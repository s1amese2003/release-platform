package com.huatai.release.engine.config;

import com.huatai.release.model.dto.ConfigDiffItem;
import com.huatai.release.model.dto.ConfigDiffReport;
import com.huatai.release.model.enums.DiffType;
import com.huatai.release.model.enums.RiskLevel;
import org.springframework.stereotype.Component;

import java.util.ArrayList;
import java.util.LinkedHashMap;
import java.util.LinkedHashSet;
import java.util.List;
import java.util.Map;
import java.util.Set;

@Component
public class ConfigDiffEngine {

    private static final Set<String> PROD_FORBIDDEN_TRUE = Set.of(
            "swagger.enabled",
            "springdoc.api-docs.enabled",
            "debug.iceCmdb.systemNames",
            "debug.workflowMaintenace.adjustParticipant"
    );

    public ConfigDiffReport diff(Map<String, String> packageConfig,
                                 Map<String, String> baseline,
                                 String targetEnv,
                                 String configFile) {
        List<ConfigDiffItem> items = new ArrayList<>();

        Set<String> allKeys = new LinkedHashSet<>();
        allKeys.addAll(packageConfig.keySet());
        allKeys.addAll(baseline.keySet());

        for (String key : allKeys) {
            String packageValue = packageConfig.get(key);
            String baselineValue = baseline.get(key);

            DiffType diffType;
            RiskLevel risk = RiskLevel.SAFE;
            String reason = "";

            if (packageValue == null) {
                diffType = DiffType.MISSING;
                risk = RiskLevel.P2;
                reason = "Package is missing a baseline key";
            } else if (baselineValue == null) {
                diffType = DiffType.NEW;
                risk = RiskLevel.P2;
                reason = "New config key; confirm whether to include in baseline";
            } else if (packageValue.equals(baselineValue)) {
                diffType = DiffType.MATCH;
            } else {
                diffType = DiffType.MISMATCH;
                risk = isCriticalKey(key) ? RiskLevel.P1 : RiskLevel.P2;
                reason = "Value differs from target environment baseline";
            }

            if ("prod".equalsIgnoreCase(targetEnv)
                    && PROD_FORBIDDEN_TRUE.contains(key)
                    && "true".equalsIgnoreCase(packageValue)) {
                diffType = DiffType.FORBIDDEN;
                risk = RiskLevel.P1;
                reason = "Prod forbids this debug or API-doc switch";
            }

            if ("spring.profiles.active".equals(key)
                    && "prod".equalsIgnoreCase(targetEnv)
                    && packageValue != null
                    && !"prod".equalsIgnoreCase(packageValue)) {
                diffType = DiffType.FORBIDDEN;
                risk = RiskLevel.P0;
                reason = "Prod deployment must use spring.profiles.active=prod";
            }

            if (containsDevAddress(packageValue) && "prod".equalsIgnoreCase(targetEnv)) {
                if (risk != RiskLevel.P0) {
                    risk = RiskLevel.P1;
                }
                if (reason.isBlank()) {
                    reason = "Detected possible dev address in prod package";
                }
            }

            if (isSensitive(key) && packageValue != null && !packageValue.isBlank()) {
                if (risk == RiskLevel.SAFE) {
                    risk = RiskLevel.P2;
                }
                if (reason.isBlank()) {
                    reason = "Sensitive value should be encrypted or externalized";
                }
            }

            items.add(new ConfigDiffItem(
                    configFile,
                    key,
                    maskSensitive(key, packageValue),
                    maskSensitive(key, baselineValue),
                    diffType,
                    risk,
                    categoryOf(key),
                    reason
            ));
        }

        Map<String, Long> summary = new LinkedHashMap<>();
        summary.put("P0", items.stream().filter(i -> i.riskLevel() == RiskLevel.P0).count());
        summary.put("P1", items.stream().filter(i -> i.riskLevel() == RiskLevel.P1).count());
        summary.put("P2", items.stream().filter(i -> i.riskLevel() == RiskLevel.P2).count());
        summary.put("MATCH", items.stream().filter(i -> i.diffType() == DiffType.MATCH).count());

        return new ConfigDiffReport(targetEnv, items, summary);
    }

    private boolean isCriticalKey(String key) {
        return key.startsWith("spring.datasource")
                || key.startsWith("spring.redis")
                || key.startsWith("nacos")
                || key.equals("spring.profiles.active")
                || key.equals("server.port");
    }

    private boolean containsDevAddress(String value) {
        return value != null && (value.contains("172.16.10.234") || value.contains("dev"));
    }

    private String categoryOf(String key) {
        if (key.startsWith("spring.datasource")) {
            return "DATABASE";
        }
        if (key.startsWith("spring.redis")) {
            return "REDIS";
        }
        if (key.startsWith("nacos")) {
            return "NACOS";
        }
        if (key.startsWith("debug") || key.contains("swagger") || key.contains("springdoc")) {
            return "SWITCH";
        }
        if (key.contains("workflow") || key.contains("cmdb") || key.contains("file-url")) {
            return "API";
        }
        return "OTHER";
    }

    private boolean isSensitive(String key) {
        String lower = key.toLowerCase();
        return lower.contains("password") || lower.contains("secret") || lower.endsWith(".key") || lower.contains("token");
    }

    private String maskSensitive(String key, String value) {
        if (value == null) {
            return null;
        }
        if (!isSensitive(key)) {
            return value;
        }
        if (value.length() <= 4) {
            return "****";
        }
        return value.substring(0, 2) + "****" + value.substring(value.length() - 2);
    }
}
