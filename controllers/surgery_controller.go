package controllers

import (
	"errors"
	"net/http"

	"CRUD-hospital-go/config"
	"CRUD-hospital-go/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func ScheduleSurgery(c *gin.Context) {
	var request models.SurgeryScheduleRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var surgery models.SurgerySchedule

	err := config.DB.Transaction(func(tx *gorm.DB) error {
		var ot models.OperatingTheater
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("status = ?", models.OTStatusAvailable).
			First(&ot).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("no available Operating Theater found")
			}
			return err
		}

		ot.Status = models.OTStatusOccupied
		if err := tx.Save(&ot).Error; err != nil {
			return errors.New("failed to update Operating Theater status")
		}

		var doctor models.Doctor
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ?", request.DoctorID).
			First(&doctor).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("doctor not found")
			}
			return err
		}

		if !doctor.IsAvailable {
			return errors.New("doctor is not available - currently in another surgery")
		}

		doctor.IsAvailable = false
		if err := tx.Save(&doctor).Error; err != nil {
			return errors.New("failed to update doctor availability")
		}

		var patient models.Patient
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ?", request.PatientID).
			First(&patient).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("patient not found")
			}
			return err
		}

		if patient.Deposit < request.DepositRequired {
			return errors.New("insufficient patient deposit for surgery")
		}

		patient.Deposit -= request.DepositRequired
		if err := tx.Save(&patient).Error; err != nil {
			return errors.New("failed to deduct patient deposit")
		}

		surgery = models.SurgerySchedule{
			PatientID:          request.PatientID,
			DoctorID:           request.DoctorID,
			OperatingTheaterID: ot.ID,
			SurgeryType:        request.SurgeryType,
			ScheduledAt:        request.ScheduledAt,
			EstimatedDuration:  request.EstimatedDuration,
			DepositDeducted:    request.DepositRequired,
			Status:             models.SurgeryStatusScheduled,
			Notes:              request.Notes,
		}

		if err := tx.Create(&surgery).Error; err != nil {
			return errors.New("failed to create surgery schedule")
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to schedule surgery",
			"details": err.Error(),
		})
		return
	}

	config.DB.Preload("Patient").Preload("Doctor").Preload("OperatingTheater").First(&surgery, surgery.ID)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Surgery scheduled successfully",
		"surgery": surgery,
	})
}

func CompleteSurgery(c *gin.Context) {
	surgeryID := c.Param("id")

	err := config.DB.Transaction(func(tx *gorm.DB) error {
		var surgery models.SurgerySchedule
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ?", surgeryID).
			First(&surgery).Error; err != nil {
			return errors.New("surgery not found")
		}

		if surgery.Status != models.SurgeryStatusScheduled && surgery.Status != models.SurgeryStatusInProgress {
			return errors.New("surgery is already completed or cancelled")
		}

		var ot models.OperatingTheater
		if err := tx.First(&ot, surgery.OperatingTheaterID).Error; err == nil {
			ot.Status = models.OTStatusAvailable
			tx.Save(&ot)
		}

		var doctor models.Doctor
		if err := tx.First(&doctor, surgery.DoctorID).Error; err == nil {
			doctor.IsAvailable = true
			tx.Save(&doctor)
		}

		surgery.Status = models.SurgeryStatusCompleted
		if err := tx.Save(&surgery).Error; err != nil {
			return errors.New("failed to update surgery status")
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Surgery completed successfully"})
}

func CancelSurgery(c *gin.Context) {
	surgeryID := c.Param("id")

	err := config.DB.Transaction(func(tx *gorm.DB) error {
		var surgery models.SurgerySchedule
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ?", surgeryID).
			First(&surgery).Error; err != nil {
			return errors.New("surgery not found")
		}

		if surgery.Status != models.SurgeryStatusScheduled {
			return errors.New("can only cancel scheduled surgeries")
		}

		var ot models.OperatingTheater
		if err := tx.First(&ot, surgery.OperatingTheaterID).Error; err == nil {
			ot.Status = models.OTStatusAvailable
			tx.Save(&ot)
		}

		var doctor models.Doctor
		if err := tx.First(&doctor, surgery.DoctorID).Error; err == nil {
			doctor.IsAvailable = true
			tx.Save(&doctor)
		}

		var patient models.Patient
		if err := tx.First(&patient, surgery.PatientID).Error; err == nil {
			patient.Deposit += surgery.DepositDeducted
			tx.Save(&patient)
		}

		surgery.Status = models.SurgeryStatusCancelled
		if err := tx.Save(&surgery).Error; err != nil {
			return errors.New("failed to update surgery status")
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Surgery cancelled and deposit refunded"})
}

func GetSurgeryByID(c *gin.Context) {
	var surgery models.SurgerySchedule

	if err := config.DB.Preload("Patient").Preload("Doctor").Preload("OperatingTheater").
		Where("id = ?", c.Param("id")).
		First(&surgery).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Surgery not found!"})
		return
	}

	c.JSON(http.StatusOK, surgery)
}

func GetAllSurgeries(c *gin.Context) {
	var surgeries []models.SurgerySchedule

	if err := config.DB.Preload("Patient").Preload("Doctor").Preload("OperatingTheater").
		Find(&surgeries).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, surgeries)
}

func GetSurgeriesByDoctor(c *gin.Context) {
	var surgeries []models.SurgerySchedule

	if err := config.DB.Preload("Patient").Preload("OperatingTheater").
		Where("doctor_id = ?", c.Param("doctor_id")).
		Find(&surgeries).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, surgeries)
}

func GetSurgeriesByPatient(c *gin.Context) {
	var surgeries []models.SurgerySchedule

	if err := config.DB.Preload("Doctor").Preload("OperatingTheater").
		Where("patient_id = ?", c.Param("patient_id")).
		Find(&surgeries).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, surgeries)
}
