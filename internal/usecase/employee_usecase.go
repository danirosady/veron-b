package usecase

import (
	"context"
	"errors"
	"strings"

	"github.com/tms/tyre/internal/domain/entity"
	"github.com/tms/tyre/internal/domain/repository"
	"github.com/tms/tyre/internal/dto/request"
)

// Common errors used by the employee use case
var (
	ErrEmployeeNotFound      = errors.New("employee not found")
	ErrEmployeeNameRequired  = errors.New("name is required")
	ErrEmployeeNameTooLong   = errors.New("name must not exceed 255 characters")
	ErrEmployeeIDRequired    = errors.New("employee id is required")
	ErrEmployeeIDTooLong     = errors.New("employee id must not exceed 50 characters")
	ErrEmployeeIDExists      = errors.New("employee id already exists for this company and role")
	ErrEmployeeRoleRequired   = errors.New("role is required")
	ErrEmployeeRoleInvalid    = errors.New("role must be one of: driver, supervisor, tyreman")
)

// EmployeeUseCase handles business logic for employee management.
type EmployeeUseCase struct {
	employeeRepo repository.EmployeeRepository
	companyRepo repository.CompanyRepository
}

// NewEmployeeUseCase creates a new EmployeeUseCase instance.
func NewEmployeeUseCase(
	employeeRepo repository.EmployeeRepository,
	companyRepo repository.CompanyRepository,
) *EmployeeUseCase {
	return &EmployeeUseCase{
		employeeRepo: employeeRepo,
		companyRepo:  companyRepo,
	}
}

// List returns a paginated list of employees.
func (uc *EmployeeUseCase) List(ctx context.Context, page, perPage int, companyID uint, role, status string) ([]*entity.Employee, int64, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}
	return uc.employeeRepo.List(page, perPage, companyID, role, status)
}

// GetByID returns a single employee by ID.
func (uc *EmployeeUseCase) GetByID(ctx context.Context, id uint) (*entity.Employee, error) {
	employee, err := uc.employeeRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if employee == nil {
		return nil, ErrEmployeeNotFound
	}
	return employee, nil
}

// Create creates a new employee.
func (uc *EmployeeUseCase) Create(ctx context.Context, req *request.CreateEmployeeRequest) (*entity.Employee, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, ErrEmployeeNameRequired
	}
	if len(name) > 255 {
		return nil, ErrEmployeeNameTooLong
	}

	employeeID := strings.TrimSpace(req.EmployeeID)
	if employeeID == "" {
		return nil, ErrEmployeeIDRequired
	}
	if len(employeeID) > 50 {
		return nil, ErrEmployeeIDTooLong
	}

	role := strings.TrimSpace(req.Role)
	if role == "" {
		return nil, ErrEmployeeRoleRequired
	}
	if role != "driver" && role != "supervisor" && role != "tyreman" {
		return nil, ErrEmployeeRoleInvalid
	}

	company, err := uc.companyRepo.GetByID(req.CompanyID)
	if err != nil {
		return nil, err
	}
	if company == nil {
		return nil, ErrCompanyNotFound
	}

	if err := uc.assertEmployeeUnique(req.CompanyID, employeeID, role, 0); err != nil {
		return nil, err
	}

	employee := &entity.Employee{
		CompanyID:  req.CompanyID,
		Name:       name,
		EmployeeID: employeeID,
		Phone:      req.Phone,
		Role:       entity.EmployeeRole(role),
		Status:     "active",
	}

	if err := uc.employeeRepo.Create(employee); err != nil {
		return nil, err
	}

	return employee, nil
}

// Update updates an existing employee.
func (uc *EmployeeUseCase) Update(ctx context.Context, id uint, req *request.UpdateEmployeeRequest) (*entity.Employee, error) {
	employee, err := uc.employeeRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if employee == nil {
		return nil, ErrEmployeeNotFound
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, ErrEmployeeNameRequired
	}
	if len(name) > 255 {
		return nil, ErrEmployeeNameTooLong
	}

	employeeID := strings.TrimSpace(req.EmployeeID)
	if employeeID == "" {
		return nil, ErrEmployeeIDRequired
	}
	if len(employeeID) > 50 {
		return nil, ErrEmployeeIDTooLong
	}

	if err := uc.assertEmployeeUnique(employee.CompanyID, employeeID, string(employee.Role), id); err != nil {
		return nil, err
	}

	employee.Name = name
	employee.EmployeeID = employeeID
	employee.Phone = req.Phone
	if req.Status != "" {
		employee.Status = req.Status
	}

	if err := uc.employeeRepo.Update(employee); err != nil {
		return nil, err
	}

	return employee, nil
}

// Delete removes an employee by ID.
func (uc *EmployeeUseCase) Delete(ctx context.Context, id uint) error {
	employee, err := uc.employeeRepo.GetByID(id)
	if err != nil {
		return err
	}
	if employee == nil {
		return ErrEmployeeNotFound
	}
	return uc.employeeRepo.Delete(id)
}

// assertEmployeeUnique verifies that no other employee in the same company and role uses the given employee id.
func (uc *EmployeeUseCase) assertEmployeeUnique(companyID uint, employeeID, role string, excludeID uint) error {
	employees, _, err := uc.employeeRepo.List(1, 1000, companyID, role, "")
	if err != nil {
		return err
	}
	for _, e := range employees {
		if e.ID == excludeID {
			continue
		}
		if e.EmployeeID == employeeID {
			return ErrEmployeeIDExists
		}
	}
	return nil
}
