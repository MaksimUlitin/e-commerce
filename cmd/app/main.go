package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/maksimulitin/internal/controllers"
	"github.com/maksimulitin/internal/db"
	"github.com/maksimulitin/internal/routes"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	app := controllers.NewApplication(db.AbuotProduct(db.Client, "products"), db.AbuotUser(db.Client, "user"))
	router := gin.New()

	router.Use(gin.Logger())
	routes.UserRoutes(router)
	router.GET("/creatingcart", app.CreatingСart())
	router.GET("/deletecart", app.DeleteCart())
	router.GET("/listcart", controllers.GetProductFromCart())
	router.GET("/buycart", app.BuyProductFromCart())
	router.GET("/fustbuy", app.FastBuy())
	log.Fatal(router.Run(":" + port))

}
