package response

import "github.com/tms/tyre/internal/domain/entity"

// UnitTyrePosition represents a single tyre position on the unit canvas
type UnitTyrePosition struct {
	Position     string        `json:"position"`
	Label        string        `json:"label"`
	Side         string           `json:"side"`
	Axle         string           `json:"axle"`
	X            float64          `json:"x"`
	Y            float64          `json:"y"`
	Tyre         *TyreResponse    `json:"tyre,omitempty"`
	RTD          *float64         `json:"rtd,omitempty"`
	Status       string           `json:"status"`
}

// UnitTyresResponse is the structured response for GET /units/:id/tyres
type UnitTyresResponse struct {
	Unit            *UnitResponse              `json:"unit"`
	UnitTypeConfig  *UnitTypeConfigResponse    `json:"unit_type_config,omitempty"`
	Positions       []*UnitTyrePosition        `json:"positions"`
	SpareTyres      []*TyreResponse            `json:"spare_tyres"`
	TotalMounted    int                        `json:"total_mounted"`
	TotalSpare      int                        `json:"total_spare"`
}

// ToUnitTyrePosition converts a PositionConfig + optional tyre data into a position response
func ToUnitTyrePosition(pos entity.PositionConfig, tyre *entity.TyreMaster, rtd1, rtd2 *float64) *UnitTyrePosition {
	var rtd *float64
	if rtd1 != nil || rtd2 != nil {
		avg := 0.0
		cnt := 0
		if rtd1 != nil {
			avg += *rtd1
			cnt++
		}
		if rtd2 != nil {
			avg += *rtd2
			cnt++
		}
		if cnt > 0 {
			v := avg / float64(cnt)
			rtd = &v
		}
	}

	status := "empty"
	if tyre != nil {
		status = tyre.Status
	}

	var tyreResp *TyreResponse
	if tyre != nil {
		tyreResp = &TyreResponse{
			ID:              tyre.ID,
			CompanyID:       tyre.CompanyID,
			UnitID:          tyre.UnitID,
			MountedPosition: tyre.MountedPosition,
			Barcode:         tyre.Barcode,
			SerialNumber:    tyre.SerialNumber,
			DOTCode:         tyre.DOTCode,
			Type:            tyre.Type,
			SizeID:          tyre.SizeID,
			BrandID:         tyre.BrandID,
			PatternID:       tyre.PatternID,
			OTD:             tyre.OTD,
			RTD:             tyre.RTD,
			RTD1:            tyre.RTD1,
			RTD2:            tyre.RTD2,
			Lifetime:        tyre.Lifetime,
			PSI:             tyre.PSI,
			Status:          tyre.Status,
			Remarks:         tyre.Remarks,
			Company:         tyre.Company,
			Unit:            tyre.Unit,
			Size:            tyre.Size,
			Brand:           tyre.Brand,
			Pattern:         tyre.Pattern,
		}
	}

	return &UnitTyrePosition{
		Position: pos.Position,
		Label:    pos.Label,
		Side:     pos.Side,
		Axle:     pos.Axle,
		X:        pos.X,
		Y:        pos.Y,
		Tyre:     tyreResp,
		RTD:      rtd,
		Status:   status,
	}
}
