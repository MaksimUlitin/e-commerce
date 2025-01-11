package controllers

import (
	"context"
	"errors"
	"github.com/maksimulitin/internal/models"
	"github.com/maksimulitin/lib/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func AddAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Query("id")

		if userID == "" {
			logger.Error("User ID is empty")
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error": "Invalid ID"})
			c.Abort()
			return
		}

		address, err := primitive.ObjectIDFromHex(userID)

		if err != nil {
			logger.Error("Invalid user ID format", slog.Any("error", err))
			c.IndentedJSON(http.StatusInternalServerError, "Internal Server Error")
			return
		}

		var addresses models.Address
		addresses.AddressId = primitive.NewObjectID()

		if err := c.BindJSON(&addresses); err != nil {
			logger.Error("Failed to bind JSON to address struct", slog.Any("error", err))
			c.IndentedJSON(http.StatusNotAcceptable, err.Error())
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		matchFilter := bson.D{{Key: "$match",
			Value: bson.D{{Key: "_id",
				Value: address,
			}}}}

		unwind := bson.D{{Key: "$unwind",
			Value: bson.D{{Key: "path",
				Value: "$address",
			}}}}

		group := bson.D{{Key: "$group",
			Value: bson.D{{Key: "_id",
				Value: "$address_id"},
				{Key: "count",
					Value: bson.D{{Key: "$sum",
						Value: 1,
					}}}}}}

		pointCursor, err := UserCollection.Aggregate(ctx, mongo.Pipeline{matchFilter, unwind, group})

		if err != nil {
			logger.Error("Failed to aggregate address data", slog.Any("error", err))
			c.IndentedJSON(http.StatusInternalServerError, "Internal Server Error")
			return
		}

		var addressInfo []bson.M

		if err = pointCursor.All(ctx, &addressInfo); err != nil {
			logger.Error("Failed to decode aggregation results", slog.Any("error", err))
			c.IndentedJSON(http.StatusInternalServerError, "Internal Server Error")
			return
		}

		var size int32
		for _, addressNo := range addressInfo {
			count := addressNo["count"]
			size = count.(int32)
		}

		if size < 2 {
			filter := bson.D{{Key: "_id",
				Value: address,
			}}

			update := bson.D{{Key: "$push",
				Value: bson.D{{Key: "address",
					Value: addresses,
				}}}}

			_, err := UserCollection.UpdateOne(ctx, filter, update)

			if err != nil {
				logger.Error("Failed to update user address", slog.Any("error", err))
				c.IndentedJSON(http.StatusInternalServerError, "Internal Server Error")
				return
			}

			logger.Info("Address added successfully", slog.String("userID", userID))
			c.IndentedJSON(http.StatusCreated, "Address added successfully")

		} else {
			logger.Warn("Address limit exceeded for user", slog.String("userID", userID))
			c.IndentedJSON(http.StatusBadRequest, "Address limit exceeded")
		}
	}
}

func EditHomeAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Query("id")

		if userID == "" {
			logger.Error("User ID is empty")
			c.JSON(http.StatusNotFound, gin.H{"error": "Invalid ID"})
			c.Abort()
			return
		}

		userObjID, err := primitive.ObjectIDFromHex(userID)

		if err != nil {
			logger.Error("Invalid user ID format", slog.Any("error", err))
			c.IndentedJSON(http.StatusInternalServerError, "Internal Server Error")
			return
		}

		var editAddress models.Address

		if err := c.BindJSON(&editAddress); err != nil {
			logger.Error("Failed to bind JSON to address struct", slog.Any("error", err))
			c.IndentedJSON(http.StatusBadRequest, err.Error())
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		filter := bson.D{{Key: "_id", Value: userObjID}}
		update := bson.D{{Key: "$set", Value: bson.D{
			{Key: "address.0.house_name", Value: editAddress.House},
			{Key: "address.0.street_name", Value: editAddress.Street},
			{Key: "address.0.city_name", Value: editAddress.City},
			{Key: "address.0.pin_code", Value: editAddress.PinCode},
		}}}

		_, err = UserCollection.UpdateOne(ctx, filter, update)

		if err != nil {
			logger.Error("Failed to update home address", slog.Any("error", err))
			c.IndentedJSON(http.StatusInternalServerError, "Something went wrong")
			return
		}

		logger.Info("Home address updated successfully", slog.String("userID", userID))
		c.IndentedJSON(http.StatusOK, "Home address updated successfully")
	}
}

func EditWorkAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.Info("EditWorkAddress handler invoked")
		userId := c.Query("id")

		if userId == "" {
			logger.Error("User ID not provided")
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error": "Wrong id not provided"})
			c.Abort()
			return
		}

		usertId, err := primitive.ObjectIDFromHex(userId)

		if err != nil {
			logger.Error("Invalid User ID", slog.Any("error", err))
			c.IndentedJSON(500, err)
			return
		}

		var editAddress models.Address

		if err := c.BindJSON(&editAddress); err != nil {
			logger.Error("Failed to bind JSON", slog.Any("error", err))
			c.IndentedJSON(http.StatusBadRequest, err.Error())
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		filter := bson.D{primitive.E{Key: "_id", Value: usertId}}
		update := bson.D{{Key: "$set", Value: bson.D{
			primitive.E{Key: "address.1.house_name", Value: editAddress.House},
			{Key: "address.1.street_name", Value: editAddress.Street},
			{Key: "address.1.city_name", Value: editAddress.City},
			{Key: "address.1.pin_code", Value: editAddress.PinCode},
		}}}

		_, err = UserCollection.UpdateOne(ctx, filter, update)

		if err != nil {
			logger.Error("Failed to update work address", slog.Any("error", err))
			c.IndentedJSON(500, "Something went wrong")
			return
		}

		logger.Info("Successfully updated the work address", slog.String("userID", userId))
		c.IndentedJSON(200, "Successfully updated the Work Address")
	}
}

func DeleteAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.Info("DeleteAddress handler invoked")

		userId := c.Query("id")

		if userId == "" {
			logger.Error("Invalid Search Index", slog.Any("error", errors.New("Invalid ID")))
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error": "Invalid Search Index"})
			c.Abort()
			return
		}

		addresses := make([]models.Address, 0)
		usertId, err := primitive.ObjectIDFromHex(userId)

		if err != nil {
			logger.Error("Invalid User ID", slog.Any("error", err))
			c.IndentedJSON(500, "Internal Server Error")
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		filter := bson.D{primitive.E{Key: "_id", Value: usertId}}
		update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "address", Value: addresses}}}}
		_, err = UserCollection.UpdateOne(ctx, filter, update)

		if err != nil {
			logger.Error("Failed to delete address", slog.Any("error", err))
			c.IndentedJSON(404, "Wrong")
			return
		}

		logger.Info("Successfully deleted address", slog.String("userID", userId))
		c.IndentedJSON(200, "Successfully Deleted!")
	}
}
