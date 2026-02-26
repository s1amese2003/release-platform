package com.huatai.release.entity;

import com.baomidou.mybatisplus.annotation.IdType;
import com.baomidou.mybatisplus.annotation.TableId;
import com.baomidou.mybatisplus.annotation.TableName;

import java.time.LocalDateTime;

@TableName("release_request")
public class ReleaseRequestEntity {

    @TableId(type = IdType.AUTO)
    private Long id;
    private String requestNo;
    private String appName;
    private String appVersion;
    private String targetEnv;
    private String upgradeVersion;
    private String packagePath;
    private String packageMd5;
    private LocalDateTime buildTime;
    private String buildJdk;
    private String submitter;
    private String status;
    private LocalDateTime createdAt;
    private LocalDateTime deployedAt;

    public Long getId() { return id; }
    public void setId(Long id) { this.id = id; }
    public String getRequestNo() { return requestNo; }
    public void setRequestNo(String requestNo) { this.requestNo = requestNo; }
    public String getAppName() { return appName; }
    public void setAppName(String appName) { this.appName = appName; }
    public String getAppVersion() { return appVersion; }
    public void setAppVersion(String appVersion) { this.appVersion = appVersion; }
    public String getTargetEnv() { return targetEnv; }
    public void setTargetEnv(String targetEnv) { this.targetEnv = targetEnv; }
    public String getUpgradeVersion() { return upgradeVersion; }
    public void setUpgradeVersion(String upgradeVersion) { this.upgradeVersion = upgradeVersion; }
    public String getPackagePath() { return packagePath; }
    public void setPackagePath(String packagePath) { this.packagePath = packagePath; }
    public String getPackageMd5() { return packageMd5; }
    public void setPackageMd5(String packageMd5) { this.packageMd5 = packageMd5; }
    public LocalDateTime getBuildTime() { return buildTime; }
    public void setBuildTime(LocalDateTime buildTime) { this.buildTime = buildTime; }
    public String getBuildJdk() { return buildJdk; }
    public void setBuildJdk(String buildJdk) { this.buildJdk = buildJdk; }
    public String getSubmitter() { return submitter; }
    public void setSubmitter(String submitter) { this.submitter = submitter; }
    public String getStatus() { return status; }
    public void setStatus(String status) { this.status = status; }
    public LocalDateTime getCreatedAt() { return createdAt; }
    public void setCreatedAt(LocalDateTime createdAt) { this.createdAt = createdAt; }
    public LocalDateTime getDeployedAt() { return deployedAt; }
    public void setDeployedAt(LocalDateTime deployedAt) { this.deployedAt = deployedAt; }
}
