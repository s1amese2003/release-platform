package com.huatai.release.engine.config;

import org.springframework.stereotype.Component;
import org.yaml.snakeyaml.Yaml;

import java.util.LinkedHashMap;
import java.util.Map;

@Component
public class YamlFlattener {

    public Map<String, String> flatten(String yamlContent) {
        Map<String, String> flattened = new LinkedHashMap<>();
        if (yamlContent == null || yamlContent.isBlank()) {
            return flattened;
        }

        Yaml yaml = new Yaml();
        Object obj = yaml.load(yamlContent);
        if (obj instanceof Map<?, ?> map) {
            flattenMap("", map, flattened);
        }
        return flattened;
    }

    @SuppressWarnings("unchecked")
    private void flattenMap(String prefix, Map<?, ?> src, Map<String, String> out) {
        for (Map.Entry<?, ?> entry : src.entrySet()) {
            String key = entry.getKey().toString();
            String path = prefix.isEmpty() ? key : prefix + "." + key;
            Object value = entry.getValue();
            if (value instanceof Map<?, ?> nested) {
                flattenMap(path, nested, out);
            } else {
                out.put(path, value == null ? "null" : String.valueOf(value));
            }
        }
    }
}
