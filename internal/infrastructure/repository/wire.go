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
	Employee    repository.EmployeeRepository
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
		Employee:    NewEmployeeRepository(db),
		Tyre:        NewTyreRepository(db),
		Replacement: NewReplacementRepository(db),
		Master:      NewMasterRepository(db),
	}
}
