package com.huatai.release.entity;

import com.baomidou.mybatisplus.annotation.IdType;
import com.baomidou.mybatisplus.annotation.TableId;
import com.baomidou.mybatisplus.annotation.TableName;

import java.time.LocalDateTime;

@TableName("env_baseline")
public class EnvBaselineEntity {
    @TableId(type = IdType.AUTO)
    private Long id;
    private String envCode;
    private String appName;
    private String configKey;
    private String configValue;
    private Integer isSensitive;
    private Long sourceRequest;
    private LocalDateTime updatedAt;

    public Long getId() { return id; }
    public void setId(Long id) { this.id = id; }
    public String getEnvCode() { return envCode; }
    public void setEnvCode(String envCode) { this.envCode = envCode; }
    public String getAppName() { return appName; }
    public void setAppName(String appName) { this.appName = appName; }
    public String getConfigKey() { return configKey; }
    public void setConfigKey(String configKey) { this.configKey = configKey; }
    public String getConfigValue() { return configValue; }
    public void setConfigValue(String configValue) { this.configValue = configValue; }
    public Integer getIsSensitive() { return isSensitive; }
    public void setIsSensitive(Integer isSensitive) { this.isSensitive = isSensitive; }
    public Long getSourceRequest() { return sourceRequest; }
    public void setSourceRequest(Long sourceRequest) { this.sourceRequest = sourceRequest; }
    public LocalDateTime getUpdatedAt() { return updatedAt; }
    public void setUpdatedAt(LocalDateTime updatedAt) { this.updatedAt = updatedAt; }
}
