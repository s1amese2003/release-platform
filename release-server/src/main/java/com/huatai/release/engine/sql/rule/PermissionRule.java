package com.huatai.release.engine.sql.rule;

import com.huatai.release.engine.sql.SqlRule;
import com.huatai.release.engine.sql.SqlRuleContext;
import com.huatai.release.model.dto.SqlAuditIssue;
import com.huatai.release.model.enums.RiskLevel;

import java.util.Optional;

public class PermissionRule implements SqlRule {

    @Override
    public Optional<SqlAuditIssue> apply(SqlRuleContext context) {
        String normalized = context.normalized();
        if (normalized.startsWith("GRANT ") || normalized.startsWith("REVOKE ")) {
            return Optional.of(new SqlAuditIssue(context.line(), context.sql(), RiskLevel.P1,
                    "Permission statement detected", "Verify least-privilege and keep approval record"));
        }
        return Optional.empty();
    }
}
