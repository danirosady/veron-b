package repository

import "github.com/tms/tyre/internal/domain/entity"

type EmployeeRepository interface {
	Create(employee *entity.Employee) error
	GetByID(id uint) (*entity.Employee, error)
	Update(employee *entity.Employee) error
	Delete(id uint) error
	List(page, perPage int, companyID uint, role, status string) ([]*entity.Employee, int64, error)
}
