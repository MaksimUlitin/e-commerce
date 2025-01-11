package main

import (
	"github.com/gin-gonic/gin"
	"github.com/maksimulitin/internal/controllers"
	"github.com/maksimulitin/internal/database"
	"github.com/maksimulitin/internal/routes"
	"github.com/maksimulitin/lib/logger"
	"log"
	"log/slog"
)

func main() {
	logger.Info("Starting application initialization")

	app := controllers.NewApplication(
		database.ProductData(database.Client, "Products"),
		database.UserData(database.Client, "Users"),
	)

	logger.Info("Application controllers initialized successfully")

	router := gin.New()
	router.Use(gin.Logger())
	routes.SetupRoutes(router, app)

	logger.Info("Router configured successfully", slog.String("port", "8084"), slog.Any("routes", router.Routes()))
	logger.Info("Starting server on port 8084")

	if err := router.Run(":8084"); err != nil {
		logger.Error("Server failed to start", slog.Any("error", err))
		log.Fatal(err)
	}
}
