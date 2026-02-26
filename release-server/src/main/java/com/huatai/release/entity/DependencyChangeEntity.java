package com.huatai.release.entity;

import com.baomidou.mybatisplus.annotation.IdType;
import com.baomidou.mybatisplus.annotation.TableId;
import com.baomidou.mybatisplus.annotation.TableName;

@TableName("dependency_change")
public class DependencyChangeEntity {
    @TableId(type = IdType.AUTO)
    private Long id;
    private Long requestId;
    private String groupId;
    private String artifactId;
    private String oldVersion;
    private String newVersion;
    private String changeType;
    private String cveIds;
    private String riskLevel;

    public Long getId() { return id; }
    public void setId(Long id) { this.id = id; }
    public Long getRequestId() { return requestId; }
    public void setRequestId(Long requestId) { this.requestId = requestId; }
    public String getGroupId() { return groupId; }
    public void setGroupId(String groupId) { this.groupId = groupId; }
    public String getArtifactId() { return artifactId; }
    public void setArtifactId(String artifactId) { this.artifactId = artifactId; }
    public String getOldVersion() { return oldVersion; }
    public void setOldVersion(String oldVersion) { this.oldVersion = oldVersion; }
    public String getNewVersion() { return newVersion; }
    public void setNewVersion(String newVersion) { this.newVersion = newVersion; }
    public String getChangeType() { return changeType; }
    public void setChangeType(String changeType) { this.changeType = changeType; }
    public String getCveIds() { return cveIds; }
    public void setCveIds(String cveIds) { this.cveIds = cveIds; }
    public String getRiskLevel() { return riskLevel; }
    public void setRiskLevel(String riskLevel) { this.riskLevel = riskLevel; }
}
