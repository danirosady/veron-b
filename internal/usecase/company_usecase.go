package usecase

import (
	"context"
	"errors"
	"net/mail"
	"strings"

	"github.com/tms/tyre/internal/domain/entity"
	"github.com/tms/tyre/internal/domain/repository"
	"github.com/tms/tyre/internal/dto/request"
)

// Common errors used by the company use case
var (
	ErrCompanyNameRequired = errors.New("company name is required")
	ErrCompanyNameTooLong  = errors.New("company name must not exceed 255 characters")
	ErrCompanyNameExists   = errors.New("company name already exists")
	ErrCompanyEmailInvalid = errors.New("company email is invalid")
	ErrCompanyHasActive    = errors.New("company has active projects, units, or drivers")
)

// CompanyUseCase handles business logic for company management.
type CompanyUseCase struct {
	companyRepo repository.CompanyRepository
	projectRepo repository.ProjectRepository
	unitRepo    repository.UnitRepository
	driverRepo  repository.DriverRepository
}

// NewCompanyUseCase creates a new CompanyUseCase instance.
func NewCompanyUseCase(
	companyRepo repository.CompanyRepository,
	projectRepo repository.ProjectRepository,
	unitRepo repository.UnitRepository,
	driverRepo repository.DriverRepository,
) *CompanyUseCase {
	return &CompanyUseCase{
		companyRepo: companyRepo,
		projectRepo: projectRepo,
		unitRepo:    unitRepo,
		driverRepo:  driverRepo,
	}
}

// List returns a paginated list of companies.
func (uc *CompanyUseCase) List(ctx context.Context, page, perPage int, status string) ([]*entity.Company, int64, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}
	return uc.companyRepo.List(page, perPage, status)
}

// GetByID returns a single company by ID.
func (uc *CompanyUseCase) GetByID(ctx context.Context, id uint) (*entity.Company, error) {
	company, err := uc.companyRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if company == nil {
		return nil, ErrCompanyNotFound
	}
	return company, nil
}

// Create creates a new company.
func (uc *CompanyUseCase) Create(ctx context.Context, req *request.CreateCompanyRequest) (*entity.Company, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, ErrCompanyNameRequired
	}
	if len(name) > 255 {
		return nil, ErrCompanyNameTooLong
	}

	if req.Email != "" {
		if _, err := mail.ParseAddress(req.Email); err != nil {
			return nil, ErrCompanyEmailInvalid
		}
	}

	existing, err := uc.companyRepo.GetByName(name)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrCompanyNameExists
	}

	company := &entity.Company{
		Name:          name,
		Address:       req.Address,
		ContactPerson: req.ContactPerson,
		Phone:         req.Phone,
		Email:         req.Email,
		Status:        "active",
	}

	if err := uc.companyRepo.Create(company); err != nil {
		return nil, err
	}

	return company, nil
}

// Update updates an existing company.
func (uc *CompanyUseCase) Update(ctx context.Context, id uint, req *request.UpdateCompanyRequest) (*entity.Company, error) {
	company, err := uc.companyRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if company == nil {
		return nil, ErrCompanyNotFound
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, ErrCompanyNameRequired
	}
	if len(name) > 255 {
		return nil, ErrCompanyNameTooLong
	}

	if req.Email != "" {
		if _, err := mail.ParseAddress(req.Email); err != nil {
			return nil, ErrCompanyEmailInvalid
		}
	}

	existing, err := uc.companyRepo.GetByName(name)
	if err != nil {
		return nil, err
	}
	if existing != nil && existing.ID != id {
		return nil, ErrCompanyNameExists
	}

	company.Name = name
	company.Address = req.Address
	company.ContactPerson = req.ContactPerson
	company.Phone = req.Phone
	company.Email = req.Email
	if req.Status != "" {
		company.Status = req.Status
	}

	if err := uc.companyRepo.Update(company); err != nil {
		return nil, err
	}

	return company, nil
}

// Delete removes a company by ID.
func (uc *CompanyUseCase) Delete(ctx context.Context, id uint) error {
	company, err := uc.companyRepo.GetByID(id)
	if err != nil {
		return err
	}
	if company == nil {
		return ErrCompanyNotFound
	}

	hasActive, err := uc.companyRepo.HasActiveData(id)
	if err != nil {
		return err
	}
	if hasActive {
		return ErrCompanyHasActive
	}

	return uc.companyRepo.Delete(id)
}
