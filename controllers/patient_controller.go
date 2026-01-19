package controllers

import (
	"log"
	"net/http"
	"time"

	"CRUD-hospital-go/config"
	"CRUD-hospital-go/models"

	"github.com/gin-gonic/gin"
)

func CreatePatient(c *gin.Context) {
	log.Println("CreatePatient: Request received")

	var input models.Patient

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("CreatePatient: Invalid request body - %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config.DB.Create(&input)

	log.Printf("CreatePatient: Patient created successfully with ID %d", input.ID)
	c.JSON(http.StatusOK, input)
}

func GetAllPatients(c *gin.Context) {
	log.Println("GetAllPatients: Request received")

	var patients []models.Patient

	if err := config.DB.Find(&patients).Error; err != nil {
		log.Printf("GetAllPatients: Error fetching patients - %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("GetAllPatients: Found %d patients", len(patients))
	c.JSON(http.StatusOK, patients)
}

func GetPatientByID(c *gin.Context) {
	log.Printf("GetPatientByID: Request received for ID %s", c.Param("id"))

	var patient models.Patient

	if err := config.DB.Where("id = ?", c.Param("id")).First(&patient).Error; err != nil {
		log.Printf("GetPatientByID: Patient not found with ID %s", c.Param("id"))
		c.JSON(http.StatusNotFound, gin.H{"error": "Patient not found!"})
		return
	}

	log.Printf("GetPatientByID: Patient found with ID %d", patient.ID)
	c.JSON(http.StatusOK, patient)
}

func UpdatePatient(c *gin.Context) {
	log.Printf("UpdatePatient: Request received for ID %s", c.Param("id"))

	var patient models.Patient
	id := c.Param("id")

	if err := config.DB.First(&patient, "id = ?", id).Error; err != nil {
		log.Printf("UpdatePatient: Patient not found with ID %s", id)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Patient not found!"})
		return
	}

	var input struct {
		Name      *string  `json:"name"`
		ContactNo *string  `json:"contact_no"`
		Address   *string  `json:"address"`
		DoctorID  *uint    `json:"doctor_id"`
		Deposit   *float64 `json:"deposit"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("UpdatePatient: Invalid request body - %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Name != nil {
		patient.Name = *input.Name
	}
	if input.ContactNo != nil {
		patient.ContactNo = *input.ContactNo
	}
	if input.Address != nil {
		patient.Address = *input.Address
	}
	if input.DoctorID != nil {
		patient.DoctorID = *input.DoctorID
	}
	if input.Deposit != nil {
		patient.Deposit = *input.Deposit
	}
	patient.UpdatedAt = time.Now()

	config.DB.Save(&patient)
	log.Printf("UpdatePatient: Patient updated successfully with ID %s", id)
	c.JSON(http.StatusOK, patient)
}

func DeletePatient(c *gin.Context) {
	log.Printf("DeletePatient: Request received for ID %s", c.Param("id"))

	var patient models.Patient
	id := c.Param("id")

	if err := config.DB.First(&patient, "id = ?", id).Error; err != nil {
		log.Printf("DeletePatient: Patient not found with ID %s", id)
		c.JSON(http.StatusNotFound, gin.H{"error": "Patient not found!"})
		return
	}

	config.DB.Delete(&patient)

	log.Printf("DeletePatient: Patient deleted successfully with ID %s", id)
	c.JSON(http.StatusOK, gin.H{"message": "Patient deleted successfully"})
}

func GetPatientsByDoctorID(c *gin.Context) {
	log.Printf("GetPatientsByDoctorID: Request received for doctor_id %s", c.Param("doctor_id"))

	var patients []models.Patient

	if err := config.DB.Where("doctor_id = ?", c.Param("doctor_id")).Find(&patients).Error; err != nil {
		log.Printf("GetPatientsByDoctorID: Error fetching patients - %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Patients not found!"})
		return
	}

	log.Printf("GetPatientsByDoctorID: Found %d patients for doctor_id %s", len(patients), c.Param("doctor_id"))
	c.JSON(http.StatusOK, patients)
}

func SearchPatientByName(c *gin.Context) {
	name := c.Query("name")
	log.Printf("SearchPatientByName: Request received for name %s", name)

	var patients []models.Patient

	if err := config.DB.Where("name LIKE ?", "%"+name+"%").Find(&patients).Error; err != nil {
		log.Printf("SearchPatientByName: Error searching patients - %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("SearchPatientByName: Found %d patients matching name %s", len(patients), name)
	c.JSON(http.StatusOK, patients)
}
