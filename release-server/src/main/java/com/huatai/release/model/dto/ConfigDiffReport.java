package com.huatai.release.model.dto;

import java.util.List;
import java.util.Map;

public record ConfigDiffReport(
        String targetEnv,
        List<ConfigDiffItem> items,
        Map<String, Long> summary
) {
}
