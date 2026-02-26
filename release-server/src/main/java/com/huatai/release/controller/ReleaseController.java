package com.huatai.release.controller;

import com.huatai.release.model.dto.ApiResponse;
import com.huatai.release.model.dto.ApprovalActionRequest;
import com.huatai.release.model.dto.ReleaseCheckReport;
import com.huatai.release.model.enums.RequestStatus;
import com.huatai.release.service.ApprovalService;
import com.huatai.release.service.ReleaseService;
import jakarta.validation.Valid;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RequestPart;
import org.springframework.web.bind.annotation.RestController;
import org.springframework.web.multipart.MultipartFile;

@RestController
@RequestMapping("/api/releases")
public class ReleaseController {

    private final ReleaseService releaseService;
    private final ApprovalService approvalService;

    public ReleaseController(ReleaseService releaseService, ApprovalService approvalService) {
        this.releaseService = releaseService;
        this.approvalService = approvalService;
    }

    @PostMapping("/analyze")
    public ApiResponse<ReleaseCheckReport> analyze(@RequestPart("file") MultipartFile file,
                                                   @RequestParam("appName") String appName,
                                                   @RequestParam("targetEnv") String targetEnv,
                                                   @RequestParam(value = "submitter", defaultValue = "anonymous") String submitter) {
        return ApiResponse.ok(releaseService.analyze(file, appName, targetEnv, submitter));
    }

    @GetMapping("/{requestNo}")
    public ApiResponse<ReleaseCheckReport> findReport(@PathVariable String requestNo) {
        ReleaseCheckReport report = releaseService.findReport(requestNo);
        if (report == null) {
            return ApiResponse.fail("鏈壘鍒扮敵璇峰崟: " + requestNo);
        }
        return ApiResponse.ok(report);
    }

    @PostMapping("/{requestNo}/approval")
    public ApiResponse<RequestStatus> approval(@PathVariable String requestNo,
                                               @Valid @org.springframework.web.bind.annotation.RequestBody ApprovalActionRequest request) {
        return ApiResponse.ok(approvalService.action(requestNo, request));
    }
}
