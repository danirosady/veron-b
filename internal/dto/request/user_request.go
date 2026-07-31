package request

type CreateUserRequest struct {
	Name      string `json:"name" binding:"required,max=255"`
	Email     string `json:"email" binding:"required,email,max=255"`
	Password  string `json:"password" binding:"required,min=8,max=72"`
	Role      string `json:"role" binding:"required,oneof=superadmin admin_company"`
	CompanyID *uint  `json:"company_id"`
}

type UpdateUserRequest struct {
	Name      string `json:"name" binding:"required,max=255"`
	Email     string `json:"email" binding:"required,email,max=255"`
	Role      string `json:"role" binding:"required,oneof=superadmin admin_company"`
	CompanyID *uint  `json:"company_id"`
	Status    string `json:"status" binding:"omitempty,oneof=active inactive"`
}
