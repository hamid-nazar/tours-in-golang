package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/hamid-nazari/tours-in-go/internal/controllers"
	"github.com/hamid-nazari/tours-in-go/internal/middleware"
)

func SetupTourRoutes(router *gin.RouterGroup) {

	router.POST("/", controllers.ProtectHandler, controllers.RestrictTo("admin", "lead-guide"), controllers.CreateTourHandler)
	router.GET("/", controllers.GetAllToursHandler)

	router.GET("/id", controllers.ProtectHandler, controllers.RestrictTo("admin", "lead-guide"), controllers.GetTourHandler)
	router.PATCH("/id", controllers.ProtectHandler, controllers.RestrictTo("admin", "lead-guide"), controllers.UpdateTourHandler)
	router.DELETE("/id", controllers.ProtectHandler, controllers.RestrictTo("admin", "lead-guide"), controllers.DeleteTourHandler)

	router.GET("/top-5-cheap", middleware.AliasTopTours, controllers.GetAllToursHandler)

	// router.GET("/tours-within/:distance/center/:latlng/unit/:unit", controllers.GetToursWithinHandler)
	// router.GET("/distances/:latlng/unit/:unit", controllers.GetDistancesHandler)
}
