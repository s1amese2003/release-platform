package com.huatai.release.engine.sql.rule;

import com.huatai.release.engine.sql.SqlRule;
import com.huatai.release.engine.sql.SqlRuleContext;
import com.huatai.release.model.dto.SqlAuditIssue;
import com.huatai.release.model.enums.RiskLevel;

import java.util.Optional;

public class CreateIfNotExistsRule implements SqlRule {

    @Override
    public Optional<SqlAuditIssue> apply(SqlRuleContext context) {
        String normalized = context.normalized();
        if (normalized.startsWith("CREATE TABLE") && normalized.contains("IF NOT EXISTS")) {
            return Optional.of(new SqlAuditIssue(context.line(), context.sql(), RiskLevel.P2,
                    "CREATE TABLE IF NOT EXISTS", "纭骞傜瓑鍒涘缓閫昏緫鏄惁绗﹀悎棰勬湡"));
        }
        return Optional.empty();
    }
}
