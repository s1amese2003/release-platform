package com.huatai.release.model.dto;

import com.huatai.release.model.enums.DiffType;
import com.huatai.release.model.enums.RiskLevel;

public record ConfigDiffItem(
        String configFile,
        String configKey,
        String packageValue,
        String baselineValue,
        DiffType diffType,
        RiskLevel riskLevel,
        String category,
        String reason
) {
}
