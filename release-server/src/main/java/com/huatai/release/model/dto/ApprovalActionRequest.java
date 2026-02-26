package com.huatai.release.model.dto;

import jakarta.validation.constraints.NotBlank;

public class ApprovalActionRequest {

    @NotBlank
    private String approver;
    @NotBlank
    private String action;
    private String comment;

    public String getApprover() {
        return approver;
    }

    public void setApprover(String approver) {
        this.approver = approver;
    }

    public String getAction() {
        return action;
    }

    public void setAction(String action) {
        this.action = action;
    }

    public String getComment() {
        return comment;
    }

    public void setComment(String comment) {
        this.comment = comment;
    }
}
