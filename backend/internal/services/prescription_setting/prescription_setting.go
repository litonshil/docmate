package prescription_setting

import (
	"context"
	"docmate/internal/model"
	"docmate/types"
	"log/slog"
)

type Service struct {
	repo model.PrescriptionSettingRepo
}

func NewService(repo model.PrescriptionSettingRepo) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Upsert(ctx context.Context, req types.PrescriptionSettingReq) (types.PrescriptionSettingResp, error) {
	payload := model.PrescriptionSetting{
		DoctorID:           req.DoctorID,
		ChamberID:          req.ChamberID,
		HeaderLeftBangla:   req.HeaderLeftBangla,
		HeaderRightEnglish: req.HeaderRightEnglish,
		FooterInfoBangla:   req.FooterInfoBangla,
		FooterInfoEnglish:  req.FooterInfoEnglish,
		TemplateType:       req.TemplateType,
	}

	setting, err := s.repo.Upsert(payload)
	if err != nil {
		slog.Error("failed to upsert prescription setting", "error", err.Error())
		return types.PrescriptionSettingResp{}, err
	}

	return mapToResponse(setting), nil
}

func (s *Service) GetByChamber(ctx context.Context, doctorID, chamberID int) (types.PrescriptionSettingResp, error) {
	setting, err := s.repo.GetByChamber(doctorID, chamberID)
	if err != nil {
		slog.Error("failed to get prescription setting", "doctor_id", doctorID, "chamber_id", chamberID, "error", err.Error())
		return types.PrescriptionSettingResp{}, err
	}

	return mapToResponse(setting), nil
}

func mapToResponse(setting model.PrescriptionSetting) types.PrescriptionSettingResp {
	return types.PrescriptionSettingResp{
		ID:                 setting.ID,
		DoctorID:           setting.DoctorID,
		ChamberID:          setting.ChamberID,
		HeaderLeftBangla:   setting.HeaderLeftBangla,
		HeaderRightEnglish: setting.HeaderRightEnglish,
		FooterInfoBangla:   setting.FooterInfoBangla,
		FooterInfoEnglish:  setting.FooterInfoEnglish,
		TemplateType:       setting.TemplateType,
		CreatedAt:          setting.CreatedAt,
		UpdatedAt:          setting.UpdatedAt,
	}
}
