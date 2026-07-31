package repository

import "github.com/tms/tyre/internal/domain/entity"

type DriverRepository interface {
	Create(driver *entity.Driver) error
	GetByID(id uint) (*entity.Driver, error)
	Update(driver *entity.Driver) error
	Delete(id uint) error
	List(page, perPage int, companyID uint, status string) ([]*entity.Driver, int64, error)
}
