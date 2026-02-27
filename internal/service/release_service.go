package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"release-platform/internal/model"
	"release-platform/internal/queue"

	"github.com/hibiken/asynq"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type ReleaseService struct {
	db          *gorm.DB
	artifactSvc *ArtifactService
	scannerSvc  *ScannerService
	queueClient *asynq.Client
	queueName   string
}

type CreateReleaseInput struct {
	Application string `json:"application"`
	Owner       string `json:"owner"`
	Environment string `json:"environment"`
	Version     string `json:"version"`
	Uploader    string `json:"uploader"`
}

type ApproveInput struct {
	Approver string `json:"approver"`
	Action   string `json:"action"`
	Comment  string `json:"comment"`
	Level    int    `json:"level"`
}

type DeployInput struct {
	Operator string `json:"operator"`
}

func NewReleaseService(db *gorm.DB, artifactSvc *ArtifactService, scannerSvc *ScannerService, queueClient *asynq.Client, queueName string) *ReleaseService {
	return &ReleaseService{
		db:          db,
		artifactSvc: artifactSvc,
		scannerSvc:  scannerSvc,
		queueClient: queueClient,
		queueName:   queueName,
	}
}

func (s *ReleaseService) CreateRelease(ctx context.Context, in CreateReleaseInput) (*model.ReleaseTicket, error) {
	if in.Application == "" || in.Environment == "" || in.Version == "" || in.Uploader == "" {
		return nil, errors.New("application, environment, version and uploader are required")
	}

	var app model.Application
	err := s.db.WithContext(ctx).Where("name = ?", in.Application).First(&app).Error
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("query application: %w", err)
		}
		app = model.Application{Name: in.Application, Owner: in.Owner}
		if app.Owner == "" {
			app.Owner = in.Uploader
		}
		if err := s.db.WithContext(ctx).Create(&app).Error; err != nil {
			return nil, fmt.Errorf("create application: %w", err)
		}
	}

	ticket := model.ReleaseTicket{
		ApplicationID:   app.ID,
		Environment:     in.Environment,
		Version:         in.Version,
		Uploader:        in.Uploader,
		ScanStatus:      model.ScanStatusPending,
		ApprovalStatus:  model.ApprovalStatusPending,
		AdmissionResult: model.AdmissionWarn,
		DeployStatus:    model.DeployStatusNone,
	}

	policy, err := s.resolvePolicy(ctx, app.ID, in.Environment)
	if err != nil {
		return nil, err
	}
	if policy != nil {
		ticket.PolicyVersionID = &policy.ID
	}

	if err := s.db.WithContext(ctx).Create(&ticket).Error; err != nil {
		return nil, fmt.Errorf("create release ticket: %w", err)
	}

	s.audit(ctx, in.Uploader, "create_release", "release_ticket", ticket.ID, fmt.Sprintf("app=%s env=%s version=%s", in.Application, in.Environment, in.Version))

	if err := s.db.WithContext(ctx).Preload("Application").First(&ticket, ticket.ID).Error; err != nil {
		return nil, fmt.Errorf("reload release ticket: %w", err)
	}
	return &ticket, nil
}

func (s *ReleaseService) ListReleases(ctx context.Context, environment, admissionResult string) ([]model.ReleaseTicket, error) {
	query := s.db.WithContext(ctx).Model(&model.ReleaseTicket{}).
		Preload("Application").
		Order("created_at desc")

	if environment != "" {
		query = query.Where("environment = ?", environment)
	}
	if admissionResult != "" {
		query = query.Where("admission_result = ?", admissionResult)
	}

	var releases []model.ReleaseTicket
	if err := query.Find(&releases).Error; err != nil {
		return nil, fmt.Errorf("list releases: %w", err)
	}
	return releases, nil
}

func (s *ReleaseService) UploadArtifact(ctx context.Context, releaseID uint, fileName string, storagePath string, extractedPath string, size int64, sha string) (*model.Artifact, error) {
	var ticket model.ReleaseTicket
	if err := s.db.WithContext(ctx).First(&ticket, releaseID).Error; err != nil {
		return nil, fmt.Errorf("query release ticket: %w", err)
	}

	artifact := model.Artifact{
		ReleaseTicketID: releaseID,
		FileName:        fileName,
		StoragePath:     storagePath,
		ExtractedPath:   extractedPath,
		SHA256:          sha,
		SizeBytes:       size,
	}

	if err := s.db.WithContext(ctx).Where("release_ticket_id = ?", releaseID).Delete(&model.Artifact{}).Error; err != nil {
		return nil, fmt.Errorf("remove old artifact: %w", err)
	}
	if err := s.db.WithContext(ctx).Create(&artifact).Error; err != nil {
		return nil, fmt.Errorf("save artifact: %w", err)
	}

	s.audit(ctx, ticket.Uploader, "upload_artifact", "release_ticket", releaseID, fmt.Sprintf("artifact=%s sha=%s", fileName, sha))
	return &artifact, nil
}

func (s *ReleaseService) EnqueueScan(ctx context.Context, releaseID uint, operator string) error {
	if s.queueClient == nil {
		if err := s.RunScan(ctx, releaseID); err != nil {
			return err
		}
		s.audit(ctx, operator, "trigger_scan_sync", "release_ticket", releaseID, "scan executed synchronously")
		return nil
	}

	var ticket model.ReleaseTicket
	if err := s.db.WithContext(ctx).First(&ticket, releaseID).Error; err != nil {
		return fmt.Errorf("query release ticket: %w", err)
	}

	task, err := queue.NewScanTask(releaseID, s.queueName)
	if err != nil {
		return err
	}
	if _, err := s.queueClient.EnqueueContext(ctx, task, asynq.MaxRetry(3), asynq.Timeout(10*time.Minute)); err != nil {
		return fmt.Errorf("enqueue scan task: %w", err)
	}

	if err := s.db.WithContext(ctx).Model(&model.ReleaseTicket{}).
		Where("id = ?", releaseID).
		Updates(map[string]any{"scan_status": model.ScanStatusPending}).Error; err != nil {
		return fmt.Errorf("update release scan status: %w", err)
	}

	s.audit(ctx, operator, "trigger_scan", "release_ticket", releaseID, "scan task enqueued")
	return nil
}

func (s *ReleaseService) RunScan(ctx context.Context, releaseID uint) error {
	var ticket model.ReleaseTicket
	if err := s.db.WithContext(ctx).Preload("Artifact").First(&ticket, releaseID).Error; err != nil {
		return fmt.Errorf("load ticket for scan: %w", err)
	}
	if ticket.Artifact == nil {
		return errors.New("artifact not uploaded")
	}

	if err := s.db.WithContext(ctx).Model(&model.ReleaseTicket{}).Where("id = ?", releaseID).
		Updates(map[string]any{"scan_status": model.ScanStatusRunning}).Error; err != nil {
		return fmt.Errorf("mark scan running: %w", err)
	}

	sqlPath := ""
	if located, err := s.artifactSvc.LocateLatestUpgradeSQL(ticket.Artifact.ExtractedPath); err == nil {
		sqlPath = located
	}
	configPath := ""
	if located, err := s.artifactSvc.LocateBootstrapDevYAML(ticket.Artifact.ExtractedPath); err == nil {
		configPath = located
	}

	outcome, err := s.scannerSvc.ScanRelease(ctx, &ticket, ticket.Artifact.ExtractedPath, sqlPath, configPath)
	if err != nil {
		s.db.WithContext(ctx).Model(&model.ReleaseTicket{}).Where("id = ?", releaseID).
			Updates(map[string]any{"scan_status": model.ScanStatusFailed})
		return err
	}

	if sqlPath == "" {
		outcome.SQLIssues = append(outcome.SQLIssues, model.SQLIssue{
			Sequence:   0,
			Snippet:    "upgrade sql not found",
			RuleName:   "MISSING_UPGRADE_SQL",
			RiskLevel:  model.RiskMedium,
			Suggestion: "Provide BOOT-INF/classes/upgrade/YYYYMMDD/upgrade.sql",
		})
		outcome.MediumCount++
		if outcome.Result == model.AdmissionPass {
			outcome.Result = model.AdmissionWarn
			outcome.BlockReason = "No upgrade.sql found"
		}
	}

	summaryObj := map[string]any{
		"candidate_config": outcome.CandidateConfig,
		"scanned_at":       time.Now().UTC().Format(time.RFC3339),
	}
	summaryBytes, _ := json.Marshal(summaryObj)

	if err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var existing model.ScanReport
		if err := tx.Where("release_ticket_id = ?", releaseID).First(&existing).Error; err == nil {
			if err := tx.Where("scan_report_id = ?", existing.ID).Delete(&model.SQLIssue{}).Error; err != nil {
				return err
			}
			if err := tx.Where("scan_report_id = ?", existing.ID).Delete(&model.ConfigDiff{}).Error; err != nil {
				return err
			}
			if err := tx.Delete(&existing).Error; err != nil {
				return err
			}
		}

		report := model.ScanReport{
			ReleaseTicketID: releaseID,
			SQLPath:         outcome.SQLPath,
			ConfigPath:      outcome.ConfigPath,
			HighCount:       outcome.HighCount,
			MediumCount:     outcome.MediumCount,
			LowCount:        outcome.LowCount,
			ConfigAdded:     outcome.AddedCount,
			ConfigDeleted:   outcome.DeletedCount,
			ConfigModified:  outcome.ModifiedCount,
			Result:          outcome.Result,
			BlockReason:     outcome.BlockReason,
			Summary:         datatypes.JSON(summaryBytes),
		}
		if err := tx.Create(&report).Error; err != nil {
			return err
		}

		for i := range outcome.SQLIssues {
			outcome.SQLIssues[i].ScanReportID = report.ID
		}
		for i := range outcome.ConfigDiffs {
			outcome.ConfigDiffs[i].ScanReportID = report.ID
		}
		if len(outcome.SQLIssues) > 0 {
			if err := tx.Create(&outcome.SQLIssues).Error; err != nil {
				return err
			}
		}
		if len(outcome.ConfigDiffs) > 0 {
			if err := tx.Create(&outcome.ConfigDiffs).Error; err != nil {
				return err
			}
		}

		approvalStatus := model.ApprovalStatusPending
		if outcome.Result == model.AdmissionPass && !outcome.NeedsSecondaryApprove {
			approvalStatus = model.ApprovalStatusNotRequired
		}

		if err := tx.Model(&model.ReleaseTicket{}).Where("id = ?", releaseID).
			Updates(map[string]any{
				"scan_status":             model.ScanStatusDone,
				"admission_result":        outcome.Result,
				"approval_status":         approvalStatus,
				"needs_secondary_approve": outcome.NeedsSecondaryApprove,
			}).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		s.db.WithContext(ctx).Model(&model.ReleaseTicket{}).Where("id = ?", releaseID).
			Updates(map[string]any{"scan_status": model.ScanStatusFailed})
		return fmt.Errorf("persist scan report: %w", err)
	}

	s.audit(ctx, ticket.Uploader, "scan_release", "release_ticket", releaseID, fmt.Sprintf("result=%s", outcome.Result))
	return nil
}

func (s *ReleaseService) GetReleaseReport(ctx context.Context, releaseID uint) (*model.ReleaseTicket, error) {
	var ticket model.ReleaseTicket
	if err := s.db.WithContext(ctx).
		Preload("Application").
		Preload("Artifact").
		Preload("ScanReport").
		Preload("ScanReport.SQLIssues").
		Preload("ScanReport.ConfigDiffs").
		Preload("PolicyVersion").
		First(&ticket, releaseID).Error; err != nil {
		return nil, fmt.Errorf("get report: %w", err)
	}
	return &ticket, nil
}

func (s *ReleaseService) Approve(ctx context.Context, releaseID uint, in ApproveInput) error {
	if in.Approver == "" || in.Action == "" {
		return errors.New("approver and action are required")
	}

	var ticket model.ReleaseTicket
	if err := s.db.WithContext(ctx).First(&ticket, releaseID).Error; err != nil {
		return fmt.Errorf("query ticket: %w", err)
	}

	action := normalizeApprovalAction(in.Action)
	if action == "" {
		return errors.New("action must be approve or reject")
	}
	if ticket.AdmissionResult == model.AdmissionBlock && action == "APPROVE" {
		return errors.New("blocked release cannot be approved directly")
	}

	newStatus := model.ApprovalStatusRejected
	if action == "APPROVE" {
		if ticket.NeedsSecondaryApprove && in.Level < 2 {
			return errors.New("secondary approval is required (level>=2)")
		}
		newStatus = model.ApprovalStatusApproved
	}

	record := model.ApprovalRecord{
		ReleaseTicketID: releaseID,
		Approver:        in.Approver,
		Action:          action,
		Comment:         in.Comment,
		Level:           in.Level,
	}
	if err := s.db.WithContext(ctx).Create(&record).Error; err != nil {
		return fmt.Errorf("save approval record: %w", err)
	}

	if err := s.db.WithContext(ctx).Model(&model.ReleaseTicket{}).
		Where("id = ?", releaseID).
		Updates(map[string]any{"approval_status": newStatus}).Error; err != nil {
		return fmt.Errorf("update approval status: %w", err)
	}

	s.audit(ctx, in.Approver, "approve_release", "release_ticket", releaseID, fmt.Sprintf("action=%s level=%d", action, in.Level))
	return nil
}

func (s *ReleaseService) Deploy(ctx context.Context, releaseID uint, in DeployInput) error {
	if in.Operator == "" {
		return errors.New("operator is required")
	}

	var ticket model.ReleaseTicket
	if err := s.db.WithContext(ctx).
		Preload("Application").
		Preload("ScanReport").
		First(&ticket, releaseID).Error; err != nil {
		return fmt.Errorf("query ticket: %w", err)
	}
	if ticket.ScanStatus != model.ScanStatusDone {
		return errors.New("scan not completed")
	}
	if ticket.AdmissionResult == model.AdmissionBlock {
		return errors.New("blocked release cannot be deployed")
	}
	if ticket.ApprovalStatus != model.ApprovalStatusApproved && ticket.ApprovalStatus != model.ApprovalStatusNotRequired {
		return errors.New("release is not approved")
	}

	if err := s.db.WithContext(ctx).Model(&model.ReleaseTicket{}).Where("id = ?", releaseID).
		Updates(map[string]any{"deploy_status": model.DeployStatusRunning}).Error; err != nil {
		return fmt.Errorf("mark deploy running: %w", err)
	}

	if err := s.persistSnapshotFromReport(ctx, &ticket); err != nil {
		s.db.WithContext(ctx).Model(&model.ReleaseTicket{}).Where("id = ?", releaseID).
			Updates(map[string]any{"deploy_status": model.DeployStatusFailed})
		return err
	}

	if err := s.db.WithContext(ctx).Model(&model.ReleaseTicket{}).Where("id = ?", releaseID).
		Updates(map[string]any{"deploy_status": model.DeployStatusSuccess}).Error; err != nil {
		return fmt.Errorf("mark deploy success: %w", err)
	}

	s.audit(ctx, in.Operator, "deploy_release", "release_ticket", releaseID, "deployment callback simulated")
	return nil
}

func (s *ReleaseService) persistSnapshotFromReport(ctx context.Context, ticket *model.ReleaseTicket) error {
	if ticket.ScanReport == nil {
		return errors.New("scan report missing")
	}

	summary := map[string]any{}
	if len(ticket.ScanReport.Summary) > 0 {
		if err := json.Unmarshal(ticket.ScanReport.Summary, &summary); err != nil {
			return fmt.Errorf("decode scan summary: %w", err)
		}
	}

	candidate := map[string]string{}
	if raw, ok := summary["candidate_config"]; ok {
		if mapValue, castOK := raw.(map[string]any); castOK {
			for key, value := range mapValue {
				candidate[key] = fmt.Sprint(value)
			}
		}
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.ConfigSnapshot{}).
			Where("application_id = ? AND environment = ? AND is_active = ?", ticket.ApplicationID, ticket.Environment, true).
			Updates(map[string]any{"is_active": false}).Error; err != nil {
			return err
		}
		payload, _ := json.Marshal(candidate)
		snapshot := model.ConfigSnapshot{
			ApplicationID:   ticket.ApplicationID,
			Environment:     ticket.Environment,
			SourceReleaseID: &ticket.ID,
			Snapshot:        datatypes.JSON(payload),
			IsActive:        true,
		}
		if err := tx.Create(&snapshot).Error; err != nil {
			return err
		}
		return nil
	})
}

func (s *ReleaseService) GetPolicies(ctx context.Context, applicationID *uint, environment string) ([]model.PolicyVersion, error) {
	query := s.db.WithContext(ctx).Model(&model.PolicyVersion{}).Order("created_at desc")
	if applicationID != nil {
		query = query.Where("application_id = ?", *applicationID)
	}
	if environment != "" {
		query = query.Where("environment = ?", environment)
	}
	var policies []model.PolicyVersion
	if err := query.Find(&policies).Error; err != nil {
		return nil, fmt.Errorf("query policies: %w", err)
	}
	return policies, nil
}

func (s *ReleaseService) UpdatePolicy(ctx context.Context, policyID uint, rules json.RawMessage, isActive bool, operator string) error {
	if len(rules) == 0 {
		return errors.New("rules cannot be empty")
	}

	updates := map[string]any{
		"rules":     datatypes.JSON(rules),
		"is_active": isActive,
	}
	if err := s.db.WithContext(ctx).Model(&model.PolicyVersion{}).Where("id = ?", policyID).Updates(updates).Error; err != nil {
		return fmt.Errorf("update policy: %w", err)
	}

	s.audit(ctx, operator, "update_policy", "policy_version", policyID, string(rules))
	return nil
}

func (s *ReleaseService) resolvePolicy(ctx context.Context, applicationID uint, environment string) (*model.PolicyVersion, error) {
	var policy model.PolicyVersion
	err := s.db.WithContext(ctx).
		Where("is_active = ? AND (application_id = ? OR application_id IS NULL) AND (environment = ? OR environment = ?)", true, applicationID, environment, "default").
		Order("application_id desc, created_at desc").
		First(&policy).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("resolve policy: %w", err)
	}
	return &policy, nil
}

func (s *ReleaseService) audit(ctx context.Context, actor, action, entityType string, entityID uint, detail string) {
	_ = s.db.WithContext(ctx).Create(&model.AuditLog{
		Actor:      actor,
		Action:     action,
		EntityType: entityType,
		EntityID:   fmt.Sprintf("%d", entityID),
		Detail:     detail,
	}).Error
}

func normalizeApprovalAction(action string) string {
	switch strings.ToLower(action) {
	case "approve", "approved", "pass":
		return "APPROVE"
	case "reject", "rejected", "deny":
		return "REJECT"
	default:
		return ""
	}
}
