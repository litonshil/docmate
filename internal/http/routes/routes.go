package routes

import (
	"docmate/internal/consts"
	"docmate/internal/http/controllers"
	"docmate/internal/http/middlewares"
	"os"

	"github.com/labstack/echo/v4"
)

type Routes struct {
	echo                          *echo.Echo
	userController                *controllers.UserController
	doctorController              *controllers.DoctorController
	patientController             *controllers.PatientController
	chamberController             *controllers.ChamberController
	medicineController            *controllers.MedicineController
	prescriptionController        *controllers.PrescriptionController
	prescriptionSettingController *controllers.PrescriptionSettingController
	uploadController              *controllers.UploadController
	dashboardController           *controllers.DashboardController
	aiSettingController           *controllers.AISuggestionController
	appointmentController         *controllers.AppointmentController
}

func New(
	e *echo.Echo,
	userController *controllers.UserController,
	doctorController *controllers.DoctorController,
	patientController *controllers.PatientController,
	chamberController *controllers.ChamberController,
	medicineController *controllers.MedicineController,
	prescriptionController *controllers.PrescriptionController,
	prescriptionSettingController *controllers.PrescriptionSettingController,
	dashboardController *controllers.DashboardController,
	aiSettingController *controllers.AISuggestionController,
	appointmentController *controllers.AppointmentController,
) *Routes {
	return &Routes{
		echo:                          e,
		userController:                userController,
		doctorController:              doctorController,
		patientController:             patientController,
		chamberController:             chamberController,
		medicineController:            medicineController,
		prescriptionController:        prescriptionController,
		prescriptionSettingController: prescriptionSettingController,
		uploadController:              controllers.NewUploadController(),
		dashboardController:           dashboardController,
		aiSettingController:           aiSettingController,
		appointmentController:         appointmentController,
	}
}

func (r *Routes) Init() {
	e := r.echo

	// Public routes
	v1 := e.Group("/v1")

	user := v1.Group("/users")
	{
		user.GET("", r.userController.List, middlewares.AuthRoles(consts.RoleAdmin))
		user.POST("/register", r.userController.Create)
		user.POST("/login", r.userController.Login)
	}

	doctors := v1.Group("/doctors", middlewares.AuthRoles(consts.RoleAdmin, consts.RoleDoctor))
	{
		doctors.GET("/profile", r.doctorController.GetProfile)
		doctors.GET("", r.doctorController.List)
		doctors.POST("", r.doctorController.Create)
		doctors.GET("/:id", r.doctorController.Get)
		doctors.PUT("/:id", r.doctorController.Update)
		doctors.PATCH("/:id/ai-settings", r.aiSettingController.AdminUpdateSettings, middlewares.AuthRoles(consts.RoleAdmin))
	}

	patients := v1.Group("/patients", middlewares.AuthRoles(consts.RoleAdmin, consts.RoleDoctor))
	{
		patients.GET("", r.patientController.List)
		patients.POST("", r.patientController.Create)
		patients.GET("/:id", r.patientController.Get)
		patients.PUT("/:id", r.patientController.Update)
	}

	chambers := v1.Group("/doctors/:doctor_id/chambers", middlewares.AuthRoles(consts.RoleAdmin, consts.RoleDoctor))
	{
		chambers.GET("", r.chamberController.List)
		chambers.POST("", r.chamberController.Create)
		chambers.GET("/:id", r.chamberController.Get)
		chambers.PUT("/:id", r.chamberController.Update)
	}

	medicines := v1.Group("/medicines", middlewares.AuthRoles(consts.RoleAdmin, consts.RoleDoctor))
	{
		medicines.GET("", r.medicineController.List)
		medicines.POST("", r.medicineController.Create)
		medicines.GET("/:id", r.medicineController.Get)
		medicines.PUT("/:id", r.medicineController.Update)
		medicines.DELETE("/:id", r.medicineController.Delete)
	}

	prescriptions := v1.Group("/prescriptions", middlewares.AuthRoles(consts.RoleAdmin, consts.RoleDoctor))
	{
		prescriptions.GET("", r.prescriptionController.List)
		prescriptions.POST("", r.prescriptionController.Create)
		prescriptions.GET("/:id", r.prescriptionController.Get)
		prescriptions.PUT("/:id", r.prescriptionController.Update)
	}

	prescriptionSettings := v1.Group("/doctors/:doctor_id/prescription-settings", middlewares.AuthRoles(consts.RoleAdmin, consts.RoleDoctor))
	{
		prescriptionSettings.GET("", r.prescriptionSettingController.GetByChamber)
		prescriptionSettings.POST("", r.prescriptionSettingController.Upsert)
	}

	dashboard := v1.Group("/dashboard", middlewares.AuthRoles(consts.RoleAdmin, consts.RoleDoctor))
	{
		dashboard.GET("/summary", r.dashboardController.GetSummary)
	}

	appointments := v1.Group("/appointments", middlewares.AuthRoles(consts.RoleAdmin, consts.RoleDoctor))
	{
		appointments.GET("", r.appointmentController.List)
		appointments.POST("", r.appointmentController.Book)
		appointments.GET("/:id", r.appointmentController.GetDetails)
		appointments.PATCH("/:id/status", r.appointmentController.UpdateStatus)
		appointments.PATCH("/:id/collect-fee", r.appointmentController.CollectFee)
	}

	aiSuggestions := v1.Group("/suggestions", middlewares.AuthRoles(consts.RoleDoctor))
	{
		aiSuggestions.POST("", r.aiSettingController.GetSuggestions)
	}

	aiSettings := v1.Group("/settings/ai", middlewares.AuthRoles(consts.RoleDoctor))
	{
		aiSettings.GET("", r.aiSettingController.GetSettings)
		aiSettings.POST("", r.aiSettingController.UpsertSettings)
	}

	// Upload route
	v1.POST("/upload", r.uploadController.UploadFile, middlewares.AuthRoles(consts.RoleAdmin, consts.RoleDoctor))

	// Static file serving
	uploadPath := "uploads"
	if _, err := os.Stat("/project"); err == nil {
		uploadPath = "/project/uploads"
	}
	e.Static("/uploads", uploadPath)
}
