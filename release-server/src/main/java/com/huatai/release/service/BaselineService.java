package com.huatai.release.service;

import java.util.Map;

public interface BaselineService {

    Map<String, String> getBaseline(String env, String appName);

    void upsertBaseline(String env, String appName, Map<String, String> values);

    void snapshotFromPackage(String env, String appName, Map<String, String> values);
}
