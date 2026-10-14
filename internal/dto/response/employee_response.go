package response

import "github.com/tms/tyre/internal/domain/entity"

type EmployeeResponse struct {
	ID         uint              `json:"id"`
	CompanyID  uint              `json:"company_id"`
	Name       string            `json:"name"`
	EmployeeID string           `json:"employee_id"`
	Phone      string           `json:"phone,omitempty"`
	Role       string           `json:"role"`
	Status     string           `json:"status"`
	Company    *CompanyResponse `json:"company,omitempty"`
}

func ToEmployeeResponse(e *entity.Employee) *EmployeeResponse {
	if e == nil {
		return nil
	}
	resp := &EmployeeResponse{
		ID:         e.ID,
		CompanyID:  e.CompanyID,
		Name:       e.Name,
		EmployeeID: e.EmployeeID,
		Phone:      e.Phone,
		Role:       string(e.Role),
		Status:     e.Status,
	}
	if e.Company != nil {
		resp.Company = &CompanyResponse{
			ID:   e.Company.ID,
			Name: e.Company.Name,
		}
	}
	return resp
}

func ToEmployeeResponses(entities []*entity.Employee) []*EmployeeResponse {
	responses := make([]*EmployeeResponse, 0, len(entities))
	for _, e := range entities {
		responses = append(responses, ToEmployeeResponse(e))
	}
	return responses
}

