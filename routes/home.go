package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lmwcode/ggservices-backend/controllers"
)

func RegisterHomeRoutes(router *gin.Engine, hc *controllers.HomeController) {
	// Serve images from the local images directory
	router.Static("/images", "./images")

	api := router.Group("/api")
	{
		api.GET("/", hc.GetHome)
		api.POST("/contact", hc.HandleContactForm)
	}
}
