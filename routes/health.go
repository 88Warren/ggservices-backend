package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lmwcode/ggservices-backend/controllers"
)

func RegisterHealthRoutes(router *gin.Engine, hc *controllers.HealthController) {
	health := router.Group("/api/health")
	{
		health.GET("/live", hc.LivenessCheck)
		health.GET("/ready", hc.ReadinessCheck)
	}
}
