package repository

import (
	"errors"

	"github.com/tms/tyre/internal/domain/entity"
	"github.com/tms/tyre/internal/domain/repository"
	"gorm.io/gorm"
)

type employeeRepository struct {
	*BaseRepository
}

func NewEmployeeRepository(db *gorm.DB) repository.EmployeeRepository {
	return &employeeRepository{NewBaseRepository(db)}
}

func (r *employeeRepository) Create(employee *entity.Employee) error {
	return r.db.Create(employee).Error
}

func (r *employeeRepository) GetByID(id uint) (*entity.Employee, error) {
	var employee entity.Employee
	err := r.db.Preload("Company").First(&employee, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &employee, nil
}

func (r *employeeRepository) Update(employee *entity.Employee) error {
	return r.db.Save(employee).Error
}

func (r *employeeRepository) Delete(id uint) error {
	return r.db.Delete(&entity.Employee{}, id).Error
}

func (r *employeeRepository) List(page, perPage int, companyID uint, role, status string) ([]*entity.Employee, int64, error) {
	var employees []*entity.Employee
	var total int64

	query := r.db.Model(&entity.Employee{})
	if companyID > 0 {
		query = query.Where("company_id = ?", companyID)
	}
	if role != "" {
		query = query.Where("role = ?", role)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage
	if err := query.Preload("Company").Offset(offset).Limit(perPage).Order("id DESC").Find(&employees).Error; err != nil {
		return nil, 0, err
	}

	return employees, total, nil
}
