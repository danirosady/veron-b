package response

import "github.com/tms/tyre/internal/domain/entity"

type CompanyResponse struct {
	ID            uint   `json:"id"`
	Name          string `json:"name"`
	Address       string `json:"address,omitempty"`
	ContactPerson string `json:"contact_person,omitempty"`
	Phone         string `json:"phone,omitempty"`
	Email         string `json:"email,omitempty"`
	Status        string `json:"status"`
}

func ToCompanyResponse(e *entity.Company) *CompanyResponse {
	if e == nil {
		return nil
	}
	return &CompanyResponse{
		ID:            e.ID,
		Name:          e.Name,
		Address:       e.Address,
		ContactPerson: e.ContactPerson,
		Phone:         e.Phone,
		Email:         e.Email,
		Status:        e.Status,
	}
}

func ToCompanyResponses(entities []*entity.Company) []*CompanyResponse {
	responses := make([]*CompanyResponse, 0, len(entities))
	for _, e := range entities {
		responses = append(responses, ToCompanyResponse(e))
	}
	return responses
}
