package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hamid-nazari/tours-in-go/internal/models"
	"github.com/hamid-nazari/tours-in-go/internal/services"
	"github.com/stripe/stripe-go/v79"
)

func GetCheckoutSessionHandler(c *gin.Context) {
	tourId := c.Param("id")
	if tourId == "" {
		c.JSON(http.StatusBadRequest, models.CustomResponse{
			Status:  "Failed",
			Message: "Tour ID is required",
			Data:    nil,
		})
		return
	}

	tour := services.FindTourById(c, tourId)
	if tour == nil {
		c.JSON(http.StatusNotFound, models.CustomResponse{
			Status:  "Failed",
			Message: "Tour not found",
			Data:    nil,
		})
		return
	}

	user, ok := c.Get("user")
	if !ok {
		c.JSON(http.StatusInternalServerError, models.CustomResponse{
			Status:  "Failed",
			Message: "You need to be logged in to proceed",
			Data:    nil,
		})
		return
	}

	session := stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		Mode:               stripe.String("payment"),
		SuccessURL:         stripe.String("https://example.com/success"),
		CancelURL:          stripe.String("https://example.com/canceled"),
		ClientReferenceID:  stripe.String(tourId),
		CustomerEmail:      stripe.String(user.(*models.User).Email),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price: stripe.String("price_1L3Zb2HlB6kSf1sR4t6H4Hqo"),
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency:   stripe.String("usd"),
					UnitAmount: stripe.Int64(int64(tour.Price)),
					Product:    stripe.String(tourId),
				},
				Quantity: stripe.Int64(1),
			},
		},
	}

	// sessionParams := &stripe.CheckoutSessionParams{
	// 	Params: stripe.Params{
	// 		Metadata: map[string]string{
	// 			"tour_id": tourId,
	// 		},
	// 	},
	// }

	// sessionParams.AddExpand("line_items")
	// sessionParams.AddExpand("customer")

	c.JSON(http.StatusOK, models.CustomResponse{
		Status:  "Success",
		Message: "Checkout session created successfully",
		Data:    session,
	})

}

func WebhookHandler(c *gin.Context) {

}

func CreateBookingHandler(c *gin.Context) {

}

func GetAllBookingsHandler(c *gin.Context) {

}

func GetBookingHandler(c *gin.Context) {

}

func UpdateBookingHandler(c *gin.Context) {

}

func DeleteBookingHandler(c *gin.Context) {

}
