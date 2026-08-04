package repository

import "github.com/tms/tyre/internal/domain/entity"

type MasterRepository interface {
	// Brands
	ListBrands() ([]*entity.MasterBrand, error)
	GetBrandByID(id uint) (*entity.MasterBrand, error)
	CreateBrand(name string) (*entity.MasterBrand, error)
	UpdateBrand(id uint, name string) (*entity.MasterBrand, error)
	DeleteBrand(id uint) error

	// Sizes
	ListSizes() ([]*entity.MasterSize, error)
	GetSizeByID(id uint) (*entity.MasterSize, error)
	CreateSize(name string) (*entity.MasterSize, error)
	UpdateSize(id uint, name string) (*entity.MasterSize, error)
	DeleteSize(id uint) error

	// Types
	ListTypes() ([]*entity.MasterType, error)
	GetTypeByID(id uint) (*entity.MasterType, error)
	CreateType(name string) (*entity.MasterType, error)
	UpdateType(id uint, name string) (*entity.MasterType, error)
	DeleteType(id uint) error

	// Patterns
	ListPatterns(brandID uint) ([]*entity.MasterPattern, error)
	GetPatternByID(id uint) (*entity.MasterPattern, error)
	CreatePattern(brandID uint, name string) (*entity.MasterPattern, error)
	UpdatePattern(id uint, name string) (*entity.MasterPattern, error)
	DeletePattern(id uint) error

	// Reasons
	ListReasons() ([]*entity.MasterReason, error)
	GetReasonByID(id uint) (*entity.MasterReason, error)
	CreateReason(name string) (*entity.MasterReason, error)
	UpdateReason(id uint, name string) (*entity.MasterReason, error)
	DeleteReason(id uint) error

	// Actions
	ListActions() ([]*entity.MasterAction, error)
	GetActionByID(id uint) (*entity.MasterAction, error)
	CreateAction(name string) (*entity.MasterAction, error)
	UpdateAction(id uint, name string) (*entity.MasterAction, error)
	DeleteAction(id uint) error

	// Remarks
	ListRemarks() ([]*entity.MasterRemark, error)
	GetRemarkByID(id uint) (*entity.MasterRemark, error)
	CreateRemark(name string) (*entity.MasterRemark, error)
	UpdateRemark(id uint, name string) (*entity.MasterRemark, error)
	DeleteRemark(id uint) error

	// Unit Type Configs
	GetUnitTypeConfig(unitType string) (*entity.UnitTypeConfig, error)
	ListUnitTypeConfigs() ([]*entity.UnitTypeConfig, error)
	GetUnitTypeConfigByID(id uint) (*entity.UnitTypeConfig, error)
	CreateUnitTypeConfig(unitType, displayName string, maxPosition int, positionConfigs entity.PositionConfigs) (*entity.UnitTypeConfig, error)
	UpdateUnitTypeConfig(id uint, displayName string, maxPosition int, positionConfigs entity.PositionConfigs) (*entity.UnitTypeConfig, error)
	DeleteUnitTypeConfig(id uint) error
}
