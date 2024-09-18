package services

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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

func GetAllTours(ctx *gin.Context) []models.Tour {
	collection := utils.GetCollection(TourDatabaseClient, "tours")
	var tours []models.Tour
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil
	}
	for cursor.Next(ctx) {
		var tour models.Tour
		cursor.Decode(&tour)
		tours = append(tours, tour)
	}
	return tours
}

func FindTourById(ctx *gin.Context, id string) *models.Tour {
	collection := utils.GetCollection(TourDatabaseClient, "tours")

	var tour models.Tour

	err := collection.FindOne(ctx, bson.M{"id": id}).Decode(&tour)
	if err != nil {
		return nil
	}
	return &tour
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

func ValidateTour(tour models.Tour) error {
	err := validator.New().Struct(tour)
	if err != nil {
		return err.(validator.ValidationErrors)
	}
	return nil
}
