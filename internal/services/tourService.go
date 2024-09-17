package services

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/hamid-nazari/tours-in-go/internal/models"
	"github.com/hamid-nazari/tours-in-go/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var TourDatabaseClient *mongo.Client

func CreateTour(ctx *gin.Context, tour *models.Tour) error {
	collection := utils.GetCollection(TourDatabaseClient, "tours")
	_, err := collection.InsertOne(ctx, tour)
	if err != nil {
		return err
	}
	return nil
}

func GetAllTours(ctx *gin.Context) ([]models.Tour, error) {
	collection := utils.GetCollection(TourDatabaseClient, "tours")
	var tours []models.Tour
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to find tours: %v", err)
	}
	for cursor.Next(ctx) {
		var tour models.Tour
		cursor.Decode(&tour)
		tours = append(tours, tour)
	}
	return tours, nil
}

func GetTour(ctx *gin.Context, tour *models.Tour) (*models.Tour, error) {
	collection := utils.GetCollection(TourDatabaseClient, "tours")
	err := collection.FindOne(ctx, bson.M{"id": tour.Id}).Decode(&tour)
	if err != nil {
		return nil, fmt.Errorf("failed to find tour: %v", err)
	}
	return tour, nil
}

func UpdateTour(ctx *gin.Context, tour *models.Tour) error {
	collection := utils.GetCollection(TourDatabaseClient, "tours")
	_, err := collection.UpdateOne(ctx, bson.M{"id": tour.Id}, bson.M{"$set": tour})
	if err != nil {
		return err
	}
	return nil
}

func DeleteTour(ctx *gin.Context, tour *models.Tour) error {
	collection := utils.GetCollection(TourDatabaseClient, "tours")
	_, err := collection.DeleteOne(ctx, bson.M{"id": tour.Id})
	if err != nil {
		return err
	}
	return nil
}
