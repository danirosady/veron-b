package response

import (
	"time"

	"github.com/tms/tyre/internal/domain/entity"
)

type ReplacementResponse struct {
	ID               uint                    `json:"id"`
	CompanyID       uint                    `json:"company_id"`
	ProjectID       uint                    `json:"project_id"`
	UnitID          uint                    `json:"unit_id"`
	DriverID        uint                    `json:"driver_id"`
	Date            time.Time               `json:"date"`
	HMUpdate        float64                 `json:"hm_update"`
	CurrentLifeHM   float64                 `json:"current_life_hm"`
	HMPlan          float64                 `json:"hm_plan"`
	Remarks         string                  `json:"remarks,omitempty"`
	CreatedBy       uint                    `json:"created_by"`
	Details         []ReplacementDetailDTO   `json:"details,omitempty"`
	ReplacementDate time.Time               `json:"replacement_date"`
	Action          string                  `json:"action"`
	Position        string                  `json:"position"`
	OldTyre         *TyreMiniDTO            `json:"old_tyre,omitempty"`
	NewTyre         *TyreMiniDTO            `json:"new_tyre,omitempty"`
	HM              float64                 `json:"hm"`
	Unit            *ReplacementUnitDTO     `json:"unit,omitempty"`
	Operator        *OperatorDTO            `json:"operator,omitempty"`
}

type TyreMiniDTO struct {
	ID           uint   `json:"id"`
	SerialNumber string `json:"serial_number"`
}

type ReplacementUnitDTO struct {
	ID          uint   `json:"id"`
	UnitID      string `json:"unit_id"`
	PlateNumber string `json:"plate_number"`
}

type OperatorDTO struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type ReplacementDetailDTO struct {
	ID                 uint    `json:"id"`
	ReplacementID      uint    `json:"replacement_id"`
	Position           string  `json:"position"`
	Action             string  `json:"action"`
	OldTyreID         *uint   `json:"old_tyre_id,omitempty"`
	OldTyreSerialNum  string  `json:"old_tyre_serial_number,omitempty"`
	OldTyrePattern    string  `json:"old_tyre_pattern,omitempty"`
	OldTyreSize       string  `json:"old_tyre_size,omitempty"`
	OldTyreTread1     *float64 `json:"old_tyre_tread_1,omitempty"`
	OldTyreTread2     *float64 `json:"old_tyre_tread_2,omitempty"`
	OldTyreLifetime   *float64 `json:"old_tyre_lifetime,omitempty"`
	OldTyreStatus     string  `json:"old_tyre_status,omitempty"`
	FailureReasonID   *uint   `json:"failure_reason_id,omitempty"`
	FromUnitID        string  `json:"from_unit_id,omitempty"`
	NewTyreID         *uint   `json:"new_tyre_id,omitempty"`
	NewTyreSerialNum  string  `json:"new_tyre_serial_number,omitempty"`
	NewTyrePattern    string  `json:"new_tyre_pattern,omitempty"`
	NewTyreSize       string  `json:"new_tyre_size,omitempty"`
	NewTyreTread1     *float64 `json:"new_tyre_tread_1,omitempty"`
	NewTyreTread2     *float64 `json:"new_tyre_tread_2,omitempty"`
	NewTyreCurrentLife float64 `json:"new_tyre_current_lifetime"`
	NewTyreStatus     string  `json:"new_tyre_status,omitempty"`
	Remark            string  `json:"remark,omitempty"`
}

func ToReplacementResponse(e *entity.Replacement) *ReplacementResponse {
	if e == nil {
		return nil
	}
	resp := &ReplacementResponse{
		ID:               e.ID,
		CompanyID:        e.CompanyID,
		ProjectID:        e.ProjectID,
		UnitID:           e.UnitID,
		DriverID:         e.DriverID,
		Date:             e.Date,
		HMUpdate:         e.HMUpdate,
		CurrentLifeHM:    e.CurrentLifeHM,
		HMPlan:           e.HMPlan,
		Remarks:          e.Remarks,
		CreatedBy:        e.CreatedBy,
		Details:          make([]ReplacementDetailDTO, 0, len(e.Details)),
		ReplacementDate:  e.Date,
		HM:               e.HMUpdate,
	}
	if e.Unit != nil {
		resp.Unit = &ReplacementUnitDTO{
			ID:          e.Unit.ID,
			UnitID:      e.Unit.UnitID,
			PlateNumber: e.Unit.PlateNumber,
		}
	}
	if e.Creator != nil {
		resp.Operator = &OperatorDTO{
			ID:   e.Creator.ID,
			Name: e.Creator.Name,
		}
	}
	for _, d := range e.Details {
		dto := ReplacementDetailDTO{
			ID:                d.ID,
			ReplacementID:     d.ReplacementID,
			Position:          d.Position,
			Action:            d.Action,
			OldTyreID:        d.OldTyreID,
			OldTyreSerialNum: d.OldTyreSerialNum,
			OldTyrePattern:   d.OldTyrePattern,
			OldTyreSize:      d.OldTyreSize,
			OldTyreTread1:    d.OldTyreTread1,
			OldTyreTread2:    d.OldTyreTread2,
			OldTyreLifetime:  d.OldTyreLifetime,
			OldTyreStatus:    d.OldTyreStatus,
			FailureReasonID:  d.FailureReasonID,
			FromUnitID:       d.FromUnitID,
			NewTyreID:        d.NewTyreID,
			NewTyreSerialNum: d.NewTyreSerialNum,
			NewTyrePattern:   d.NewTyrePattern,
			NewTyreSize:      d.NewTyreSize,
			NewTyreTread1:    d.NewTyreTread1,
			NewTyreTread2:    d.NewTyreTread2,
			NewTyreCurrentLife: d.NewTyreCurrentLife,
			NewTyreStatus:   d.NewTyreStatus,
			Remark:           d.Remark,
		}
		resp.Details = append(resp.Details, dto)

		if resp.Action == "" {
			resp.Action = d.Action
			resp.Position = d.Position
			if d.OldTyreID != nil && d.OldTyreSerialNum != "" {
				resp.OldTyre = &TyreMiniDTO{ID: *d.OldTyreID, SerialNumber: d.OldTyreSerialNum}
			} else if d.OldTyreSerialNum != "" {
				resp.OldTyre = &TyreMiniDTO{SerialNumber: d.OldTyreSerialNum}
			}
			if d.NewTyreID != nil && d.NewTyreSerialNum != "" {
				resp.NewTyre = &TyreMiniDTO{ID: *d.NewTyreID, SerialNumber: d.NewTyreSerialNum}
			} else if d.NewTyreSerialNum != "" {
				resp.NewTyre = &TyreMiniDTO{SerialNumber: d.NewTyreSerialNum}
			}
		}
	}
	return resp
}

func ToReplacementResponses(entities []*entity.Replacement) []*ReplacementResponse {
	responses := make([]*ReplacementResponse, 0, len(entities))
	for _, e := range entities {
		responses = append(responses, ToReplacementResponse(e))
	}
	return responses
}
