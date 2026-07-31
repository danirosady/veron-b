package repository

import "github.com/tms/tyre/internal/domain/entity"

type ProjectRepository interface {
	Create(project *entity.Project) error
	GetByID(id uint) (*entity.Project, error)
	Update(project *entity.Project) error
	Delete(id uint) error
	List(page, perPage int, companyID uint, status string) ([]*entity.Project, int64, error)
	HasActiveUnits(id uint) (bool, error)
}
