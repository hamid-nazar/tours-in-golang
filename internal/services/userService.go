package services

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/hamid-nazari/tours-in-go/internal/models"
	"github.com/hamid-nazari/tours-in-go/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var UserDatabaseClient *mongo.Client

func CreateUser(cxt *gin.Context, user *models.User) error {
	collection := utils.GetCollection(UserDatabaseClient, "users")
	_, err := collection.InsertOne(cxt, user)
	if err != nil {
		return err
	}

	return nil
}

func GetAllUsers(ctx *gin.Context) ([]models.User, error) {
	collection := utils.GetCollection(UserDatabaseClient, "users")

	var users []models.User

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to find users: %v", err)
	}

	for cursor.Next(ctx) {
		var user models.User
		cursor.Decode(&user)
		users = append(users, user)
	}

	return users, nil
}

func FindUserByEmail(ctx *gin.Context, email string) *models.User {
	collection := utils.GetCollection(UserDatabaseClient, "users")

	var user models.User

	err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil
	}

	return &user
}

func FindUserById(ctx *gin.Context, id string) *models.User {
	collection := utils.GetCollection(UserDatabaseClient, "users")

	var user models.User

	err := collection.FindOne(ctx, bson.M{"id": id}).Decode(&user)
	if err != nil {
		return nil
	}

	return &user
}

func UpdateUser(ctx *gin.Context, user *models.User) error {
	collection := utils.GetCollection(UserDatabaseClient, "users")

	_, err := collection.UpdateOne(ctx, bson.M{"id": user.Id}, bson.M{"$set": user})
	if err != nil {
		return fmt.Errorf("failed to update user: %v", err)
	}

	return nil
}

func DeleteUser(ctx *gin.Context, id string) error {
	collection := utils.GetCollection(UserDatabaseClient, "users")

	_, err := collection.DeleteOne(ctx, bson.M{"id": id})
	if err != nil {
		return fmt.Errorf("failed to delete user: %v", err)
	}

	return nil
}

func DeleteAllUsers(ctx *gin.Context) error {
	collection := utils.GetCollection(UserDatabaseClient, "users")

	_, err := collection.DeleteMany(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("failed to delete all users: %v", err)
	}

	return nil
}

func ValidateUser(user models.User) error {
	err := validator.New().Struct(user)
	if err != nil {
		return fmt.Errorf("invalid user: %v", err.(validator.ValidationErrors))
	}
	return nil
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %v", err)
	}
	return string(hashedPassword), nil
}

func ComparePassword(providedPassword string, actualPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(actualPassword))
}

func FindActiveUsers(users []models.User) []models.User {

	var activeUsers []models.User

	for _, user := range users {
		if user.Active {
			activeUsers = append(activeUsers, user)
		}
	}

	return activeUsers
}
