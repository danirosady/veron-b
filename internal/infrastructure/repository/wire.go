package repository

import (
	"github.com/tms/tyre/internal/domain/repository"
	"gorm.io/gorm"
)

type Repositories struct {
	User        repository.UserRepository
	Company     repository.CompanyRepository
	Project     repository.ProjectRepository
	Unit        repository.UnitRepository
	Driver      repository.DriverRepository
	Tyre        repository.TyreRepository
	Replacement repository.ReplacementRepository
	Master      repository.MasterRepository
}

func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		User:        NewUserRepository(db),
		Company:     NewCompanyRepository(db),
		Project:     NewProjectRepository(db),
		Unit:        NewUnitRepository(db),
		Driver:      NewDriverRepository(db),
		Tyre:        NewTyreRepository(db),
		Replacement: NewReplacementRepository(db),
		Master:      NewMasterRepository(db),
	}
}
