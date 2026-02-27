package api

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

func NewRouter(handler *Handler, webRoot string) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())
	router.Use(auditActorMiddleware())

	router.GET("/healthz", handler.Healthz)
	router.Static("/assets", filepath.Join(webRoot, "assets"))
	router.StaticFile("/", filepath.Join(webRoot, "index.html"))

	api := router.Group("/api")
	{
		api.GET("/releases", handler.ListReleases)
		api.POST("/releases", handler.CreateRelease)
		api.POST("/releases/:id/artifact", handler.UploadArtifact)
		api.POST("/releases/:id/scan", handler.TriggerScan)
		api.GET("/releases/:id/report", handler.GetReport)
		api.POST("/releases/:id/approve", requireRoles("owner", "ops", "dba"), handler.Approve)
		api.POST("/releases/:id/deploy", requireRoles("owner", "ops"), handler.Deploy)

		api.GET("/policies", handler.GetPolicies)
		api.PUT("/policies/:id", requireRoles("dba", "audit"), handler.UpdatePolicy)
	}

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "route not found"})
	})

	return router
}

func auditActorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("X-User") == "" {
			c.Request.Header.Set("X-User", "anonymous")
		}
		c.Next()
	}
}

func requireRoles(allowed ...string) gin.HandlerFunc {
	allowMap := map[string]struct{}{}
	for _, role := range allowed {
		allowMap[role] = struct{}{}
	}

	return func(c *gin.Context) {
		role := strings.TrimSpace(strings.ToLower(c.GetHeader("X-Role")))
		if role == "" {
			c.JSON(http.StatusForbidden, gin.H{"error": "missing role header"})
			c.Abort()
			return
		}
		if _, ok := allowMap[role]; !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "role not allowed"})
			c.Abort()
			return
		}
		c.Next()
	}
}
