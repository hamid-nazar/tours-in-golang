package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hamid-nazari/tours-in-go/internal/models"
	"github.com/hamid-nazari/tours-in-go/internal/services"
)

func CreateTourHandler(c *gin.Context) {
	tour := models.NewTour()

	if err := c.ShouldBindJSON(&tour); err != nil {
		c.JSON(http.StatusBadRequest, models.CustomResponse{
			Status:  "Failed",
			Message: fmt.Errorf("Failed to bind JSON: %v", err).Error(),
			Data:    nil,
		})
		return
	}

	if err := services.ValidateTour(*tour); err != nil {
		c.JSON(http.StatusBadRequest, models.CustomResponse{
			Status:  "Failed",
			Message: fmt.Errorf("Validation failed: %v", err).Error(),
			Data:    nil,
		})
		return
	}

	if err := services.CreateTour(c, tour); err != nil {
		c.JSON(http.StatusInternalServerError, models.CustomResponse{
			Status:  "Failed",
			Message: fmt.Errorf("Failed to create tour: %v", err).Error(),
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, models.CustomResponse{
		Status:  "Success",
		Message: "Tour created successfully",
		Data:    tour,
	})

}

func GetAllToursHandler(c *gin.Context) {
	tours := services.GetAllTours(c)
	if tours == nil {
		c.JSON(http.StatusNotFound, models.CustomResponse{
			Status:  "Failed",
			Message: "No tours not found",
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, models.CustomResponse{
		Status:  "Success",
		Message: "Tours retrieved successfully",
		Data:    tours,
	})

}

func GetTourHandler(c *gin.Context) {
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

	c.JSON(http.StatusOK, models.CustomResponse{
		Status:  "Success",
		Message: "Tour retrieved successfully",
		Data:    tour,
	})

}

func UpdateTourHandler(c *gin.Context) {
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

	if err := c.ShouldBindJSON(&tour); err != nil {
		c.JSON(http.StatusBadRequest, models.CustomResponse{
			Status:  "Failed",
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	if err := services.ValidateTour(*tour); err != nil {
		c.JSON(http.StatusBadRequest, models.CustomResponse{
			Status:  "Failed",
			Message: fmt.Sprintf("Validation failed: %v", err),
			Data:    nil,
		})
		return
	}

	if err := services.UpdateTour(c, tour); err != nil {
		c.JSON(http.StatusInternalServerError, models.CustomResponse{
			Status:  "Failed",
			Message: "Failed to update tour",
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, models.CustomResponse{
		Status:  "Success",
		Message: "Tour updated successfully",
		Data:    tour,
	})

}

func DeleteTourHandler(c *gin.Context) {

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

	if err := services.DeleteTour(c, tour); err != nil {
		c.JSON(http.StatusInternalServerError, models.CustomResponse{
			Status:  "Failed",
			Message: fmt.Sprintf("Failed to delete tour: %v", err),
			Data:    nil,
		})
		return
	}

}
