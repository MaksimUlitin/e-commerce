package main

import (
	"github.com/maksimulitin/internal/controllers"
	"github.com/maksimulitin/internal/database"
	"github.com/maksimulitin/internal/middleware"
	"github.com/maksimulitin/internal/routes"
	"log"

	"github.com/gin-gonic/gin"
)

/*
TODO
 1. ADD SWAGGER
 2 CHANGE LOGIC
 3 REWRITE README
 4. ADD LOGGER
*/
func main() {
	app := controllers.NewApplication(database.ProductData(database.Client, "Products"), database.UserData(database.Client, "Users"))

	router := gin.New()
	router.Use(gin.Logger())
	routes.UserRoutes(router)
	router.Use(middleware.Authentication())
	router.GET("/addtocart", app.AddToCart())
	router.GET("/removeitem", app.RemoveItem())
	router.GET("/listcart", controllers.GetItemFromCart())
	router.POST("/addaddress", controllers.AddAddress())
	router.PUT("/edithomeaddress", controllers.EditHomeAddress())
	router.PUT("/editworkaddress", controllers.EditWorkAddress())
	router.GET("/deleteaddresses", controllers.DeleteAddress())
	router.GET("/cartcheckout", app.BuyFromCart())
	router.GET("/instantbuy", app.InstantBuy())
	log.Fatal(router.Run(":" + "8084"))
}
