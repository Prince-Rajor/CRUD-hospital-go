package routers

import (
	"CRUD-hospital-go/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Welcome to Hospital API"})
	})

	// Doctor Routes
	router.POST("/doctor/", controllers.CreateDoctor)
	router.GET("/doctor/:id", controllers.GetDoctorByID)
	router.PATCH("/doctor/:id", controllers.UpdateDoctor)
	router.DELETE("/doctor/:id", controllers.DeleteDoctor)
	router.GET("/searchDoctorByName", controllers.SearchDoctorByName)

	// Patient Routes
	router.POST("/patient/", controllers.CreatePatient)
	router.GET("/patient/:id", controllers.GetPatientByID)
	router.PATCH("/patient/:id", controllers.UpdatePatient)
	router.GET("/fetchPatientByDoctorId/:doctor_id", controllers.GetPatientsByDoctorID)
	router.DELETE("/patient/:id", controllers.DeletePatient)
	router.GET("/searchPatientByName", controllers.SearchPatientByName)

	return router
}
