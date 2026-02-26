package com.huatai.release.service;

import com.huatai.release.model.dto.PackageParseResult;
import org.springframework.web.multipart.MultipartFile;

public interface PackageParseService {

    PackageParseResult parse(MultipartFile file, String appName, String targetEnv, String submitter);
}
