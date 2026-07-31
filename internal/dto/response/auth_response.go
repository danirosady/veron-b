package response

type LoginResponse struct {
	AccessToken  string  `json:"access_token"`
	RefreshToken string  `json:"refresh_token"`
	User         *UserDTO `json:"user"`
}

type UserDTO struct {
	ID        uint  `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	CompanyID *uint  `json:"company_id,omitempty"`
	Status    string `json:"status"`
}
