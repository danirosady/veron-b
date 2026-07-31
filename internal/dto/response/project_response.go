package response

import (
	"time"

	"github.com/tms/tyre/internal/domain/entity"
)

type ProjectResponse struct {
	ID        uint       `json:"id"`
	CompanyID uint       `json:"company_id"`
	Name      string     `json:"name"`
	Location  string     `json:"location,omitempty"`
	StartDate *time.Time `json:"start_date,omitempty"`
	EndDate   *time.Time `json:"end_date,omitempty"`
	Status    string     `json:"status"`
}

func ToProjectResponse(e *entity.Project) *ProjectResponse {
	if e == nil {
		return nil
	}
	return &ProjectResponse{
		ID:        e.ID,
		CompanyID: e.CompanyID,
		Name:      e.Name,
		Location:  e.Location,
		StartDate: e.StartDate,
		EndDate:   e.EndDate,
		Status:    e.Status,
	}
}

func ToProjectResponses(entities []*entity.Project) []*ProjectResponse {
	responses := make([]*ProjectResponse, 0, len(entities))
	for _, e := range entities {
		responses = append(responses, ToProjectResponse(e))
	}
	return responses
}
