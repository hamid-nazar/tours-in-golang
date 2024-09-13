package controllers

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/hamid-nazari/tours-in-go/internal/models"
	"github.com/hamid-nazari/tours-in-go/internal/services"
)

func CreateJwtTokenAndSend(c *gin.Context, user *models.User, message string) {

	claims := models.CustomClaims{
		UserId: user.Id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.CustomResponse{
			Status:  "Failed",
			Message: "Failed to create JWT token",
			Data:    nil,
		})
		return

	}

	cookieOptions := &http.Cookie{
		Name:    "jwt",
		Value:   token,
		Expires: time.Now().Add(time.Hour * 24),
	}

	http.SetCookie(c.Writer, cookieOptions)

	c.JSON(http.StatusOK, models.CustomResponse{
		Status:  "Success",
		Message: message,
		Data: gin.H{
			"token": token,
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

	CreateJwtTokenAndSend(c, user, "User created successfully")
}
func LoginHandler(c *gin.Context) {
	var jsonData map[string]string

	if err := c.ShouldBindJSON(&jsonData); err != nil {
		c.JSON(http.StatusBadRequest, models.CustomResponse{
			Status:  "Failed",
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	email, password := jsonData["email"], jsonData["password"]

	if email == "" || password == "" {
		c.JSON(http.StatusBadRequest, models.CustomResponse{
			Status:  "Failed",
			Message: "Email and password are required",
			Data:    nil,
		})
		return
	}

	user := services.FindUserByEmail(c, email)
	if user == nil {
		c.JSON(http.StatusNotFound, models.CustomResponse{
			Status:  "Failed",
			Message: "User not found",
			Data:    nil,
		})
		return
	}
	if !services.VerifyPassword(password, user.Password) {
		c.JSON(http.StatusUnauthorized, models.CustomResponse{
			Status:  "Failed",
			Message: "Password is incorrect",
			Data:    nil,
		})
		return
	}

	CreateJwtTokenAndSend(c, user, "User logged in successfully")

}

func LogoutHandler(c *gin.Context) {
	cookie := &http.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	}

	http.SetCookie(c.Writer, cookie)

	c.JSON(http.StatusOK, models.CustomResponse{
		Status:  "Success",
		Message: "User logged out successfully",
		Data:    nil,
	})
}

func ProtectHandler(c *gin.Context) {
	token := c.GetHeader("Authorization")

	if token == "" {
		c.JSON(http.StatusUnauthorized, models.CustomResponse{
			Status:  "Failed",
			Message: "Unauthorized",
			Data:    nil,
		})
		return
	}

	claims, err := extractAndvalidateToken(c)

	if err != nil {
		c.JSON(http.StatusUnauthorized, models.CustomResponse{
			Status:  "Failed",
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	currentUser := services.FindUserById(c, claims.UserId)
	if currentUser == nil {
		c.JSON(http.StatusUnauthorized, models.CustomResponse{
			Status:  "Failed",
			Message: "User assigned to token not found",
			Data:    nil,
		})
		return
	}

	if currentUser.PasswordChangedAt.UTC().After(claims.IssuedAt.Time) {
		c.JSON(http.StatusUnauthorized, models.CustomResponse{
			Status:  "Failed",
			Message: "User recently changed password. Please login again",
			Data:    nil,
		})
		return
	}

	c.Set("user", currentUser)

	c.Next()
}

func RestrictTo(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {

		currentUser, ok := c.Get("user")

		if !ok {
			c.JSON(http.StatusUnauthorized, models.CustomResponse{
				Status:  "Failed",
				Message: "Unauthorized",
				Data:    nil,
			})
			return
		}

		for _, role := range roles {
			if role == currentUser.(*models.User).Role {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusUnauthorized, models.CustomResponse{
			Status:  "Failed",
			Message: "You are not authorized to access this resource",
			Data:    nil,
		})
	}
}

func ForgotPasswordHandler(c *gin.Context) {
	var jsonData map[string]string

	if err := c.ShouldBindJSON(&jsonData); err != nil {
		c.JSON(http.StatusBadRequest, models.CustomResponse{
			Status:  "Failed",
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	email := jsonData["email"]

	user := services.FindUserByEmail(c, email)
	if user == nil {
		c.JSON(http.StatusNotFound, models.CustomResponse{
			Status:  "Failed",
			Message: "User associated with email not found",
			Data:    nil,
		})
		return
	}

	user.PasswordResetToken = generatePasswordResetToken()
	user.PasswordResetTokenExpiry = time.Now().Add(10 * time.Minute)

	services.UpdateUser(c, user)

	resetURL := fmt.Sprintf("%s://%s/api/users/reset-password/%s", c.Request.Proto, c.Request.Host, user.PasswordResetToken)

	c.JSON(http.StatusOK, models.CustomResponse{
		Status:  "Success",
		Message: "Password reset token sent to email",
		Data:    map[string]string{"reset_url": resetURL},
	})

}

func ResetPasswordHandler(c *gin.Context) {

	hashedResetToken := c.Param("token")

	user := services.FindUserByPasswordResetToken(c, hashedResetToken)

	if user == nil {
		c.JSON(http.StatusNotFound, models.CustomResponse{
			Status:  "Failed",
			Message: "User associated with token not found",
			Data:    nil,
		})
		return
	}

	if time.Now().After(user.PasswordResetTokenExpiry) {
		c.JSON(http.StatusNotFound, models.CustomResponse{
			Status:  "Failed",
			Message: "Password reset token expired",
			Data:    nil,
		})
		return
	}

	var jsonData map[string]string

	if err := c.ShouldBindJSON(&jsonData); err != nil {
		c.JSON(http.StatusBadRequest, models.CustomResponse{
			Status:  "Failed",
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	password := jsonData["password"]
	if password == "" {
		c.JSON(http.StatusBadRequest, models.CustomResponse{
			Status:  "Failed",
			Message: "Password is required",
			Data:    nil,
		})
		return
	}

	user.Password = password
	user.PasswordChangedAt = time.Now()
	user.PasswordResetToken = ""
	user.PasswordResetTokenExpiry = time.Time{}

	services.UpdateUser(c, user)

	CreateJwtTokenAndSend(c, user, "Password reset successful")
}
func UpdatePasswordHandler(c *gin.Context) {
	var jsonData map[string]string

	if err := c.ShouldBindJSON(&jsonData); err != nil {
		c.JSON(http.StatusBadRequest, models.CustomResponse{
			Status:  "Failed",
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	currentPassword, newPassword, newPasswordConfirm := jsonData["currentPassword"], jsonData["newPassword"], jsonData["newPasswordConfirm"]

	if currentPassword == "" || newPassword == "" || newPasswordConfirm == "" {
		c.JSON(http.StatusBadRequest, models.CustomResponse{
			Status:  "Failed",
			Message: "All fields are required",
			Data:    nil,
		})
		return
	}

	currentUser, ok := c.Get("user")

	if !ok {
		c.JSON(http.StatusUnauthorized, models.CustomResponse{
			Status:  "Failed",
			Message: "Unauthorized",
			Data:    nil,
		})
		return
	}

	isPasswordCorrect := services.VerifyPassword(currentUser.(*models.User).Password, currentPassword)

	if !isPasswordCorrect {
		c.JSON(http.StatusUnauthorized, models.CustomResponse{
			Status:  "Failed",
			Message: "Incorrect password",
			Data:    nil,
		})
		return
	}

	if newPassword != newPasswordConfirm {
		c.JSON(http.StatusBadRequest, models.CustomResponse{
			Status:  "Failed",
			Message: "New password and confirm password do not match",
			Data:    nil,
		})
		return
	}

	hashedPassword, err := services.HashPassword(newPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.CustomResponse{
			Status:  "Failed",
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	currentUser.(*models.User).Password = hashedPassword
	currentUser.(*models.User).PasswordChangedAt = time.Now()

	services.UpdateUser(c, currentUser.(*models.User))

	CreateJwtTokenAndSend(c, currentUser.(*models.User), "Password updated successfully")
}

func extractAndvalidateToken(c *gin.Context) (*models.CustomClaims, error) {
	authHeader := c.GetHeader("Authorization")

	if len(strings.Split(authHeader, " ")) != 2 || strings.Split(authHeader, " ")[0] != "Bearer" {
		return nil, errors.New("invalid token")
	}

	bearerToken := strings.Split(authHeader, " ")[1]

	token, err := jwt.ParseWithClaims(bearerToken, &models.CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*models.CustomClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func generatePasswordResetToken() string {
	resetToken := make([]byte, 32)

	rand.Read(resetToken)

	hashedToken := sha256.Sum256(resetToken)
	hashedResetToken := hex.EncodeToString(hashedToken[:])

	return hashedResetToken
}
