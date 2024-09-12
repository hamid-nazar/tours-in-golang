package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/hamid-nazari/tours-in-go/internal/routes"
	"github.com/hamid-nazari/tours-in-go/internal/services"
	"github.com/hamid-nazari/tours-in-go/internal/utils"
)

func main() {
	godotenv.Load("../.env")

	var databaseClient *mongo.Client = utils.ConnectDB()

	services.UserDatabaseClient = databaseClient

	router := gin.Default()

	usersRouter := router.Group("api/v1/users")
	// toursRouter := router.Group("api/tours")

	routes.SetupRoutes(usersRouter)

	if err := router.Run(":8000"); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Server started on port 8080")

	// defer databaseClient.Disconnect(context.TODO())
}
