package usecase

import (
	"context"
	"errors"
	"strings"

	"github.com/tms/tyre/internal/domain/entity"
	"github.com/tms/tyre/internal/domain/repository"
	"github.com/tms/tyre/internal/dto/request"
)

// Common errors used by the driver use case
var (
	ErrDriverNotFound      = errors.New("driver not found")
	ErrDriverNameRequired  = errors.New("driver name is required")
	ErrDriverNameTooLong   = errors.New("driver name must not exceed 255 characters")
	ErrDriverEmployeeReq   = errors.New("employee id is required")
	ErrDriverEmployeeLong  = errors.New("employee id must not exceed 50 characters")
	ErrDriverEmployeeExist = errors.New("employee id already exists for this company")
)

// DriverUseCase handles business logic for driver management.
type DriverUseCase struct {
	driverRepo repository.DriverRepository
	companyRepo repository.CompanyRepository
}

// NewDriverUseCase creates a new DriverUseCase instance.
func NewDriverUseCase(
	driverRepo repository.DriverRepository,
	companyRepo repository.CompanyRepository,
) *DriverUseCase {
	return &DriverUseCase{
		driverRepo:  driverRepo,
		companyRepo: companyRepo,
	}
}

// List returns a paginated list of drivers.
func (uc *DriverUseCase) List(ctx context.Context, page, perPage int, companyID uint, status string) ([]*entity.Driver, int64, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}
	return uc.driverRepo.List(page, perPage, companyID, status)
}

// GetByID returns a single driver by ID.
func (uc *DriverUseCase) GetByID(ctx context.Context, id uint) (*entity.Driver, error) {
	driver, err := uc.driverRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if driver == nil {
		return nil, ErrDriverNotFound
	}
	return driver, nil
}

// Create creates a new driver.
func (uc *DriverUseCase) Create(ctx context.Context, req *request.CreateDriverRequest) (*entity.Driver, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, ErrDriverNameRequired
	}
	if len(name) > 255 {
		return nil, ErrDriverNameTooLong
	}

	employeeID := strings.TrimSpace(req.EmployeeID)
	if employeeID == "" {
		return nil, ErrDriverEmployeeReq
	}
	if len(employeeID) > 50 {
		return nil, ErrDriverEmployeeLong
	}

	company, err := uc.companyRepo.GetByID(req.CompanyID)
	if err != nil {
		return nil, err
	}
	if company == nil {
		return nil, ErrCompanyNotFound
	}

	if err := uc.assertEmployeeUnique(req.CompanyID, employeeID, 0); err != nil {
		return nil, err
	}

	driver := &entity.Driver{
		CompanyID:  req.CompanyID,
		Name:       name,
		EmployeeID: employeeID,
		Phone:      req.Phone,
		Status:     "active",
	}

	if err := uc.driverRepo.Create(driver); err != nil {
		return nil, err
	}

	return driver, nil
}

// Update updates an existing driver.
func (uc *DriverUseCase) Update(ctx context.Context, id uint, req *request.UpdateDriverRequest) (*entity.Driver, error) {
	driver, err := uc.driverRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if driver == nil {
		return nil, ErrDriverNotFound
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, ErrDriverNameRequired
	}
	if len(name) > 255 {
		return nil, ErrDriverNameTooLong
	}

	employeeID := strings.TrimSpace(req.EmployeeID)
	if employeeID == "" {
		return nil, ErrDriverEmployeeReq
	}
	if len(employeeID) > 50 {
		return nil, ErrDriverEmployeeLong
	}

	if err := uc.assertEmployeeUnique(driver.CompanyID, employeeID, id); err != nil {
		return nil, err
	}

	driver.Name = name
	driver.EmployeeID = employeeID
	driver.Phone = req.Phone
	if req.Status != "" {
		driver.Status = req.Status
	}

	if err := uc.driverRepo.Update(driver); err != nil {
		return nil, err
	}

	return driver, nil
}

// Delete removes a driver by ID.
func (uc *DriverUseCase) Delete(ctx context.Context, id uint) error {
	driver, err := uc.driverRepo.GetByID(id)
	if err != nil {
		return err
	}
	if driver == nil {
		return ErrDriverNotFound
	}
	return uc.driverRepo.Delete(id)
}

// assertEmployeeUnique verifies that no other driver in the same company uses the given employee id.
func (uc *DriverUseCase) assertEmployeeUnique(companyID uint, employeeID string, excludeID uint) error {
	drivers, _, err := uc.driverRepo.List(1, 1000, companyID, "")
	if err != nil {
		return err
	}
	for _, d := range drivers {
		if d.ID == excludeID {
			continue
		}
		if d.EmployeeID == employeeID {
			return ErrDriverEmployeeExist
		}
	}
	return nil
}