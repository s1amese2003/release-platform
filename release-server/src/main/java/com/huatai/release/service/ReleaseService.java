package com.huatai.release.service;

import com.huatai.release.model.dto.ReleaseCheckReport;
import org.springframework.web.multipart.MultipartFile;

public interface ReleaseService {

    ReleaseCheckReport analyze(MultipartFile file, String appName, String targetEnv, String submitter);

    ReleaseCheckReport findReport(String requestNo);
}
