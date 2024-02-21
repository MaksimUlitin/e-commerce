package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/maksimulitin/internal/controllers"
)

func UserRoutes(Router *gin.Engine) {
	Router.POST("/users/signup", controllers.Signup())
	Router.POST("/users/login", controllers.Login())
	Router.POST("/admin/addproduct", controllers.AddProduct())
	Router.GET("/users/view", controllers.ViewProduct())
	Router.GET("/users/serch", controllers.SerchProduct())
}
