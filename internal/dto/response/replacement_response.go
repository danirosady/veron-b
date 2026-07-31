package response

import (
	"time"

	"github.com/tms/tyre/internal/domain/entity"
)

type ReplacementResponse struct {
	ID            uint                   `json:"id"`
	CompanyID     uint                   `json:"company_id"`
	ProjectID     uint                   `json:"project_id"`
	UnitID        uint                   `json:"unit_id"`
	DriverID      uint                   `json:"driver_id"`
	Date          time.Time              `json:"date"`
	HMUpdate      float64                `json:"hm_update"`
	CurrentLifeHM float64                `json:"current_life_hm"`
	HMPlan        float64                `json:"hm_plan"`
	Remarks       string                 `json:"remarks,omitempty"`
	CreatedBy     uint                   `json:"created_by"`
	Details       []ReplacementDetailDTO `json:"details,omitempty"`
}

type ReplacementDetailDTO struct {
	ID                  uint    `json:"id"`
	ReplacementID       uint    `json:"replacement_id"`
	Position            int     `json:"position"`
	Action              string  `json:"action"`
	OldTyreID           *uint   `json:"old_tyre_id,omitempty"`
	OldTyreSerialNum    string  `json:"old_tyre_serial_number,omitempty"`
	OldTyrePattern      string  `json:"old_tyre_pattern,omitempty"`
	OldTyreSize         string  `json:"old_tyre_size,omitempty"`
	OldTyreTread1       *float64 `json:"old_tyre_tread_1,omitempty"`
	OldTyreTread2       *float64 `json:"old_tyre_tread_2,omitempty"`
	OldTyreLifetime     *float64 `json:"old_tyre_lifetime,omitempty"`
	OldTyreStatus       string  `json:"old_tyre_status,omitempty"`
	FailureReasonID     *uint   `json:"failure_reason_id,omitempty"`
	FromUnitID          string  `json:"from_unit_id,omitempty"`
	NewTyreID           *uint   `json:"new_tyre_id,omitempty"`
	NewTyreSerialNum    string  `json:"new_tyre_serial_number,omitempty"`
	NewTyrePattern      string  `json:"new_tyre_pattern,omitempty"`
	NewTyreSize         string  `json:"new_tyre_size,omitempty"`
	NewTyreTread1       *float64 `json:"new_tyre_tread_1,omitempty"`
	NewTyreTread2       *float64 `json:"new_tyre_tread_2,omitempty"`
	NewTyreCurrentLife  float64 `json:"new_tyre_current_lifetime"`
	NewTyreStatus       string  `json:"new_tyre_status,omitempty"`
	Remark              string  `json:"remark,omitempty"`
}

func ToReplacementResponse(e *entity.Replacement) *ReplacementResponse {
	if e == nil {
		return nil
	}
	resp := &ReplacementResponse{
		ID:            e.ID,
		CompanyID:     e.CompanyID,
		ProjectID:     e.ProjectID,
		UnitID:        e.UnitID,
		DriverID:      e.DriverID,
		Date:          e.Date,
		HMUpdate:      e.HMUpdate,
		CurrentLifeHM: e.CurrentLifeHM,
		HMPlan:        e.HMPlan,
		Remarks:       e.Remarks,
		CreatedBy:     e.CreatedBy,
		Details:       make([]ReplacementDetailDTO, 0, len(e.Details)),
	}
	for _, d := range e.Details {
		resp.Details = append(resp.Details, ReplacementDetailDTO{
			ID:                 d.ID,
			ReplacementID:      d.ReplacementID,
			Position:           d.Position,
			Action:             d.Action,
			OldTyreID:          d.OldTyreID,
			OldTyreSerialNum:   d.OldTyreSerialNum,
			OldTyrePattern:     d.OldTyrePattern,
			OldTyreSize:        d.OldTyreSize,
			OldTyreTread1:      d.OldTyreTread1,
			OldTyreTread2:      d.OldTyreTread2,
			OldTyreLifetime:    d.OldTyreLifetime,
			OldTyreStatus:      d.OldTyreStatus,
			FailureReasonID:   d.FailureReasonID,
			FromUnitID:        d.FromUnitID,
			NewTyreID:          d.NewTyreID,
			NewTyreSerialNum:   d.NewTyreSerialNum,
			NewTyrePattern:     d.NewTyrePattern,
			NewTyreSize:        d.NewTyreSize,
			NewTyreTread1:      d.NewTyreTread1,
			NewTyreTread2:      d.NewTyreTread2,
			NewTyreCurrentLife: d.NewTyreCurrentLife,
			NewTyreStatus:      d.NewTyreStatus,
			Remark:             d.Remark,
		})
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
