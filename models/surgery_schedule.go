package models

import (
	"time"

	"gorm.io/gorm"
)

type SurgeryStatus string

const (
	SurgeryStatusScheduled  SurgeryStatus = "Scheduled"
	SurgeryStatusInProgress SurgeryStatus = "In Progress"
	SurgeryStatusCompleted  SurgeryStatus = "Completed"
	SurgeryStatusCancelled  SurgeryStatus = "Cancelled"
)

type SurgerySchedule struct {
	gorm.Model
	PatientID          uint             `json:"patient_id"`
	Patient            Patient          `json:"patient" gorm:"foreignKey:PatientID"`
	DoctorID           uint             `json:"doctor_id"`
	Doctor             Doctor           `json:"doctor" gorm:"foreignKey:DoctorID"`
	OperatingTheaterID uint             `json:"operating_theater_id"`
	OperatingTheater   OperatingTheater `json:"operating_theater" gorm:"foreignKey:OperatingTheaterID"`
	SurgeryType        string           `json:"surgery_type"`
	ScheduledAt        time.Time        `json:"scheduled_at"`
	EstimatedDuration  int              `json:"estimated_duration"`
	DepositDeducted    float64          `json:"deposit_deducted"`
	Status             SurgeryStatus    `json:"status" gorm:"default:'Scheduled'"`
	Notes              string           `json:"notes"`
}

type SurgeryScheduleRequest struct {
	PatientID         uint      `json:"patient_id" binding:"required"`
	DoctorID          uint      `json:"doctor_id" binding:"required"`
	SurgeryType       string    `json:"surgery_type" binding:"required"`
	ScheduledAt       time.Time `json:"scheduled_at" binding:"required"`
	EstimatedDuration int       `json:"estimated_duration" binding:"required"`
	DepositRequired   float64   `json:"deposit_required" binding:"required"`
	Notes             string    `json:"notes"`
}
