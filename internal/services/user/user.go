package user

import (
	"context"
	"docmate/config"
	"docmate/internal/consts"
	"docmate/internal/model"
	"docmate/types"
	"fmt"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
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
		UserName: req.UserName,
		Password: string(hashedPass),
		Email:    req.Email,
		Role:     req.Role,
	}
	user, err := service.repo.CreateUser(payload)
	if err != nil {
		slog.Error("failed to create user", "error", err.Error())

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
		slog.Error("failed to get user", "error", err.Error())

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
		slog.Error("failed to list users", "error", err)

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

func (service *Service) Login(ctx context.Context, req types.LoginReq) (types.LoginResp, error) {
	user, err := service.repo.GetUserByEmail(req.Email)
	if err != nil {
		slog.Error("failed to get user by email", "error", err.Error())

		return types.LoginResp{}, fmt.Errorf("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		slog.Error("invalid password", "error", err.Error())

		return types.LoginResp{}, fmt.Errorf("invalid credentials")
	}

	token, err := service.generateJWT(user)
	if err != nil {
		slog.Error("failed to generate token", "error", err.Error())

		return types.LoginResp{}, fmt.Errorf("failed to generate token")
	}

	return types.LoginResp{
		Token: token,
		User:  mapToUserResponse(user),
	}, nil
}

func (service *Service) generateJWT(user model.UserResp) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(config.App().JWTSecret))
}

func mapToUserResponse(user model.UserResp) types.UserResp {
	resp := types.UserResp{
		ID:        user.ID,
		UserName:  user.UserName,
		Email:     user.Email,
		Password:  user.Password,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return resp
}
