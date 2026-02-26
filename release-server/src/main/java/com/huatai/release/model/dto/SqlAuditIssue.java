package com.huatai.release.model.dto;

import com.huatai.release.model.enums.RiskLevel;

public record SqlAuditIssue(
        int line,
        String sql,
        RiskLevel level,
        String reason,
        String suggestion
) {
}
