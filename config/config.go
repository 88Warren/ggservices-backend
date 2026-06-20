package config

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/lmwcode/ggservices-backend/controllers"
	"github.com/lmwcode/ggservices-backend/middleware"
	"github.com/lmwcode/ggservices-backend/routes"
)

func LoadEnv() {
	env := os.Getenv("GO_ENV")
	var envFile string

	if env == "production" {
		envFile = ".env.production"
	} else {
		envFile = ".env.development"
	}

	err := godotenv.Load(envFile)
	if err != nil {
		log.Printf("Warning: No %s file found, relying on system environment variables", envFile)
	}
}

func SetGinMode() {
	gin.SetMode(gin.ReleaseMode)
}

func GetEnv(key string, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}
	return value
}

func SetupServer() *gin.Engine {
	isTestEnv := os.Getenv("GO_ENV") == "test" ||
		os.Getenv("CI") != "" ||
		os.Getenv("HARNESS_BUILD_ID") != "" ||
		os.Getenv("GITHUB_ACTIONS") != ""

	var router *gin.Engine
	if isTestEnv {
		gin.SetMode(gin.ReleaseMode)
		router = gin.New()
		router.Use(gin.Recovery())
	} else {
		router = gin.Default()
	}

	// Disable 405 Method Not Allowed responses so that the CORS middleware
	// can intercept OPTIONS preflight requests before Gin rejects them.
	router.HandleMethodNotAllowed = false

	router.Use(middleware.CORSMiddleware())

	return router
}

func SetupHandlers(router *gin.Engine) {
	homeController := controllers.NewHomeController()
	sitemapController := controllers.NewSitemapController()

	routes.RegisterHomeRoutes(router, homeController)
	routes.RegisterSitemapRoutes(router, sitemapController)
}
