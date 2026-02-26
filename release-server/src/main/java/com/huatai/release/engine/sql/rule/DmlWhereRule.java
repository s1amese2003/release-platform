package com.huatai.release.engine.sql.rule;

import com.huatai.release.engine.sql.SqlRule;
import com.huatai.release.engine.sql.SqlRuleContext;
import com.huatai.release.model.dto.SqlAuditIssue;
import com.huatai.release.model.enums.RiskLevel;

import java.util.Optional;

public class DmlWhereRule implements SqlRule {

    @Override
    public Optional<SqlAuditIssue> apply(SqlRuleContext context) {
        String normalized = context.normalized();
        if (normalized.startsWith("DELETE FROM") && !normalized.contains(" WHERE ")) {
            return Optional.of(new SqlAuditIssue(context.line(), context.sql(), RiskLevel.P0,
                    "DELETE statement missing WHERE clause", "Add WHERE condition to avoid full-table delete"));
        }
        if (normalized.startsWith("UPDATE ") && !normalized.contains(" WHERE ")) {
            return Optional.of(new SqlAuditIssue(context.line(), context.sql(), RiskLevel.P0,
                    "UPDATE statement missing WHERE clause", "Add WHERE condition to avoid full-table update"));
        }
        return Optional.empty();
    }
}
