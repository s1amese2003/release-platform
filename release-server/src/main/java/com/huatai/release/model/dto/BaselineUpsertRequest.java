package com.huatai.release.model.dto;

import java.util.Map;

public record BaselineUpsertRequest(
        Map<String, String> values
) {
}
