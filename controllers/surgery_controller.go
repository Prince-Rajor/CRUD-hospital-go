package controllers

import (
	"errors"
	"log"
	"net/http"
	"time"

	"CRUD-hospital-go/config"
	"CRUD-hospital-go/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func ScheduleSurgery(c *gin.Context) {
	log.Println("ScheduleSurgery: Request received")

	var request models.SurgeryScheduleRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		log.Printf("ScheduleSurgery: Invalid request body - %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("ScheduleSurgery: Scheduling surgery for patient_id=%d, doctor_id=%d", request.PatientID, request.DoctorID)

	var surgery models.SurgerySchedule

	err := config.DB.Transaction(func(tx *gorm.DB) error {
		var ot models.OperatingTheater
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("status = ?", models.OTStatusAvailable).
			First(&ot).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				log.Println("ScheduleSurgery: No available Operating Theater found")
				return errors.New("no available Operating Theater found")
			}
			return err
		}
		log.Printf("ScheduleSurgery: Found available OT with ID %d", ot.ID)

		ot.Status = models.OTStatusOccupied
		if err := tx.Save(&ot).Error; err != nil {
			log.Printf("ScheduleSurgery: Failed to update OT status - %v", err)
			return errors.New("failed to update Operating Theater status")
		}

		var doctor models.Doctor
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ?", request.DoctorID).
			First(&doctor).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				log.Printf("ScheduleSurgery: Doctor not found with ID %d", request.DoctorID)
				return errors.New("doctor not found")
			}
			return err
		}

		surgeryDate := request.ScheduledAt.Truncate(24 * time.Hour)
		nextDay := surgeryDate.Add(24 * time.Hour)

		var existingSurgery models.SurgerySchedule
		err := tx.Where("doctor_id = ? AND scheduled_at >= ? AND scheduled_at < ? AND status NOT IN ?",
			request.DoctorID, surgeryDate, nextDay,
			[]models.SurgeryStatus{models.SurgeryStatusCompleted, models.SurgeryStatusCancelled}).
			First(&existingSurgery).Error

		if err == nil {
			log.Printf("ScheduleSurgery: Doctor %d already has surgery on %s", request.DoctorID, surgeryDate.Format("2006-01-02"))
			return errors.New("doctor already has a surgery scheduled on this date")
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		var patient models.Patient
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ?", request.PatientID).
			First(&patient).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				log.Printf("ScheduleSurgery: Patient not found with ID %d", request.PatientID)
				return errors.New("patient not found")
			}
			return err
		}

		if patient.Deposit < request.DepositRequired {
			log.Printf("ScheduleSurgery: Insufficient deposit for patient %d (has %.2f, needs %.2f)", request.PatientID, patient.Deposit, request.DepositRequired)
			return errors.New("insufficient patient deposit for surgery")
		}

		patient.Deposit -= request.DepositRequired
		if err := tx.Save(&patient).Error; err != nil {
			log.Printf("ScheduleSurgery: Failed to deduct patient deposit - %v", err)
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
			log.Printf("ScheduleSurgery: Failed to create surgery schedule - %v", err)
			return errors.New("failed to create surgery schedule")
		}

		return nil
	})

	if err != nil {
		log.Printf("ScheduleSurgery: Transaction failed - %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to schedule surgery",
			"details": err.Error(),
		})
		return
	}

	config.DB.Preload("Patient").Preload("Doctor").Preload("OperatingTheater").First(&surgery, surgery.ID)

	log.Printf("ScheduleSurgery: Surgery scheduled successfully with ID %d", surgery.ID)
	c.JSON(http.StatusCreated, gin.H{
		"message": "Surgery scheduled successfully",
		"surgery": surgery,
	})
}

func CompleteSurgery(c *gin.Context) {
	surgeryID := c.Param("id")
	log.Printf("CompleteSurgery: Request received for surgery ID %s", surgeryID)

	err := config.DB.Transaction(func(tx *gorm.DB) error {
		var surgery models.SurgerySchedule
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ?", surgeryID).
			First(&surgery).Error; err != nil {
			log.Printf("CompleteSurgery: Surgery not found with ID %s", surgeryID)
			return errors.New("surgery not found")
		}

		if surgery.Status != models.SurgeryStatusScheduled && surgery.Status != models.SurgeryStatusInProgress {
			log.Printf("CompleteSurgery: Surgery %s is already completed or cancelled", surgeryID)
			return errors.New("surgery is already completed or cancelled")
		}

		var ot models.OperatingTheater
		if err := tx.First(&ot, surgery.OperatingTheaterID).Error; err == nil {
			ot.Status = models.OTStatusAvailable
			tx.Save(&ot)
			log.Printf("CompleteSurgery: OT %d marked as available", ot.ID)
		}

		surgery.Status = models.SurgeryStatusCompleted
		if err := tx.Save(&surgery).Error; err != nil {
			log.Printf("CompleteSurgery: Failed to update surgery status - %v", err)
			return errors.New("failed to update surgery status")
		}

		return nil
	})

	if err != nil {
		log.Printf("CompleteSurgery: Transaction failed - %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("CompleteSurgery: Surgery %s completed successfully", surgeryID)
	c.JSON(http.StatusOK, gin.H{"message": "Surgery completed successfully"})
}

func CancelSurgery(c *gin.Context) {
	surgeryID := c.Param("id")
	log.Printf("CancelSurgery: Request received for surgery ID %s", surgeryID)

	err := config.DB.Transaction(func(tx *gorm.DB) error {
		var surgery models.SurgerySchedule
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ?", surgeryID).
			First(&surgery).Error; err != nil {
			log.Printf("CancelSurgery: Surgery not found with ID %s", surgeryID)
			return errors.New("surgery not found")
		}

		if surgery.Status != models.SurgeryStatusScheduled {
			log.Printf("CancelSurgery: Cannot cancel surgery %s - status is %s", surgeryID, surgery.Status)
			return errors.New("can only cancel scheduled surgeries")
		}

		var ot models.OperatingTheater
		if err := tx.First(&ot, surgery.OperatingTheaterID).Error; err == nil {
			ot.Status = models.OTStatusAvailable
			tx.Save(&ot)
			log.Printf("CancelSurgery: OT %d marked as available", ot.ID)
		}

		var patient models.Patient
		if err := tx.First(&patient, surgery.PatientID).Error; err == nil {
			patient.Deposit += surgery.DepositDeducted
			tx.Save(&patient)
			log.Printf("CancelSurgery: Refunded %.2f to patient %d", surgery.DepositDeducted, patient.ID)
		}

		surgery.Status = models.SurgeryStatusCancelled
		if err := tx.Save(&surgery).Error; err != nil {
			log.Printf("CancelSurgery: Failed to update surgery status - %v", err)
			return errors.New("failed to update surgery status")
		}

		return nil
	})

	if err != nil {
		log.Printf("CancelSurgery: Transaction failed - %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("CancelSurgery: Surgery %s cancelled and deposit refunded", surgeryID)
	c.JSON(http.StatusOK, gin.H{"message": "Surgery cancelled and deposit refunded"})
}

func GetSurgeryByID(c *gin.Context) {
	log.Printf("GetSurgeryByID: Request received for ID %s", c.Param("id"))

	var surgery models.SurgerySchedule

	if err := config.DB.Preload("Patient").Preload("Doctor").Preload("OperatingTheater").
		Where("id = ?", c.Param("id")).
		First(&surgery).Error; err != nil {
		log.Printf("GetSurgeryByID: Surgery not found with ID %s", c.Param("id"))
		c.JSON(http.StatusNotFound, gin.H{"error": "Surgery not found!"})
		return
	}

	log.Printf("GetSurgeryByID: Surgery found with ID %d", surgery.ID)
	c.JSON(http.StatusOK, surgery)
}

func GetAllSurgeries(c *gin.Context) {
	log.Println("GetAllSurgeries: Request received")

	var surgeries []models.SurgerySchedule

	if err := config.DB.Preload("Patient").Preload("Doctor").Preload("OperatingTheater").
		Find(&surgeries).Error; err != nil {
		log.Printf("GetAllSurgeries: Error fetching surgeries - %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("GetAllSurgeries: Found %d surgeries", len(surgeries))
	c.JSON(http.StatusOK, surgeries)
}

func GetSurgeriesByDoctor(c *gin.Context) {
	log.Printf("GetSurgeriesByDoctor: Request received for doctor_id %s", c.Param("doctor_id"))

	var surgeries []models.SurgerySchedule

	if err := config.DB.Preload("Patient").Preload("OperatingTheater").
		Where("doctor_id = ?", c.Param("doctor_id")).
		Find(&surgeries).Error; err != nil {
		log.Printf("GetSurgeriesByDoctor: Error fetching surgeries - %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("GetSurgeriesByDoctor: Found %d surgeries for doctor_id %s", len(surgeries), c.Param("doctor_id"))
	c.JSON(http.StatusOK, surgeries)
}

func GetSurgeriesByPatient(c *gin.Context) {
	log.Printf("GetSurgeriesByPatient: Request received for patient_id %s", c.Param("patient_id"))

	var surgeries []models.SurgerySchedule

	if err := config.DB.Preload("Doctor").Preload("OperatingTheater").
		Where("patient_id = ?", c.Param("patient_id")).
		Find(&surgeries).Error; err != nil {
		log.Printf("GetSurgeriesByPatient: Error fetching surgeries - %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("GetSurgeriesByPatient: Found %d surgeries for patient_id %s", len(surgeries), c.Param("patient_id"))
	c.JSON(http.StatusOK, surgeries)
}
