package repository

import (
	"errors"

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
		Preload("Unit").
		Preload("Size").
		Preload("Brand").
		Preload("Pattern").
		First(&tyre, id).Error
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

	query := r.db.Model(&entity.TyreMaster{}).Where("company_id = ?", companyID)
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if brandID != "" {
		query = query.Where("brand_id = ?", brandID)
	}
	if sizeID != "" {
		query = query.Where("size_id = ?", sizeID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage
	if err := query.
		Preload("Company").
		Preload("Unit").
		Preload("Size").
		Preload("Brand").
		Preload("Pattern").
		Offset(offset).Limit(perPage).
		Order("id DESC").
		Find(&tyres).Error; err != nil {
		return nil, 0, err
	}

	return tyres, total, nil
}

func (r *tyreRepository) GetByUnitID(unitID uint) ([]*entity.TyreMaster, error) {
	var tyres []*entity.TyreMaster
	err := r.db.
		Where("unit_id = ? AND status = 'mounted'", unitID).
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
		Preload("Size").
		Preload("Brand").
		Preload("Pattern").
		Order("id DESC").
		Find(&tyres).Error
	return tyres, err
}

func (r *tyreRepository) Mount(tyreID uint, unitID uint, position int) error {
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
