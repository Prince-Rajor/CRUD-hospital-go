package controllers

import (
	"net/http"
	"time"

	"CRUD-hospital-go/config"
	"CRUD-hospital-go/models"

	"github.com/gin-gonic/gin"
)

func CreateOperatingTheater(c *gin.Context) {
	var input models.OperatingTheater

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Status == "" {
		input.Status = models.OTStatusAvailable
	}

	config.DB.Create(&input)
	c.JSON(http.StatusCreated, input)
}

func GetOperatingTheaterByID(c *gin.Context) {
	var ot models.OperatingTheater

	if err := config.DB.Where("id = ?", c.Param("id")).First(&ot).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Operating Theater not found!"})
		return
	}

	c.JSON(http.StatusOK, ot)
}

func GetAllOperatingTheaters(c *gin.Context) {
	var ots []models.OperatingTheater

	if err := config.DB.Find(&ots).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ots)
}

func GetAvailableOperatingTheaters(c *gin.Context) {
	var ots []models.OperatingTheater

	if err := config.DB.Where("status = ?", models.OTStatusAvailable).Find(&ots).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ots)
}

func UpdateOperatingTheater(c *gin.Context) {
	var ot models.OperatingTheater
	id := c.Param("id")

	if err := config.DB.First(&ot, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Operating Theater not found!"})
		return
	}

	var input struct {
		Name     *string          `json:"name"`
		Floor    *int             `json:"floor"`
		Status   *models.OTStatus `json:"status"`
		Capacity *int             `json:"capacity"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Name != nil {
		ot.Name = *input.Name
	}
	if input.Floor != nil {
		ot.Floor = *input.Floor
	}
	if input.Status != nil {
		ot.Status = *input.Status
	}
	if input.Capacity != nil {
		ot.Capacity = *input.Capacity
	}
	ot.UpdatedAt = time.Now()

	config.DB.Save(&ot)
	c.JSON(http.StatusOK, ot)
}

func DeleteOperatingTheater(c *gin.Context) {
	var ot models.OperatingTheater
	id := c.Param("id")

	if err := config.DB.First(&ot, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Operating Theater not found!"})
		return
	}

	config.DB.Delete(&ot)
	c.JSON(http.StatusOK, gin.H{"message": "Operating Theater deleted successfully"})
}
