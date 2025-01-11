package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/maksimulitin/internal/controllers"
)

func SetupRoutes(router *gin.Engine, app *controllers.Application) {
	setupUserRoutes(router)
	setupCartRoutes(router, app)
	setupAddressRoutes(router)
	setupAdminRoutes(router)
}
