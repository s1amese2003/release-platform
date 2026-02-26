package com.huatai.release.engine.sql.rule;

import com.huatai.release.engine.sql.SqlRule;
import com.huatai.release.engine.sql.SqlRuleContext;
import com.huatai.release.model.dto.SqlAuditIssue;
import com.huatai.release.model.enums.RiskLevel;

import java.util.Optional;

public class AlterRiskRule implements SqlRule {

    @Override
    public Optional<SqlAuditIssue> apply(SqlRuleContext context) {
        String normalized = context.normalized();
        if (!normalized.startsWith("ALTER TABLE")) {
            return Optional.empty();
        }
        if (normalized.contains(" DROP COLUMN ")) {
            return Optional.of(new SqlAuditIssue(context.line(), context.sql(), RiskLevel.P1,
                    "ALTER TABLE DROP COLUMN is irreversible", "Prepare backup and rollback plan"));
        }
        if (normalized.contains(" MODIFY COLUMN ") || normalized.contains(" CHANGE COLUMN ")) {
            return Optional.of(new SqlAuditIssue(context.line(), context.sql(), RiskLevel.P1,
                    "Column definition/type change may truncate data", "Validate compatibility and add rollback SQL"));
        }
        if (normalized.contains(" ADD INDEX ") && !normalized.contains("ALGORITHM=INPLACE")) {
            return Optional.of(new SqlAuditIssue(context.line(), context.sql(), RiskLevel.P1,
                    "ADD INDEX without explicit ALGORITHM=INPLACE", "Use online DDL options to reduce lock risk"));
        }
        if (normalized.contains(" ADD COLUMN ") && !normalized.contains(" DEFAULT ")) {
            return Optional.of(new SqlAuditIssue(context.line(), context.sql(), RiskLevel.P2,
                    "ADD COLUMN has no DEFAULT value", "Check if null compatibility is acceptable"));
        }
        return Optional.empty();
    }
}
