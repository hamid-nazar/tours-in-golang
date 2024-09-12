package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/hamid-nazari/tours-in-go/internal/controllers"
)

func SetupRoutes(router *gin.RouterGroup) {

	router.POST("/", controllers.RegisterUserHandler)

	router.GET("/", controllers.GetAllUsersHandler)

	router.GET("/:id", controllers.GetUserHandler)

	router.PATCH("/:id", controllers.UpdateUserHandler)

	router.DELETE("/", controllers.DeleteAllUsersHandler)

	router.DELETE("/:id", controllers.DeleteUserdHandler)

}
