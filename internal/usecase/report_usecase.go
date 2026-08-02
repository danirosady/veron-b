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
	Date         string
	UnitCode     string
	PlateNumber  string
	CompanyName  string
	Position     string
	Action       string
	OldTyreSN    string
	NewTyreSN    string
	Brand        string
	Size         string
	Pattern      string
	OldRTD       float64
	NewRTD       float64
	HM           float64
	DriverName   string
	OperatorName string
	Remarks      string
}

// InventoryReportRow represents a row in the inventory report.
type InventoryReportRow struct {
	Barcode       string
	SerialNumber  string
	Brand         string
	Size          string
	Pattern       string
	Status        string
	OTD           float64
	RTD           float64
	PercentWorn   float64
	UnitCode      string
	PlateNumber   string
	CompanyName   string
	LastPosition  string
	LastMountDate string
	TotalDepth    float64
	CurrentDepth  float64
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

// GetReplacementReport generates a replacement report based on filters.
func (uc *ReportUseCase) GetReplacementReport(ctx context.Context, filters map[string]interface{}) ([]*ReplacementReportRow, error) {
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
		Joins("LEFT JOIN drivers d ON r.driver_id = d.id").
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

	var results []*ReplacementReportRow
	if err := query.Scan(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}

// GetInventoryReport generates an inventory report based on filters.
func (uc *ReportUseCase) GetInventoryReport(ctx context.Context, filters map[string]interface{}) ([]*InventoryReportRow, error) {
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
			'' AS last_mount_date,
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

	var results []*InventoryReportRow
	if err := query.Scan(&results).Error; err != nil {
		return nil, err
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
