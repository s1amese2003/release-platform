package com.huatai.release.service.impl;

import com.huatai.release.engine.config.YamlFlattener;
import com.huatai.release.model.dto.DependencyReport;
import com.huatai.release.model.dto.ManualOperationItem;
import com.huatai.release.model.dto.PackageParseResult;
import com.huatai.release.model.dto.ReleaseCheckReport;
import com.huatai.release.model.dto.SqlAuditReport;
import com.huatai.release.model.enums.RequestStatus;
import com.huatai.release.repository.ReleaseMemoryStore;
import com.huatai.release.service.BaselineService;
import com.huatai.release.service.ConfigDiffService;
import com.huatai.release.service.DependencyScanService;
import com.huatai.release.service.PackageParseService;
import com.huatai.release.service.ReleaseService;
import com.huatai.release.service.SqlAuditService;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Service;
import org.springframework.web.multipart.MultipartFile;

import java.util.ArrayList;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Locale;
import java.util.Map;

@Service
public class ReleaseServiceImpl implements ReleaseService {

    private final PackageParseService packageParseService;
    private final SqlAuditService sqlAuditService;
    private final ConfigDiffService configDiffService;
    private final DependencyScanService dependencyScanService;
    private final BaselineService baselineService;
    private final YamlFlattener yamlFlattener;
    private final ReleaseMemoryStore store;

    @Value("${release.default-business-db:itsm_data}")
    private String businessDb;

    public ReleaseServiceImpl(PackageParseService packageParseService,
                              SqlAuditService sqlAuditService,
                              ConfigDiffService configDiffService,
                              DependencyScanService dependencyScanService,
                              BaselineService baselineService,
                              YamlFlattener yamlFlattener,
                              ReleaseMemoryStore store) {
        this.packageParseService = packageParseService;
        this.sqlAuditService = sqlAuditService;
        this.configDiffService = configDiffService;
        this.dependencyScanService = dependencyScanService;
        this.baselineService = baselineService;
        this.yamlFlattener = yamlFlattener;
        this.store = store;
    }

    @Override
    public ReleaseCheckReport analyze(MultipartFile file, String appName, String targetEnv, String submitter) {
        PackageParseResult parse = packageParseService.parse(file, appName, targetEnv, submitter);
        store.saveStatus(parse.requestNo(), RequestStatus.CHECKING);

        SqlAuditReport sqlAudit = sqlAuditService.audit(parse.sqlContent(), parse.latestUpgradeVersion(), businessDb);

        Map<String, String> mergedPackageConfig = mergeConfig(parse.configFiles());
        Map<String, String> baseline = baselineService.getBaseline(targetEnv, appName);
        var configDiff = configDiffService.diff(mergedPackageConfig, baseline, targetEnv, "merged-config");

        DependencyReport depReport = dependencyScanService.scan(parse.dependencies(), List.of());
        List<ManualOperationItem> checklist = buildChecklist(parse.manualOperations());

        boolean rejectedByP0 = sqlAudit.summary().getOrDefault("P0", 0L) > 0
                || configDiff.summary().getOrDefault("P0", 0L) > 0;

        if (rejectedByP0) {
            store.saveStatus(parse.requestNo(), RequestStatus.REJECTED);
        } else {
            store.saveStatus(parse.requestNo(), RequestStatus.REVIEW);
        }

        ReleaseCheckReport report = new ReleaseCheckReport(parse, sqlAudit, configDiff, depReport, checklist, rejectedByP0);
        store.saveReport(parse.requestNo(), report);
        return report;
    }

    @Override
    public ReleaseCheckReport findReport(String requestNo) {
        return store.getReport(requestNo);
    }

    private Map<String, String> mergeConfig(Map<String, String> configFiles) {
        Map<String, String> merged = new LinkedHashMap<>();
        for (String content : configFiles.values()) {
            merged.putAll(yamlFlattener.flatten(content));
        }
        return merged;
    }

    private List<ManualOperationItem> buildChecklist(List<String> operations) {
        List<ManualOperationItem> list = new ArrayList<>();
        for (String operation : operations) {
            String type = "MANUAL";
            String lower = operation.toLowerCase(Locale.ROOT);
            if (lower.contains("nacos")) {
                type = "NACOS_CONFIG";
            } else if (lower.contains("sql") || lower.contains("db")) {
                type = "DB_SCRIPT";
            }
            list.add(new ManualOperationItem(operation, type, "PENDING"));
        }
        return list;
    }
}
