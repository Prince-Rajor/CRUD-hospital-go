package main

import (
	"CRUD-hospital-go/database"
	"CRUD-hospital-go/routers"
)

func main() {
	database.InitializeDatabase()
	router := routers.SetupRouter()
	router.Run(":8080")
}
