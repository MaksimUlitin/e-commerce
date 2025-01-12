package main

import (
	"github.com/gin-gonic/gin"
	"github.com/maksimulitin/config"
	"github.com/maksimulitin/internal/controllers"
	"github.com/maksimulitin/internal/database"
	"github.com/maksimulitin/internal/routes"
	"github.com/maksimulitin/lib/logger"
	"github.com/maksimulitin/lib/serverutils"
	"log"
	"log/slog"
	"os"
)

func main() {
	config.LoadConfigEnv()

	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = "8084"
	}

	serverPortFallback := os.Getenv("SERVER_PORT_FALLBACK")
	if serverPortFallback == "" {
		serverPortFallback = "8085"
	}

	logger.Info("Starting application initialization")

	app := controllers.NewApplication(
		database.ProductData(database.Client, "Products"),
		database.UserData(database.Client, "Users"),
	)

	logger.Info("Application controllers initialized successfully")

	router := gin.New()
	router.Use(gin.Logger())
	routes.SetupRoutes(router, app)

	logger.Info("Router configured successfully", slog.String("port", serverPort), slog.Any("routes", router.Routes()))
	logger.Info("Attempting to start server on port " + serverPort)

	if err := serverutils.TryRunServer(router, serverPort); err != nil {
		logger.Warn("Main port is occupied, trying fallback port", slog.String("fallbackPort", serverPortFallback))
		if err := serverutils.TryRunServer(router, serverPortFallback); err != nil {
			logger.Error("Server failed to start on both main and fallback ports", slog.Any("error", err))
			log.Fatal(err)
		}
	}
}
