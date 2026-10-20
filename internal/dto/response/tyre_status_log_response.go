package response

import (
	"time"

	"github.com/tms/tyre/internal/domain/entity"
)

type TyreStatusLogResponse struct {
	ID           uint                   `json:"id"`
	TyreID       uint                   `json:"tyre_id"`
	OldStatus    string                 `json:"old_status,omitempty"`
	OldRemarksID *uint                 `json:"old_remarks_id,omitempty"`
	OldRemarks   *MasterRemarkResponse  `json:"old_remarks,omitempty"`
	NewStatus    string                 `json:"new_status"`
	NewRemarksID *uint                 `json:"new_remarks_id,omitempty"`
	NewRemarks   *MasterRemarkResponse  `json:"new_remarks,omitempty"`
	ChangedBy    *uint                 `json:"changed_by,omitempty"`
	Reason       string                 `json:"reason,omitempty"`
	ChangedAt    time.Time              `json:"changed_at"`
}

func ToTyreStatusLogResponse(e *entity.TyreStatusLog) *TyreStatusLogResponse {
	if e == nil {
		return nil
	}
	return &TyreStatusLogResponse{
		ID:           e.ID,
		TyreID:       e.TyreID,
		OldStatus:    e.OldStatus,
		OldRemarksID: e.OldRemarksID,
		OldRemarks:   ToMasterRemarkResponse(e.OldRemarks),
		NewStatus:    e.NewStatus,
		NewRemarksID: e.NewRemarksID,
		NewRemarks:   ToMasterRemarkResponse(e.NewRemarks),
		ChangedBy:    e.ChangedBy,
		Reason:       e.Reason,
		ChangedAt:    e.ChangedAt,
	}
}

func ToTyreStatusLogResponses(entities []*entity.TyreStatusLog) []*TyreStatusLogResponse {
	responses := make([]*TyreStatusLogResponse, 0, len(entities))
	for _, e := range entities {
		responses = append(responses, ToTyreStatusLogResponse(e))
	}
	return responses
}
