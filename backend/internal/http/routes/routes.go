package routes

import (
	"docmate/internal/consts"
	"docmate/internal/http/controllers"
	"docmate/internal/http/middlewares"

	"github.com/labstack/echo/v4"
)

type Routes struct {
	echo              *echo.Echo
	userController    *controllers.UserController
	doctorController  *controllers.DoctorController
	patientController *controllers.PatientController
	chamberController *controllers.ChamberController
}

func New(
	e *echo.Echo,
	userController *controllers.UserController,
	doctorController *controllers.DoctorController,
	patientController *controllers.PatientController,
	chamberController *controllers.ChamberController,
) *Routes {
	return &Routes{
		echo:              e,
		userController:    userController,
		doctorController:  doctorController,
		patientController: patientController,
		chamberController: chamberController,
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
	}

	patients := v1.Group("/patients", middlewares.AuthRoles(consts.RoleDoctor))
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
}
