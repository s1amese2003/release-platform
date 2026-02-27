package service

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"

	"release-platform/internal/model"
	"release-platform/internal/util"

	"gopkg.in/yaml.v3"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type ScannerService struct {
	db *gorm.DB
}

type ScanOutcome struct {
	SQLPath               string
	ConfigPath            string
	SQLIssues             []model.SQLIssue
	ConfigDiffs           []model.ConfigDiff
	CandidateConfig       map[string]string
	HighCount             int
	MediumCount           int
	LowCount              int
	AddedCount            int
	DeletedCount          int
	ModifiedCount         int
	Result                string
	BlockReason           string
	NeedsSecondaryApprove bool
}

func NewScannerService(db *gorm.DB) *ScannerService {
	return &ScannerService{db: db}
}

func (s *ScannerService) ScanRelease(ctx context.Context, release *model.ReleaseTicket, extractedPath, sqlPath, configPath string) (ScanOutcome, error) {
	_ = ctx
	outcome := ScanOutcome{
		SQLPath:    sqlPath,
		ConfigPath: configPath,
	}

	if sqlPath != "" {
		issues, err := scanSQLFile(sqlPath)
		if err != nil {
			return outcome, fmt.Errorf("scan SQL file: %w", err)
		}
		outcome.SQLIssues = issues
		for _, issue := range issues {
			switch issue.RiskLevel {
			case model.RiskHigh:
				outcome.HighCount++
			case model.RiskMedium:
				outcome.MediumCount++
			case model.RiskLow:
				outcome.LowCount++
			}
		}
	}

	candidateConfig := map[string]string{}
	if configPath != "" {
		cfg, err := loadFlattenYAML(configPath)
		if err != nil {
			return outcome, fmt.Errorf("read bootstrap-dev.yml: %w", err)
		}
		candidateConfig = cfg
	}
	outcome.CandidateConfig = candidateConfig

	onlineConfig, err := s.loadOnlineSnapshot(release.ApplicationID, release.Environment)
	if err != nil {
		return outcome, err
	}

	diffs := buildConfigDiffs(onlineConfig, candidateConfig)
	outcome.ConfigDiffs = diffs
	for _, diff := range diffs {
		switch diff.DiffType {
		case model.DiffAdded:
			outcome.AddedCount++
		case model.DiffDeleted:
			outcome.DeletedCount++
		case model.DiffModified:
			outcome.ModifiedCount++
		}
	}

	result, blockReason, secondary := evaluatePolicy(outcome.SQLIssues, diffs)
	outcome.Result = result
	outcome.BlockReason = blockReason
	outcome.NeedsSecondaryApprove = secondary

	_ = extractedPath
	return outcome, nil
}

func (s *ScannerService) loadOnlineSnapshot(applicationID uint, environment string) (map[string]string, error) {
	var snapshot model.ConfigSnapshot
	err := s.db.Where("application_id = ? AND environment = ? AND is_active = ?", applicationID, environment, true).
		Order("created_at desc").
		First(&snapshot).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return map[string]string{}, nil
		}
		return nil, fmt.Errorf("query active config snapshot: %w", err)
	}

	result := map[string]string{}
	if len(snapshot.Snapshot) == 0 {
		return result, nil
	}
	if err := json.Unmarshal(snapshot.Snapshot, &result); err != nil {
		return nil, fmt.Errorf("decode snapshot JSON: %w", err)
	}
	return result, nil
}

func scanSQLFile(path string) ([]model.SQLIssue, error) {
	contentBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read SQL file: %w", err)
	}

	content := string(contentBytes)
	statements := splitSQLStatements(content)
	issues := make([]model.SQLIssue, 0)

	for i, stmt := range statements {
		issue := evaluateSQLStatement(i+1, stmt)
		if issue == nil {
			continue
		}
		issues = append(issues, *issue)
	}

	if !strings.Contains(strings.ToLower(content), "rollback") {
		issues = append(issues, model.SQLIssue{
			Sequence:   0,
			Snippet:    "-- no rollback description",
			RuleName:   "MISSING_ROLLBACK_NOTE",
			RiskLevel:  model.RiskLow,
			Suggestion: "Document rollback strategy in SQL comments.",
		})
	}

	return issues, nil
}

func splitSQLStatements(content string) []string {
	lines := strings.Split(content, "\n")
	cleaned := make([]string, 0, len(lines))
	for _, line := range lines {
		cleanLine := strings.ReplaceAll(line, "\uFEFF", "")
		trimmed := strings.TrimSpace(cleanLine)
		if strings.HasPrefix(trimmed, "--") || strings.HasPrefix(trimmed, "#") {
			continue
		}
		cleaned = append(cleaned, cleanLine)
	}
	joined := strings.Join(cleaned, "\n")
	rawStatements := strings.Split(joined, ";")
	result := make([]string, 0, len(rawStatements))
	for _, stmt := range rawStatements {
		trimmed := strings.TrimSpace(stmt)
		if trimmed == "" {
			continue
		}
		result = append(result, trimmed)
	}
	return result
}

func evaluateSQLStatement(seq int, stmt string) *model.SQLIssue {
	cleanStmt := strings.ReplaceAll(stmt, "\uFEFF", "")
	normalized := strings.ToLower(strings.TrimSpace(cleanStmt))
	if normalized == "" {
		return nil
	}
	if strings.HasPrefix(normalized, "--") || strings.HasPrefix(normalized, "#") {
		return nil
	}

	highRules := []struct {
		name       string
		pattern    *regexp.Regexp
		suggestion string
	}{
		{name: "DROP_DATABASE", pattern: regexp.MustCompile(`(?i)\bdrop\s+database\b`), suggestion: "Do not run DROP DATABASE in release scripts."},
		{name: "DROP_TABLE", pattern: regexp.MustCompile(`(?i)\bdrop\s+table\b`), suggestion: "Avoid DROP TABLE. Prefer phased decommission and backup."},
		{name: "TRUNCATE_TABLE", pattern: regexp.MustCompile(`(?i)\btruncate\s+table\b`), suggestion: "TRUNCATE TABLE is destructive and should be reviewed manually."},
		{name: "ALTER_DROP_COLUMN", pattern: regexp.MustCompile(`(?i)\balter\s+table\b[\s\S]*\bdrop\s+column\b`), suggestion: "Drop-column migrations require strict review and fallback plan."},
	}
	for _, rule := range highRules {
		if rule.pattern.MatchString(stmt) {
			return &model.SQLIssue{Sequence: seq, Snippet: stmt, RuleName: rule.name, RiskLevel: model.RiskHigh, Suggestion: rule.suggestion}
		}
	}

	if regexp.MustCompile(`(?i)^\s*delete\s+from\b`).MatchString(stmt) && !regexp.MustCompile(`(?i)\bwhere\b`).MatchString(stmt) {
		return &model.SQLIssue{Sequence: seq, Snippet: stmt, RuleName: "DELETE_WITHOUT_WHERE", RiskLevel: model.RiskHigh, Suggestion: "Add WHERE clause or convert to batched deletes."}
	}
	if regexp.MustCompile(`(?i)^\s*update\s+\S+\s+set\b`).MatchString(stmt) && !regexp.MustCompile(`(?i)\bwhere\b`).MatchString(stmt) {
		return &model.SQLIssue{Sequence: seq, Snippet: stmt, RuleName: "UPDATE_WITHOUT_WHERE", RiskLevel: model.RiskHigh, Suggestion: "Add WHERE clause and validate affected row scope."}
	}

	mediumRules := []struct {
		name       string
		pattern    *regexp.Regexp
		suggestion string
	}{
		{name: "RENAME_TABLE", pattern: regexp.MustCompile(`(?i)\brename\s+table\b`), suggestion: "Validate dependent services and scheduled jobs before rename."},
		{name: "REBUILD_INDEX", pattern: regexp.MustCompile(`(?i)\b(create|drop)\s+index\b`), suggestion: "Large index rebuild can lock table and impact SLA."},
	}
	for _, rule := range mediumRules {
		if rule.pattern.MatchString(stmt) {
			return &model.SQLIssue{Sequence: seq, Snippet: stmt, RuleName: rule.name, RiskLevel: model.RiskMedium, Suggestion: rule.suggestion}
		}
	}

	if regexp.MustCompile(`(?i)^\s*create\s+table\b`).MatchString(stmt) && !regexp.MustCompile(`(?i)\bif\s+not\s+exists\b`).MatchString(stmt) {
		return &model.SQLIssue{Sequence: seq, Snippet: stmt, RuleName: "CREATE_TABLE_WITHOUT_IF_NOT_EXISTS", RiskLevel: model.RiskLow, Suggestion: "Use IF NOT EXISTS for idempotent migrations."}
	}

	if regexp.MustCompile(`(?i)^\s*drop\s+table\b`).MatchString(stmt) && !regexp.MustCompile(`(?i)\bif\s+exists\b`).MatchString(stmt) {
		return &model.SQLIssue{Sequence: seq, Snippet: stmt, RuleName: "DROP_TABLE_WITHOUT_IF_EXISTS", RiskLevel: model.RiskLow, Suggestion: "Use IF EXISTS and attach explicit backup plan."}
	}

	return nil
}

func loadFlattenYAML(path string) (map[string]string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	raw := map[string]any{}
	if err := yaml.Unmarshal(content, &raw); err != nil {
		return nil, err
	}

	return util.FlattenMap(raw), nil
}

func buildConfigDiffs(online map[string]string, candidate map[string]string) []model.ConfigDiff {
	allKeys := map[string]struct{}{}
	for key := range online {
		allKeys[key] = struct{}{}
	}
	for key := range candidate {
		allKeys[key] = struct{}{}
	}

	keys := make([]string, 0, len(allKeys))
	for key := range allKeys {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	diffs := make([]model.ConfigDiff, 0)
	for _, key := range keys {
		onlineVal, onlineOK := online[key]
		candidateVal, candidateOK := candidate[key]
		if onlineOK && candidateOK && onlineVal == candidateVal {
			continue
		}

		diffType := model.DiffModified
		switch {
		case !onlineOK && candidateOK:
			diffType = model.DiffAdded
		case onlineOK && !candidateOK:
			diffType = model.DiffDeleted
		}

		isSensitive := util.IsSensitiveKey(key)
		onlineDisplay := onlineVal
		candidateDisplay := candidateVal
		if isSensitive {
			onlineDisplay = util.MaskValue(onlineVal)
			candidateDisplay = util.MaskValue(candidateVal)
		}

		diffs = append(diffs, model.ConfigDiff{
			ConfigKey:      key,
			OnlineValue:    onlineDisplay,
			CandidateValue: candidateDisplay,
			DiffType:       diffType,
			IsCritical:     util.IsCriticalConfigKey(key),
			IsMasked:       isSensitive,
		})
	}

	return diffs
}

func evaluatePolicy(sqlIssues []model.SQLIssue, configDiffs []model.ConfigDiff) (result, blockReason string, secondary bool) {
	hasHigh := false
	hasMedium := false
	hasLowOnly := true
	for _, issue := range sqlIssues {
		switch issue.RiskLevel {
		case model.RiskHigh:
			hasHigh = true
			hasLowOnly = false
		case model.RiskMedium:
			hasMedium = true
			hasLowOnly = false
		case model.RiskLow:
		default:
			hasLowOnly = false
		}
	}

	criticalConfigChanged := false
	for _, diff := range configDiffs {
		if diff.IsCritical {
			criticalConfigChanged = true
			break
		}
	}

	if hasHigh {
		return model.AdmissionBlock, "HIGH risk SQL operation detected", false
	}
	if criticalConfigChanged {
		return model.AdmissionWarn, "Critical configuration changed, secondary approval required", true
	}
	if hasMedium {
		return model.AdmissionWarn, "Medium risk SQL operation requires manual review", false
	}
	if hasLowOnly || len(sqlIssues) == 0 {
		return model.AdmissionPass, "", false
	}
	return model.AdmissionWarn, "Policy fallback to WARN", false
}

func encodeConfigSnapshot(config map[string]string) datatypes.JSON {
	if config == nil {
		return datatypes.JSON([]byte("{}"))
	}
	payload, err := json.Marshal(config)
	if err != nil {
		return datatypes.JSON([]byte("{}"))
	}
	return datatypes.JSON(payload)
}
