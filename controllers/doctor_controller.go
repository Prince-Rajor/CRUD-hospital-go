package controllers

import (
	"net/http"
	"time"

	"CRUD-hospital-go/config"
	"CRUD-hospital-go/models"

	"github.com/gin-gonic/gin"
)

func CreateDoctor(c *gin.Context) {
	var input models.Doctor

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config.DB.Create(&input)

	c.JSON(http.StatusOK, input)
}

func GetDoctorByID(c *gin.Context) {
	var doctor models.Doctor

	if err := config.DB.Where("id = ?", c.Param("id")).First(&doctor).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Doctor not found!"})
		return
	}

	c.JSON(http.StatusOK, doctor)
}

func DeleteDoctor(c *gin.Context) {
	var doctor models.Doctor
	id := c.Param("id")

	if err := config.DB.First(&doctor, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Doctor not found!"})
		return
	}

	config.DB.Delete(&doctor)

	c.JSON(http.StatusOK, gin.H{"message": "Doctor deleted successfully"})
}

func SearchDoctorByName(c *gin.Context) {
	name := c.Query("name")
	var doctors []models.Doctor

	if err := config.DB.Where("name LIKE ?", "%"+name+"%").Find(&doctors).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, doctors)
}

func UpdateDoctor(c *gin.Context) {
	var doctor models.Doctor
	id := c.Param("id")

	if err := config.DB.First(&doctor, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Doctor not found!"})
		return
	}

	var input struct {
		Name      *string `json:"name"`
		ContactNo *string `json:"contact_no"`
		Address   *string `json:"address"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Name != nil {
		doctor.Name = *input.Name
	}
	if input.ContactNo != nil {
		doctor.ContactNo = *input.ContactNo
	}
	if input.Address != nil {
		doctor.Address = *input.Address
	}
	doctor.UpdatedAt = time.Now()

	config.DB.Save(&doctor)
	c.JSON(http.StatusOK, doctor)
}
