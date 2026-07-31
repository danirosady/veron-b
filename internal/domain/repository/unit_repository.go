package repository

import "github.com/tms/tyre/internal/domain/entity"

type UnitRepository interface {
	Create(unit *entity.Unit) error
	GetByID(id uint) (*entity.Unit, error)
	GetByUnitID(unitID string) (*entity.Unit, error)
	Update(unit *entity.Unit) error
	Delete(id uint) error
	List(page, perPage int, companyID, projectID uint, status string) ([]*entity.Unit, int64, error)
	UpdateHM(id uint, hm float64) error
	HasMountedTyres(id uint) (bool, error)
	HasReplacements(id uint) (bool, error)
}
