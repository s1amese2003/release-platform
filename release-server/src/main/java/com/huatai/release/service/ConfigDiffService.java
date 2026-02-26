package com.huatai.release.service;

import com.huatai.release.model.dto.ConfigDiffReport;

import java.util.Map;

public interface ConfigDiffService {

    ConfigDiffReport diff(Map<String, String> packageConfig,
                          Map<String, String> baseline,
                          String targetEnv,
                          String configFile);
}
