package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/maksimulitin/internal/controllers"
	"github.com/maksimulitin/internal/middleware"
)

func setupAdminRoutes(router *gin.Engine) {
	admin := router.Group("/admin")
	admin.Use(middleware.Authentication())
	{
		admin.POST("/products/add", controllers.ProductViewerAdmin())
	}
}
