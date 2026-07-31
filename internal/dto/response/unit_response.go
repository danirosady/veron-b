package response

import "github.com/tms/tyre/internal/domain/entity"

type UnitResponse struct {
	ID              uint    `json:"id"`
	CompanyID       uint    `json:"company_id"`
	ProjectID       uint    `json:"project_id"`
	UnitID          string  `json:"unit_id"`
	UnitModel       string  `json:"unit_model"`
	PlateNumber     string  `json:"plate_number,omitempty"`
	TyreSizeDefault string  `json:"tyre_size_default"`
	UnitType        string  `json:"unit_type"`
	MaxPosition     int     `json:"max_position"`
	CurrentHM       float64 `json:"current_hm"`
	Status          string  `json:"status"`
}

func ToUnitResponse(e *entity.Unit) *UnitResponse {
	if e == nil {
		return nil
	}
	return &UnitResponse{
		ID:              e.ID,
		CompanyID:       e.CompanyID,
		ProjectID:       e.ProjectID,
		UnitID:          e.UnitID,
		UnitModel:       e.UnitModel,
		PlateNumber:     e.PlateNumber,
		TyreSizeDefault: e.TyreSizeDefault,
		UnitType:        e.UnitType,
		MaxPosition:     e.MaxPosition,
		CurrentHM:       e.CurrentHM,
		Status:          e.Status,
	}
}

func ToUnitResponses(entities []*entity.Unit) []*UnitResponse {
	responses := make([]*UnitResponse, 0, len(entities))
	for _, e := range entities {
		responses = append(responses, ToUnitResponse(e))
	}
	return responses
}
