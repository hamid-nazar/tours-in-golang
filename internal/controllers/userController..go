package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hamid-nazari/tours-in-go/internal/models"
	"github.com/hamid-nazari/tours-in-go/internal/services"
)

func RegisterUserHandler(c *gin.Context) {
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

	c.JSON(http.StatusCreated, models.CustomResponse{
		Status:  "Success",
		Message: "User created successfully",
		Data:    user,
	})
}

func GetAllUsersHandler(c *gin.Context) {

	users, err := services.GetAllUsers(c)
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
		Message: fmt.Sprint(len(users)) + " users found",
		Data:    users,
	})
}

func GetUserHandler(c *gin.Context) {

	id := c.Param("id")

	user := services.FindUserById(c, id)
	if user == nil {
		c.JSON(http.StatusNotFound, models.CustomResponse{
			Status:  "Failed",
			Message: "User not found",
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, models.CustomResponse{
		Status:  "Success",
		Message: "User found",
		Data:    user,
	})
}

func UpdateUserHandler(c *gin.Context) {

	id := c.Param("id")

	user := services.FindUserById(c, id)
	if user == nil {
		c.JSON(http.StatusNotFound, models.CustomResponse{
			Status:  "Failed",
			Message: "User not found",
			Data:    nil,
		})
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, models.CustomResponse{
			Status:  "Failed",
			Message: err.Error(),
			Data:    nil,
		})
		return
	}
	name := form.Value["name"][0]
	photo := form.File["photo"][0].Filename
	if name == "" {
		c.JSON(http.StatusBadRequest, models.CustomResponse{
			Status:  "Failed",
			Message: "Name is required",
			Data:    nil,
		})
		return
	}
	user.Name = name
	user.Photo = photo

	if err := services.UpdateUser(c, user); err != nil {
		c.JSON(http.StatusInternalServerError, models.CustomResponse{
			Status:  "Failed",
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, models.CustomResponse{
		Status:  "Success",
		Message: "User updated successfully",
		Data:    user,
	})
}

func DeleteAllUsersHandler(c *gin.Context) {
	users, err := services.GetAllUsers(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.CustomResponse{
			Status:  "Failed",
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	if err := services.DeleteAllUsers(c); err != nil {
		c.JSON(http.StatusInternalServerError, models.CustomResponse{
			Status:  "Failed",
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, models.CustomResponse{
		Status:  "Success",
		Message: fmt.Sprint(len(users)) + " users deleted",
		Data:    nil,
	})
}

func DeleteUserdHandler(c *gin.Context) {

	id := c.Param("id")

	if err := services.DeleteUser(c, id); err != nil {
		c.JSON(http.StatusInternalServerError, models.CustomResponse{
			Status:  "Failed",
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, models.CustomResponse{
		Status:  "Success",
		Message: "User deleted successfully",
		Data:    nil,
	})
}

func DeleteMeHandler(c *gin.Context) {
}
