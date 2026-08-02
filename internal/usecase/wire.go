package usecase

import (
	"github.com/tms/tyre/internal/infrastructure/repository"
)

type UseCases struct {
	Auth    *AuthUseCase
	Company *CompanyUseCase
	Project *ProjectUseCase
	Unit    *UnitUseCase
	Driver  *DriverUseCase
	Tyre    *TyreUseCase
	Master  *MasterUseCase
}

func NewUseCases(repos *repository.Repositories) *UseCases {
	return &UseCases{
		Company: NewCompanyUseCase(repos.Company, repos.Project, repos.Unit, repos.Driver),
		Project: NewProjectUseCase(repos.Project, repos.Company, repos.Unit),
		Unit:    NewUnitUseCase(repos.Unit, repos.Project, repos.Company, repos.Tyre, repos.Master),
		Driver:  NewDriverUseCase(repos.Driver, repos.Company),
		Tyre:    NewTyreUseCase(repos.Tyre, repos.Company, repos.Master),
		Master:  NewMasterUseCase(repos.Master),
	}
}