package response

import "github.com/tms/tyre/internal/domain/entity"

type UserResponse struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	CompanyID *uint  `json:"company_id,omitempty"`
	Status    string `json:"status"`
}

func ToUserResponse(e *entity.User) *UserResponse {
	if e == nil {
		return nil
	}
	return &UserResponse{
		ID:        e.ID,
		Name:      e.Name,
		Email:     e.Email,
		Role:      e.Role,
		CompanyID: e.CompanyID,
		Status:    e.Status,
	}
}

func ToUserResponses(entities []*entity.User) []*UserResponse {
	responses := make([]*UserResponse, 0, len(entities))
	for _, e := range entities {
		responses = append(responses, ToUserResponse(e))
	}
	return responses
}
