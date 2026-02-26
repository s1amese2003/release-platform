package com.huatai.release.controller;

import com.huatai.release.model.dto.ApiResponse;
import com.huatai.release.model.dto.BaselineUpsertRequest;
import com.huatai.release.service.BaselineService;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import java.util.Map;

@RestController
@RequestMapping("/api/baselines")
public class BaselineController {

    private final BaselineService baselineService;

    public BaselineController(BaselineService baselineService) {
        this.baselineService = baselineService;
    }

    @GetMapping("/{env}/{appName}")
    public ApiResponse<Map<String, String>> getBaseline(@PathVariable String env, @PathVariable String appName) {
        return ApiResponse.ok(baselineService.getBaseline(env, appName));
    }

    @PostMapping("/{env}/{appName}")
    public ApiResponse<String> upsertBaseline(@PathVariable String env,
                                              @PathVariable String appName,
                                              @RequestBody BaselineUpsertRequest request) {
        baselineService.upsertBaseline(env, appName, request.values());
        return ApiResponse.ok("UPDATED");
    }
}
