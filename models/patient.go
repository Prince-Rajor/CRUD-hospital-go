package models

import "gorm.io/gorm"

type Patient struct {
	gorm.Model
	Name      string  `json:"name"`
	ContactNo string  `json:"contact_no"`
	Address   string  `json:"address"`
	DoctorID  uint    `json:"doctor_id"`
	Deposit   float64 `json:"deposit" gorm:"default:0"`
}
