package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/maksimulitin/internal/controllers"
	"github.com/maksimulitin/internal/db"
)

func main() {
	err := os.Getenv("PORT")
	if err == "" {
		port := 8000
	}

	app := controllers.NewApplication(db.AbuotProduct(db.Client, "products"), db.AbuotUser(db.Client, "user"))
	routes := gin.Default()

	routes
}
