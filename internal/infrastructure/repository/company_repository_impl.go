package repository

import (
	"errors"

	"github.com/tms/tyre/internal/domain/entity"
	"github.com/tms/tyre/internal/domain/repository"
	"gorm.io/gorm"
)

type companyRepository struct {
	*BaseRepository
}

func NewCompanyRepository(db *gorm.DB) repository.CompanyRepository {
	return &companyRepository{NewBaseRepository(db)}
}

func (r *companyRepository) Create(company *entity.Company) error {
	return r.db.Create(company).Error
}

func (r *companyRepository) GetByID(id uint) (*entity.Company, error) {
	var company entity.Company
	err := r.db.First(&company, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &company, nil
}

func (r *companyRepository) GetByName(name string) (*entity.Company, error) {
	var company entity.Company
	err := r.db.Where("name = ?", name).First(&company).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &company, nil
}

func (r *companyRepository) Update(company *entity.Company) error {
	return r.db.Save(company).Error
}

func (r *companyRepository) Delete(id uint) error {
	return r.db.Delete(&entity.Company{}, id).Error
}

func (r *companyRepository) List(page, perPage int, status string) ([]*entity.Company, int64, error) {
	var companies []*entity.Company
	var total int64

	query := r.db.Model(&entity.Company{})
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("id DESC").Find(&companies).Error; err != nil {
		return nil, 0, err
	}

	return companies, total, nil
}

func (r *companyRepository) HasActiveData(id uint) (bool, error) {
	var count int64

	if err := r.db.Model(&entity.Project{}).Where("company_id = ? AND status = 'active'", id).Count(&count).Error; err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	}

	if err := r.db.Model(&entity.Unit{}).Where("company_id = ? AND status = 'active'", id).Count(&count).Error; err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	}

	if err := r.db.Model(&entity.Driver{}).Where("company_id = ? AND status = 'active'", id).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}
