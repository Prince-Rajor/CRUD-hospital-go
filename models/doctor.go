package models

import "gorm.io/gorm"

type Doctor struct {
	gorm.Model
	Name        string `json:"name"`
	ContactNo   string `json:"contact_no"`
	Address     string `json:"address"`
	IsAvailable bool   `json:"is_available" gorm:"default:true"`
}
