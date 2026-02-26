package com.huatai.release.model.dto;

import java.util.List;

public record ReleaseCheckReport(
        PackageParseResult parse,
        SqlAuditReport sqlAudit,
        ConfigDiffReport configDiff,
        DependencyReport dependency,
        List<ManualOperationItem> manualChecklist,
        boolean rejectedByP0
) {
}
