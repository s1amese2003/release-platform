package com.huatai.release.engine.dependency;

import com.huatai.release.model.dto.DependencyChangeItem;
import com.huatai.release.model.dto.DependencyCoordinate;
import com.huatai.release.model.dto.DependencyReport;
import com.huatai.release.model.enums.RiskLevel;
import org.springframework.stereotype.Component;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

@Component
public class DependencyDiffEngine {

    private static final Map<String, List<String>> KNOWN_CVE = Map.of(
            "com.alibaba:fastjson", List.of("CVE-2022-25845"),
            "org.apache.logging.log4j:log4j-core", List.of("CVE-2021-44228")
    );

    public DependencyReport scan(List<DependencyCoordinate> current, List<DependencyCoordinate> previous) {
        Map<String, DependencyCoordinate> currentMap = toMap(current);
        Map<String, DependencyCoordinate> previousMap = toMap(previous);

        List<DependencyChangeItem> changes = new ArrayList<>();

        for (Map.Entry<String, DependencyCoordinate> entry : currentMap.entrySet()) {
            String key = entry.getKey();
            DependencyCoordinate now = entry.getValue();
            DependencyCoordinate old = previousMap.remove(key);

            if (old == null) {
                changes.add(toChange(now, null, "ADD"));
            } else if (!now.version().equals(old.version())) {
                String type = compareVersion(old.version(), now.version()) <= 0 ? "UPGRADE" : "DOWNGRADE";
                changes.add(toChange(now, old.version(), type));
            }
        }

        for (DependencyCoordinate removed : previousMap.values()) {
            changes.add(new DependencyChangeItem(
                    removed.groupId(),
                    removed.artifactId(),
                    removed.version(),
                    null,
                    "REMOVE",
                    List.of(),
                    RiskLevel.SAFE
            ));
        }

        int highRisk = (int) changes.stream().filter(item -> item.riskLevel() == RiskLevel.P1).count();
        return new DependencyReport(current.size(), highRisk, changes);
    }

    private DependencyChangeItem toChange(DependencyCoordinate now, String oldVersion, String type) {
        String key = now.groupId() + ":" + now.artifactId();
        List<String> cves = KNOWN_CVE.getOrDefault(key, List.of());
        RiskLevel risk = cves.isEmpty() ? RiskLevel.SAFE : RiskLevel.P1;

        if ("fastjson".equalsIgnoreCase(now.artifactId()) && now.version().startsWith("1.2.")) {
            risk = RiskLevel.P1;
            cves = List.of("CVE-2022-25845", "Legacy 1.2.x branch");
        }

        return new DependencyChangeItem(
                now.groupId(),
                now.artifactId(),
                oldVersion,
                now.version(),
                type,
                cves,
                risk
        );
    }

    private Map<String, DependencyCoordinate> toMap(List<DependencyCoordinate> coordinates) {
        Map<String, DependencyCoordinate> map = new HashMap<>();
        for (DependencyCoordinate coordinate : coordinates) {
            String key = coordinate.groupId() + ":" + coordinate.artifactId();
            map.put(key, coordinate);
        }
        return map;
    }

    private int compareVersion(String oldVersion, String newVersion) {
        return oldVersion.compareTo(newVersion);
    }
}
