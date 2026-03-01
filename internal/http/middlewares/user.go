package middlewares

import (
	"docmate/internal/model"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func authorizeUser(config userConfig) echo.MiddlewareFunc {
	if config.Skipper == nil {
		config.Skipper = middleware.DefaultSkipper
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			headers := c.Request().Header

			id, _ := strconv.Atoi(headers.Get(headerUserID))

			user := &model.User{
				ID:       id,
				UserName: headers.Get(headerUserName),
				Email:    headers.Get(headerUserEmail),
			}

			c.Set("user", user)

			return next(c)
		}
	}
}
