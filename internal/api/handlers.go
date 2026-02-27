package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"release-platform/internal/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	releaseSvc  *service.ReleaseService
	artifactSvc *service.ArtifactService
}

func NewHandler(releaseSvc *service.ReleaseService, artifactSvc *service.ArtifactService) *Handler {
	return &Handler{releaseSvc: releaseSvc, artifactSvc: artifactSvc}
}

func (h *Handler) Healthz(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *Handler) CreateRelease(c *gin.Context) {
	var req service.CreateReleaseInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	release, err := h.releaseSvc.CreateRelease(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, release)
}

func (h *Handler) ListReleases(c *gin.Context) {
	env := c.Query("environment")
	result := c.Query("result")
	list, err := h.releaseSvc.ListReleases(c.Request.Context(), env, result)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": list, "count": len(list)})
}

func (h *Handler) UploadArtifact(c *gin.Context) {
	releaseID, err := parseID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fileHeader, err := c.FormFile("artifact")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing artifact file"})
		return
	}

	storagePath, size, sha, err := h.artifactSvc.SaveUploadedArtifact(releaseID, fileHeader)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	extractedPath, err := h.artifactSvc.ExtractArtifact(releaseID, storagePath)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	artifact, err := h.releaseSvc.UploadArtifact(c.Request.Context(), releaseID, fileHeader.Filename, storagePath, extractedPath, size, sha)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, artifact)
}

func (h *Handler) TriggerScan(c *gin.Context) {
	releaseID, err := parseID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	operator := c.GetHeader("X-User")
	if operator == "" {
		operator = "unknown"
	}

	if err := h.releaseSvc.EnqueueScan(c.Request.Context(), releaseID, operator); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"message": "scan queued"})
}

func (h *Handler) GetReport(c *gin.Context) {
	releaseID, err := parseID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	report, err := h.releaseSvc.GetReleaseReport(c.Request.Context(), releaseID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, report)
}

func (h *Handler) Approve(c *gin.Context) {
	releaseID, err := parseID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var req service.ApproveInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.releaseSvc.Approve(c.Request.Context(), releaseID, req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "approval recorded"})
}

func (h *Handler) Deploy(c *gin.Context) {
	releaseID, err := parseID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var req service.DeployInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.releaseSvc.Deploy(c.Request.Context(), releaseID, req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deployment simulated and baseline snapshot updated"})
}

func (h *Handler) GetPolicies(c *gin.Context) {
	var appID *uint
	if appStr := c.Query("application_id"); appStr != "" {
		parsed, err := parseID(appStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid application_id"})
			return
		}
		appID = &parsed
	}
	env := c.Query("environment")
	policies, err := h.releaseSvc.GetPolicies(c.Request.Context(), appID, env)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": policies, "count": len(policies)})
}

func (h *Handler) UpdatePolicy(c *gin.Context) {
	policyID, err := parseID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var req struct {
		Rules    json.RawMessage `json:"rules"`
		IsActive bool            `json:"is_active"`
		Operator string          `json:"operator"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.Operator == "" {
		req.Operator = "unknown"
	}
	if err := h.releaseSvc.UpdatePolicy(c.Request.Context(), policyID, req.Rules, req.IsActive, req.Operator); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "policy updated"})
}

func parseID(raw string) (uint, error) {
	if raw == "" {
		return 0, errors.New("missing id")
	}
	parsed, err := strconv.ParseUint(raw, 10, 64)
	if err != nil {
		return 0, errors.New("invalid numeric id")
	}
	return uint(parsed), nil
}
