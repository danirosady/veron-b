package response

import "github.com/tms/tyre/internal/domain/entity"

type TyreResponse struct {
	ID              uint           `json:"id"`
	CompanyID       uint           `json:"company_id"`
	UnitID          *uint          `json:"unit_id,omitempty"`
	MountedPosition *int           `json:"mounted_position,omitempty"`
	Barcode         string         `json:"barcode"`
	SerialNumber    string         `json:"serial_number"`
	DOTCode         string         `json:"dot_code,omitempty"`
	Type            string         `json:"type"`
	SizeID          uint           `json:"size_id"`
	BrandID         uint           `json:"brand_id"`
	PatternID       uint           `json:"pattern_id"`
	OTD             float64        `json:"otd"`
	RTD             float64        `json:"rtd"`
	RTD1            *float64       `json:"rtd_1,omitempty"`
	RTD2            *float64       `json:"rtd_2,omitempty"`
	Lifetime        float64        `json:"lifetime"`
	PSI             *float64       `json:"psi,omitempty"`
	Status          string         `json:"status"`
	Remarks         string         `json:"remarks,omitempty"`
	Company         *entity.Company       `json:"company,omitempty"`
	Unit            *entity.Unit         `json:"unit,omitempty"`
	Size            *entity.MasterSize   `json:"size,omitempty"`
	Brand           *entity.MasterBrand  `json:"brand,omitempty"`
	Pattern         *entity.MasterPattern `json:"pattern,omitempty"`
}

func ToTyreResponse(e *entity.TyreMaster) *TyreResponse {
	if e == nil {
		return nil
	}
	return &TyreResponse{
		ID:              e.ID,
		CompanyID:       e.CompanyID,
		UnitID:          e.UnitID,
		MountedPosition: e.MountedPosition,
		Barcode:         e.Barcode,
		SerialNumber:    e.SerialNumber,
		DOTCode:         e.DOTCode,
		Type:            e.Type,
		SizeID:          e.SizeID,
		BrandID:         e.BrandID,
		PatternID:       e.PatternID,
		OTD:             e.OTD,
		RTD:             e.RTD,
		RTD1:            e.RTD1,
		RTD2:            e.RTD2,
		Lifetime:         e.Lifetime,
		PSI:             e.PSI,
		Status:          e.Status,
		Remarks:         e.Remarks,
		Company:         e.Company,
		Unit:            e.Unit,
		Size:            e.Size,
		Brand:           e.Brand,
		Pattern:         e.Pattern,
	}
}

func ToTyreResponses(entities []*entity.TyreMaster) []*TyreResponse {
	responses := make([]*TyreResponse, 0, len(entities))
	for _, e := range entities {
		responses = append(responses, ToTyreResponse(e))
	}
	return responses
}
