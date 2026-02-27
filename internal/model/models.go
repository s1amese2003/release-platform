package model

import (
	"time"

	"gorm.io/datatypes"
)

const (
	ScanStatusPending = "PENDING"
	ScanStatusRunning = "RUNNING"
	ScanStatusDone    = "DONE"
	ScanStatusFailed  = "FAILED"
)

const (
	ApprovalStatusPending    = "PENDING"
	ApprovalStatusApproved   = "APPROVED"
	ApprovalStatusRejected   = "REJECTED"
	ApprovalStatusNotRequired = "NOT_REQUIRED"
)

const (
	AdmissionPass  = "PASS"
	AdmissionWarn  = "WARN"
	AdmissionBlock = "BLOCK"
)

const (
	DeployStatusNone     = "NOT_TRIGGERED"
	DeployStatusRunning  = "RUNNING"
	DeployStatusSuccess  = "SUCCESS"
	DeployStatusFailed   = "FAILED"
)

const (
	RiskHigh   = "HIGH"
	RiskMedium = "MEDIUM"
	RiskLow    = "LOW"
)

const (
	DiffAdded    = "ADDED"
	DiffDeleted  = "DELETED"
	DiffModified = "MODIFIED"
)

type Application struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:128;uniqueIndex" json:"name"`
	Owner     string    `gorm:"size:128" json:"owner"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ReleaseTicket struct {
	ID                    uint      `gorm:"primaryKey" json:"id"`
	ApplicationID         uint      `json:"application_id"`
	Application           Application
	Environment           string    `gorm:"size:64;index:idx_release_app_env" json:"environment"`
	Version               string    `gorm:"size:128" json:"version"`
	Uploader              string    `gorm:"size:128" json:"uploader"`
	ScanStatus            string    `gorm:"size:32" json:"scan_status"`
	ApprovalStatus        string    `gorm:"size:32" json:"approval_status"`
	AdmissionResult       string    `gorm:"size:16" json:"admission_result"`
	DeployStatus          string    `gorm:"size:32" json:"deploy_status"`
	NeedsSecondaryApprove bool      `json:"needs_secondary_approve"`
	PolicyVersionID       *uint     `json:"policy_version_id"`
	PolicyVersion         *PolicyVersion
	Artifact              *Artifact
	ScanReport            *ScanReport
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

type Artifact struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	ReleaseTicketID uint      `gorm:"uniqueIndex" json:"release_ticket_id"`
	FileName        string    `gorm:"size:255" json:"file_name"`
	StoragePath     string    `gorm:"size:1024" json:"storage_path"`
	ExtractedPath   string    `gorm:"size:1024" json:"extracted_path"`
	SHA256          string    `gorm:"size:64;index" json:"sha256"`
	SizeBytes       int64     `json:"size_bytes"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type ScanReport struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	ReleaseTicketID uint     `gorm:"uniqueIndex" json:"release_ticket_id"`
	SQLPath        string    `gorm:"size:1024" json:"sql_path"`
	ConfigPath     string    `gorm:"size:1024" json:"config_path"`
	HighCount      int       `json:"high_count"`
	MediumCount    int       `json:"medium_count"`
	LowCount       int       `json:"low_count"`
	ConfigAdded    int       `json:"config_added"`
	ConfigDeleted  int       `json:"config_deleted"`
	ConfigModified int       `json:"config_modified"`
	Result         string    `gorm:"size:16" json:"result"`
	BlockReason    string    `gorm:"type:text" json:"block_reason"`
	Summary        datatypes.JSON `gorm:"type:jsonb" json:"summary"`
	SQLIssues      []SQLIssue `json:"sql_issues"`
	ConfigDiffs    []ConfigDiff `json:"config_diffs"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type SQLIssue struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	ScanReportID uint     `gorm:"index" json:"scan_report_id"`
	Sequence    int       `json:"sequence"`
	Snippet     string    `gorm:"type:text" json:"snippet"`
	RuleName    string    `gorm:"size:128" json:"rule_name"`
	RiskLevel   string    `gorm:"size:16" json:"risk_level"`
	Suggestion  string    `gorm:"type:text" json:"suggestion"`
	Whitelisted bool      `json:"whitelisted"`
	CreatedAt   time.Time `json:"created_at"`
}

type ConfigSnapshot struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	ApplicationID   uint      `gorm:"index:idx_snapshot_app_env" json:"application_id"`
	Environment     string    `gorm:"size:64;index:idx_snapshot_app_env" json:"environment"`
	SourceReleaseID *uint     `json:"source_release_id"`
	Snapshot        datatypes.JSON `gorm:"type:jsonb" json:"snapshot"`
	IsActive        bool      `gorm:"index" json:"is_active"`
	CreatedAt       time.Time `json:"created_at"`
}

type ConfigDiff struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	ScanReportID   uint      `gorm:"index" json:"scan_report_id"`
	ConfigKey      string    `gorm:"size:255" json:"config_key"`
	OnlineValue    string    `gorm:"type:text" json:"online_value"`
	CandidateValue string    `gorm:"type:text" json:"candidate_value"`
	DiffType       string    `gorm:"size:16" json:"diff_type"`
	IsCritical     bool      `json:"is_critical"`
	IsMasked       bool      `json:"is_masked"`
	CreatedAt      time.Time `json:"created_at"`
}

type ApprovalRecord struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	ReleaseTicketID uint     `gorm:"index" json:"release_ticket_id"`
	Approver       string    `gorm:"size:128" json:"approver"`
	Action         string    `gorm:"size:32" json:"action"`
	Comment        string    `gorm:"type:text" json:"comment"`
	Level          int       `json:"level"`
	CreatedAt      time.Time `json:"created_at"`
}

type PolicyVersion struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	ApplicationID *uint     `gorm:"index" json:"application_id"`
	Environment   string    `gorm:"size:64;index" json:"environment"`
	Version       string    `gorm:"size:64;uniqueIndex" json:"version"`
	Rules         datatypes.JSON `gorm:"type:jsonb" json:"rules"`
	IsActive      bool      `gorm:"index" json:"is_active"`
	CreatedBy     string    `gorm:"size:128" json:"created_by"`
	CreatedAt     time.Time `json:"created_at"`
}

type AuditLog struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Actor      string    `gorm:"size:128;index" json:"actor"`
	Action     string    `gorm:"size:128;index" json:"action"`
	EntityType string    `gorm:"size:64;index" json:"entity_type"`
	EntityID   string    `gorm:"size:128;index" json:"entity_id"`
	Detail     string    `gorm:"type:text" json:"detail"`
	CreatedAt  time.Time `json:"created_at"`
}
