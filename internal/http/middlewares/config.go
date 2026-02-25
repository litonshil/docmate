package middlewares

import (
	"github.com/labstack/echo/v4/middleware"
)

type userConfig struct {
	Skipper middleware.Skipper
}

const (
	headerUserID        = "User-Id"
	headerAdmin         = "Admin"
	headerUserFirstName = "User-Firstname"
	headerUserLastName  = "User-Lastname"
	headerUserEmail     = "User-Email"
	headerServiceName   = "Service-Name"
)
