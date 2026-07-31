package usecase

import (
	"time"

	"gorm.io/gorm"
)

type DashboardUseCase struct {
	db *gorm.DB
}

func NewDashboardUseCase(db *gorm.DB) *DashboardUseCase {
	return &DashboardUseCase{db: db}
}

type DashboardStats struct {
	TotalCompanies         int64 `json:"total_companies"`
	TotalProjects         int64 `json:"total_projects"`
	TotalUnits            int64 `json:"total_units"`
	TotalDrivers          int64 `json:"total_drivers"`
	TotalTyres            int64 `json:"total_tyres"`
	MountedTyres          int64 `json:"mounted_tyres"`
	SpareTyres            int64 `json:"spare_tyres"`
	ScrapTyres            int64 `json:"scrap_tyres"`
	TotalReplacements      int64 `json:"total_replacements"`
	ReplacementsThisMonth int64 `json:"replacements_this_month"`
}

func (uc *DashboardUseCase) GetStats(companyID *uint) (*DashboardStats, error) {
	stats := &DashboardStats{}

	var total int64

	if companyID != nil {
		// Company-scoped stats
		if err := uc.db.Table("companies").Where("id = ? AND status = 'active'", *companyID).Count(&total).Error; err != nil {
			return nil, err
		}
		stats.TotalCompanies = total

		if err := uc.db.Table("projects").Where("company_id = ? AND status = 'active'", *companyID).Count(&stats.TotalProjects).Error; err != nil {
			return nil, err
		}
		if err := uc.db.Table("units").Where("company_id = ? AND status = 'active'", *companyID).Count(&stats.TotalUnits).Error; err != nil {
			return nil, err
		}
		if err := uc.db.Table("drivers").Where("company_id = ? AND status = 'active'", *companyID).Count(&stats.TotalDrivers).Error; err != nil {
			return nil, err
		}
		if err := uc.db.Table("tyre_master").Where("company_id = ?", *companyID).Count(&stats.TotalTyres).Error; err != nil {
			return nil, err
		}
		if err := uc.db.Table("tyre_master").Where("company_id = ? AND status = 'mounted'", *companyID).Count(&stats.MountedTyres).Error; err != nil {
			return nil, err
		}
		if err := uc.db.Table("tyre_master").Where("company_id = ? AND status = 'spare'", *companyID).Count(&stats.SpareTyres).Error; err != nil {
			return nil, err
		}
		if err := uc.db.Table("tyre_master").Where("company_id = ? AND status = 'scrap'", *companyID).Count(&stats.ScrapTyres).Error; err != nil {
			return nil, err
		}
		if err := uc.db.Table("replacements").Where("company_id = ?", *companyID).Count(&stats.TotalReplacements).Error; err != nil {
			return nil, err
		}
	} else {
		// Global stats
		if err := uc.db.Table("companies").Where("status = 'active'").Count(&stats.TotalCompanies).Error; err != nil {
			return nil, err
		}
		if err := uc.db.Table("projects").Where("status = 'active'").Count(&stats.TotalProjects).Error; err != nil {
			return nil, err
		}
		if err := uc.db.Table("units").Where("status = 'active'").Count(&stats.TotalUnits).Error; err != nil {
			return nil, err
		}
		if err := uc.db.Table("drivers").Where("status = 'active'").Count(&stats.TotalDrivers).Error; err != nil {
			return nil, err
		}
		if err := uc.db.Table("tyre_master").Count(&stats.TotalTyres).Error; err != nil {
			return nil, err
		}
		if err := uc.db.Table("tyre_master").Where("status = 'mounted'").Count(&stats.MountedTyres).Error; err != nil {
			return nil, err
		}
		if err := uc.db.Table("tyre_master").Where("status = 'spare'").Count(&stats.SpareTyres).Error; err != nil {
			return nil, err
		}
		if err := uc.db.Table("tyre_master").Where("status = 'scrap'").Count(&stats.ScrapTyres).Error; err != nil {
			return nil, err
		}
		if err := uc.db.Table("replacements").Count(&stats.TotalReplacements).Error; err != nil {
			return nil, err
		}
	}

	monthStart := time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.UTC)
	if err := uc.db.Table("replacements").Where("date >= ?", monthStart).Count(&stats.ReplacementsThisMonth).Error; err != nil {
		return nil, err
	}

	return stats, nil
}
