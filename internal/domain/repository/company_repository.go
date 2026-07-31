package repository

import "github.com/tms/tyre/internal/domain/entity"

type CompanyRepository interface {
	Create(company *entity.Company) error
	GetByID(id uint) (*entity.Company, error)
	GetByName(name string) (*entity.Company, error)
	Update(company *entity.Company) error
	Delete(id uint) error
	List(page, perPage int, status string) ([]*entity.Company, int64, error)
	HasActiveData(id uint) (bool, error)
}
