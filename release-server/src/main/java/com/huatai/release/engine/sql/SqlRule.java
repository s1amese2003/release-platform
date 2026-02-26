package com.huatai.release.engine.sql;

import com.huatai.release.model.dto.SqlAuditIssue;

import java.util.Optional;

public interface SqlRule {

    Optional<SqlAuditIssue> apply(SqlRuleContext context);
}
