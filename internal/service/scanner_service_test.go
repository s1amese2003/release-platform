package service

import (
	"testing"

	"release-platform/internal/model"
)

func TestEvaluateSQLStatementHighRisk(t *testing.T) {
	issue := evaluateSQLStatement(1, "DELETE FROM users")
	if issue == nil {
		t.Fatalf("expected issue, got nil")
	}
	if issue.RiskLevel != model.RiskHigh {
		t.Fatalf("expected HIGH, got %s", issue.RiskLevel)
	}
}

func TestEvaluateSQLStatementLowRisk(t *testing.T) {
	issue := evaluateSQLStatement(1, "CREATE TABLE t(id bigint)")
	if issue == nil {
		t.Fatalf("expected issue, got nil")
	}
	if issue.RiskLevel != model.RiskLow {
		t.Fatalf("expected LOW, got %s", issue.RiskLevel)
	}
}
