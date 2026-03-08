package controllers

import (
	"context"
	"docmate/internal/model"
	"docmate/response"
	"docmate/types"
	"docmate/utils/contextutil"
	"strconv"

	"github.com/labstack/echo/v4"
)

type MedicineController struct {
	ctx             context.Context
	medicineUseCase model.MedicineUseCase
}

func NewMedicineController(ctx context.Context, medicineUseCase model.MedicineUseCase) *MedicineController {
	return &MedicineController{
		ctx:             ctx,
		medicineUseCase: medicineUseCase,
	}
}

func (c *MedicineController) Create(ctx echo.Context) error {
	var req types.MedicineReq
	if err := ctx.Bind(&req); err != nil {
		return response.BadRequest(ctx, "invalid request body")
	}

	user, err := contextutil.GetUserFromContext(ctx)
	if err != nil {
		return response.Unauthorized(ctx, "unauthorized")
	}
	req.CreatedBy = user.ID

	medicine, err := c.medicineUseCase.Create(ctx.Request().Context(), req)
	if err != nil {
		return response.InternalServerError(ctx, err.Error())
	}

	return response.Success(ctx, "medicine created successfully", medicine)
}

func (c *MedicineController) Get(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return response.BadRequest(ctx, "invalid medicine id")
	}

	medicine, err := c.medicineUseCase.Get(ctx.Request().Context(), id)
	if err != nil {
		return response.NotFound(ctx, "medicine not found")
	}

	return response.Success(ctx, "medicine retrieved successfully", medicine)
}

func (c *MedicineController) Update(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return response.BadRequest(ctx, "invalid medicine id")
	}

	var req types.MedicineUpdateReq
	if err := ctx.Bind(&req); err != nil {
		return response.BadRequest(ctx, "invalid request body")
	}
	req.ID = id

	medicine, err := c.medicineUseCase.Update(ctx.Request().Context(), req)
	if err != nil {
		return response.InternalServerError(ctx, err.Error())
	}

	return response.Success(ctx, "medicine updated successfully", medicine)
}

func (c *MedicineController) Delete(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return response.BadRequest(ctx, "invalid medicine id")
	}

	if err := c.medicineUseCase.Delete(ctx.Request().Context(), id); err != nil {
		return response.InternalServerError(ctx, err.Error())
	}

	return response.Success(ctx, "medicine deleted successfully", nil)
}

func (c *MedicineController) List(ctx echo.Context) error {
	var req types.MedicineListReq
	if err := ctx.Bind(&req); err != nil {
		return response.BadRequest(ctx, "invalid query parameters")
	}

	medicines, err := c.medicineUseCase.List(ctx.Request().Context(), req)
	if err != nil {
		return response.InternalServerError(ctx, err.Error())
	}

	return response.Success(ctx, "medicines retrieved successfully", medicines)
}
