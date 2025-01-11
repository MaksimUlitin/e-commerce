package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/maksimulitin/internal/controllers"
)

func UserRoutes(r *gin.Engine) {
	r.POST("/users/signup", controllers.SignUp())
	r.POST("/users/login", controllers.Login())
	r.POST("/admin/addproduct", controllers.ProductViewerAdmin())
	r.GET("/users/productview", controllers.SearchProduct())
	r.GET("/users/search", controllers.SearchProductByQuery())
}
