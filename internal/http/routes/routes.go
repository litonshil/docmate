package routes

import (
	"docmate/internal/http/controllers"

	"github.com/labstack/echo/v4"
)

type Routes struct {
	echo           *echo.Echo
	userController *controllers.UserController
}

func New(
	e *echo.Echo,
	userController *controllers.UserController,
) *Routes {
	return &Routes{
		echo:           e,
		userController: userController,
	}
}

func (r *Routes) Init() {
	e := r.echo

	// Public routes
	v1 := e.Group("/v1")

	user := v1.Group("/users")
	{
		user.GET("", r.userController.ListUsers)
		user.POST("/register", r.userController.CreateUser)
		user.GET("/:id", r.userController.GetUser)
	}

}
