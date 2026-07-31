package usecase

import (
	"context"
	"errors"
	"strings"

	"github.com/tms/tyre/internal/domain/entity"
	"github.com/tms/tyre/internal/domain/repository"
	"github.com/tms/tyre/internal/dto/request"
)

// Common errors used by the unit use case
var (
	ErrUnitNotFound          = errors.New("unit not found")
	ErrUnitIDRequired        = errors.New("unit id is required")
	ErrUnitIDTooLong         = errors.New("unit id must not exceed 50 characters")
	ErrUnitIDExists          = errors.New("unit id already exists")
	ErrUnitModelRequired     = errors.New("unit model is required")
	ErrUnitTyreSizeRequired  = errors.New("tyre size default is required")
	ErrUnitTypeRequired      = errors.New("unit type is required")
	ErrUnitInvalidPosition   = errors.New("max position must be between 1 and 20")
	ErrUnitHasMountedTyres   = errors.New("unit has mounted tyres")
	ErrUnitHasReplacements   = errors.New("unit has replacement history")
	ErrUnitHMDecrease        = errors.New("current HM cannot be lower than existing value")
)

// UnitUseCase handles business logic for unit (vehicle) management.
type UnitUseCase struct {
	unitRepo    repository.UnitRepository
	projectRepo repository.ProjectRepository
	companyRepo repository.CompanyRepository
	tyreRepo    repository.TyreRepository
}

// NewUnitUseCase creates a new UnitUseCase instance.
func NewUnitUseCase(
	unitRepo repository.UnitRepository,
	projectRepo repository.ProjectRepository,
	companyRepo repository.CompanyRepository,
	tyreRepo repository.TyreRepository,
) *UnitUseCase {
	return &UnitUseCase{
		unitRepo:    unitRepo,
		projectRepo: projectRepo,
		companyRepo: companyRepo,
		tyreRepo:    tyreRepo,
	}
}

// List returns a paginated list of units.
func (uc *UnitUseCase) List(ctx context.Context, page, perPage int, companyID, projectID uint, status string) ([]*entity.Unit, int64, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}
	return uc.unitRepo.List(page, perPage, companyID, projectID, status)
}

// GetByID returns a single unit by ID.
func (uc *UnitUseCase) GetByID(ctx context.Context, id uint) (*entity.Unit, error) {
	unit, err := uc.unitRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if unit == nil {
		return nil, ErrUnitNotFound
	}
	return unit, nil
}

// Create creates a new unit.
func (uc *UnitUseCase) Create(ctx context.Context, req *request.CreateUnitRequest) (*entity.Unit, error) {
	unitID := strings.TrimSpace(req.UnitID)
	if unitID == "" {
		return nil, ErrUnitIDRequired
	}
	if len(unitID) > 50 {
		return nil, ErrUnitIDTooLong
	}

	if strings.TrimSpace(req.UnitModel) == "" {
		return nil, ErrUnitModelRequired
	}
	if strings.TrimSpace(req.TyreSizeDefault) == "" {
		return nil, ErrUnitTyreSizeRequired
	}
	if strings.TrimSpace(req.UnitType) == "" {
		return nil, ErrUnitTypeRequired
	}
	if req.MaxPosition < 1 || req.MaxPosition > 20 {
		return nil, ErrUnitInvalidPosition
	}

	company, err := uc.companyRepo.GetByID(req.CompanyID)
	if err != nil {
		return nil, err
	}
	if company == nil {
		return nil, ErrCompanyNotFound
	}

	project, err := uc.projectRepo.GetByID(req.ProjectID)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, ErrProjectNotFound
	}

	existing, err := uc.unitRepo.GetByUnitID(unitID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrUnitIDExists
	}

	unit := &entity.Unit{
		CompanyID:       req.CompanyID,
		ProjectID:       req.ProjectID,
		UnitID:          unitID,
		UnitModel:       req.UnitModel,
		PlateNumber:     req.PlateNumber,
		TyreSizeDefault: req.TyreSizeDefault,
		UnitType:        req.UnitType,
		MaxPosition:     req.MaxPosition,
		CurrentHM:       0,
		Status:          "active",
	}

	if err := uc.unitRepo.Create(unit); err != nil {
		return nil, err
	}

	return unit, nil
}

// Update updates an existing unit.
func (uc *UnitUseCase) Update(ctx context.Context, id uint, req *request.UpdateUnitRequest) (*entity.Unit, error) {
	unit, err := uc.unitRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if unit == nil {
		return nil, ErrUnitNotFound
	}

	if strings.TrimSpace(req.UnitModel) == "" {
		return nil, ErrUnitModelRequired
	}
	if strings.TrimSpace(req.TyreSizeDefault) == "" {
		return nil, ErrUnitTyreSizeRequired
	}
	if strings.TrimSpace(req.UnitType) == "" {
		return nil, ErrUnitTypeRequired
	}
	if req.MaxPosition < 1 || req.MaxPosition > 20 {
		return nil, ErrUnitInvalidPosition
	}

	project, err := uc.projectRepo.GetByID(req.ProjectID)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, ErrProjectNotFound
	}

	unit.ProjectID = req.ProjectID
	unit.UnitModel = req.UnitModel
	unit.PlateNumber = req.PlateNumber
	unit.TyreSizeDefault = req.TyreSizeDefault
	unit.UnitType = req.UnitType
	unit.MaxPosition = req.MaxPosition
	if req.Status != "" {
		unit.Status = req.Status
	}

	if err := uc.unitRepo.Update(unit); err != nil {
		return nil, err
	}

	return unit, nil
}

// Delete removes a unit by ID.
func (uc *UnitUseCase) Delete(ctx context.Context, id uint) error {
	unit, err := uc.unitRepo.GetByID(id)
	if err != nil {
		return err
	}
	if unit == nil {
		return ErrUnitNotFound
	}

	hasMounted, err := uc.unitRepo.HasMountedTyres(id)
	if err != nil {
		return err
	}
	if hasMounted {
		return ErrUnitHasMountedTyres
	}

	hasReplacements, err := uc.unitRepo.HasReplacements(id)
	if err != nil {
		return err
	}
	if hasReplacements {
		return ErrUnitHasReplacements
	}

	return uc.unitRepo.Delete(id)
}

// UpdateHM updates the current HM (hour meter) of a unit.
func (uc *UnitUseCase) UpdateHM(ctx context.Context, id uint, req *request.UpdateUnitHMRequest) (*entity.Unit, error) {
	unit, err := uc.unitRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if unit == nil {
		return nil, ErrUnitNotFound
	}

	if req.CurrentHM < unit.CurrentHM {
		return nil, ErrUnitHMDecrease
	}

	if err := uc.unitRepo.UpdateHM(id, req.CurrentHM); err != nil {
		return nil, err
	}

	unit.CurrentHM = req.CurrentHM
	return unit, nil
}

// GetTyres returns all tyres mounted on the unit.
func (uc *UnitUseCase) GetTyres(ctx context.Context, id uint) ([]*entity.TyreMaster, error) {
	unit, err := uc.unitRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if unit == nil {
		return nil, ErrUnitNotFound
	}
	return uc.tyreRepo.GetByUnitID(id)
}
