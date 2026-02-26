package com.huatai.release.entity;

import com.baomidou.mybatisplus.annotation.IdType;
import com.baomidou.mybatisplus.annotation.TableId;
import com.baomidou.mybatisplus.annotation.TableName;

@TableName("config_diff_result")
public class ConfigDiffResultEntity {
    @TableId(type = IdType.AUTO)
    private Long id;
    private Long requestId;
    private String configFile;
    private String configKey;
    private String packageValue;
    private String baselineValue;
    private String diffType;
    private String riskLevel;
    private String category;

    public Long getId() { return id; }
    public void setId(Long id) { this.id = id; }
    public Long getRequestId() { return requestId; }
    public void setRequestId(Long requestId) { this.requestId = requestId; }
    public String getConfigFile() { return configFile; }
    public void setConfigFile(String configFile) { this.configFile = configFile; }
    public String getConfigKey() { return configKey; }
    public void setConfigKey(String configKey) { this.configKey = configKey; }
    public String getPackageValue() { return packageValue; }
    public void setPackageValue(String packageValue) { this.packageValue = packageValue; }
    public String getBaselineValue() { return baselineValue; }
    public void setBaselineValue(String baselineValue) { this.baselineValue = baselineValue; }
    public String getDiffType() { return diffType; }
    public void setDiffType(String diffType) { this.diffType = diffType; }
    public String getRiskLevel() { return riskLevel; }
    public void setRiskLevel(String riskLevel) { this.riskLevel = riskLevel; }
    public String getCategory() { return category; }
    public void setCategory(String category) { this.category = category; }
}
