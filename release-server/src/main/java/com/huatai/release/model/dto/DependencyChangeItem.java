package com.huatai.release.model.dto;

import com.huatai.release.model.enums.RiskLevel;

import java.util.List;

public record DependencyChangeItem(
        String groupId,
        String artifactId,
        String oldVersion,
        String newVersion,
        String changeType,
        List<String> cveIds,
        RiskLevel riskLevel
) {
}
