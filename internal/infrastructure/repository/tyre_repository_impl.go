package repository

import (
	"errors"
	"strconv"
	"time"

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
	// Unit fields
	UnitID_Tyre   *uint  `gorm:"column:u_id"`
	UnitCode      string `gorm:"column:u_unit_id"`
	UnitModel     string `gorm:"column:u_unit_model"`
	UnitPlateNum  string `gorm:"column:u_plate_number"`
	UnitType      string `gorm:"column:u_unit_type"`
	UnitStatus    string `gorm:"column:u_status"`
	UnitCompanyID uint   `gorm:"column:u_company_id"`
	UnitProjectID uint   `gorm:"column:u_project_id"`
	// Size fields
	SizeID2  uint   `gorm:"column:s_id"`
	SizeName string `gorm:"column:s_name"`
	// Brand fields
	BrandID2  uint   `gorm:"column:b_id"`
	BrandName string `gorm:"column:b_name"`
	// Pattern fields
	PatternID2  uint   `gorm:"column:p_id"`
	PatternName string `gorm:"column:p_name"`
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

// loadCurrentMount fetches the most recent mount record for a tyre from replacement_details
func (r *tyreRepository) loadCurrentMount(tyre *entity.TyreMaster) {
	if tyre.Status != "mounted" {
		return
	}
	var row struct {
		Position  string    `gorm:"column:position"`
		MountDate time.Time `gorm:"column:mount_date"`
		MountHM   float64   `gorm:"column:hm_update"`
		UnitID    uint      `gorm:"column:unit_id"`
		UnitCode  string    `gorm:"column:unit_id_str"`
		UnitPlate string    `gorm:"column:plate_number"`
		UnitType  string    `gorm:"column:unit_type"`
		UnitStat  string    `gorm:"column:u_status"`
	}
	err := r.db.Raw(`
		SELECT
			rd.position,
			rd.created_at AS mount_date,
			r.hm_update,
			u.id AS unit_id,
			u.unit_id AS unit_id_str,
			u.plate_number,
			u.unit_type,
			u.status AS u_status
		FROM replacement_details rd
		JOIN replacements r ON r.id = rd.replacement_id
		JOIN units u ON u.id = r.unit_id
		WHERE rd.new_tyre_id = ? AND rd.action = 'mount'
		ORDER BY rd.created_at DESC
		LIMIT 1
	`, tyre.ID).Scan(&row).Error
	if err != nil || row.UnitID == 0 {
		return
	}
	tyre.CurrentMount = &entity.CurrentMount{
		Position:  row.Position,
		MountDate: row.MountDate,
		MountHM:   row.MountHM,
		Unit: &entity.Unit{
			ID:          row.UnitID,
			UnitID:      row.UnitCode,
			PlateNumber: row.UnitPlate,
			UnitType:    row.UnitType,
			Status:      row.UnitStat,
		},
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
	r.loadCurrentMount(&tyre)
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
	r.loadCurrentMount(&tyre)
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
	r.loadCurrentMount(&tyre)
	return &tyre, nil
}

func (r *tyreRepository) Update(tyre *entity.TyreMaster) error {
	return r.db.Model(tyre).Select("brand_id", "size_id", "pattern_id", "type", "rtd", "rtd1", "rtd2", "psi", "remarks", "status", "dot_code").Updates(tyre).Error
}

func (r *tyreRepository) Delete(id uint) error {
	return r.db.Delete(&entity.TyreMaster{}, id).Error
}

func (r *tyreRepository) List(page, perPage int, companyID uint, status, brandID, sizeID, search string) ([]*entity.TyreMaster, int64, error) {
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
	if search != "" {
		db = db.Where("barcode LIKE ? OR serial_number LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Build the full query string to use with Scan (avoids GORM preload type mismatch)
	query := `SELECT t.*,
		u.id AS u_id, u.unit_id AS u_unit_id, u.unit_model AS u_unit_model,
		u.plate_number AS u_plate_number, u.unit_type AS u_unit_type,
		u.status AS u_status, u.company_id AS u_company_id, u.project_id AS u_project_id,
		s.id AS s_id, s.name AS s_name,
		b.id AS b_id, b.name AS b_name,
		p.id AS p_id, p.name AS p_name
		FROM tyre_master t
		LEFT JOIN units u ON u.id = t.unit_id
		LEFT JOIN master_sizes s ON s.id = t.size_id
		LEFT JOIN master_brands b ON b.id = t.brand_id
		LEFT JOIN master_patterns p ON p.id = t.pattern_id`

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
	if search != "" {
		if where != "" {
			where += " AND"
		}
		where += " (t.barcode LIKE ? OR t.serial_number LIKE ?)"
		args = append(args, "%"+search+"%", "%"+search+"%")
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
		if row.UnitID_Tyre != nil {
			row.Unit = &entity.Unit{
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
		if row.SizeID2 > 0 {
			row.Size = &entity.MasterSize{ID: row.SizeID2, Name: row.SizeName}
		}
		if row.BrandID2 > 0 {
			row.Brand = &entity.MasterBrand{ID: row.BrandID2, Name: row.BrandName}
		}
		if row.PatternID2 > 0 {
			row.Pattern = &entity.MasterPattern{ID: row.PatternID2, Name: row.PatternName}
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

func (r *tyreRepository) Mount(tyreID uint, unitID uint, position string, rtd1, rtd2 *float64) error {
	updates := map[string]interface{}{
		"unit_id":           unitID,
		"mounted_position":  position,
		"status":            "mounted",
	}
	if rtd1 != nil {
		updates["rtd1"] = rtd1
		updates["rtd"] = *rtd1
	}
	if rtd2 != nil {
		updates["rtd2"] = rtd2
	}
	return r.db.Model(&entity.TyreMaster{}).
		Where("id = ?", tyreID).
		Updates(updates).Error
}

func (r *tyreRepository) Dismount(tyreID uint, status string, rtd1, rtd2 *float64) error {
	updates := map[string]interface{}{
		"unit_id":           nil,
		"mounted_position":  nil,
		"status":            status,
	}
	if rtd1 != nil {
		updates["rtd1"] = rtd1
		updates["rtd"] = *rtd1
	}
	if rtd2 != nil {
		updates["rtd2"] = rtd2
	}
	return r.db.Model(&entity.TyreMaster{}).
		Where("id = ?", tyreID).
		Updates(updates).Error
}

func (r *tyreRepository) GetTyreHistory(tyreID uint) ([]*entity.TyreHistoryItem, error) {
	rows, err := r.db.Raw(`
		SELECT
			rd.id,
			rd.created_at AS replacement_date,
			r.hm_update AS hm,
			rd.position,
			rd.action,
			COALESCE(rd.new_tyre_tread_1, rd.new_tyre_tread_2) AS rtd,
			rd.remark AS remarks,
			u.id AS unit_id,
			u.unit_id AS unit_code,
			u.plate_number AS unit_plate_number,
			u.unit_type AS unit_type,
			u.status AS unit_status
		FROM replacement_details rd
		JOIN replacements r ON r.id = rd.replacement_id
		JOIN units u ON u.id = r.unit_id
		WHERE rd.new_tyre_id = ? OR rd.old_tyre_id = ?
		ORDER BY rd.created_at DESC
	`, tyreID, tyreID).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*entity.TyreHistoryItem
	for rows.Next() {
		var item entity.TyreHistoryItem
		var unitID *uint
		var unitCode, unitPlate, unitType, unitStatus string

		rows.Scan(
			&item.ID, &item.ReplacementDate, &item.HM,
			&item.Position, &item.Action, &item.RTD, &item.Remarks,
			&unitID, &unitCode, &unitPlate, &unitType, &unitStatus,
		)

		if unitID != nil {
			item.Unit = &entity.Unit{
				ID:          *unitID,
				UnitID:      unitCode,
				PlateNumber: unitPlate,
				UnitType:    unitType,
				Status:      unitStatus,
			}
		}
		items = append(items, &item)
	}
	return items, nil
}
