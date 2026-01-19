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
	router.GET("/doctors/", controllers.GetAllDoctors)
	router.POST("/doctor/", controllers.CreateDoctor)
	router.GET("/doctor/:id", controllers.GetDoctorByID)
	router.GET("/doctor/:id/availability", controllers.CheckDoctorAvailability)
	router.PATCH("/doctor/:id", controllers.UpdateDoctor)
	router.DELETE("/doctor/:id", controllers.DeleteDoctor)
	router.GET("/searchDoctorByName", controllers.SearchDoctorByName)

	// Patient Routes
	router.GET("/patients/", controllers.GetAllPatients)
	router.POST("/patient/", controllers.CreatePatient)
	router.GET("/patient/:id", controllers.GetPatientByID)
	router.PATCH("/patient/:id", controllers.UpdatePatient)
	router.GET("/fetchPatientByDoctorId/:doctor_id", controllers.GetPatientsByDoctorID)
	router.DELETE("/patient/:id", controllers.DeletePatient)
	router.GET("/searchPatientByName", controllers.SearchPatientByName)

	// Operating Theater Routes
	router.POST("/operating-theater/", controllers.CreateOperatingTheater)
	router.GET("/operating-theater/:id", controllers.GetOperatingTheaterByID)
	router.GET("/operating-theaters/", controllers.GetAllOperatingTheaters)
	router.GET("/operating-theaters/available", controllers.GetAvailableOperatingTheaters)
	router.PATCH("/operating-theater/:id", controllers.UpdateOperatingTheater)
	router.DELETE("/operating-theater/:id", controllers.DeleteOperatingTheater)

	// Surgery Scheduling Routes (Transactional)
	router.POST("/surgery/schedule", controllers.ScheduleSurgery)           // Schedule a new surgery (THE MAIN TRANSACTION)
	router.POST("/surgery/:id/complete", controllers.CompleteSurgery)       // Mark surgery as completed
	router.POST("/surgery/:id/cancel", controllers.CancelSurgery)           // Cancel surgery and refund deposit
	router.GET("/surgery/:id", controllers.GetSurgeryByID)                  // Get surgery details
	router.GET("/surgeries/", controllers.GetAllSurgeries)                  // Get all surgeries
	router.GET("/surgeries/doctor/:doctor_id", controllers.GetSurgeriesByDoctor)   // Get surgeries by doctor
	router.GET("/surgeries/patient/:patient_id", controllers.GetSurgeriesByPatient) // Get surgeries by patient

	return router
}
