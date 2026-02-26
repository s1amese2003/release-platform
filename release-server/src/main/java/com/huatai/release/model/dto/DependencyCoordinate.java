package com.huatai.release.model.dto;

public record DependencyCoordinate(
        String groupId,
        String artifactId,
        String version,
        String rawName
) {
}
