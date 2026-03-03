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
}

func New(
	e *echo.Echo,
	userController *controllers.UserController,
	doctorController *controllers.DoctorController,
	patientController *controllers.PatientController,
) *Routes {
	return &Routes{
		echo:              e,
		userController:    userController,
		doctorController:  doctorController,
		patientController: patientController,
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
		// General CRUD Endpoints
		doctors.GET("", r.doctorController.List)
		doctors.POST("", r.doctorController.Create)
		doctors.GET("/:id", r.doctorController.Get)
		doctors.PUT("/:id", r.doctorController.Update)
	}

	patients := v1.Group("/patients", middlewares.AuthRoles(consts.RoleAdmin, consts.RoleDoctor))
	{
		// General CRUD Endpoints
		patients.GET("", r.patientController.List)
		patients.POST("", r.patientController.Create)
		patients.GET("/:id", r.patientController.Get)
		patients.PUT("/:id", r.patientController.Update)
	}
}
