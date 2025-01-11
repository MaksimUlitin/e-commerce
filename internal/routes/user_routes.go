package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/maksimulitin/internal/controllers"
)

func setupUserRoutes(router *gin.Engine) {
	public := router.Group("/users")
	{
		public.POST("/signup", controllers.SignUp())
		public.POST("/login", controllers.Login())
		public.GET("/productview", controllers.SearchProduct())
		public.GET("/search", controllers.SearchProductByQuery())
	}
}
