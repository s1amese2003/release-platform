package com.huatai.release.engine.sql.rule;

import com.huatai.release.engine.sql.SqlRule;
import com.huatai.release.engine.sql.SqlRuleContext;
import com.huatai.release.model.dto.SqlAuditIssue;
import com.huatai.release.model.enums.RiskLevel;

import java.util.Optional;

public class SystemSchemaRule implements SqlRule {

    @Override
    public Optional<SqlAuditIssue> apply(SqlRuleContext context) {
        String normalized = context.normalized();
        if (normalized.contains("MYSQL.") || normalized.contains("INFORMATION_SCHEMA.") || normalized.contains("PERFORMANCE_SCHEMA.")) {
            return Optional.of(new SqlAuditIssue(context.line(), context.sql(), RiskLevel.P0,
                    "System schema operation is forbidden", "Only operate business schemas"));
        }
        return Optional.empty();
    }
}
