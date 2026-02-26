package com.huatai.release.service;

import com.huatai.release.model.dto.SqlAuditReport;

public interface SqlAuditService {

    SqlAuditReport audit(String sqlContent, String version, String businessDb);
}
