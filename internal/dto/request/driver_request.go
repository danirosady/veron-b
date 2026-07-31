package request

type CreateDriverRequest struct {
	CompanyID  uint   `json:"company_id" binding:"required"`
	Name       string `json:"name" binding:"required,max=255"`
	EmployeeID string `json:"employee_id" binding:"required,max=50"`
	Phone      string `json:"phone" binding:"max=50"`
}

type UpdateDriverRequest struct {
	Name       string `json:"name" binding:"required,max=255"`
	EmployeeID string `json:"employee_id" binding:"required,max=50"`
	Phone      string `json:"phone" binding:"max=50"`
	Status     string `json:"status" binding:"omitempty,oneof=active inactive"`
}
