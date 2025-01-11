package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/maksimulitin/internal/controllers"
	"github.com/maksimulitin/internal/middleware"
)

func setupCartRoutes(router *gin.Engine, app *controllers.Application) {
	cart := router.Group("/cart")
	cart.Use(middleware.Authentication())
	{
		cart.GET("/add", app.AddToCart())
		cart.GET("/remove", app.RemoveItem())
		cart.GET("/list", controllers.GetItemFromCart())
		cart.GET("/checkout", app.BuyFromCart())
		cart.GET("/buy", app.InstantBuy())
	}
}
