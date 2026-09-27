package response

import (
	"time"

	"github.com/tms/tyre/internal/domain/entity"
)

type TyreHistoryItemResponse struct {
	ID              uint             `json:"id"`
	ReplacementDate time.Time        `json:"replacement_date"`
	HM              float64          `json:"hm"`
	Position        string           `json:"position"`
	Action          string           `json:"action"`
	RTD             *float64         `json:"rtd,omitempty"`
	Remarks         string           `json:"remarks,omitempty"`
	Unit            *entity.Unit     `json:"unit,omitempty"`
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
		HM:              e.HM,
		Position:        e.Position,
		Action:          e.Action,
		RTD:             e.RTD,
		Remarks:         e.Remarks,
		Unit:            e.Unit,
	}
}
