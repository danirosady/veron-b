package usecase

import (
	"context"
	"time"

	"github.com/tms/tyre/internal/domain/repository"
	"gorm.io/gorm"
)

// ReportUseCase handles business logic for report generation.
type ReportUseCase struct {
	db              *gorm.DB
	tyreRepo        repository.TyreRepository
	replacementRepo repository.ReplacementRepository
	unitRepo        repository.UnitRepository
	companyRepo     repository.CompanyRepository
}

// NewReportUseCase creates a new ReportUseCase instance.
func NewReportUseCase(
	db *gorm.DB,
	tyreRepo repository.TyreRepository,
	replacementRepo repository.ReplacementRepository,
	unitRepo repository.UnitRepository,
	companyRepo repository.CompanyRepository,
) *ReportUseCase {
	return &ReportUseCase{
		db:              db,
		tyreRepo:        tyreRepo,
		replacementRepo: replacementRepo,
		unitRepo:        unitRepo,
		companyRepo:     companyRepo,
	}
}

// ReplacementReportRow represents a row in the replacement report.
type ReplacementReportRow struct {
	Date         string  `json:"replacement_date"`
	Unit         UnitInfo `json:"unit"`
	Position     string  `json:"position"`
	Action       string  `json:"action"`
	OldTyre      TyreInfo `json:"old_tyre"`
	NewTyre      TyreInfo `json:"new_tyre"`
	TyreSpec     TyreSpecInfo `json:"tyre_spec"`
	RTDBefore    float64 `json:"rtd_before"`
	RTDAfter     float64 `json:"rtd_after"`
	HM           float64 `json:"hm"`
	DriverName   string  `json:"driver_name"`
	OperatorName string  `json:"operator_name"`
	Remarks      string  `json:"remarks"`
	CompanyName  string  `json:"company_name"`
}

type UnitInfo struct {
	PlateNumber string `json:"plate_number"`
	Code        string `json:"code"`
}

type TyreInfo struct {
	SerialNumber string `json:"serial_number"`
}

type TyreSpecInfo struct {
	BrandName string `json:"brand_name"`
	SizeName  string `json:"size_name"`
	Pattern   string `json:"pattern"`
}

// InventoryReportRow represents a row in the inventory report.
type InventoryReportRow struct {
	Barcode      string       `json:"barcode"`
	SerialNumber string       `json:"serial_number"`
	Brand        BrandInfo    `json:"brand"`
	Size         SizeInfo     `json:"size"`
	Pattern      string       `json:"pattern"`
	Status       string       `json:"status"`
	OTD          float64      `json:"otd"`
	RTD          float64      `json:"rtd"`
	PercentWorn  float64      `json:"percent_worn"`
	Unit         UnitInfo     `json:"unit"`
	CompanyName  string       `json:"company_name"`
	LastPosition string       `json:"last_position"`
	TotalDepth   float64      `json:"total_depth"`
	CurrentDepth float64      `json:"current_depth"`
}

type BrandInfo struct {
	Name string `json:"name"`
}

type SizeInfo struct {
	Name string `json:"name"`
}

// ScheduleReportRow represents a row in the schedule report.
type ScheduleReportRow struct {
	SerialNumber       string
	Barcode            string
	Brand              string
	Size               string
	Pattern            string
	UnitCode           string
	PlateNumber        string
	Position           string
	CurrentRTD         float64
	EstimatedEnd       float64
	DaysRemaining      float64
	RecommendedAction  string
}

// TyreSummaryReportRow represents a row in the tyre summary report (Cost/HM per tyre).
type TyreSummaryReportRow struct {
	SerialNumber    string  `json:"serial_number"`
	Barcode         string  `json:"barcode"`
	Brand           string  `json:"brand"`
	Size            string  `json:"size"`
	Pattern         string  `json:"pattern"`
	OTD             float64 `json:"otd"`
	RTD             float64 `json:"rtd"`
	Cost            float64 `json:"cost"`
	TotalLifeHM     float64 `json:"total_life_hm"`
	CostPerHM       float64 `json:"cost_per_hm"`
	TotalTreadUsed  float64 `json:"total_tread_used"`
	PotentialHM     float64 `json:"potential_hm"`
	Status          string  `json:"status"`
	CurrentUnit     string  `json:"current_unit"`
	CurrentPosition string  `json:"current_position"`
}

// PotentialHMReportRow represents a row in the Potential HM report.
type PotentialHMReportRow struct {
	SerialNumber      string  `json:"serial_number"`
	Barcode          string  `json:"barcode"`
	Brand            string  `json:"brand"`
	Size             string  `json:"size"`
	Pattern          string  `json:"pattern"`
	UnitCode         string  `json:"unit_code"`
	PlateNumber      string  `json:"plate_number"`
	Position         string  `json:"position"`
	OTD              float64 `json:"otd"`
	RTDATMount       float64 `json:"rtd_at_mount"`
	RTDCurrent       float64 `json:"rtd_current"`
	RTDMinimum      float64 `json:"rtd_minimum"`
	TotalLifeHM      float64 `json:"total_life_hm"`
	WearRate         float64 `json:"wear_rate"`
	RemainingRTD     float64 `json:"remaining_rtd"`
	PotentialHM      float64 `json:"potential_hm"`
	TotalPotentialHM float64 `json:"total_potential_hm"`
	Status           string  `json:"status"`
}

// GetReplacementReport generates a replacement report based on filters.
func (uc *ReportUseCase) GetReplacementReport(ctx context.Context, filters map[string]interface{}) ([]*ReplacementReportRow, error) {
	type rawRow struct {
		Date           string
		UnitCode       string
		PlateNumber    string
		CompanyName    string
		Position       string
		Action         string
		OldTyreSN      string
		NewTyreSN      string
		Brand          string
		Size           string
		Pattern        string
		OldRTD         float64
		NewRTD         float64
		HM             float64
		DriverName     string
		OperatorName   string
		Remarks        string
	}

	query := uc.db.Table("replacement_details rd").
		Select(`
			r.date,
			u.unit_id AS unit_code,
			u.plate_number,
			c.name AS company_name,
			rd.position,
			rd.action,
			rd.old_tyre_serial_number AS old_tyre_sn,
			rd.new_tyre_serial_number AS new_tyre_sn,
			COALESCE(b.name, '') AS brand,
			COALESCE(s.name, '') AS size,
			COALESCE(p.name, '') AS pattern,
			COALESCE(t_old.rtd, 0) AS old_rtd,
			COALESCE(t_new.rtd, 0) AS new_rtd,
			r.hm_update AS hm,
			COALESCE(d.name, '') AS driver_name,
			COALESCE(usr.name, '') AS operator_name,
			rd.remark AS remarks
		`).
		Joins("JOIN replacements r ON rd.replacement_id = r.id").
		Joins("LEFT JOIN units u ON r.unit_id = u.id").
		Joins("LEFT JOIN companies c ON r.company_id = c.id").
		Joins("LEFT JOIN employees d ON r.driver_id = d.id AND d.role = 'driver'").
		Joins("LEFT JOIN users usr ON r.created_by = usr.id").
		Joins("LEFT JOIN tyre_master t_old ON rd.old_tyre_id = t_old.id").
		Joins("LEFT JOIN tyre_master t_new ON rd.new_tyre_id = t_new.id").
		Joins("LEFT JOIN master_brands b ON COALESCE(t_old.brand_id, t_new.brand_id) = b.id").
		Joins("LEFT JOIN master_sizes s ON COALESCE(t_old.size_id, t_new.size_id) = s.id").
		Joins("LEFT JOIN master_patterns p ON COALESCE(t_old.pattern_id, t_new.pattern_id) = p.id")

	// Apply filters
	if companyID, ok := filters["company_id"].(uint); ok && companyID > 0 {
		query = query.Where("r.company_id = ?", companyID)
	}
	if projectID, ok := filters["project_id"].(uint); ok && projectID > 0 {
		query = query.Where("r.project_id = ?", projectID)
	}
	if unitID, ok := filters["unit_id"].(uint); ok && unitID > 0 {
		query = query.Where("r.unit_id = ?", unitID)
	}
	if dateFrom, ok := filters["date_from"].(string); ok && dateFrom != "" {
		query = query.Where("r.date >= ?", dateFrom)
	}
	if dateTo, ok := filters["date_to"].(string); ok && dateTo != "" {
		query = query.Where("r.date <= ?", dateTo)
	}
	if action, ok := filters["action"].(string); ok && action != "" {
		query = query.Where("rd.action = ?", action)
	}

	query = query.Order("r.date DESC, rd.position ASC")

	var rawResults []rawRow
	if err := query.Scan(&rawResults).Error; err != nil {
		return nil, err
	}

	results := make([]*ReplacementReportRow, len(rawResults))
	for i, r := range rawResults {
		results[i] = &ReplacementReportRow{
			Date:       r.Date,
			Unit:       UnitInfo{PlateNumber: r.PlateNumber, Code: r.UnitCode},
			Position:   r.Position,
			Action:     r.Action,
			OldTyre:    TyreInfo{SerialNumber: r.OldTyreSN},
			NewTyre:    TyreInfo{SerialNumber: r.NewTyreSN},
			TyreSpec:   TyreSpecInfo{BrandName: r.Brand, SizeName: r.Size, Pattern: r.Pattern},
			RTDBefore:  r.OldRTD,
			RTDAfter:   r.NewRTD,
			HM:         r.HM,
			DriverName: r.DriverName,
			OperatorName: r.OperatorName,
			Remarks:    r.Remarks,
			CompanyName: r.CompanyName,
		}
	}

	return results, nil
}

// GetInventoryReport generates an inventory report based on filters.
func (uc *ReportUseCase) GetInventoryReport(ctx context.Context, filters map[string]interface{}) ([]*InventoryReportRow, error) {
	type rawRow struct {
		Barcode      string
		SerialNumber string
		Brand        string
		Size         string
		Pattern      string
		Status       string
		OTD          float64
		RTD          float64
		PercentWorn  float64
		UnitCode     string
		PlateNumber  string
		CompanyName  string
		LastPosition string
		TotalDepth   float64
		CurrentDepth float64
	}

	query := uc.db.Table("tyre_master t").
		Select(`
			t.barcode,
			t.serial_number,
			COALESCE(b.name, '') AS brand,
			COALESCE(s.name, '') AS size,
			COALESCE(p.name, '') AS pattern,
			t.status,
			t.otd,
			t.rtd,
			CASE WHEN t.otd > 0 THEN ROUND((t.otd - t.rtd) / t.otd * 100, 2) ELSE 0 END AS percent_worn,
			COALESCE(u.unit_id, '') AS unit_code,
			COALESCE(u.plate_number, '') AS plate_number,
			COALESCE(c.name, '') AS company_name,
			CASE WHEN t.mounted_position IS NOT NULL THEN CAST(t.mounted_position AS TEXT) ELSE '' END AS last_position,
			t.otd AS total_depth,
			t.rtd AS current_depth
		`).
		Joins("LEFT JOIN units u ON t.unit_id = u.id").
		Joins("LEFT JOIN companies c ON t.company_id = c.id").
		Joins("LEFT JOIN master_brands b ON t.brand_id = b.id").
		Joins("LEFT JOIN master_sizes s ON t.size_id = s.id").
		Joins("LEFT JOIN master_patterns p ON t.pattern_id = p.id")

	// Apply filters
	if companyID, ok := filters["company_id"].(uint); ok && companyID > 0 {
		query = query.Where("t.company_id = ?", companyID)
	}
	if status, ok := filters["status"].(string); ok && status != "" {
		query = query.Where("t.status = ?", status)
	}
	if brandID, ok := filters["brand_id"].(uint); ok && brandID > 0 {
		query = query.Where("t.brand_id = ?", brandID)
	}
	if sizeID, ok := filters["size_id"].(uint); ok && sizeID > 0 {
		query = query.Where("t.size_id = ?", sizeID)
	}

	query = query.Order("c.name, t.status, t.serial_number")

	var rawResults []rawRow
	if err := query.Scan(&rawResults).Error; err != nil {
		return nil, err
	}

	results := make([]*InventoryReportRow, len(rawResults))
	for i, r := range rawResults {
		results[i] = &InventoryReportRow{
			Barcode:      r.Barcode,
			SerialNumber: r.SerialNumber,
			Brand:        BrandInfo{Name: r.Brand},
			Size:         SizeInfo{Name: r.Size},
			Pattern:      r.Pattern,
			Status:       r.Status,
			OTD:          r.OTD,
			RTD:          r.RTD,
			PercentWorn:  r.PercentWorn,
			Unit:         UnitInfo{Code: r.UnitCode, PlateNumber: r.PlateNumber},
			CompanyName:  r.CompanyName,
			LastPosition: r.LastPosition,
			TotalDepth:   r.TotalDepth,
			CurrentDepth: r.CurrentDepth,
		}
	}

	return results, nil
}

// GetScheduleReport generates a tyre replacement schedule report based on RTD analysis.
func (uc *ReportUseCase) GetScheduleReport(ctx context.Context, filters map[string]interface{}) ([]*ScheduleReportRow, error) {
	// Get mounted tyres with their mount dates for wear rate calculation
	type TyreMountInfo struct {
		TyreID        uint
		UnitID        string
		PlateNumber   string
		MountDate     *time.Time
		BrandName     string
		SizeName      string
		PatternName   string
		OTD           float64
		RTD           float64
		MountedPos    *string
	}

	var mountedTyres []TyreMountInfo
	query := uc.db.Table("tyre_master t").
		Select(`
			t.id AS tyre_id,
			COALESCE(u.unit_id, '') AS unit_id,
			COALESCE(u.plate_number, '') AS plate_number,
			NULL::timestamp AS mount_date,
			COALESCE(b.name, '') AS brand_name,
			COALESCE(s.name, '') AS size_name,
			COALESCE(p.name, '') AS pattern_name,
			t.otd,
			t.rtd,
			t.mounted_position
		`).
		Joins("LEFT JOIN units u ON t.unit_id = u.id").
		Joins("LEFT JOIN master_brands b ON t.brand_id = b.id").
		Joins("LEFT JOIN master_sizes s ON t.size_id = s.id").
		Joins("LEFT JOIN master_patterns p ON t.pattern_id = p.id").
		Where("t.status = ?", "mounted")

	if companyID, ok := filters["company_id"].(uint); ok && companyID > 0 {
		query = query.Where("t.company_id = ?", companyID)
	}

	// Get last mount date for each tyre from replacements
	now := time.Now()

	if err := query.Scan(&mountedTyres).Error; err != nil {
		return nil, err
	}

	results := make([]*ScheduleReportRow, 0, len(mountedTyres))

	for _, tyre := range mountedTyres {
		// Get the last mount date for this tyre from replacement_details
		var lastMountDate *time.Time
		uc.db.Table("replacement_details rd").
			Select("r.date").
			Joins("JOIN replacements r ON rd.replacement_id = r.id").
			Where("rd.new_tyre_id = ? AND rd.action = 'mount'", tyre.TyreID).
			Order("r.date DESC").
			Limit(1).
			Scan(&lastMountDate)

		var daysRemaining float64
		var recommendedAction string
		var estimatedEnd float64 = 0

		if tyre.OTD > 0 && tyre.RTD >= 0 {
			if lastMountDate != nil {
				monthsSinceMount := now.Sub(*lastMountDate).Hours() / 24.0 / 30.0
				if monthsSinceMount > 0 {
					wearPerMonth := (tyre.OTD - tyre.RTD) / monthsSinceMount
					if wearPerMonth > 0 {
						remainingMonths := tyre.RTD / wearPerMonth
						daysRemaining = remainingMonths * 30
						estimatedEnd = tyre.RTD - (wearPerMonth * remainingMonths)
					}
				}
			}

			// Determine recommended action based on RTD
			if tyre.RTD < 2.0 {
				recommendedAction = "URGENT: Replace immediately - RTD below minimum"
			} else if tyre.RTD < 5.0 {
				recommendedAction = "Schedule replacement soon"
			} else if tyre.RTD < tyre.OTD*0.3 {
				recommendedAction = "Monitor closely - approaching replacement"
			} else {
				recommendedAction = "Continue operation"
			}
		}

		position := ""
		if tyre.MountedPos != nil {
			position = *tyre.MountedPos
		}

		row := &ScheduleReportRow{
			SerialNumber:      "",
			Barcode:            "",
			Brand:              tyre.BrandName,
			Size:               tyre.SizeName,
			Pattern:            tyre.PatternName,
			UnitCode:           tyre.UnitID,
			PlateNumber:        tyre.PlateNumber,
			Position:           position,
			CurrentRTD:         tyre.RTD,
			EstimatedEnd:       estimatedEnd,
			DaysRemaining:      daysRemaining,
			RecommendedAction:  recommendedAction,
		}

		// Get tyre details
		var tyreDetail struct {
			SerialNumber string
			Barcode      string
		}
		uc.db.Table("tyre_master").Select("serial_number, barcode").Where("id = ?", tyre.TyreID).Scan(&tyreDetail)
		row.SerialNumber = tyreDetail.SerialNumber
		row.Barcode = tyreDetail.Barcode

		results = append(results, row)
	}

	// Sort by days remaining (ascending, urgent first)
	// Simple insertion sort for small datasets
	for i := 1; i < len(results); i++ {
		for j := i; j > 0 && results[j].DaysRemaining < results[j-1].DaysRemaining; j-- {
			results[j], results[j-1] = results[j-1], results[j]
		}
	}

	return results, nil
}

// GetTyreSummaryReport generates a tyre summary report (Cost/HM per tyre).
// Total Life HM = SUM(old_tyre_lifetime from all dismount records) + (current_unit_hm - new_tyre_hm from last mount)
func (uc *ReportUseCase) GetTyreSummaryReport(ctx context.Context, filters map[string]interface{}) ([]*TyreSummaryReportRow, error) {
	type resultRow struct {
		TyreID           uint
		SerialNumber     string
		Barcode          string
		Brand            string
		Size             string
		Pattern          string
		OTD              float64
		RTD              float64
		Cost             float64
		Status           string
		TotalLifeHM      float64
		CurrentUnit      string
		CurrentPosition  string
		WearRate         float64
		NewTyreHM        float64
	}

	// CTE: sum old_tyre_lifetime per tyre from all dismount records
	// CTE: last mount info (new_tyre_hm, new_tyre_tread_1) per tyre
	withSQL := `
	WITH tyre_dismount_life AS (
		SELECT old_tyre_id AS tyre_id, SUM(old_tyre_lifetime) AS dismount_life
		FROM replacement_details
		WHERE old_tyre_id IS NOT NULL AND action = 'dismount' AND old_tyre_lifetime IS NOT NULL
		GROUP BY old_tyre_id
	),
	tyre_wear AS (
		SELECT
			rd.old_tyre_id AS tyre_id,
			CASE WHEN SUM(rd.old_tyre_lifetime) > 0
				THEN SUM(rd.new_tyre_tread_1 - rd.old_tyre_tread_1) / SUM(rd.old_tyre_lifetime)
				ELSE 0 END AS wear_rate
		FROM replacement_details rd
		WHERE rd.old_tyre_id IS NOT NULL AND rd.action = 'dismount'
			AND rd.old_tyre_lifetime IS NOT NULL AND rd.old_tyre_lifetime > 0
		GROUP BY rd.old_tyre_id
	),
	last_mount AS (
		SELECT DISTINCT ON (rd.new_tyre_id)
			rd.new_tyre_id AS tyre_id,
			rd.new_tyre_hm,
			rd.new_tyre_tread_1 AS rtd_at_mount
		FROM replacement_details rd
		JOIN replacements r ON rd.replacement_id = r.id
		WHERE rd.new_tyre_id IS NOT NULL AND rd.action = 'mount' AND rd.new_tyre_hm IS NOT NULL
		ORDER BY rd.new_tyre_id, r.id DESC
	)
	SELECT
		t.id AS tyre_id,
		t.serial_number,
		t.barcode,
		COALESCE(b.name, '') AS brand,
		COALESCE(s.name, '') AS size,
		COALESCE(p.name, '') AS pattern,
		t.otd,
		t.rtd,
		COALESCE(t.cost, 0) AS cost,
		t.status,
		COALESCE(u.unit_id, '') AS current_unit,
		COALESCE(t.mounted_position, '') AS current_position,
		COALESCE(tdl.dismount_life, 0) AS dismount_life,
		COALESCE(tw.wear_rate, 0) AS wear_rate,
		lm.new_tyre_hm,
		lm.rtd_at_mount
	FROM tyre_master t
	LEFT JOIN master_brands b ON t.brand_id = b.id
	LEFT JOIN master_sizes s ON t.size_id = s.id
	LEFT JOIN master_patterns p ON t.pattern_id = p.id
	LEFT JOIN units u ON t.unit_id = u.id
	LEFT JOIN tyre_dismount_life tdl ON tdl.tyre_id = t.id
	LEFT JOIN tyre_wear tw ON tw.tyre_id = t.id
	LEFT JOIN last_mount lm ON lm.tyre_id = t.id
	WHERE (? = 0 OR t.company_id = ?)
	`

	var companyID uint = 0
	if cid, ok := filters["company_id"].(uint); ok {
		companyID = cid
	}

	var rows []resultRow
	if err := uc.db.Raw(withSQL, companyID, companyID).Scan(&rows).Error; err != nil {
		return nil, err
	}

	results := make([]*TyreSummaryReportRow, 0, len(rows))
	for _, r := range rows {
		// Total life HM = SUM dismount life + current mounted lifetime
		// Current mounted lifetime = unit.current_hm - new_tyre_hm (if mounted)
		totalLifeHM := r.TotalLifeHM
		if r.Status == "mounted" && r.NewTyreHM > 0 {
			var currentHM float64
			uc.db.Table("units").Select("current_hm").Where("unit_id = ?", r.CurrentUnit).Scan(&currentHM)
			currentLifeHM := currentHM - r.NewTyreHM
			if currentLifeHM > 0 {
				totalLifeHM += currentLifeHM
			}
		}

		// Cost per HM
		costPerHM := 0.0
		if totalLifeHM > 0 && r.Cost > 0 {
			costPerHM = r.Cost / totalLifeHM
		}

		// Potential HM: (rtd - rtd_minimum) / wear_rate
		potentialHM := 0.0
		if r.WearRate > 0 {
			rtdMin := 10.0
			var rtdMinFromDB float64
			uc.db.Table("master_sizes").
				Select("COALESCE(rtd_minimum, 10.0)").
				Where("id = (SELECT size_id FROM tyre_master WHERE id = ?)", r.TyreID).
				Scan(&rtdMinFromDB)
			if rtdMinFromDB > 0 {
				rtdMin = rtdMinFromDB
			}
			remainingRTD := r.RTD - rtdMin
			if remainingRTD > 0 {
				potentialHM = remainingRTD / r.WearRate
			}
		}

		results = append(results, &TyreSummaryReportRow{
			SerialNumber:    r.SerialNumber,
			Barcode:        r.Barcode,
			Brand:          r.Brand,
			Size:           r.Size,
			Pattern:        r.Pattern,
			OTD:            r.OTD,
			RTD:            r.RTD,
			Cost:           r.Cost,
			TotalLifeHM:    totalLifeHM,
			CostPerHM:      costPerHM,
			TotalTreadUsed: r.OTD - r.RTD,
			PotentialHM:     potentialHM,
			Status:         r.Status,
			CurrentUnit:    r.CurrentUnit,
			CurrentPosition: r.CurrentPosition,
		})
	}

	return results, nil
}

// GetPotentialHMReport generates a Potential HM report for mounted tyres.
func (uc *ReportUseCase) GetPotentialHMReport(ctx context.Context, filters map[string]interface{}) ([]*PotentialHMReportRow, error) {
	type tyreData struct {
		TyreID       uint
		SerialNumber string
		Barcode      string
		Brand        string
		Size         string
		Pattern      string
		UnitCode     string
		PlateNumber  string
		MountedPos   string
		OTD          float64
		RTDCurrent   float64
	}

	companyID := uint(0)
	if cid, ok := filters["company_id"].(uint); ok {
		companyID = cid
	}

	var tyres []tyreData
	query := uc.db.Table("tyre_master t").
		Select(`
			t.id AS tyre_id,
			t.serial_number,
			t.barcode,
			COALESCE(b.name, '') AS brand,
			COALESCE(s.name, '') AS size,
			COALESCE(p.name, '') AS pattern,
			COALESCE(u.unit_id, '') AS unit_code,
			COALESCE(u.plate_number, '') AS plate_number,
			COALESCE(t.mounted_position, '') AS mounted_pos,
			t.otd,
			t.rtd AS rtd_current
		`).
		Joins("LEFT JOIN units u ON t.unit_id = u.id").
		Joins("LEFT JOIN master_brands b ON t.brand_id = b.id").
		Joins("LEFT JOIN master_sizes s ON t.size_id = s.id").
		Joins("LEFT JOIN master_patterns p ON t.pattern_id = p.id").
		Where("t.status = 'mounted'")

	if companyID > 0 {
		query = query.Where("t.company_id = ?", companyID)
	}

	if err := query.Scan(&tyres).Error; err != nil {
		return nil, err
	}

	results := make([]*PotentialHMReportRow, 0, len(tyres))
	for _, tyre := range tyres {
		// Get total life HM: SUM(old_tyre_lifetime) + (current_unit_hm - new_tyre_hm)
		var sumDismountLife float64
		uc.db.Table("replacement_details").
			Select("COALESCE(SUM(old_tyre_lifetime), 0)").
			Where("old_tyre_id = ? AND action = 'dismount'", tyre.TyreID).
			Scan(&sumDismountLife)

		totalLifeHM := sumDismountLife

		// Get last mount info
		var lastMount struct {
			NewTyreHM  float64
			RTDAtMount float64
		}
		uc.db.Table("replacement_details rd").
			Select("rd.new_tyre_hm, rd.new_tyre_tread_1 AS rtd_at_mount").
			Joins("JOIN replacements r ON rd.replacement_id = r.id").
			Where("rd.new_tyre_id = ? AND rd.action = 'mount'", tyre.TyreID).
			Order("r.id DESC").Limit(1).
			Scan(&lastMount)

		if lastMount.NewTyreHM > 0 && tyre.UnitCode != "" {
			var currentHM float64
			uc.db.Table("units").Select("current_hm").Where("unit_id = ?", tyre.UnitCode).Scan(&currentHM)
			if currentHM > lastMount.NewTyreHM {
				totalLifeHM += (currentHM - lastMount.NewTyreHM)
			}
		}

		// Calculate wear rate from all dismount cycles
		var totalTreadConsumed, totalHMFromDismount float64
		var wearRow struct {
			Tread float64
			HM    float64
		}
		uc.db.Table("replacement_details").
			Select("COALESCE(SUM(new_tyre_tread_1 - old_tyre_tread_1), 0) AS tread, COALESCE(SUM(old_tyre_lifetime), 0) AS hm").
			Where("old_tyre_id = ? AND action = 'dismount' AND old_tyre_lifetime IS NOT NULL AND old_tyre_lifetime > 0", tyre.TyreID).
			Scan(&wearRow)
		totalTreadConsumed = wearRow.Tread
		totalHMFromDismount = wearRow.HM

		wearRate := 0.0
		if totalHMFromDismount > 0 {
			wearRate = totalTreadConsumed / totalHMFromDismount
		}

		// Get rtd_minimum from master_sizes
		rtdMin := 10.0
		uc.db.Table("master_sizes").
			Select("COALESCE(rtd_minimum, 10.0)").
			Where("id = (SELECT size_id FROM tyre_master WHERE id = ?)", tyre.TyreID).
			Scan(&rtdMin)

		// Potential HM calculation
		remainingRTD := tyre.RTDCurrent - rtdMin
		if remainingRTD < 0 {
			remainingRTD = 0
		}

		potentialHM := 0.0
		if wearRate > 0 && remainingRTD > 0 {
			potentialHM = remainingRTD / wearRate
		}
		totalPotentialHM := totalLifeHM + potentialHM

		// Status urgency
		status := "good"
		if tyre.RTDCurrent < 5.0 {
			status = "critical"
		} else if tyre.RTDCurrent < 10.0 {
			status = "warning"
		}

		results = append(results, &PotentialHMReportRow{
			SerialNumber:     tyre.SerialNumber,
			Barcode:          tyre.Barcode,
			Brand:            tyre.Brand,
			Size:             tyre.Size,
			Pattern:          tyre.Pattern,
			UnitCode:         tyre.UnitCode,
			PlateNumber:      tyre.PlateNumber,
			Position:         tyre.MountedPos,
			OTD:              tyre.OTD,
			RTDATMount:       lastMount.RTDAtMount,
			RTDCurrent:       tyre.RTDCurrent,
			RTDMinimum:       rtdMin,
			TotalLifeHM:      totalLifeHM,
			WearRate:         wearRate,
			RemainingRTD:     remainingRTD,
			PotentialHM:      potentialHM,
			TotalPotentialHM:  totalPotentialHM,
			Status:           status,
		})
	}

	// Sort by TotalPotentialHM ascending (urgent first)
	for i := 1; i < len(results); i++ {
		for j := i; j > 0 && results[j].TotalPotentialHM < results[j-1].TotalPotentialHM; j-- {
			results[j], results[j-1] = results[j-1], results[j]
		}
	}

	return results, nil
}
