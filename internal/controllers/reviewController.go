package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hamid-nazari/tours-in-go/internal/models"
	"github.com/hamid-nazari/tours-in-go/internal/services"
)

func CreateReviewHandler(c *gin.Context) {
	review := models.NewReview()

	if err := c.ShouldBindJSON(&review); err != nil {
		c.JSON(http.StatusBadRequest, models.CustomResponse{
			Status:  "Failed",
			Message: fmt.Errorf("Failed to bind JSON: %v", err).Error(),
			Data:    nil,
		})
		return
	}

	if err := services.ValidateReview(*review); err != nil {
		c.JSON(http.StatusBadRequest, models.CustomResponse{
			Status:  "Failed",
			Message: fmt.Errorf("Validation failed: %v", err).Error(),
			Data:    nil,
		})
		return
	}

	if err := services.CreateReview(c, review); err != nil {
		c.JSON(http.StatusInternalServerError, models.CustomResponse{
			Status:  "Failed",
			Message: fmt.Errorf("Failed to create review: %v", err).Error(),
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, models.CustomResponse{
		Status:  "Success",
		Message: "Review created successfully",
		Data:    review,
	})

}

func GetAllReviewsHandler(c *gin.Context) {
	tours := services.GetAllReviews(c)
	if tours == nil {
		c.JSON(http.StatusNotFound, models.CustomResponse{
			Status:  "Failed",
			Message: "No reviews not found",
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, models.CustomResponse{
		Status:  "Success",
		Message: "Reviews retrieved successfully",
		Data:    tours,
	})
}

func GetReviewHandler(c *gin.Context) {
	reviewId := c.Param("id")

	if reviewId == "" {
		c.JSON(http.StatusBadRequest, models.CustomResponse{
			Status:  "Failed",
			Message: "Review ID is required",
			Data:    nil,
		})
		return
	}

	review := services.GetReviewById(c, reviewId)
	if review == nil {
		c.JSON(http.StatusNotFound, models.CustomResponse{
			Status:  "Failed",
			Message: "Review not found",
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, models.CustomResponse{
		Status:  "Success",
		Message: "Review retrieved successfully",
		Data:    review,
	})

}

func UpdateReviewHandler(c *gin.Context) {
	reviewId := c.Param("id")

	if reviewId == "" {
		c.JSON(http.StatusBadRequest, models.CustomResponse{
			Status:  "Failed",
			Message: "Review ID is required",
			Data:    nil,
		})
		return
	}

	review := services.GetReviewById(c, reviewId)
	if review == nil {
		c.JSON(http.StatusNotFound, models.CustomResponse{
			Status:  "Failed",
			Message: "Review not found",
			Data:    nil,
		})
		return
	}

	if err := c.ShouldBindJSON(&review); err != nil {
		c.JSON(http.StatusBadRequest, models.CustomResponse{
			Status:  "Failed",
			Message: fmt.Errorf("Failed to bind JSON: %v", err).Error(),
			Data:    nil,
		})
		return
	}

	if err := services.ValidateReview(*review); err != nil {
		c.JSON(http.StatusBadRequest, models.CustomResponse{
			Status:  "Failed",
			Message: fmt.Errorf("Validation failed: %v", err).Error(),
			Data:    nil,
		})
		return
	}

	if err := services.UpdateReview(c, review); err != nil {
		c.JSON(http.StatusInternalServerError, models.CustomResponse{
			Status:  "Failed",
			Message: fmt.Errorf("Failed to update review: %v", err).Error(),
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, models.CustomResponse{
		Status:  "Success",
		Message: "Review updated successfully",
		Data:    review,
	})

}

func DeleteReviewHandler(c *gin.Context) {

	reviewId := c.Param("id")

	if reviewId == "" {
		c.JSON(http.StatusBadRequest, models.CustomResponse{
			Status:  "Failed",
			Message: "Review ID is required",
			Data:    nil,
		})
		return
	}

	review := services.GetReviewById(c, reviewId)
	if review == nil {
		c.JSON(http.StatusNotFound, models.CustomResponse{
			Status:  "Failed",
			Message: "Review not found",
			Data:    nil,
		})
		return
	}

	if err := services.DeleteReview(c, reviewId); err != nil {
		c.JSON(http.StatusInternalServerError, models.CustomResponse{
			Status:  "Failed",
			Message: fmt.Errorf("Failed to delete review: %v", err).Error(),
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, models.CustomResponse{
		Status:  "Success",
		Message: "Review deleted successfully",
		Data:    nil,
	})
}
