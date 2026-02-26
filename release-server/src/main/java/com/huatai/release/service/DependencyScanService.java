package com.huatai.release.service;

import com.huatai.release.model.dto.DependencyCoordinate;
import com.huatai.release.model.dto.DependencyReport;

import java.util.List;

public interface DependencyScanService {

    DependencyReport scan(List<DependencyCoordinate> current,
                          List<DependencyCoordinate> previous);
}
