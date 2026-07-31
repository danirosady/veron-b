package repository

import (
	"errors"

	"github.com/tms/tyre/internal/domain/entity"
	"github.com/tms/tyre/internal/domain/repository"
	"gorm.io/gorm"
)

type driverRepository struct {
	*BaseRepository
}

func NewDriverRepository(db *gorm.DB) repository.DriverRepository {
	return &driverRepository{NewBaseRepository(db)}
}

func (r *driverRepository) Create(driver *entity.Driver) error {
	return r.db.Create(driver).Error
}

func (r *driverRepository) GetByID(id uint) (*entity.Driver, error) {
	var driver entity.Driver
	err := r.db.Preload("Company").First(&driver, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &driver, nil
}

func (r *driverRepository) Update(driver *entity.Driver) error {
	return r.db.Save(driver).Error
}

func (r *driverRepository) Delete(id uint) error {
	return r.db.Delete(&entity.Driver{}, id).Error
}

func (r *driverRepository) List(page, perPage int, companyID uint, status string) ([]*entity.Driver, int64, error) {
	var drivers []*entity.Driver
	var total int64

	query := r.db.Model(&entity.Driver{}).Where("company_id = ?", companyID)
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage
	if err := query.Preload("Company").Offset(offset).Limit(perPage).Order("id DESC").Find(&drivers).Error; err != nil {
		return nil, 0, err
	}

	return drivers, total, nil
}
