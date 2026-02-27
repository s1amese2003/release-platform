package com.huatai.release.service.impl;

import com.huatai.release.model.dto.ApprovalActionRequest;
import com.huatai.release.model.dto.ReleaseCheckReport;
import com.huatai.release.model.enums.RequestStatus;
import com.huatai.release.engine.config.YamlFlattener;
import com.huatai.release.repository.ReleaseMemoryStore;
import com.huatai.release.service.ApprovalService;
import com.huatai.release.service.BaselineService;
import org.springframework.stereotype.Service;

import java.util.LinkedHashMap;
import java.util.Map;

@Service
public class ApprovalServiceImpl implements ApprovalService {

    private final ReleaseMemoryStore store;
    private final BaselineService baselineService;
    private final YamlFlattener yamlFlattener;

    public ApprovalServiceImpl(ReleaseMemoryStore store,
                               BaselineService baselineService,
                               YamlFlattener yamlFlattener) {
        this.store = store;
        this.baselineService = baselineService;
        this.yamlFlattener = yamlFlattener;
    }

    @Override
    public RequestStatus action(String requestNo, ApprovalActionRequest request) {
        RequestStatus current = store.getStatus(requestNo);
        if (current == null) {
            throw new IllegalArgumentException("申请单不存在: " + requestNo);
        }

        String action = request.getAction().toUpperCase();
        RequestStatus next;
        switch (action) {
            case "APPROVE" -> next = RequestStatus.APPROVED;
            case "REJECT" -> next = RequestStatus.REJECTED;
            case "DEPLOY" -> {
                next = RequestStatus.DEPLOYED;
                snapshotBaseline(requestNo);
            }
            default -> throw new IllegalArgumentException("不支持的审批动作: " + request.getAction());
        }

        store.saveStatus(requestNo, next);
        return next;
    }

    private void snapshotBaseline(String requestNo) {
        ReleaseCheckReport report = store.getReport(requestNo);
        if (report == null || report.parse() == null) {
            return;
        }
        Map<String, String> merged = new LinkedHashMap<>();
        report.parse().configFiles().values().forEach(content -> merged.putAll(yamlFlattener.flatten(content)));
        baselineService.snapshotFromPackage(report.parse().targetEnv(), report.parse().appName(), merged);
    }
}
