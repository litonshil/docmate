package user

import (
	"context"
	"docmate/internal/consts"
	"docmate/internal/model"
	"docmate/types"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
)

type Service struct {
	repo model.UserRepo
}

func NewService(userRepo model.UserRepo) *Service {
	return &Service{
		repo: userRepo,
	}
}

func (service *Service) CreateUser(ctx context.Context, req types.UserReq) (types.UserResp, error) {

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return types.UserResp{}, err
	}

	payload := model.User{
		Type:      req.Type,
		UserName:  req.UserName,
		Password:  string(hashedPass),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Phone:     req.Phone,
	}
	user, err := service.repo.CreateUser(payload)
	if err != nil {
		slog.Error("failed to create user", err)
		return types.UserResp{}, err
	}

	resp := types.UserResp{
		ID: user.ID,
	}
	return resp, nil
}
func (service *Service) GetUser(ctx context.Context, userID int) (types.UserResp, error) {
	user, err := service.repo.GetUser(userID)
	if err != nil {
		slog.Error("failed to get user", err)
		return types.UserResp{}, fmt.Errorf("failed to get user: %w", err)
	}
	resp := mapToUserResponse(user)
	return resp, nil
}

func (service *Service) ListUsers(ctx context.Context, req types.UserListReq) (types.PaginatedResponse, error) {

	if req.Page == 0 || req.Limit == 0 {
		req.Page = consts.Page
		req.Limit = consts.Limit
	}

	offset := (req.Page - 1) * req.Limit

	users, total, err := service.repo.ListUsers(offset, req.Limit)
	if err != nil {
		slog.Error("failed to list users", err)
		return types.PaginatedResponse{}, fmt.Errorf("failed to list users: %w", err)
	}

	userResp := make([]types.UserResp, len(users))
	for i, user := range users {
		userResp[i] = mapToUserResponse(user)
	}

	lastPage := (total + req.Limit - 1) / req.Limit

	resp := types.PaginatedResponse{
		Pagination: types.Pagination{
			Page:     req.Page,
			Limit:    req.Limit,
			Total:    total,
			LastPage: lastPage,
		},
		Records: userResp,
	}

	return resp, nil
}

func mapToUserResponse(user model.UserResp) types.UserResp {
	resp := types.UserResp{
		ID:        user.ID,
		Type:      user.Type,
		UserName:  user.UserName,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Phone:     user.Phone,
		Password:  user.Password,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
	return resp
}
