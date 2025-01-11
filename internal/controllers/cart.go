package controllers

import (
	"context"
	"errors"
	"github.com/maksimulitin/internal/database"
	"github.com/maksimulitin/internal/models"
	"github.com/maksimulitin/lib/logger"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Application struct {
	prodCollection *mongo.Collection
	userCollection *mongo.Collection
}

func NewApplication(prodCollection, userCollection *mongo.Collection) *Application {
	return &Application{
		prodCollection: prodCollection,
		userCollection: userCollection,
	}
}

func (app *Application) AddToCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		productQueryID := c.Query("id")
		if productQueryID == "" {
			logger.Error("Product ID is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("product id is empty"))
			return
		}
		userQueryID := c.Query("userID")
		if userQueryID == "" {
			logger.Error("User ID is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("user id is empty"))
			return
		}
		productID, err := primitive.ObjectIDFromHex(productQueryID)
		if err != nil {
			logger.Error("Invalid product ID", slog.Any("error", err))
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err = database.AddProductToCart(ctx, app.prodCollection, app.userCollection, productID, userQueryID)
		if err != nil {
			logger.Error("Failed to add product to cart", slog.Any("error", err))
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}
		logger.Info("Product successfully added to cart", slog.String("productID", productQueryID), slog.String("userID", userQueryID))
		c.IndentedJSON(200, "Successfully added to the cart")
	}
}

func (app *Application) RemoveItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		productQueryID := c.Query("id")
		if productQueryID == "" {
			logger.Error("Product ID is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("product id is empty"))
			return
		}

		userQueryID := c.Query("userID")
		if userQueryID == "" {
			logger.Error("User ID is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("user id is empty"))
			return
		}

		ProductID, err := primitive.ObjectIDFromHex(productQueryID)
		if err != nil {
			logger.Error("Invalid product ID", slog.Any("error", err))
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err = database.RemoveCartItem(ctx, app.prodCollection, app.userCollection, ProductID, userQueryID)
		if err != nil {
			logger.Error("Failed to remove product from cart", slog.Any("error", err))
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}
		logger.Info("Product successfully removed from cart", slog.String("productID", productQueryID), slog.String("userID", userQueryID))
		c.IndentedJSON(200, "Successfully removed from cart")
	}
}

func GetItemFromCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Query("id")
		if userId == "" {
			logger.Error("User ID is empty")
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error": "invalid id"})
			c.Abort()
			return
		}

		usertId, _ := primitive.ObjectIDFromHex(userId)

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var filledCart models.User
		err := UserCollection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: usertId}}).Decode(&filledCart)
		if err != nil {
			logger.Error("Failed to find user cart", slog.Any("error", err))
			c.IndentedJSON(500, "not id found")
			return
		}

		filterMatch := bson.D{{Key: "$match", Value: bson.D{primitive.E{Key: "_id", Value: usertId}}}}
		unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$usercart"}}}}
		grouping := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$_id"}, {Key: "total", Value: bson.D{primitive.E{Key: "$sum", Value: "$usercart.price"}}}}}}
		pointCursor, err := UserCollection.Aggregate(ctx, mongo.Pipeline{filterMatch, unwind, grouping})
		if err != nil {
			logger.Error("Failed to aggregate cart data", slog.Any("error", err))
		}
		var listing []bson.M
		if err = pointCursor.All(ctx, &listing); err != nil {
			logger.Error("Failed to decode aggregated data", slog.Any("error", err))
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		for _, json := range listing {
			logger.Info("Cart data retrieved successfully", slog.String("userID", userId))
			c.IndentedJSON(200, json["total"])
			c.IndentedJSON(200, filledCart.UserCart)
		}
		ctx.Done()
	}
}

func (app *Application) BuyFromCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		userQueryID := c.Query("id")
		if userQueryID == "" {
			logger.Error("User ID is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("user id is empty"))
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		err := database.BuyItemFromCart(ctx, app.userCollection, userQueryID)
		if err != nil {
			logger.Error("Failed to buy items from cart", slog.Any("error", err))
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}
		logger.Info("Items successfully purchased from cart", slog.String("userID", userQueryID))
		c.IndentedJSON(200, "Successfully placed the order")
	}
}

func (app *Application) InstantBuy() gin.HandlerFunc {
	return func(c *gin.Context) {
		UserQueryID := c.Query("userid")
		if UserQueryID == "" {
			logger.Error("User ID is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("user id is empty"))
			return
		}
		ProductQueryID := c.Query("pid")
		if ProductQueryID == "" {
			logger.Error("Product ID is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("product id is empty"))
			return
		}
		productID, err := primitive.ObjectIDFromHex(ProductQueryID)
		if err != nil {
			logger.Error("Invalid product ID", slog.Any("error", err))
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err = database.InstantBuyer(ctx, app.prodCollection, app.userCollection, productID, UserQueryID)
		if err != nil {
			logger.Error("Failed to place instant buy order", slog.Any("error", err))
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}
		logger.Info("Instant buy order placed successfully", slog.String("productID", ProductQueryID), slog.String("userID", UserQueryID))
		c.IndentedJSON(200, "Successfully placed the order")
	}
}
