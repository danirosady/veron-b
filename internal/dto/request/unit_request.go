package request

type CreateUnitRequest struct {
	CompanyID       uint    `json:"company_id" binding:"required"`
	ProjectID       uint    `json:"project_id" binding:"required"`
	UnitID          string  `json:"unit_id" binding:"required,max=50"`
	UnitModel       string  `json:"unit_model" binding:"required,max=255"`
	PlateNumber     string  `json:"plate_number" binding:"max=50"`
	TyreSizeDefault string  `json:"tyre_size_default" binding:"required,max=50"`
	UnitType        string  `json:"unit_type" binding:"required,max=50"`
	MaxPosition     int     `json:"max_position" binding:"required,min=1,max=20"`
}

type UpdateUnitRequest struct {
	ProjectID       uint    `json:"project_id" binding:"required"`
	UnitModel       string  `json:"unit_model" binding:"required,max=255"`
	PlateNumber     string  `json:"plate_number" binding:"max=50"`
	TyreSizeDefault string  `json:"tyre_size_default" binding:"required,max=50"`
	UnitType        string  `json:"unit_type" binding:"required,max=50"`
	MaxPosition     int     `json:"max_position" binding:"required,min=1,max=20"`
	Status          string  `json:"status" binding:"omitempty,oneof=active inactive"`
}

type UpdateUnitHMRequest struct {
	CurrentHM float64 `json:"current_hm" binding:"required,gte=0"`
}
