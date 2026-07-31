package response

import "github.com/tms/tyre/internal/domain/entity"

type DriverResponse struct {
	ID         uint   `json:"id"`
	CompanyID  uint   `json:"company_id"`
	Name       string `json:"name"`
	EmployeeID string `json:"employee_id"`
	Phone      string `json:"phone,omitempty"`
	Status     string `json:"status"`
}

func ToDriverResponse(e *entity.Driver) *DriverResponse {
	if e == nil {
		return nil
	}
	return &DriverResponse{
		ID:         e.ID,
		CompanyID:  e.CompanyID,
		Name:       e.Name,
		EmployeeID: e.EmployeeID,
		Phone:      e.Phone,
		Status:     e.Status,
	}
}

func ToDriverResponses(entities []*entity.Driver) []*DriverResponse {
	responses := make([]*DriverResponse, 0, len(entities))
	for _, e := range entities {
		responses = append(responses, ToDriverResponse(e))
	}
	return responses
}
