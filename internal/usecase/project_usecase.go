package usecase

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/tms/tyre/internal/domain/entity"
	"github.com/tms/tyre/internal/domain/repository"
	"github.com/tms/tyre/internal/dto/request"
)

// Common errors used by the project use case
var (
	ErrProjectNameRequired = errors.New("project name is required")
	ErrProjectNameTooLong  = errors.New("project name must not exceed 255 characters")
	ErrProjectNotFound     = errors.New("project not found")
	ErrProjectHasUnits     = errors.New("project has active units")
	ErrProjectInvalidDates = errors.New("start date must be before or equal to end date")
)

// ProjectUseCase handles business logic for project management.
type ProjectUseCase struct {
	projectRepo repository.ProjectRepository
	companyRepo repository.CompanyRepository
	unitRepo    repository.UnitRepository
}

// NewProjectUseCase creates a new ProjectUseCase instance.
func NewProjectUseCase(
	projectRepo repository.ProjectRepository,
	companyRepo repository.CompanyRepository,
	unitRepo repository.UnitRepository,
) *ProjectUseCase {
	return &ProjectUseCase{
		projectRepo: projectRepo,
		companyRepo: companyRepo,
		unitRepo:    unitRepo,
	}
}

// List returns a paginated list of projects.
func (uc *ProjectUseCase) List(ctx context.Context, page, perPage int, companyID uint, status string) ([]*entity.Project, int64, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}
	return uc.projectRepo.List(page, perPage, companyID, status)
}

// GetByID returns a single project by ID.
func (uc *ProjectUseCase) GetByID(ctx context.Context, id uint) (*entity.Project, error) {
	project, err := uc.projectRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, ErrProjectNotFound
	}
	return project, nil
}

// Create creates a new project.
func (uc *ProjectUseCase) Create(ctx context.Context, req *request.CreateProjectRequest) (*entity.Project, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, ErrProjectNameRequired
	}
	if len(name) > 255 {
		return nil, ErrProjectNameTooLong
	}

	company, err := uc.companyRepo.GetByID(req.CompanyID)
	if err != nil {
		return nil, err
	}
	if company == nil {
		return nil, ErrCompanyNotFound
	}

	var startDate *time.Time
	if strings.TrimSpace(req.StartDate) != "" {
		t, err := parseDate(req.StartDate)
		if err != nil {
			return nil, err
		}
		startDate = &t
	}

	var endDate *time.Time
	if strings.TrimSpace(req.EndDate) != "" {
		t, err := parseDate(req.EndDate)
		if err != nil {
			return nil, err
		}
		endDate = &t
	}

	if startDate != nil && endDate != nil && startDate.After(*endDate) {
		return nil, ErrProjectInvalidDates
	}

	project := &entity.Project{
		CompanyID: req.CompanyID,
		Name:      name,
		Location:  req.Location,
		StartDate: startDate,
		EndDate:   endDate,
		Status:    "active",
	}

	if err := uc.projectRepo.Create(project); err != nil {
		return nil, err
	}

	return project, nil
}

// Update updates an existing project.
func (uc *ProjectUseCase) Update(ctx context.Context, id uint, req *request.UpdateProjectRequest) (*entity.Project, error) {
	project, err := uc.projectRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, ErrProjectNotFound
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, ErrProjectNameRequired
	}
	if len(name) > 255 {
		return nil, ErrProjectNameTooLong
	}

	var startDate *time.Time
	if strings.TrimSpace(req.StartDate) != "" {
		t, err := parseDate(req.StartDate)
		if err != nil {
			return nil, err
		}
		startDate = &t
	}

	var endDate *time.Time
	if strings.TrimSpace(req.EndDate) != "" {
		t, err := parseDate(req.EndDate)
		if err != nil {
			return nil, err
		}
		endDate = &t
	}

	if startDate != nil && endDate != nil && startDate.After(*endDate) {
		return nil, ErrProjectInvalidDates
	}

	project.Name = name
	project.Location = req.Location
	project.StartDate = startDate
	project.EndDate = endDate
	if req.Status != "" {
		project.Status = req.Status
	}

	if err := uc.projectRepo.Update(project); err != nil {
		return nil, err
	}

	return project, nil
}

// Delete removes a project by ID.
func (uc *ProjectUseCase) Delete(ctx context.Context, id uint) error {
	project, err := uc.projectRepo.GetByID(id)
	if err != nil {
		return err
	}
	if project == nil {
		return ErrProjectNotFound
	}

	hasUnits, err := uc.projectRepo.HasActiveUnits(id)
	if err != nil {
		return err
	}
	if hasUnits {
		return ErrProjectHasUnits
	}

	return uc.projectRepo.Delete(id)
}

// parseDate parses a date string in YYYY-MM-DD format (date only).
func parseDate(s string) (time.Time, error) {
	return time.Parse("2006-01-02", s)
}