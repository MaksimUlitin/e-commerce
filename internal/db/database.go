package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DbSet() *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://development:testpassword@localhost:27017"))
	if err != nil {
		log.Fatal()
	}

	ctx, cancle := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal()
	}

	defer cancle()

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Println("failed to connection to mongo")

		return nil
	}

	fmt.Println("Successfully Connected to the mongodb")

	return client
}

var Client *mongo.Client = DbSet()

func SignupDb(Client *mongo.Client, colName string) *mongo.Collection {
	var collectiondb *mongo.Collection = Client.Database("Ecommerce").Collection(colName)

	return collectiondb
}
