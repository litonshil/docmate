package routes

import (
	"docmate/internal/http/controllers"
	"docmate/internal/http/middlewares"

	"github.com/labstack/echo/v4"
)

type Routes struct {
	echo             *echo.Echo
	userController   *controllers.UserController
	doctorController *controllers.DoctorController
}

func New(
	e *echo.Echo,
	userController *controllers.UserController,
	doctorController *controllers.DoctorController,
) *Routes {
	return &Routes{
		echo:             e,
		userController:   userController,
		doctorController: doctorController,
	}
}

func (r *Routes) Init() {
	e := r.echo

	// Public routes
	v1 := e.Group("/v1")

	user := v1.Group("/users")
	{
		user.GET("", r.userController.List, middlewares.AuthRoles("admin"))
		user.POST("/register", r.userController.Create)
		user.POST("/login", r.userController.Login)
	}

	doctors := v1.Group("/doctors", middlewares.AuthRoles("admin", "doctor"))
	{
		// General CRUD Endpoints
		doctors.GET("", r.doctorController.List)
		doctors.POST("", r.doctorController.Create)
		doctors.GET("/:id", r.doctorController.Get)
		doctors.PUT("/:id", r.doctorController.Update)
	}
}
