package repository

import (
	"errors"

	"github.com/tms/tyre/internal/domain/entity"
	"github.com/tms/tyre/internal/domain/repository"
	"gorm.io/gorm"
)

type projectRepository struct {
	*BaseRepository
}

func NewProjectRepository(db *gorm.DB) repository.ProjectRepository {
	return &projectRepository{NewBaseRepository(db)}
}

func (r *projectRepository) Create(project *entity.Project) error {
	return r.db.Create(project).Error
}

func (r *projectRepository) GetByID(id uint) (*entity.Project, error) {
	var project entity.Project
	err := r.db.Preload("Company").First(&project, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &project, nil
}

func (r *projectRepository) Update(project *entity.Project) error {
	return r.db.Save(project).Error
}

func (r *projectRepository) Delete(id uint) error {
	return r.db.Delete(&entity.Project{}, id).Error
}

func (r *projectRepository) List(page, perPage int, companyID uint, status string) ([]*entity.Project, int64, error) {
	var projects []*entity.Project
	var total int64

	query := r.db.Model(&entity.Project{}).Where("company_id = ?", companyID)
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage
	if err := query.Preload("Company").Offset(offset).Limit(perPage).Order("id DESC").Find(&projects).Error; err != nil {
		return nil, 0, err
	}

	return projects, total, nil
}

func (r *projectRepository) HasActiveUnits(id uint) (bool, error) {
	var count int64
	err := r.db.Model(&entity.Unit{}).Where("project_id = ? AND status = 'active'", id).Count(&count).Error
	return count > 0, err
}
