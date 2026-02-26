package com.huatai.release.engine.sql.rule;

import com.huatai.release.engine.sql.SqlRule;
import com.huatai.release.engine.sql.SqlRuleContext;
import com.huatai.release.model.dto.SqlAuditIssue;
import com.huatai.release.model.enums.RiskLevel;

import java.util.Optional;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

public class CrossDatabaseRule implements SqlRule {

    private static final Pattern INSERT_DB_PATTERN = Pattern.compile("(?i)^\\s*INSERT\\s+INTO\\s+`?([a-zA-Z0-9_]+)`?\\.");

    @Override
    public Optional<SqlAuditIssue> apply(SqlRuleContext context) {
        Matcher matcher = INSERT_DB_PATTERN.matcher(context.sql());
        if (!matcher.find()) {
            return Optional.empty();
        }
        String targetDb = matcher.group(1);
        if (targetDb.equalsIgnoreCase(context.businessDb())) {
            return Optional.empty();
        }
        return Optional.of(new SqlAuditIssue(context.line(), context.sql(), RiskLevel.P2,
                "Cross-database write: target DB " + targetDb + " != business DB " + context.businessDb(),
                "Confirm privilege and release with the target-service owner"));
    }
}
