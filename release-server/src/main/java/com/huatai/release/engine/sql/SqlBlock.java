package com.huatai.release.engine.sql;

public record SqlBlock(
        String sql,
        int line
) {
}
