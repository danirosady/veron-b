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

func NewTyreRepository(db *gorm.DB) repository.TyreRepository {
	return &tyreRepository{NewBaseRepository(db)}
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
	return &tyre, nil
}

func (r *tyreRepository) GetByBarcode(barcode string) (*entity.TyreMaster, error) {
	var tyre entity.TyreMaster
	err := r.db.
		Preload("Company").
		Preload("Unit").
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
	return &tyre, nil
}

func (r *tyreRepository) GetBySerialNumber(sn string) (*entity.TyreMaster, error) {
	var tyre entity.TyreMaster
	err := r.db.
		Preload("Company").
		Preload("Unit").
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
	return &tyre, nil
}

func (r *tyreRepository) Update(tyre *entity.TyreMaster) error {
	return r.db.Save(tyre).Error
}

func (r *tyreRepository) Delete(id uint) error {
	return r.db.Delete(&entity.TyreMaster{}, id).Error
}

func (r *tyreRepository) List(page, perPage int, companyID uint, status, brandID, sizeID string) ([]*entity.TyreMaster, int64, error) {
	var tyres []*entity.TyreMaster
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

	offset := (page - 1) * perPage
	err := db.
		Select("tyre_master.*").
		Preload("Company").
		Preload("Size").
		Preload("Brand").
		Preload("Pattern").
		Offset(offset).Limit(perPage).
		Order("id DESC").
		Find(&tyres).Error
	if err != nil {
		return nil, 0, err
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
