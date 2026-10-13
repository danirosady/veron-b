package response

import (
	"time"

	"github.com/tms/tyre/internal/domain/entity"
)

type TyreHistoryItemResponse struct {
	ID              uint         `json:"id"`
	ReplacementDate time.Time    `json:"replacement_date"`
	HM              float64      `json:"hm"`
	CurrentLifeHM   float64      `json:"current_life_hm"`
	Position        string       `json:"position"`
	Action          string       `json:"action"`
	RTD             *float64    `json:"rtd,omitempty"`
	RTDATMount      *float64    `json:"rtd_at_mount,omitempty"`
	RTDAtDismount   *float64    `json:"rtd_at_dismount,omitempty"`
	LifeHM          *float64    `json:"life_hm,omitempty"`
	Remarks         string       `json:"remarks,omitempty"`
	Unit            *entity.Unit `json:"unit,omitempty"`
}

func ToTyreHistoryResponses(items []*entity.TyreHistoryItem) []*TyreHistoryItemResponse {
	responses := make([]*TyreHistoryItemResponse, 0, len(items))
	for _, item := range items {
		responses = append(responses, ToTyreHistoryItemResponse(item))
	}
	return responses
}

func ToTyreHistoryItemResponse(e *entity.TyreHistoryItem) *TyreHistoryItemResponse {
	if e == nil {
		return nil
	}
	return &TyreHistoryItemResponse{
		ID:              e.ID,
		ReplacementDate: e.ReplacementDate,
		HM:              e.CurrentLifeHM,
		CurrentLifeHM:   e.CurrentLifeHM,
		Position:        e.Position,
		Action:          e.Action,
		RTD:             e.RTD,
		RTDATMount:      e.RTDATMount,
		RTDAtDismount:   e.RTDAtDismount,
		LifeHM:          e.LifeHM,
		Remarks:         e.Remarks,
		Unit:            e.Unit,
	}
}
