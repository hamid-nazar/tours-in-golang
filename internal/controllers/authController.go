package controllers

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/hamid-nazari/tours-in-go/internal/models"
	"github.com/hamid-nazari/tours-in-go/internal/services"
)

func CreateJwtTokenAndSend(c *gin.Context, user *models.User) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": user.Id,
		"exp":    time.Now().Add(time.Hour * 24).Unix(),
	})

	signedToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.CustomResponse{
			Status:  "Failed",
			Message: err.Error(),
			Data:    nil,
		})
		return

	}

	c.JSON(http.StatusOK, models.CustomResponse{
		Status:  "Success",
		Message: "User created successfully",
		Data: gin.H{
			"token": signedToken,
			"user":  user,
		},
	})

}
func SignupHandler(c *gin.Context) {
	user := models.NewUser()

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, models.CustomResponse{
			Status:  "Failed",
			Message: err.Error(),
			Data:    nil,
		})
		return
	}
	if existingUser := services.FindUserByEmail(c, user.Email); existingUser != nil {
		c.JSON(http.StatusConflict, models.CustomResponse{
			Status:  "Failed",
			Message: "A user with this email already exists",
			Data:    nil,
		})
		return
	}

	if err := services.ValidateUser(*user); err != nil {
		c.JSON(http.StatusBadRequest, models.CustomResponse{
			Status:  "Failed",
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	hashedPassword, err := services.HashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.CustomResponse{
			Status:  "Failed",
			Message: err.Error(),
			Data:    nil,
		})
		return
	}
	user.Password = hashedPassword
	user.PasswordConfirm = ""

	if err := services.CreateUser(c, user); err != nil {
		c.JSON(http.StatusInternalServerError, models.CustomResponse{
			Status:  "Failed",
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	CreateJwtTokenAndSend(c, user)
}
func LoginHandler(c *gin.Context) {

}

func UpdatePasswordHandler(c *gin.Context) {
}

func ForgotPasswordHandler(c *gin.Context) {
}
