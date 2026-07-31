package repository

import (
	"time"

	"github.com/tms/tyre/internal/domain/entity"
)

type ReplacementRepository interface {
	Create(replacement *entity.Replacement) (*entity.Replacement, error)
	Update(replacement *entity.Replacement) error
	Delete(id uint) error
	GetByID(id uint) (*entity.Replacement, error)
	List(page, perPage int, companyID, projectID, unitID uint, dateFrom, dateTo *time.Time) ([]*entity.Replacement, int64, error)
	GetLastReplacementByUnitID(unitID uint) (*entity.Replacement, error)
}
