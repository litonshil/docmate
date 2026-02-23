package middlewares

import (
	"bytes"
	"context"
	"docmate/config"
	"docmate/internal/model"
	"docmate/utils"
	"docmate/utils/consts"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"io/ioutil"
	"strconv"
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
			//isAdmin, _ := strconv.ParseBool(headers.Get(headerAdmin))

			user := &model.User{
				ID:        id,
				FirstName: headers.Get(headerUserFirstName),
				LastName:  headers.Get(headerUserLastName),
				Email:     headers.Get(headerUserEmail),
			}

			c.Set("user", user)

			return next(c)
		}
	}
}

// BindBody binds request body contents to bindable object
func BindBody(c echo.Context, i interface{}) error {
	// read origin body bytes
	var bodyBytes []byte
	if !utils.IsEmpty(c.Request().Body) {
		bodyBytes, _ = ioutil.ReadAll(c.Request().Body)
		// write back to request body
		c.Request().Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		// parse json data
		err := json.Unmarshal(bodyBytes, &i)
		if err != nil {
			return err
		}
	}
	return nil
}

func GenerateMetadata(c echo.Context, user *model.User) *model.User {
	if user == nil {
		user = &model.User{}
	}

	var body interface{}
	_ = BindBody(c, &body)
	appKey := c.Request().Header.Get(config.App().AppKeyHeader)
	if appKey != "" {
		appKey = "internal call (app key provided)"
	}
	serviceName := c.Request().Header.Get(headerServiceName)
	_ = serviceName
	// metadata will be passed as slack logger metadata

	return user
}

// ContextWithValue returns a new Context that carries value u.
func ContextWithValue(seedCtx context.Context, key consts.ContextKey, u interface{}) context.Context {
	switch key {
	case consts.ContextKeyUser:
		return context.WithValue(seedCtx, consts.ContextKeyUser, u.(*model.User))
	}

	return seedCtx
}
