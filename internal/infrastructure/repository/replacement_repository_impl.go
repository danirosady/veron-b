package repository

import (
	"errors"
	"time"

	"github.com/tms/tyre/internal/domain/entity"
	"github.com/tms/tyre/internal/domain/repository"
	"gorm.io/gorm"
)

type replacementRepository struct {
	*BaseRepository
}

func NewReplacementRepository(db *gorm.DB) repository.ReplacementRepository {
	return &replacementRepository{NewBaseRepository(db)}
}

func (r *replacementRepository) Create(replacement *entity.Replacement) (*entity.Replacement, error) {
	err := r.db.Create(replacement).Error
	if err != nil {
		return nil, err
	}
	var created entity.Replacement
	err = r.db.
		Preload("Company").
		Preload("Project").
		Preload("Unit").
		Preload("Driver").
		Preload("Creator").
		Preload("Details").
		Preload("Details.FailureReason").
		First(&created, replacement.ID).Error
	if err != nil {
		return nil, err
	}
	return &created, nil
}

func (r *replacementRepository) GetByID(id uint) (*entity.Replacement, error) {
	var replacement entity.Replacement
	err := r.db.
		Preload("Company").
		Preload("Project").
		Preload("Unit").
		Preload("Driver").
		Preload("Creator").
		Preload("Details").
		Preload("Details.FailureReason").
		First(&replacement, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &replacement, nil
}

func (r *replacementRepository) List(page, perPage int, companyID, projectID, unitID uint, dateFrom, dateTo *time.Time) ([]*entity.Replacement, int64, error) {
	var replacements []*entity.Replacement
	var total int64

	base := r.db.Model(&entity.Replacement{})
	if companyID > 0 {
		base = base.Where("company_id = ?", companyID)
	}
	if projectID > 0 {
		base = base.Where("project_id = ?", projectID)
	}
	if unitID > 0 {
		base = base.Where("unit_id = ?", unitID)
	}
	if dateFrom != nil {
		base = base.Where("date >= ?", dateFrom)
	}
	if dateTo != nil {
		base = base.Where("date <= ?", dateTo)
	}

	if err := base.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage
	if err := base.
		Preload("Company").
		Preload("Project").
		Preload("Driver").
		Preload("Creator").
		Preload("Details").
		Preload("Details.FailureReason").
		Offset(offset).Limit(perPage).
		Order("date DESC, id DESC").
		Find(&replacements).Error; err != nil {
		return nil, 0, err
	}

	if len(replacements) == 0 {
		return replacements, total, nil
	}

	unitIDs := make([]uint, 0, len(replacements))
	seen := make(map[uint]bool)
	for _, rep := range replacements {
		if !seen[rep.UnitID] {
			unitIDs = append(unitIDs, rep.UnitID)
			seen[rep.UnitID] = true
		}
	}

	var units []*entity.Unit
	if err := r.db.Where("id IN ?", unitIDs).Find(&units).Error; err != nil {
		return replacements, total, nil
	}

	unitMap := make(map[uint]*entity.Unit, len(units))
	for _, u := range units {
		unitMap[u.ID] = u
	}
	for _, rep := range replacements {
		rep.Unit = unitMap[rep.UnitID]
	}

	return replacements, total, nil
}

func (r *replacementRepository) GetLastReplacementByUnitID(unitID uint) (*entity.Replacement, error) {
	var replacement entity.Replacement
	err := r.db.Where("unit_id = ?", unitID).Order("id DESC").First(&replacement).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &replacement, nil
}

func (r *replacementRepository) Update(replacement *entity.Replacement) error {
	return r.db.Save(replacement).Error
}

func (r *replacementRepository) Delete(id uint) error {
	return r.db.Delete(&entity.Replacement{}, id).Error
}
