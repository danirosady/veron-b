package request

type CreateReplacementRequest struct {
	UnitID        uint                         `json:"unit_id" binding:"required"`
	DriverID      uint                         `json:"driver_id" binding:"required"`
	Date          string                       `json:"date" binding:"required"`
	HMUpdate      float64                      `json:"hm_update" binding:"gte=0"`
	CurrentLifeHM float64                      `json:"current_life_hm" binding:"required,gte=0"`
	HMPlan        float64                      `json:"hm_plan" binding:"required,gtfield=CurrentLifeHM"`
	Remarks       string                       `json:"remarks"`
	Details       []ReplacementDetailRequest   `json:"details" binding:"required,min=1,dive"`
}

type ReplacementDetailRequest struct {
	Position      int      `json:"position" binding:"required,min=1,max=20"`
	Action        string   `json:"action" binding:"required,oneof=mount dismount swap"`
	OldTyreID     *uint    `json:"old_tyre_id"`
	NewTyreID     *uint    `json:"new_tyre_id"`
	ReasonID      *uint    `json:"failure_reason_id"`
	OldTyreTread1 *float64 `json:"old_tyre_tread_1"`
	OldTyreTread2 *float64 `json:"old_tyre_tread_2"`
	NewTyreTread1 *float64 `json:"new_tyre_tread_1"`
	NewTyreTread2 *float64 `json:"new_tyre_tread_2"`
	NewTyreStatus string   `json:"new_tyre_status"`
	Remark        string   `json:"remark"`
}

type ListReplacementRequest struct {
	CompanyID uint   `form:"company_id"`
	ProjectID uint   `form:"project_id"`
	UnitID    uint   `form:"unit_id"`
	DateFrom  string `form:"date_from"`
	DateTo    string `form:"date_to"`
	Page      int    `form:"page,default=1"`
	PerPage   int    `form:"per_page,default=20"`
}

type UpdateReplacementRequest struct {
	DriverID      *uint   `json:"driver_id"`
	Date          string  `json:"date"`
	HMUpdate      float64 `json:"hm_update"`
	CurrentLifeHM float64 `json:"current_life_hm"`
	HMPlan        float64 `json:"hm_plan"`
	Remarks       string  `json:"remarks"`
}
