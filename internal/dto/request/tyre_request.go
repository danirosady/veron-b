package request

type CreateTyreRequest struct {
	Barcode      string   `json:"barcode" binding:"required,max=100"`
	SerialNumber string   `json:"serial_number" binding:"required,max=100"`
	DOTCode      string   `json:"dot_code" binding:"max=100"`
	Type         string   `json:"type" binding:"max=50"`
	SizeID       uint     `json:"size_id" binding:"required"`
	BrandID      uint     `json:"brand_id" binding:"required"`
	PatternID    uint     `json:"pattern_id" binding:"required"`
	OTD          float64  `json:"otd" binding:"required,gte=0"`
	RTD          float64  `json:"rtd" binding:"omitempty,gte=0"`
	RTD1         *float64 `json:"rtd_1" binding:"omitempty,gte=0"`
	RTD2         *float64 `json:"rtd_2" binding:"omitempty,gte=0"`
	Lifetime     float64  `json:"lifetime" binding:"gte=0"`
	PSI          *float64 `json:"psi" binding:"omitempty,gt=0"`
	Remarks      string   `json:"remarks"`
}

type UpdateTyreRequest struct {
	DOTCode      string   `json:"dot_code" binding:"max=100"`
	Type         string   `json:"type" binding:"max=50"`
	RTD          float64  `json:"rtd" binding:"required,gte=0"`
	RTD1         *float64 `json:"rtd_1" binding:"omitempty,gte=0"`
	RTD2         *float64 `json:"rtd_2" binding:"omitempty,gte=0"`
	PSI          *float64 `json:"psi" binding:"omitempty,gt=0"`
	Remarks      string   `json:"remarks"`
	Status       string   `json:"status" binding:"omitempty,oneof=spare mounted dismounted scrap"`
}
