package repository

import (
	"errors"
	"strconv"

	"github.com/tms/tyre/internal/domain/entity"
	"github.com/tms/tyre/internal/domain/repository"
	"gorm.io/gorm"
)

type tyreRepository struct {
	*BaseRepository
}

// tyreWithUnit holds the result of the LEFT JOIN for the list query
type tyreWithUnit struct {
	entity.TyreMaster
	// Unit fields selected from the JOIN
	UnitID_Tyre   *uint   `gorm:"column:u_id"`
	UnitCode      string  `gorm:"column:u_unit_id"`
	UnitModel     string  `gorm:"column:u_unit_model"`
	UnitPlateNum  string  `gorm:"column:u_plate_number"`
	UnitType      string  `gorm:"column:u_unit_type"`
	UnitStatus    string  `gorm:"column:u_status"`
	UnitCompanyID uint    `gorm:"column:u_company_id"`
	UnitProjectID uint    `gorm:"column:u_project_id"`
}

func NewTyreRepository(db *gorm.DB) repository.TyreRepository {
	return &tyreRepository{NewBaseRepository(db)}
}

// loadUnit fetches the related Unit for a tyre via LEFT JOIN and sets tyre.Unit
func (r *tyreRepository) loadUnit(tyre *entity.TyreMaster) {
	var row tyreWithUnit
	if err := r.db.Raw(`
		SELECT t.*, u.id AS u_id, u.unit_id AS u_unit_id, u.unit_model AS u_unit_model,
			u.plate_number AS u_plate_number, u.unit_type AS u_unit_type,
			u.status AS u_status, u.company_id AS u_company_id, u.project_id AS u_project_id
		FROM tyre_master t LEFT JOIN units u ON u.id = t.unit_id WHERE t.id = ?`, tyre.ID).Scan(&row).Error; err == nil {
		if row.UnitID_Tyre != nil {
			tyre.Unit = &entity.Unit{
				ID:          *row.UnitID_Tyre,
				UnitID:      row.UnitCode,
				UnitModel:   row.UnitModel,
				PlateNumber: row.UnitPlateNum,
				UnitType:    row.UnitType,
				Status:      row.UnitStatus,
				CompanyID:   row.UnitCompanyID,
				ProjectID:   row.UnitProjectID,
			}
		}
	}
}

func (r *tyreRepository) Create(tyre *entity.TyreMaster) error {
	return r.db.Create(tyre).Error
}

func (r *tyreRepository) GetByID(id uint) (*entity.TyreMaster, error) {
	var tyre entity.TyreMaster
	err := r.db.
		Preload("Company").
		Preload("Size").
		Preload("Brand").
		Preload("Pattern").
		Where("id = ?", int64(id)).First(&tyre).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	r.loadUnit(&tyre)
	return &tyre, nil
}

func (r *tyreRepository) GetByBarcode(barcode string) (*entity.TyreMaster, error) {
	var tyre entity.TyreMaster
	err := r.db.
		Preload("Company").
		Preload("Size").
		Preload("Brand").
		Preload("Pattern").
		Where("barcode = ?", barcode).First(&tyre).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	r.loadUnit(&tyre)
	return &tyre, nil
}

func (r *tyreRepository) GetBySerialNumber(sn string) (*entity.TyreMaster, error) {
	var tyre entity.TyreMaster
	err := r.db.
		Preload("Company").
		Preload("Size").
		Preload("Brand").
		Preload("Pattern").
		Where("serial_number = ?", sn).First(&tyre).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	r.loadUnit(&tyre)
	return &tyre, nil
}

func (r *tyreRepository) Update(tyre *entity.TyreMaster) error {
	return r.db.Save(tyre).Error
}

func (r *tyreRepository) Delete(id uint) error {
	return r.db.Delete(&entity.TyreMaster{}, id).Error
}

func (r *tyreRepository) List(page, perPage int, companyID uint, status, brandID, sizeID string) ([]*entity.TyreMaster, int64, error) {
	var total int64

	db := r.db.Model(&entity.TyreMaster{})
	if companyID > 0 {
		db = db.Where("company_id = ?", companyID)
	}
	if status != "" {
		db = db.Where("status = ?", status)
	}
	if brandID != "" {
		if n, err := strconv.ParseUint(brandID, 10, 64); err == nil {
			db = db.Where("brand_id = ?", n)
		}
	}
	if sizeID != "" {
		if n, err := strconv.ParseUint(sizeID, 10, 64); err == nil {
			db = db.Where("size_id = ?", n)
		}
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Build the full query string to use with Scan (avoids GORM preload type mismatch)
	query := `SELECT t.*,
		u.id AS u_id, u.unit_id AS u_unit_id, u.unit_model AS u_unit_model,
		u.plate_number AS u_plate_number, u.unit_type AS u_unit_type,
		u.status AS u_status, u.company_id AS u_company_id, u.project_id AS u_project_id
		FROM tyre_master t
		LEFT JOIN units u ON u.id = t.unit_id`

	args := []interface{}{}
	where := ""
	if companyID > 0 {
		where += " t.company_id = ?"
		args = append(args, companyID)
	}
	if status != "" {
		if where != "" {
			where += " AND"
		}
		where += " t.status = ?"
		args = append(args, status)
	}
	if brandID != "" {
		if n, err := strconv.ParseUint(brandID, 10, 64); err == nil {
			if where != "" {
				where += " AND"
			}
			where += " t.brand_id = ?"
			args = append(args, n)
		}
	}
	if sizeID != "" {
		if n, err := strconv.ParseUint(sizeID, 10, 64); err == nil {
			if where != "" {
				where += " AND"
			}
			where += " t.size_id = ?"
			args = append(args, n)
		}
	}
	if where != "" {
		query += " WHERE" + where
	}

	query += " ORDER BY t.id DESC LIMIT ? OFFSET ?"
	offset := (page - 1) * perPage
	args = append(args, perPage, offset)

	var rows []tyreWithUnit
	if err := r.db.Raw(query, args...).Scan(&rows).Error; err != nil {
		return nil, 0, err
	}

	tyres := make([]*entity.TyreMaster, len(rows))
	for i, row := range rows {
		row.Unit = nil
		if row.UnitID_Tyre != nil {
			row.Unit = &entity.Unit{
				ID:         *row.UnitID_Tyre,
				UnitID:     row.UnitCode,
				UnitModel:  row.UnitModel,
				PlateNumber: row.UnitPlateNum,
				UnitType:   row.UnitType,
				Status:     row.UnitStatus,
				CompanyID:  row.UnitCompanyID,
				ProjectID:  row.UnitProjectID,
			}
		}
		tyres[i] = &row.TyreMaster
	}

	return tyres, total, nil
}

func (r *tyreRepository) GetByUnitID(unitID uint) ([]*entity.TyreMaster, error) {
	var tyres []*entity.TyreMaster
	err := r.db.
		Where("unit_id = ? AND status = 'mounted'", unitID).
		Preload("Company").
		Preload("Size").
		Preload("Brand").
		Preload("Pattern").
		Order("mounted_position ASC").
		Find(&tyres).Error
	return tyres, err
}

func (r *tyreRepository) GetSpareTyres(companyID uint) ([]*entity.TyreMaster, error) {
	var tyres []*entity.TyreMaster
	err := r.db.
		Where("company_id = ? AND status IN ('spare', 'dismounted')", companyID).
		Preload("Company").
		Preload("Size").
		Preload("Brand").
		Preload("Pattern").
		Order("id DESC").
		Find(&tyres).Error
	return tyres, err
}

func (r *tyreRepository) Mount(tyreID uint, unitID uint, position string) error {
	return r.db.Model(&entity.TyreMaster{}).
		Where("id = ?", tyreID).
		Updates(map[string]interface{}{
			"unit_id":           unitID,
			"mounted_position": position,
			"status":            "mounted",
		}).Error
}

func (r *tyreRepository) Dismount(tyreID uint, status string) error {
	return r.db.Model(&entity.TyreMaster{}).
		Where("id = ?", tyreID).
		Updates(map[string]interface{}{
			"unit_id":           nil,
			"mounted_position": nil,
			"status":            status,
		}).Error
}
