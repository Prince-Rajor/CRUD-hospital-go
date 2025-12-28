package config

import (
	"log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := "root:prince123@tcp(127.0.0.1:3306)/hospital_db?charset=utf8mb4&parseTime=True&loc=Local"
	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to database!", err)
	}

	log.Println("Database connected successfully!")

	DB = database
}
