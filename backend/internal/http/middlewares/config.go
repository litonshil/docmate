package middlewares

import (
	"github.com/labstack/echo/v4/middleware"
)

type userConfig struct {
	Skipper middleware.Skipper
}

const (
	headerUserID      = "User-Id"
	headerAdmin       = "Admin"
	headerUserName    = "User-Name"
	headerUserEmail   = "User-Email"
	headerServiceName = "Service-Name"
)
