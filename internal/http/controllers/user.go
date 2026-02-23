package controllers

import (
	"context"
	"docmate/internal/http/middlewares"
	"docmate/internal/model"
	"docmate/response"
	"docmate/types"
	"docmate/utils/consts"

	"github.com/labstack/echo/v4"
)

type UserController struct {
	baseCtx context.Context
	userSvc model.UserUseCase
}

func NewUserController(
	baseCtx context.Context,
	userSvc model.UserUseCase,
) *UserController {
	return &UserController{
		baseCtx: baseCtx,
		userSvc: userSvc,
	}
}
func (controller *UserController) CreateUser(c echo.Context) error {
	ctx := middlewares.ContextWithValue(controller.baseCtx, consts.ContextKeyUser, parseUser(c))
	_ = ctx

	var req types.UserReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(response.BadRequest(err.Error()))
	}

	if err := req.Validate(); err != nil {
		return c.JSON(response.BadRequest(err.Error()))
	}

	resp, err := controller.userSvc.CreateUser(ctx, req)
	if err != nil {
		return c.JSON(response.InternalServerError(err.Error()))
	}

	return c.JSON(response.Success("user created successfully", []types.UserResp{resp}))
}

func (controller *UserController) GetUser(c echo.Context) error {
	ctx := c.Request().Context()
	var req types.UserFilter
	if err := c.Bind(&req); err != nil {
		return c.JSON(response.BadRequest(err.Error()))
	}

	if err := req.Validate(); err != nil {
		return c.JSON(response.BadRequest(err.Error()))
	}

	user, err := controller.userSvc.GetUser(ctx, req.ID)
	if err != nil {
		return c.JSON(response.InternalServerError(err.Error()))
	}

	return c.JSON(response.Success("user fetched successfully", user))
}

func (controller *UserController) ListUsers(c echo.Context) error {
	ctx := c.Request().Context()
	var req types.UserListReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(response.BadRequest(err.Error()))
	}

	users, err := controller.userSvc.ListUsers(ctx, req)
	if err != nil {
		return c.JSON(response.InternalServerError(err.Error()))
	}

	return c.JSON(response.Success("users fetched successfully", users))
}
