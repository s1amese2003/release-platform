package com.huatai.release.service.impl;

import com.huatai.release.engine.sql.SqlRiskAnalyzer;
import com.huatai.release.model.dto.SqlAuditReport;
import com.huatai.release.service.SqlAuditService;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Service;

@Service
public class SqlAuditServiceImpl implements SqlAuditService {

    private final SqlRiskAnalyzer analyzer;

    @Value("${release.default-business-db:itsm_data}")
    private String defaultBusinessDb;

    public SqlAuditServiceImpl(SqlRiskAnalyzer analyzer) {
        this.analyzer = analyzer;
    }

    @Override
    public SqlAuditReport audit(String sqlContent, String version, String businessDb) {
        String db = (businessDb == null || businessDb.isBlank()) ? defaultBusinessDb : businessDb;
        return analyzer.analyze(sqlContent, version, db);
    }
}
