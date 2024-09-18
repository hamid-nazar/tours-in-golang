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

type Tour struct {
	Id             string      `json:"id"`
	Name           string      `json:"name" validate:"required" min:"10" max:"50"`
	Slug           string      `json:"slug"`
	Duration       string      `json:"duration" validate:"required"`
	Difficulty     string      `json:"difficulty"`
	Price          float64     `json:"price" validate:"required"`
	MaxGroupSize   int         `json:"maxGroupSize" validate:"required"`
	RatingsAvg     float64     `json:"ratingAvg" validate:"required" default:"4.5" min:"1" max:"5"`
	RatingQuantity int         `json:"ratingQuantity" validate:"required" default:"0"`
	ImageCover     string      `json:"imageCover"`
	Images         []string    `json:"images"`
	CreatedAt      time.Time   `json:"createdAt" default:"time.Now()"`
	StartDates     []time.Time `json:"startDates"`
	SecretTour     bool        `json:"secretTour" default:"false"`
	Summary        string      `json:"summary" validate:"required"`
	Description    string      `json:"description"`
	StartLocation  string      `json:"startLocation"`
	Locations      []string    `json:"locations"`
	Guides         []User      `json:"guides"`
}

func NewTour() *Tour {
	return &Tour{
		Id: uuid.New().String(),
	}
}

type Review struct {
	Id     string `json:"id"`
	Review string `json:"review" validate:"required" min:"10" max:"50"`
	Rating int    `json:"rating" min:"1" max:"5"`
	Tour   Tour   `json:"tour"`
	User   User   `json:"user"`
}

func NewReview() *Review {
	return &Review{
		Id: uuid.New().String(),
	}
}

type Booking struct {
	Id        string    `json:"id"`
	Tour      Tour      `json:"tour"`
	User      User      `json:"user"`
	Price     float64   `json:"price" validate:"required"`
	CreatedAt time.Time `json:"createdAt" default:"time.Now()"`
	Paid      bool      `json:"paid" default:"true"`
}

func NewBooking() *Booking {
	return &Booking{
		Id: uuid.New().String(),
	}
}
