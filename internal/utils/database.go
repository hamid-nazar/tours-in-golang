package utils

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB() *mongo.Client {
	ctx := context.Background()

	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		log.Fatalln("Failed to get DB_URL")
	}

	clientOptions := options.Client().ApplyURI(dbUrl)

	databaseClient, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalln("Failed to connect to MongoDB:", err)
	}

	err = databaseClient.Ping(ctx, nil)
	if err != nil {
		log.Fatalln("Failed to ping MongoDB:", err)
	}
	log.SetPrefix("Info: ")
	log.Println("Connected to MongoDB")
	log.Println()
	return databaseClient
}

func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	collection := client.Database("Tours").Collection(collectionName)
	return collection
}
