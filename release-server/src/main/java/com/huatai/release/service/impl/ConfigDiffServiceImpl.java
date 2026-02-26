package com.huatai.release.service.impl;

import com.huatai.release.engine.config.ConfigDiffEngine;
import com.huatai.release.model.dto.ConfigDiffReport;
import com.huatai.release.service.ConfigDiffService;
import org.springframework.stereotype.Service;

import java.util.Map;

@Service
public class ConfigDiffServiceImpl implements ConfigDiffService {

    private final ConfigDiffEngine engine;

    public ConfigDiffServiceImpl(ConfigDiffEngine engine) {
        this.engine = engine;
    }

    @Override
    public ConfigDiffReport diff(Map<String, String> packageConfig,
                                 Map<String, String> baseline,
                                 String targetEnv,
                                 String configFile) {
        return engine.diff(packageConfig, baseline, targetEnv, configFile);
    }
}
