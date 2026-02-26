package com.huatai.release.repository;

import com.huatai.release.model.dto.ReleaseCheckReport;
import com.huatai.release.model.enums.RequestStatus;
import org.springframework.stereotype.Repository;

import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;

@Repository
public class ReleaseMemoryStore {

    private final Map<String, ReleaseCheckReport> reports = new ConcurrentHashMap<>();
    private final Map<String, RequestStatus> statusMap = new ConcurrentHashMap<>();

    public void saveReport(String requestNo, ReleaseCheckReport report) {
        reports.put(requestNo, report);
    }

    public ReleaseCheckReport getReport(String requestNo) {
        return reports.get(requestNo);
    }

    public void saveStatus(String requestNo, RequestStatus status) {
        statusMap.put(requestNo, status);
    }

    public RequestStatus getStatus(String requestNo) {
        return statusMap.get(requestNo);
    }
}
