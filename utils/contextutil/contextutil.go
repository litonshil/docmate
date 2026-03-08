package contextutil

import (
	"docmate/types"
	"errors"

	"github.com/labstack/echo/v4"
)

func GetUserFromContext(c echo.Context) (*types.AuthUser, error) {
	user, ok := c.Get("user").(*types.AuthUser)
	if !ok {
		return nil, errors.New("user not found in context")
	}

	return user, nil
}
