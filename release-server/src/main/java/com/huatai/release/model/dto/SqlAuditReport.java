package com.huatai.release.model.dto;

import java.util.List;
import java.util.Map;

public record SqlAuditReport(
        String version,
        int totalStatements,
        List<SqlAuditIssue> results,
        Map<String, Long> summary
) {
}
