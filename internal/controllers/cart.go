package controllers

import (
	"github.com/gin-gonic/gin"
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

	}
}

func (app *Application) DeleteCart() gin.HandlerFunc {
	return func(ctx *gin.Context) {

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
