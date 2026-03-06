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
	userRepo   model.UserRepo
	doctorRepo model.DoctorRepo
}

func NewService(userRepo model.UserRepo, doctorRepo model.DoctorRepo) *Service {
	return &Service{
		userRepo:   userRepo,
		doctorRepo: doctorRepo,
	}
}

func (service *Service) Create(ctx context.Context, req types.UserReq) (types.UserResp, error) {
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
	user, err := service.userRepo.Create(payload)
	if err != nil {
		slog.Error("failed to create user", "error", err.Error())

		return types.UserResp{}, err
	}

	resp := types.UserResp{
		ID: user.ID,
	}

	return resp, nil
}
func (service *Service) Get(ctx context.Context, userID int) (types.UserResp, error) {
	user, err := service.userRepo.Get(userID)
	if err != nil {
		slog.Error("failed to get user", "error", err.Error())

		return types.UserResp{}, fmt.Errorf("failed to get user: %w", err)
	}
	doc, _ := service.doctorRepo.GetDoctorByUserID(user.ID)
	resp := mapToUserResponse(user, doc.ID != 0)

	return resp, nil
}

func (service *Service) List(ctx context.Context, req types.UserListReq) (types.PaginatedResponse, error) {
	if req.Page == 0 || req.Limit == 0 {
		req.Page = consts.Page
		req.Limit = consts.Limit
	}

	offset := (req.Page - 1) * req.Limit

	users, total, err := service.userRepo.List(offset, req.Limit)
	if err != nil {
		slog.Error("failed to list users", "error", err)

		return types.PaginatedResponse{}, fmt.Errorf("failed to list users: %w", err)
	}

	userResp := make([]types.UserResp, len(users))
	for i, user := range users {
		doc, _ := service.doctorRepo.GetDoctorByUserID(user.ID)
		userResp[i] = mapToUserResponse(user, doc.ID != 0)
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
	user, err := service.userRepo.GetByEmail(req.Email)
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

	doc, _ := service.doctorRepo.GetDoctorByUserID(user.ID)
	isProfileCompleted := doc.ID != 0

	return types.LoginResp{
		Token: token,
		User:  mapToUserResponse(user, isProfileCompleted),
	}, nil
}

func (service *Service) generateJWT(user model.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id":   user.ID,
		"user_name": user.UserName,
		"email":     user.Email,
		"role":      user.Role,
		"exp":       time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(config.App().JWTSecret))
}

func mapToUserResponse(user model.User, isProfileCompleted bool) types.UserResp {
	resp := types.UserResp{
		ID:                 user.ID,
		UserName:           user.UserName,
		Email:              user.Email,
		Password:           user.Password,
		Role:               user.Role,
		IsProfileCompleted: isProfileCompleted,
		CreatedAt:          user.CreatedAt,
		UpdatedAt:          user.UpdatedAt,
	}

	return resp
}
