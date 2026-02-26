package com.huatai.release.engine.sql;

public record SqlRuleContext(
        String sql,
        int line,
        String businessDb
) {

    public String normalized() {
        return sql.replaceAll("\\s+", " ").trim().toUpperCase();
    }
}
