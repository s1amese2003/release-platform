package com.huatai.release.model.dto;

public record ManualOperationItem(
        String operationDesc,
        String operationType,
        String execStatus
) {
}
