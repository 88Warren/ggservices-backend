package main

import (
	"log"
	"os"

	"github.com/lmwcode/ggservices-backend/config"
	"go.uber.org/zap"
)

func init() {
	config.LoadEnv()
	config.SetGinMode()
}

func main() {

	// Debug: Print environment variables
	log.Printf("GO_ENV: %s", os.Getenv("GO_ENV"))
	log.Printf("ALLOWED_ORIGIN: %s", os.Getenv("ALLOWED_ORIGIN"))
	log.Printf("PORT: %s", os.Getenv("PORT"))

	router := config.SetupServer()

	config.SetupHandlers(router)

	port := config.GetEnv("PORT", "8090")

	if err := router.Run("0.0.0.0:" + port); err != nil {
		zap.L().Fatal("Failed to start the server", zap.Error(err))
	}
}
