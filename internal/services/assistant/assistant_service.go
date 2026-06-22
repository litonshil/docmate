package assistant

import (
	"context"
	"docmate/internal/consts"
	"docmate/internal/model"
	"docmate/types"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type service struct {
	assistantRepo model.AssistantRepo
	userRepo      model.UserRepo
	chamberRepo   model.ChamberRepo
}

func NewService(assistantRepo model.AssistantRepo, userRepo model.UserRepo, chamberRepo model.ChamberRepo) model.AssistantUseCase {
	return &service{
		assistantRepo: assistantRepo,
		userRepo:      userRepo,
		chamberRepo:   chamberRepo,
	}
}

func (s *service) Create(ctx context.Context, req types.AssistantReq, doctorID int) (types.AssistantResp, error) {
	if err := req.Validate(); err != nil {
		return types.AssistantResp{}, err
	}

	// Check if email already exists
	existingUser, err := s.userRepo.GetByEmail(req.Email)
	if err == nil && existingUser.ID != 0 {
		return types.AssistantResp{}, fmt.Errorf("user with this email already exists")
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return types.AssistantResp{}, err
	}

	// Create user
	userPayload := model.User{
		UserName: req.Name,
		Email:    req.Email,
		Password: string(hashedPass),
		Role:     consts.RoleAssistant,
	}

	user, err := s.userRepo.Create(userPayload)
	if err != nil {
		return types.AssistantResp{}, err
	}

	// Resolve chambers
	var validChamberIDs []int
	var chambers []model.Chamber
	for _, cid := range req.ChamberIDs {
		ch, err := s.chamberRepo.GetChamberByID(cid)
		if err == nil && ch.DoctorID == doctorID {
			validChamberIDs = append(validChamberIDs, cid)
			chambers = append(chambers, ch)
		}
	}

	// Create assistant profile
	assistantPayload := model.Assistant{
		UserID:    user.ID,
		DoctorID:  doctorID,
		Name:      req.Name,
		Phone:     req.Phone,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	assistant, err := s.assistantRepo.CreateAssistant(assistantPayload, validChamberIDs)
	if err != nil {
		return types.AssistantResp{}, err
	}

	return mapToAssistantResponse(assistant, user.Email, chambers), nil
}

func (s *service) Update(ctx context.Context, id int, req types.AssistantUpdateReq, doctorID int) (types.AssistantResp, error) {
	if err := req.Validate(); err != nil {
		return types.AssistantResp{}, err
	}

	assistant, err := s.assistantRepo.GetAssistantByID(id)
	if err != nil {
		return types.AssistantResp{}, err
	}

	if assistant.DoctorID != doctorID {
		return types.AssistantResp{}, errors.New("unauthorized to update this assistant")
	}

	user, err := s.userRepo.Get(assistant.UserID)
	if err != nil {
		return types.AssistantResp{}, err
	}

	// Update user email if changed
	if req.Email != user.Email {
		existing, err := s.userRepo.GetByEmail(req.Email)
		if err == nil && existing.ID != 0 {
			return types.AssistantResp{}, fmt.Errorf("user with this email already exists")
		}
		user.Email = req.Email
	}

	// Update password if provided
	if req.Password != "" {
		hashedPass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return types.AssistantResp{}, err
		}
		user.Password = string(hashedPass)
	}

	user.UserName = req.Name
	_, err = s.userRepo.Update(user)
	if err != nil {
		return types.AssistantResp{}, err
	}

	// Update assistant profile
	assistant.Name = req.Name
	assistant.Phone = req.Phone
	if req.IsActive != nil {
		assistant.IsActive = *req.IsActive
	}
	assistant.UpdatedAt = time.Now()

	// Filter valid chambers
	var validChamberIDs []int
	for _, cid := range req.ChamberIDs {
		ch, err := s.chamberRepo.GetChamberByID(cid)
		if err == nil && ch.DoctorID == doctorID {
			validChamberIDs = append(validChamberIDs, cid)
		}
	}

	updatedAssistant, err := s.assistantRepo.UpdateAssistant(assistant, validChamberIDs)
	if err != nil {
		return types.AssistantResp{}, err
	}

	chambers, err := s.assistantRepo.GetChambersByAssistantID(updatedAssistant.ID)
	if err != nil {
		return mapToAssistantResponse(updatedAssistant, user.Email, nil), err
	}

	return mapToAssistantResponse(updatedAssistant, user.Email, chambers), nil
}

func (s *service) Get(ctx context.Context, id int, doctorID int) (types.AssistantResp, error) {
	assistant, err := s.assistantRepo.GetAssistantByID(id)
	if err != nil {
		return types.AssistantResp{}, err
	}

	if assistant.DoctorID != doctorID {
		return types.AssistantResp{}, errors.New("unauthorized to view this assistant")
	}

	chambers, err := s.assistantRepo.GetChambersByAssistantID(assistant.ID)
	if err != nil {
		return mapToAssistantResponse(assistant, assistant.User.Email, nil), err
	}

	return mapToAssistantResponse(assistant, assistant.User.Email, chambers), nil
}

func (s *service) List(ctx context.Context, doctorID int) ([]types.AssistantResp, error) {
	assistants, err := s.assistantRepo.ListAssistantsByDoctorID(doctorID)
	if err != nil {
		return nil, err
	}

	var resp []types.AssistantResp
	for _, a := range assistants {
		chambers, _ := s.assistantRepo.GetChambersByAssistantID(a.ID)
		resp = append(resp, mapToAssistantResponse(a, a.User.Email, chambers))
	}

	return resp, nil
}

func mapToAssistantResponse(a model.Assistant, email string, chambersList []model.Chamber) types.AssistantResp {
	var chambers []types.ChamberResp
	for _, c := range chambersList {
		chambers = append(chambers, types.ChamberResp{
			ID:       c.ID,
			DoctorID: c.DoctorID,
			Name:     c.Name,
			Address:  c.Address,
			Phone:    c.Phone,
			Fee:      c.Fee,
			IsActive: c.IsActive,
		})
	}

	return types.AssistantResp{
		ID:        a.ID,
		UserID:    a.UserID,
		DoctorID:  a.DoctorID,
		Name:      a.Name,
		Phone:     a.Phone,
		Email:     email,
		IsActive:  a.IsActive,
		Chambers:  chambers,
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
	}
}
