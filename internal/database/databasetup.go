package database

import (
	"context"
	"errors"
	"fmt"
	"github.com/maksimulitin/config"
	"github.com/maksimulitin/lib/logger"
	"log"
	"log/slog"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MONGO_URI=mongodb://${MONGO_USER}:${MONGO_PASSWORD}@${MONGO_HOST}:${MONGO_PORT}/${MONGO_DB}?retryWrites=true&w=majority
// "mongodb://development:testpassword@localhost:27017"
func DBSet() *mongo.Client {
	config.LoadConfigEnv()

	mongoUser := os.Getenv("MONGO_USER")
	mongoPassword := os.Getenv("MONGO_PASSWORD")
	mongoHost := os.Getenv("MONGO_HOST")
	mongoPort := os.Getenv("MONGO_PORT")

	mongoURI := fmt.Sprintf("mongodb://%s:%s@%s:%s",
		mongoUser, mongoPassword, mongoHost, mongoPort)

	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))

	if err != nil {
		logger.Error("failed to create mongo client", slog.Any("error", err))
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)

	if err != nil {
		logger.Error("failed to connect to mongo", slog.Any("error", err))
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)

	if err != nil {
		logger.Error("failed to ping mongo", slog.Any("error", err))
		return nil
	}

	logger.Info("successfully connected to mongo")
	return client
}

var Client *mongo.Client = DBSet()

func UserData(client *mongo.Client, CollectionName string) *mongo.Collection {
	if client == nil {
		logger.Error("mongo client is nil", slog.Any("error", errors.New("mongo client is nil")))
		return nil
	}

	var collection *mongo.Collection = client.Database("Ecommerce").Collection(CollectionName)
	return collection
}

func ProductData(client *mongo.Client, CollectionName string) *mongo.Collection {
	if client == nil {
		logger.Error("mongo client is nil", slog.Any("error", errors.New("mongo client is nil")))
		return nil
	}

	var productCollection *mongo.Collection = client.Database("Ecommerce").Collection(CollectionName)
	return productCollection
}
