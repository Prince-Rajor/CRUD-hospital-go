package database

import (
	"log"

	config "CRUD-hospital-go/config"
	models "CRUD-hospital-go/models"
)

func InitializeDatabase() {
	log.Println("InitializeDatabase: Connecting to database...")
	config.ConnectDatabase()
	log.Println("InitializeDatabase: Running auto migrations...")
	config.DB.AutoMigrate(
		&models.Doctor{},
		&models.Patient{},
		&models.OperatingTheater{},
		&models.SurgerySchedule{},
	)
	log.Println("InitializeDatabase: Database initialization complete")
}
