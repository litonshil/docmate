package prescription

import (
	"context"
	"docmate/internal/model"
	"docmate/types"
	"encoding/json"
	"errors"

	"gorm.io/datatypes"
)

type service struct {
	repo model.PrescriptionRepo
}

func NewService(repo model.PrescriptionRepo) model.PrescriptionUseCase {
	return &service{
		repo: repo,
	}
}

func (svc *service) Create(ctx context.Context, req types.PrescriptionReq) (types.PrescriptionResp, error) {
	// TODO: Cross-check if patient belongs to doctor
	// TODO: Cross-check if chamber belongs to doctor

	vitalsJSON, _ := json.Marshal(req.Vitals)
	complaintsJSON, _ := json.Marshal(req.ChiefComplaints)
	diagnosisJSON, _ := json.Marshal(req.Diagnosis)
	medicationsJSON, _ := json.Marshal(req.Medications)
	investigationsJSON, _ := json.Marshal(req.Investigations)

	if len(req.ChiefComplaints) == 0 {
		complaintsJSON = []byte("[]")
	}
	if len(req.Diagnosis) == 0 {
		diagnosisJSON = []byte("[]")
	}
	if len(req.Medications) == 0 {
		medicationsJSON = []byte("[]")
	}
	if len(req.Investigations) == 0 {
		investigationsJSON = []byte("[]")
	}

	p := model.Prescription{
		DoctorID:        req.DoctorID,
		PatientID:       req.PatientID,
		ChamberID:       req.ChamberID,
		Vitals:          datatypes.JSON(vitalsJSON),
		ChiefComplaints: datatypes.JSON(complaintsJSON),
		Diagnosis:       datatypes.JSON(diagnosisJSON),
		Medications:     datatypes.JSON(medicationsJSON),
		Investigations:  datatypes.JSON(investigationsJSON),
		Advice:          req.Advice,
		Status:          req.Status,
		FollowUpDate:    req.FollowUpDate,
	}

	if p.Status == "" {
		p.Status = "draft"
	}

	saved, err := svc.repo.CreatePrescription(p)
	if err != nil {
		return types.PrescriptionResp{}, err
	}

	return svc.mapToResp(saved), nil
}

func (svc *service) Update(ctx context.Context, id int, req types.PrescriptionReq) (types.PrescriptionResp, error) {
	existing, err := svc.repo.GetPrescriptionByID(id)
	if err != nil {
		return types.PrescriptionResp{}, errors.New("prescription not found")
	}

	if existing.DoctorID != req.DoctorID {
		return types.PrescriptionResp{}, errors.New("unauthorized")
	}

	if existing.Status == "finalized" {
		return types.PrescriptionResp{}, errors.New("finalized prescription cannot be updated")
	}

	vitalsJSON, _ := json.Marshal(req.Vitals)
	complaintsJSON, _ := json.Marshal(req.ChiefComplaints)
	diagnosisJSON, _ := json.Marshal(req.Diagnosis)
	medicationsJSON, _ := json.Marshal(req.Medications)
	investigationsJSON, _ := json.Marshal(req.Investigations)

	existing.Vitals = datatypes.JSON(vitalsJSON)
	existing.ChiefComplaints = datatypes.JSON(complaintsJSON)
	existing.Diagnosis = datatypes.JSON(diagnosisJSON)
	existing.Medications = datatypes.JSON(medicationsJSON)
	existing.Investigations = datatypes.JSON(investigationsJSON)
	existing.Advice = req.Advice
	existing.Status = req.Status
	existing.FollowUpDate = req.FollowUpDate

	if existing.Status == "" {
		existing.Status = "draft"
	}

	updated, err := svc.repo.UpdatePrescription(existing)
	if err != nil {
		return types.PrescriptionResp{}, err
	}

	return svc.mapToResp(updated), nil
}

func (svc *service) Get(ctx context.Context, id int, doctorID int) (types.PrescriptionResp, error) {
	p, err := svc.repo.GetPrescriptionByID(id)
	if err != nil {
		return types.PrescriptionResp{}, errors.New("prescription not found")
	}

	if p.DoctorID != doctorID {
		return types.PrescriptionResp{}, errors.New("prescription not found") // Or unauthorized
	}

	return svc.mapToResp(p), nil
}

func (svc *service) List(ctx context.Context, req types.PrescriptionListReq) (types.PaginatedPrescriptionResp, error) {
	offset := (req.Page - 1) * req.Limit

	var records []model.Prescription
	var total int
	var err error

	if req.PatientID > 0 {
		records, total, err = svc.repo.ListPrescriptionsByPatient(req.DoctorID, req.PatientID, req.Limit, offset, req.Search)
	} else {
		records, total, err = svc.repo.ListPrescriptions(req.DoctorID, req.Limit, offset, req.Search)
	}

	if err != nil {
		return types.PaginatedPrescriptionResp{}, err
	}

	resps := make([]types.PrescriptionResp, 0, len(records))
	for _, rec := range records {
		resps = append(resps, svc.mapToResp(rec))
	}

	lastPage := total / req.Limit
	if total%req.Limit != 0 {
		lastPage++
	}

	return types.PaginatedPrescriptionResp{
		Pagination: types.Pagination{
			Page:     req.Page,
			Limit:    req.Limit,
			Total:    total,
			LastPage: lastPage,
		},
		Records: resps,
	}, nil
}

func (svc *service) mapToResp(p model.Prescription) types.PrescriptionResp {
	resp := types.PrescriptionResp{
		ID:              p.ID,
		DoctorID:        p.DoctorID,
		PatientID:       p.PatientID,
		PatientName:     p.PatientName,
		ChamberID:       p.ChamberID,
		Advice:          p.Advice,
		Status:          p.Status,
		FollowUpDate:    p.FollowUpDate,
		CreatedAt:       p.CreatedAt,
		UpdatedAt:       *p.UpdatedAt,
	}

	_ = json.Unmarshal(p.Vitals, &resp.Vitals)
	_ = json.Unmarshal(p.ChiefComplaints, &resp.ChiefComplaints)
	_ = json.Unmarshal(p.Diagnosis, &resp.Diagnosis)
	_ = json.Unmarshal(p.Medications, &resp.Medications)
	_ = json.Unmarshal(p.Investigations, &resp.Investigations)

	if resp.ChiefComplaints == nil {
		resp.ChiefComplaints = []string{}
	}
	if resp.Diagnosis == nil {
		resp.Diagnosis = []string{}
	}
	if resp.Medications == nil {
		resp.Medications = []types.Medication{}
	}
	if resp.Investigations == nil {
		resp.Investigations = []string{}
	}

	return resp
}
