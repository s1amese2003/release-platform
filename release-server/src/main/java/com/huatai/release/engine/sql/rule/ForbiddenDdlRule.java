package com.huatai.release.engine.sql.rule;

import com.huatai.release.engine.sql.SqlRule;
import com.huatai.release.engine.sql.SqlRuleContext;
import com.huatai.release.model.dto.SqlAuditIssue;
import com.huatai.release.model.enums.RiskLevel;

import java.util.Optional;

public class ForbiddenDdlRule implements SqlRule {

    @Override
    public Optional<SqlAuditIssue> apply(SqlRuleContext context) {
        String normalized = context.normalized();
        if (normalized.startsWith("DROP TABLE") || normalized.startsWith("DROP DATABASE")) {
            return Optional.of(new SqlAuditIssue(context.line(), context.sql(), RiskLevel.P0,
                    "DROP is forbidden", "Use archive/logical deletion instead of destructive DDL"));
        }
        if (normalized.startsWith("TRUNCATE TABLE")) {
            return Optional.of(new SqlAuditIssue(context.line(), context.sql(), RiskLevel.P0,
                    "TRUNCATE TABLE is forbidden", "Use conditional DELETE with backup"));
        }
        return Optional.empty();
    }
}
