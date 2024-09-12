package utils

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DatabaseClient *mongo.Client

func ConnectDB() *mongo.Client {
	ctx := context.Background()

	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		log.Println("Failed to get DB_URL")
		os.Exit(1)
	}

	clientOptions := options.Client().ApplyURI(dbUrl)

	databaseClient, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Println("Failed to connect to MongoDB:", err)
	}

	err = databaseClient.Ping(ctx, nil)
	if err != nil {
		log.Println("Failed to ping MongoDB:", err)
	}

	log.Println("Connected to MongoDB")

	return databaseClient
}

func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	collection := client.Database("Tours").Collection(collectionName)
	return collection
}
