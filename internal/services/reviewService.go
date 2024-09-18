package services

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/hamid-nazari/tours-in-go/internal/models"
	"github.com/hamid-nazari/tours-in-go/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
)

func CreateReview(ctx *gin.Context, review *models.Review) error {
	collection := utils.GetCollection(TourDatabaseClient, "reviews")

	_, err := collection.InsertOne(ctx, review)
	if err != nil {
		return err
	}
	return nil
}

func GetAllReviews(ctx *gin.Context) []models.Review {
	collection := utils.GetCollection(TourDatabaseClient, "reviews")

	var reviews []models.Review

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil
	}
	for cursor.Next(ctx) {
		var review models.Review
		cursor.Decode(&review)
		reviews = append(reviews, review)
	}
	return reviews
}

func GetReviewById(ctx *gin.Context, id string) *models.Review {
	return nil
}
func ValidateReview(review models.Review) error {
	validator := validator.New()
	err := validator.Struct(review)
	if err != nil {
		return err
	}
	return nil
}

func UpdateReview(ctx *gin.Context, review *models.Review) error {
	collection := utils.GetCollection(TourDatabaseClient, "reviews")

	_, err := collection.UpdateOne(ctx, bson.M{"id": review.Id}, bson.M{"$set": review})
	if err != nil {
		return err
	}
	return nil
}
func DeleteReview(ctx *gin.Context, reviewId string) error {
	collection := utils.GetCollection(TourDatabaseClient, "reviews")

	_, err := collection.DeleteOne(ctx, bson.M{"id": reviewId})
	if err != nil {
		return err
	}
	return nil
}
