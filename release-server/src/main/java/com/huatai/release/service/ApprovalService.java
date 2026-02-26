package com.huatai.release.service;

import com.huatai.release.model.dto.ApprovalActionRequest;
import com.huatai.release.model.enums.RequestStatus;

public interface ApprovalService {

    RequestStatus action(String requestNo, ApprovalActionRequest request);
}
