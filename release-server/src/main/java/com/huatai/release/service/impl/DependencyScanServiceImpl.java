package com.huatai.release.service.impl;

import com.huatai.release.engine.dependency.DependencyDiffEngine;
import com.huatai.release.model.dto.DependencyCoordinate;
import com.huatai.release.model.dto.DependencyReport;
import com.huatai.release.service.DependencyScanService;
import org.springframework.stereotype.Service;

import java.util.List;

@Service
public class DependencyScanServiceImpl implements DependencyScanService {

    private final DependencyDiffEngine engine;

    public DependencyScanServiceImpl(DependencyDiffEngine engine) {
        this.engine = engine;
    }

    @Override
    public DependencyReport scan(List<DependencyCoordinate> current, List<DependencyCoordinate> previous) {
        return engine.scan(current, previous);
    }
}
