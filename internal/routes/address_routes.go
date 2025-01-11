package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/maksimulitin/internal/controllers"
	"github.com/maksimulitin/internal/middleware"
)

func setupAddressRoutes(router *gin.Engine) {
	address := router.Group("/address")
	address.Use(middleware.Authentication())
	{
		address.POST("/add", controllers.AddAddress())
		address.PUT("/edit/home", controllers.EditHomeAddress())
		address.PUT("/edit/work", controllers.EditWorkAddress())
		address.DELETE("/delete", controllers.DeleteAddress())
	}
}
