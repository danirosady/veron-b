package repository

import (
	"errors"

	"github.com/tms/tyre/internal/domain/entity"
	"github.com/tms/tyre/internal/domain/repository"
	"gorm.io/gorm"
)

type unitRepository struct {
	*BaseRepository
}

func NewUnitRepository(db *gorm.DB) repository.UnitRepository {
	return &unitRepository{NewBaseRepository(db)}
}

func (r *unitRepository) Create(unit *entity.Unit) error {
	return r.db.Create(unit).Error
}

func (r *unitRepository) GetByID(id uint) (*entity.Unit, error) {
	var unit entity.Unit
	err := r.db.Preload("Company").Preload("Project").First(&unit, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &unit, nil
}

func (r *unitRepository) GetByUnitID(unitID string) (*entity.Unit, error) {
	var unit entity.Unit
	err := r.db.Where("unit_id = ?", unitID).Preload("Company").Preload("Project").First(&unit).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &unit, nil
}

func (r *unitRepository) Update(unit *entity.Unit) error {
	return r.db.Save(unit).Error
}

func (r *unitRepository) Delete(id uint) error {
	return r.db.Delete(&entity.Unit{}, id).Error
}

func (r *unitRepository) List(page, perPage int, companyID, projectID uint, status string) ([]*entity.Unit, int64, error) {
	var units []*entity.Unit
	var total int64

	query := r.db.Model(&entity.Unit{}).Where("company_id = ?", companyID)
	if projectID > 0 {
		query = query.Where("project_id = ?", projectID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage
	if err := query.Preload("Company").Preload("Project").Offset(offset).Limit(perPage).Order("id DESC").Find(&units).Error; err != nil {
		return nil, 0, err
	}

	return units, total, nil
}

func (r *unitRepository) UpdateHM(id uint, hm float64) error {
	return r.db.Model(&entity.Unit{}).Where("id = ? AND current_hm < ?", id, hm).Update("current_hm", hm).Error
}

func (r *unitRepository) HasMountedTyres(id uint) (bool, error) {
	var count int64
	err := r.db.Model(&entity.TyreMaster{}).Where("unit_id = ?", id).Count(&count).Error
	return count > 0, err
}

func (r *unitRepository) HasReplacements(id uint) (bool, error) {
	var count int64
	err := r.db.Model(&entity.Replacement{}).Where("unit_id = ?", id).Count(&count).Error
	return count > 0, err
}
