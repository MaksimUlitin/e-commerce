package controllers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/maksimulitin/internal/db"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Application struct {
	ProductCollection *mongo.Collection
	UserCollection    *mongo.Collection
}

func NewApplication(CollectionProduct, CollectionUser *mongo.Collection) *Application {
	return &Application{
		ProductCollection: ProductCollection,
		UserCollection:    UserCollection,
	}
}

func (app *Application) CreatingСart() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ProductQueryId := ctx.Query("id")
		if ProductQueryId == "" {
			log.Println("product id is emty")
			_ = ctx.AbortWithError(http.StatusBadRequest, errors.New("product id is emty"))
			return
		}

		userQueryId := ctx.Query("userId")
		if userQueryId == "" {
			log.Println("user userid is emty")
			_ = ctx.AbortWithError(http.StatusBadRequest, errors.New("user id is emty"))
			return
		}

		productID, err := primitive.ObjectIDFromHex(ProductQueryId)
		if err != nil {
			log.Println(err)
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		var c, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err = db.AddProductToCart(c, app.ProductCollection, app.UserCollection, productID, userQueryId)
		if err != nil {
			ctx.IndentedJSON(500, err)

		}
		ctx.IndentedJSON(200, "Successfully Added to the cart")
	}
}

func (app *Application) DeleteCart() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		productQueryID := ctx.Query("id")
		if productQueryID == "" {
			log.Println("product id is inavalid")
			_ = ctx.AbortWithError(http.StatusBadRequest, errors.New("product id is empty"))
			return
		}

		userQueryID := ctx.Query("userID")
		if userQueryID == "" {
			log.Println("user id is empty")
			_ = ctx.AbortWithError(http.StatusBadRequest, errors.New("UserID is empty"))
		}

		ProductID, err := primitive.ObjectIDFromHex(productQueryID)
		if err != nil {
			log.Println(err)
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		var c, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err = db.RemoveCart(ctx, app.ProductCollection, app.UserCollection, ProductID, userQueryID)
		if err != nil {
			ctx.IndentedJSON(http.StatusInternalServerError, err)
			return
		}
		ctx.IndentedJSON(200, "Successfully removed from cart")
	}
}

func GetProductFromCart() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}

func (app *Application) BuyProductFromCart() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}

func (app *Application) FastBuy() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}
