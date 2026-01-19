package controllers

import (
	"log"
	"net/http"
	"time"

	"CRUD-hospital-go/config"
	"CRUD-hospital-go/models"

	"github.com/gin-gonic/gin"
)

func CreateDoctor(c *gin.Context) {
	log.Println("CreateDoctor: Request received")

	var input models.Doctor

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("CreateDoctor: Invalid request body - %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config.DB.Create(&input)

	log.Printf("CreateDoctor: Doctor created successfully with ID %d", input.ID)
	c.JSON(http.StatusOK, input)
}

func GetAllDoctors(c *gin.Context) {
	log.Println("GetAllDoctors: Request received")

	var doctors []models.Doctor

	if err := config.DB.Find(&doctors).Error; err != nil {
		log.Printf("GetAllDoctors: Error fetching doctors - %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("GetAllDoctors: Found %d doctors", len(doctors))
	c.JSON(http.StatusOK, doctors)
}

func GetDoctorByID(c *gin.Context) {
	log.Printf("GetDoctorByID: Request received for ID %s", c.Param("id"))

	var doctor models.Doctor

	if err := config.DB.Where("id = ?", c.Param("id")).First(&doctor).Error; err != nil {
		log.Printf("GetDoctorByID: Doctor not found with ID %s", c.Param("id"))
		c.JSON(http.StatusNotFound, gin.H{"error": "Doctor not found!"})
		return
	}

	log.Printf("GetDoctorByID: Doctor found with ID %d", doctor.ID)
	c.JSON(http.StatusOK, doctor)
}

func DeleteDoctor(c *gin.Context) {
	log.Printf("DeleteDoctor: Request received for ID %s", c.Param("id"))

	var doctor models.Doctor
	id := c.Param("id")

	if err := config.DB.First(&doctor, "id = ?", id).Error; err != nil {
		log.Printf("DeleteDoctor: Doctor not found with ID %s", id)
		c.JSON(http.StatusNotFound, gin.H{"error": "Doctor not found!"})
		return
	}

	config.DB.Delete(&doctor)

	log.Printf("DeleteDoctor: Doctor deleted successfully with ID %s", id)
	c.JSON(http.StatusOK, gin.H{"message": "Doctor deleted successfully"})
}

func SearchDoctorByName(c *gin.Context) {
	name := c.Query("name")
	log.Printf("SearchDoctorByName: Request received for name %s", name)

	var doctors []models.Doctor

	if err := config.DB.Where("name LIKE ?", "%"+name+"%").Find(&doctors).Error; err != nil {
		log.Printf("SearchDoctorByName: Error searching doctors - %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("SearchDoctorByName: Found %d doctors matching name %s", len(doctors), name)
	c.JSON(http.StatusOK, doctors)
}

func UpdateDoctor(c *gin.Context) {
	log.Printf("UpdateDoctor: Request received for ID %s", c.Param("id"))

	var doctor models.Doctor
	id := c.Param("id")

	if err := config.DB.First(&doctor, "id = ?", id).Error; err != nil {
		log.Printf("UpdateDoctor: Doctor not found with ID %s", id)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Doctor not found!"})
		return
	}

	var input struct {
		Name      *string `json:"name"`
		ContactNo *string `json:"contact_no"`
		Address   *string `json:"address"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("UpdateDoctor: Invalid request body - %v", err)
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
	log.Printf("UpdateDoctor: Doctor updated successfully with ID %s", id)
	c.JSON(http.StatusOK, doctor)
}

func CheckDoctorAvailability(c *gin.Context) {
	doctorID := c.Param("id")
	dateStr := c.Query("date")

	log.Printf("CheckDoctorAvailability: Request for doctor %s on date %s", doctorID, dateStr)

	var doctor models.Doctor
	if err := config.DB.First(&doctor, "id = ?", doctorID).Error; err != nil {
		log.Printf("CheckDoctorAvailability: Doctor not found with ID %s", doctorID)
		c.JSON(http.StatusNotFound, gin.H{"error": "Doctor not found!"})
		return
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		log.Printf("CheckDoctorAvailability: Invalid date format %s", dateStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
		return
	}

	nextDay := date.Add(24 * time.Hour)

	var existingSurgery models.SurgerySchedule
	err = config.DB.Where("doctor_id = ? AND scheduled_at >= ? AND scheduled_at < ? AND status NOT IN ?",
		doctorID, date, nextDay,
		[]models.SurgeryStatus{models.SurgeryStatusCompleted, models.SurgeryStatusCancelled}).
		First(&existingSurgery).Error

	isAvailable := err != nil

	log.Printf("CheckDoctorAvailability: Doctor %s availability on %s: %v", doctorID, dateStr, isAvailable)
	c.JSON(http.StatusOK, gin.H{
		"doctor_id":    doctorID,
		"doctor_name":  doctor.Name,
		"date":         dateStr,
		"is_available": isAvailable,
	})
}
