package com.huatai.release.repository;

import org.springframework.stereotype.Repository;

import java.util.Collections;
import java.util.HashMap;
import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;

@Repository
public class BaselineMemoryStore {

    private final Map<String, Map<String, String>> baselines = new ConcurrentHashMap<>();

    public BaselineMemoryStore() {
        Map<String, String> prod = new HashMap<>();
        prod.put("spring.profiles.active", "prod");
        prod.put("server.port", "8083");
        prod.put("spring.application.name", "dis-itsm");
        prod.put("spring.datasource.url", "jdbc:mysql://10.8.25.10:3306/itsm_data");
        prod.put("spring.redis.host", "10.8.25.12");
        prod.put("nacos.server-addr", "10.8.25.1:8848");
        prod.put("swagger.enabled", "false");
        prod.put("springdoc.api-docs.enabled", "false");
        prod.put("debug.iceCmdb.systemNames", "false");
        prod.put("debug.workflowMaintenace.adjustParticipant", "false");
        prod.put("email.open", "true");

        baselines.put(key("prod", "dis-modules-itsm"), prod);
    }

    public Map<String, String> getBaseline(String env, String app) {
        return baselines.getOrDefault(key(env, app), Collections.emptyMap());
    }

    public void upsert(String env, String app, Map<String, String> values) {
        baselines.computeIfAbsent(key(env, app), k -> new ConcurrentHashMap<>()).putAll(values);
    }

    public void replace(String env, String app, Map<String, String> values) {
        baselines.put(key(env, app), new ConcurrentHashMap<>(values));
    }

    private String key(String env, String app) {
        return env + "::" + app;
    }
}
