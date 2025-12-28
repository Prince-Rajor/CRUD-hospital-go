package database

import (
	config "CRUD-hospital-go/config"
	models "CRUD-hospital-go/models"
)

func InitializeDatabase() {
	config.ConnectDatabase()
	config.DB.AutoMigrate(&models.Doctor{}, &models.Patient{})
}
