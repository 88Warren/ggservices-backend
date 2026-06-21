package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthController struct{}

func NewHealthController() *HealthController {
	return &HealthController{}
}

// LivenessCheck confirms the process is alive.
func (hc *HealthController) LivenessCheck(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// ReadinessCheck confirms the service is ready to accept traffic.
func (hc *HealthController) ReadinessCheck(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"status": "ready"})
}
