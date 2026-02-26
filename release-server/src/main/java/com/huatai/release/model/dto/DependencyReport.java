package com.huatai.release.model.dto;

import java.util.List;

public record DependencyReport(
        int totalDependencies,
        int highRiskDependencies,
        List<DependencyChangeItem> changes
) {
}
