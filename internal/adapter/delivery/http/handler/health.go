package handler

import (
	netHttp "net/http"
	"todolist/internal/adapter/delivery/http"
	"todolist/internal/config"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
	appConfig config.ApplicationProvider
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(
	appConfig config.ApplicationProvider,
) *HealthHandler {
	return &HealthHandler{
		appConfig: appConfig,
	}
}

// HealthCheck godoc
// @Summary Health check
// @Description Returns the status of the application
// @Tags Health
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /health [get]
func (h *HealthHandler) HealthCheck(ctx http.RequestContext) {
	ctx.JSON(netHttp.StatusOK, gin.H{
		"status":  "ok",
		"app":     h.appConfig.GetName(),
		"version": h.appConfig.GetVersion(),
	})
}
