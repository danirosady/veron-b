package usecase

import (
	"context"
	"errors"

	"github.com/tms/tyre/internal/domain/entity"
	"github.com/tms/tyre/internal/domain/repository"
)

// Common errors used by the master use case
var (
	ErrUnitTypeConfigNotFound = errors.New("unit type config not found")
)

// MasterUseCase handles business logic for master data (brands, sizes, types, patterns, reasons, actions, remarks, unit type configs).
type MasterUseCase struct {
	masterRepo repository.MasterRepository
}

// NewMasterUseCase creates a new MasterUseCase instance.
func NewMasterUseCase(masterRepo repository.MasterRepository) *MasterUseCase {
	return &MasterUseCase{masterRepo: masterRepo}
}

// ListBrands returns all active tyre brands.
func (uc *MasterUseCase) ListBrands(ctx context.Context) ([]*entity.MasterBrand, error) {
	return uc.masterRepo.ListBrands()
}

// CreateBrand creates a new tyre brand.
func (uc *MasterUseCase) CreateBrand(ctx context.Context, name string) (*entity.MasterBrand, error) {
	return uc.masterRepo.CreateBrand(name)
}

// UpdateBrand updates an existing tyre brand.
func (uc *MasterUseCase) UpdateBrand(ctx context.Context, id uint, name string) (*entity.MasterBrand, error) {
	return uc.masterRepo.UpdateBrand(id, name)
}

// DeleteBrand deletes a tyre brand.
func (uc *MasterUseCase) DeleteBrand(ctx context.Context, id uint) error {
	return uc.masterRepo.DeleteBrand(id)
}

// ListSizes returns all active tyre sizes.
func (uc *MasterUseCase) ListSizes(ctx context.Context) ([]*entity.MasterSize, error) {
	return uc.masterRepo.ListSizes()
}

// CreateSize creates a new tyre size.
func (uc *MasterUseCase) CreateSize(ctx context.Context, name string) (*entity.MasterSize, error) {
	return uc.masterRepo.CreateSize(name)
}

// UpdateSize updates an existing tyre size.
func (uc *MasterUseCase) UpdateSize(ctx context.Context, id uint, name string) (*entity.MasterSize, error) {
	return uc.masterRepo.UpdateSize(id, name)
}

// DeleteSize deletes a tyre size.
func (uc *MasterUseCase) DeleteSize(ctx context.Context, id uint) error {
	return uc.masterRepo.DeleteSize(id)
}

// ListTypes returns all active tyre types.
func (uc *MasterUseCase) ListTypes(ctx context.Context) ([]*entity.MasterType, error) {
	return uc.masterRepo.ListTypes()
}

// CreateType creates a new tyre type.
func (uc *MasterUseCase) CreateType(ctx context.Context, name string) (*entity.MasterType, error) {
	return uc.masterRepo.CreateType(name)
}

// UpdateType updates an existing tyre type.
func (uc *MasterUseCase) UpdateType(ctx context.Context, id uint, name string) (*entity.MasterType, error) {
	return uc.masterRepo.UpdateType(id, name)
}

// DeleteType deletes a tyre type.
func (uc *MasterUseCase) DeleteType(ctx context.Context, id uint) error {
	return uc.masterRepo.DeleteType(id)
}

// ListPatterns returns active patterns, optionally filtered by brand.
func (uc *MasterUseCase) ListPatterns(ctx context.Context, brandID uint) ([]*entity.MasterPattern, error) {
	return uc.masterRepo.ListPatterns(brandID)
}

// CreatePattern creates a new tyre pattern.
func (uc *MasterUseCase) CreatePattern(ctx context.Context, brandID uint, name string) (*entity.MasterPattern, error) {
	return uc.masterRepo.CreatePattern(brandID, name)
}

// UpdatePattern updates an existing tyre pattern.
func (uc *MasterUseCase) UpdatePattern(ctx context.Context, id uint, name string) (*entity.MasterPattern, error) {
	return uc.masterRepo.UpdatePattern(id, name)
}

// DeletePattern deletes a tyre pattern.
func (uc *MasterUseCase) DeletePattern(ctx context.Context, id uint) error {
	return uc.masterRepo.DeletePattern(id)
}

// ListReasons returns all active replacement reasons.
func (uc *MasterUseCase) ListReasons(ctx context.Context) ([]*entity.MasterReason, error) {
	return uc.masterRepo.ListReasons()
}

// CreateReason creates a new replacement reason.
func (uc *MasterUseCase) CreateReason(ctx context.Context, name string) (*entity.MasterReason, error) {
	return uc.masterRepo.CreateReason(name)
}

// UpdateReason updates an existing replacement reason.
func (uc *MasterUseCase) UpdateReason(ctx context.Context, id uint, name string) (*entity.MasterReason, error) {
	return uc.masterRepo.UpdateReason(id, name)
}

// DeleteReason deletes a replacement reason.
func (uc *MasterUseCase) DeleteReason(ctx context.Context, id uint) error {
	return uc.masterRepo.DeleteReason(id)
}

// ListActions returns all active replacement actions.
func (uc *MasterUseCase) ListActions(ctx context.Context) ([]*entity.MasterAction, error) {
	return uc.masterRepo.ListActions()
}

// CreateAction creates a new replacement action.
func (uc *MasterUseCase) CreateAction(ctx context.Context, name string) (*entity.MasterAction, error) {
	return uc.masterRepo.CreateAction(name)
}

// UpdateAction updates an existing replacement action.
func (uc *MasterUseCase) UpdateAction(ctx context.Context, id uint, name string) (*entity.MasterAction, error) {
	return uc.masterRepo.UpdateAction(id, name)
}

// DeleteAction deletes a replacement action.
func (uc *MasterUseCase) DeleteAction(ctx context.Context, id uint) error {
	return uc.masterRepo.DeleteAction(id)
}

// ListRemarks returns all active remarks.
func (uc *MasterUseCase) ListRemarks(ctx context.Context) ([]*entity.MasterRemark, error) {
	return uc.masterRepo.ListRemarks()
}

// CreateRemark creates a new remark.
func (uc *MasterUseCase) CreateRemark(ctx context.Context, name string) (*entity.MasterRemark, error) {
	return uc.masterRepo.CreateRemark(name)
}

// UpdateRemark updates an existing remark.
func (uc *MasterUseCase) UpdateRemark(ctx context.Context, id uint, name string) (*entity.MasterRemark, error) {
	return uc.masterRepo.UpdateRemark(id, name)
}

// DeleteRemark deletes a remark.
func (uc *MasterUseCase) DeleteRemark(ctx context.Context, id uint) error {
	return uc.masterRepo.DeleteRemark(id)
}

// CreateUnitTypeConfig creates a new unit type configuration.
func (uc *MasterUseCase) CreateUnitTypeConfig(ctx context.Context, unitType, displayName string, maxPosition int, positionConfigs entity.PositionConfigs) (*entity.UnitTypeConfig, error) {
	return uc.masterRepo.CreateUnitTypeConfig(unitType, displayName, maxPosition, positionConfigs)
}

// UpdateUnitTypeConfig updates an existing unit type configuration.
func (uc *MasterUseCase) UpdateUnitTypeConfig(ctx context.Context, id uint, displayName string, maxPosition int, positionConfigs entity.PositionConfigs) (*entity.UnitTypeConfig, error) {
	return uc.masterRepo.UpdateUnitTypeConfig(id, displayName, maxPosition, positionConfigs)
}

// DeleteUnitTypeConfig deletes a unit type configuration.
func (uc *MasterUseCase) DeleteUnitTypeConfig(ctx context.Context, id uint) error {
	return uc.masterRepo.DeleteUnitTypeConfig(id)
}

// GetUnitTypeConfig returns the position configuration for a given unit type as a parsed array.
func (uc *MasterUseCase) GetUnitTypeConfig(ctx context.Context, unitType string) ([]entity.PositionConfig, error) {
	config, err := uc.masterRepo.GetUnitTypeConfig(unitType)
	if err != nil {
		return nil, err
	}
	if config == nil {
		return nil, ErrUnitTypeConfigNotFound
	}

	return []entity.PositionConfig(config.PositionConfig), nil
}

// ListUnitTypeConfigs returns all unit type configurations.
func (uc *MasterUseCase) ListUnitTypeConfigs(ctx context.Context) ([]*entity.UnitTypeConfig, error) {
	return uc.masterRepo.ListUnitTypeConfigs()
}
