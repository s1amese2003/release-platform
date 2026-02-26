package com.huatai.release.entity;

import com.baomidou.mybatisplus.annotation.IdType;
import com.baomidou.mybatisplus.annotation.TableId;
import com.baomidou.mybatisplus.annotation.TableName;

@TableName("sql_audit_result")
public class SqlAuditResultEntity {
    @TableId(type = IdType.AUTO)
    private Long id;
    private Long requestId;
    private String sqlFilePath;
    private Integer lineNumber;
    private String sqlContent;
    private String sqlType;
    private String riskLevel;
    private String riskReason;
    private String suggestion;
    private String reviewerAction;
    private String reviewer;

    public Long getId() { return id; }
    public void setId(Long id) { this.id = id; }
    public Long getRequestId() { return requestId; }
    public void setRequestId(Long requestId) { this.requestId = requestId; }
    public String getSqlFilePath() { return sqlFilePath; }
    public void setSqlFilePath(String sqlFilePath) { this.sqlFilePath = sqlFilePath; }
    public Integer getLineNumber() { return lineNumber; }
    public void setLineNumber(Integer lineNumber) { this.lineNumber = lineNumber; }
    public String getSqlContent() { return sqlContent; }
    public void setSqlContent(String sqlContent) { this.sqlContent = sqlContent; }
    public String getSqlType() { return sqlType; }
    public void setSqlType(String sqlType) { this.sqlType = sqlType; }
    public String getRiskLevel() { return riskLevel; }
    public void setRiskLevel(String riskLevel) { this.riskLevel = riskLevel; }
    public String getRiskReason() { return riskReason; }
    public void setRiskReason(String riskReason) { this.riskReason = riskReason; }
    public String getSuggestion() { return suggestion; }
    public void setSuggestion(String suggestion) { this.suggestion = suggestion; }
    public String getReviewerAction() { return reviewerAction; }
    public void setReviewerAction(String reviewerAction) { this.reviewerAction = reviewerAction; }
    public String getReviewer() { return reviewer; }
    public void setReviewer(String reviewer) { this.reviewer = reviewer; }
}
