package com.huatai.release.model.dto;

import java.time.LocalDateTime;
import java.util.List;
import java.util.Map;

public record PackageParseResult(
        String requestNo,
        String appName,
        String appVersion,
        String targetEnv,
        LocalDateTime buildTime,
        String buildJdk,
        String packagePath,
        String packageMd5,
        List<String> upgradeVersions,
        String latestUpgradeVersion,
        String sqlFilePath,
        String sqlContent,
        String configTxt,
        List<String> manualOperations,
        Map<String, String> configFiles,
        List<DependencyCoordinate> dependencies
) {
}
