package request

type CreateProjectRequest struct {
	CompanyID uint   `json:"company_id" binding:"required"`
	Name      string `json:"name" binding:"required,max=255"`
	Location  string `json:"location"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

type UpdateProjectRequest struct {
	Name      string `json:"name" binding:"required,max=255"`
	Location  string `json:"location"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	Status    string `json:"status" binding:"omitempty,oneof=active inactive"`
}
