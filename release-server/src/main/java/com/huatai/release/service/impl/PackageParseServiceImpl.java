package com.huatai.release.service.impl;

import com.huatai.release.engine.pkg.JarPackageParser;
import com.huatai.release.model.dto.PackageParseResult;
import com.huatai.release.service.PackageParseService;
import org.springframework.stereotype.Service;
import org.springframework.web.multipart.MultipartFile;

@Service
public class PackageParseServiceImpl implements PackageParseService {

    private final JarPackageParser parser;

    public PackageParseServiceImpl(JarPackageParser parser) {
        this.parser = parser;
    }

    @Override
    public PackageParseResult parse(MultipartFile file, String appName, String targetEnv, String submitter) {
        return parser.parse(file, appName, targetEnv, submitter);
    }
}
