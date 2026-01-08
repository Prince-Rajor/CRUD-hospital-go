package models

import "gorm.io/gorm"

type OTStatus string

const (
	OTStatusAvailable   OTStatus = "Available"
	OTStatusOccupied    OTStatus = "Occupied"
	OTStatusMaintenance OTStatus = "Maintenance"
)

type OperatingTheater struct {
	gorm.Model
	Name     string   `json:"name"`
	Floor    int      `json:"floor"`
	Status   OTStatus `json:"status" gorm:"default:'Available'"`
	Capacity int      `json:"capacity"`
}
