package controllers

import (
	m "docmate/internal/http/middlewares"
	"docmate/internal/model"

	"github.com/labstack/echo/v4"
)

func parseUser(c echo.Context) *model.User {
	if c.Get("user") == nil {
		user := m.GenerateMetadata(c, nil)

		return user
	}

	res := c.Get("user").(*model.User)

	return res
}
