package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lmwcode/ggservices-backend/controllers"
)

func RegisterSitemapRoutes(router *gin.Engine, sitemapController *controllers.SitemapController) {
	router.GET("/sitemap.xml", sitemapController.GenerateSitemap)
}
