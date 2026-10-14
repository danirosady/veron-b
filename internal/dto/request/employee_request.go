package request

type CreateEmployeeRequest struct {
	CompanyID  uint   `json:"company_id" binding:"required"`
	Name       string `json:"name" binding:"required,max=255"`
	EmployeeID string `json:"employee_id" binding:"required,max=50"`
	Phone      string `json:"phone" binding:"max=50"`
	Role       string `json:"role" binding:"required,oneof=driver supervisor tyreman"`
}

type UpdateEmployeeRequest struct {
	Name       string `json:"name" binding:"required,max=255"`
	EmployeeID string `json:"employee_id" binding:"required,max=50"`
	Phone      string `json:"phone" binding:"max=50"`
	Status     string `json:"status" binding:"omitempty,oneof=active inactive"`
}
