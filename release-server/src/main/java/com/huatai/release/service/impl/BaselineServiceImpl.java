package com.huatai.release.service.impl;

import com.huatai.release.repository.BaselineMemoryStore;
import com.huatai.release.service.BaselineService;
import org.springframework.stereotype.Service;

import java.util.Map;

@Service
public class BaselineServiceImpl implements BaselineService {

    private final BaselineMemoryStore store;

    public BaselineServiceImpl(BaselineMemoryStore store) {
        this.store = store;
    }

    @Override
    public Map<String, String> getBaseline(String env, String appName) {
        return store.getBaseline(env, appName);
    }

    @Override
    public void upsertBaseline(String env, String appName, Map<String, String> values) {
        store.upsert(env, appName, values);
    }

    @Override
    public void snapshotFromPackage(String env, String appName, Map<String, String> values) {
        store.replace(env, appName, values);
    }
}
