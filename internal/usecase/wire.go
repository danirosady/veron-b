package usecase

import (
	"github.com/tms/tyre/internal/infrastructure/repository"
)

type UseCases struct {
	Auth      *AuthUseCase
	Company   *CompanyUseCase
	Project   *ProjectUseCase
	Unit      *UnitUseCase
	Employee  *EmployeeUseCase
	Tyre      *TyreUseCase
	Master    *MasterUseCase
	Replacement *ReplacementUseCase
}

func NewUseCases(repos *repository.Repositories) *UseCases {
	return &UseCases{
		Company:   NewCompanyUseCase(repos.Company, repos.Project, repos.Unit, repos.Employee),
		Project:  NewProjectUseCase(repos.Project, repos.Company, repos.Unit),
		Unit:     NewUnitUseCase(repos.Unit, repos.Project, repos.Company, repos.Tyre, repos.Master),
		Employee: NewEmployeeUseCase(repos.Employee, repos.Company),
		Tyre:     NewTyreUseCase(repos.Tyre, repos.Company, repos.Master),
		Master:   NewMasterUseCase(repos.Master),
	}
}
