package controllers

import (
	"log"
	"net/http"
	"time"

	"CRUD-hospital-go/config"
	"CRUD-hospital-go/models"

	"github.com/gin-gonic/gin"
)

func CreateOperatingTheater(c *gin.Context) {
	log.Println("CreateOperatingTheater: Request received")

	var input models.OperatingTheater

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("CreateOperatingTheater: Invalid request body - %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Status == "" {
		input.Status = models.OTStatusAvailable
	}

	config.DB.Create(&input)
	log.Printf("CreateOperatingTheater: Operating Theater created successfully with ID %d", input.ID)
	c.JSON(http.StatusCreated, input)
}

func GetOperatingTheaterByID(c *gin.Context) {
	log.Printf("GetOperatingTheaterByID: Request received for ID %s", c.Param("id"))

	var ot models.OperatingTheater

	if err := config.DB.Where("id = ?", c.Param("id")).First(&ot).Error; err != nil {
		log.Printf("GetOperatingTheaterByID: Operating Theater not found with ID %s", c.Param("id"))
		c.JSON(http.StatusNotFound, gin.H{"error": "Operating Theater not found!"})
		return
	}

	log.Printf("GetOperatingTheaterByID: Operating Theater found with ID %d", ot.ID)
	c.JSON(http.StatusOK, ot)
}

func GetAllOperatingTheaters(c *gin.Context) {
	log.Println("GetAllOperatingTheaters: Request received")

	var ots []models.OperatingTheater

	if err := config.DB.Find(&ots).Error; err != nil {
		log.Printf("GetAllOperatingTheaters: Error fetching Operating Theaters - %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("GetAllOperatingTheaters: Found %d Operating Theaters", len(ots))
	c.JSON(http.StatusOK, ots)
}

func GetAvailableOperatingTheaters(c *gin.Context) {
	log.Println("GetAvailableOperatingTheaters: Request received")

	var ots []models.OperatingTheater

	if err := config.DB.Where("status = ?", models.OTStatusAvailable).Find(&ots).Error; err != nil {
		log.Printf("GetAvailableOperatingTheaters: Error fetching available Operating Theaters - %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("GetAvailableOperatingTheaters: Found %d available Operating Theaters", len(ots))
	c.JSON(http.StatusOK, ots)
}

func UpdateOperatingTheater(c *gin.Context) {
	log.Printf("UpdateOperatingTheater: Request received for ID %s", c.Param("id"))

	var ot models.OperatingTheater
	id := c.Param("id")

	if err := config.DB.First(&ot, "id = ?", id).Error; err != nil {
		log.Printf("UpdateOperatingTheater: Operating Theater not found with ID %s", id)
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
		log.Printf("UpdateOperatingTheater: Invalid request body - %v", err)
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
	log.Printf("UpdateOperatingTheater: Operating Theater updated successfully with ID %s", id)
	c.JSON(http.StatusOK, ot)
}

func DeleteOperatingTheater(c *gin.Context) {
	log.Printf("DeleteOperatingTheater: Request received for ID %s", c.Param("id"))

	var ot models.OperatingTheater
	id := c.Param("id")

	if err := config.DB.First(&ot, "id = ?", id).Error; err != nil {
		log.Printf("DeleteOperatingTheater: Operating Theater not found with ID %s", id)
		c.JSON(http.StatusNotFound, gin.H{"error": "Operating Theater not found!"})
		return
	}

	config.DB.Delete(&ot)
	log.Printf("DeleteOperatingTheater: Operating Theater deleted successfully with ID %s", id)
	c.JSON(http.StatusOK, gin.H{"message": "Operating Theater deleted successfully"})
}
