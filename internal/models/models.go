package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type CustomResponse struct {
	Status  string      `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type CustomClaims struct {
	UserId string `json:"userId"`
	jwt.RegisteredClaims
}

type User struct {
	Id                       string    `json:"id" validate:"required"`
	Name                     string    `json:"name" validate:"required"`
	Email                    string    `json:"email" validate:"required,email"`
	Photo                    string    `json:"photo"`
	Role                     string    `json:"role" validate:"oneof=user admin guide lead-guide"`
	Password                 string    `json:"password" validate:"required,min=8"`
	PasswordConfirm          string    `json:"passwordConfirm" validate:"required,eqfield=Password"`
	PasswordChangedAt        time.Time `json:"passwordChangedAt,omitempty"`
	PasswordResetToken       string    `json:"passwordResetToken,omitempty"`
	PasswordResetTokenExpiry time.Time `json:"passwordResetTokenExpiry,omitempty"`
	Active                   bool      `json:"active"`
}

func NewUser() *User {
	return &User{
		Id:                uuid.New().String(),
		Photo:             "https://i.pravatar.cc/300",
		Role:              "user",
		PasswordChangedAt: time.Now(),
		Active:            true,
	}
}
