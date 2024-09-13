package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/hamid-nazari/tours-in-go/internal/controllers"
)

func SetupRoutes(router *gin.RouterGroup) {

	router.POST("/signup", controllers.SignupHandler)
	router.POST("/login", controllers.LoginHandler)
	router.POST("/logout", controllers.LogoutHandler)
	router.POST("/forgot-password", controllers.ForgotPasswordHandler)
	router.POST("/reset-password", controllers.ResetPasswordHandler)

	router.Use(controllers.ProtectHandler)

	router.PATCH("/update-password", controllers.UpdatePasswordHandler)
	router.PATCH("/update-me", controllers.UpdateMeHandler)
	router.DELETE("/delete-me", controllers.DeleteMeHandler)
	router.GET("/me", controllers.GetMeHandler)

	router.Use(controllers.RestrictTo("admin"))

	router.POST("/", controllers.CreateUserHandler)
	router.GET("/", controllers.GetAllUsersHandler)
	router.DELETE("/", controllers.DeleteAllUsersHandler)
	router.GET("/:id", controllers.GetUserHandler)
	router.PATCH("/:id", controllers.UpdateUserHandler)
	router.DELETE("/:id", controllers.DeleteUserdHandler)

}
