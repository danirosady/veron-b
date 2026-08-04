package repository

import (
	"github.com/tms/tyre/internal/domain/entity"
	"github.com/tms/tyre/internal/domain/repository"
	"gorm.io/gorm"
)

type masterRepository struct {
	*BaseRepository
}

func NewMasterRepository(db *gorm.DB) repository.MasterRepository {
	return &masterRepository{NewBaseRepository(db)}
}

func (r *masterRepository) ListBrands() ([]*entity.MasterBrand, error) {
	var brands []*entity.MasterBrand
	err := r.db.Where("status = 'active'").Order("name ASC").Find(&brands).Error
	return brands, err
}

func (r *masterRepository) GetBrandByID(id uint) (*entity.MasterBrand, error) {
	var brand entity.MasterBrand
	err := r.db.First(&brand, id).Error
	if err != nil {
		return nil, err
	}
	return &brand, nil
}

func (r *masterRepository) CreateBrand(name string) (*entity.MasterBrand, error) {
	brand := &entity.MasterBrand{Name: name}
	err := r.db.Create(brand).Error
	return brand, err
}

func (r *masterRepository) UpdateBrand(id uint, name string) (*entity.MasterBrand, error) {
	var brand entity.MasterBrand
	if err := r.db.First(&brand, id).Error; err != nil {
		return nil, err
	}
	brand.Name = name
	if err := r.db.Save(&brand).Error; err != nil {
		return nil, err
	}
	return &brand, nil
}

func (r *masterRepository) DeleteBrand(id uint) error {
	return r.db.Delete(&entity.MasterBrand{}, id).Error
}

func (r *masterRepository) ListSizes() ([]*entity.MasterSize, error) {
	var sizes []*entity.MasterSize
	err := r.db.Where("status = 'active'").Order("name ASC").Find(&sizes).Error
	return sizes, err
}

func (r *masterRepository) GetSizeByID(id uint) (*entity.MasterSize, error) {
	var size entity.MasterSize
	err := r.db.First(&size, id).Error
	if err != nil {
		return nil, err
	}
	return &size, nil
}

func (r *masterRepository) CreateSize(name string) (*entity.MasterSize, error) {
	size := &entity.MasterSize{Name: name}
	err := r.db.Create(size).Error
	return size, err
}

func (r *masterRepository) UpdateSize(id uint, name string) (*entity.MasterSize, error) {
	var size entity.MasterSize
	if err := r.db.First(&size, id).Error; err != nil {
		return nil, err
	}
	size.Name = name
	if err := r.db.Save(&size).Error; err != nil {
		return nil, err
	}
	return &size, nil
}

func (r *masterRepository) DeleteSize(id uint) error {
	return r.db.Delete(&entity.MasterSize{}, id).Error
}

func (r *masterRepository) ListTypes() ([]*entity.MasterType, error) {
	var types []*entity.MasterType
	err := r.db.Where("status = 'active'").Order("name ASC").Find(&types).Error
	return types, err
}

func (r *masterRepository) GetTypeByID(id uint) (*entity.MasterType, error) {
	var t entity.MasterType
	err := r.db.First(&t, id).Error
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *masterRepository) CreateType(name string) (*entity.MasterType, error) {
	t := &entity.MasterType{Name: name}
	err := r.db.Create(t).Error
	return t, err
}

func (r *masterRepository) UpdateType(id uint, name string) (*entity.MasterType, error) {
	var t entity.MasterType
	if err := r.db.First(&t, id).Error; err != nil {
		return nil, err
	}
	t.Name = name
	if err := r.db.Save(&t).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *masterRepository) DeleteType(id uint) error {
	return r.db.Delete(&entity.MasterType{}, id).Error
}

func (r *masterRepository) ListPatterns(brandID uint) ([]*entity.MasterPattern, error) {
	var patterns []*entity.MasterPattern
	query := r.db.Where("status = 'active'")
	if brandID > 0 {
		query = query.Where("brand_id = ?", brandID)
	}
	err := query.Preload("Brand").Order("name ASC").Find(&patterns).Error
	return patterns, err
}

func (r *masterRepository) GetPatternByID(id uint) (*entity.MasterPattern, error) {
	var pattern entity.MasterPattern
	err := r.db.Preload("Brand").First(&pattern, id).Error
	if err != nil {
		return nil, err
	}
	return &pattern, nil
}

func (r *masterRepository) CreatePattern(brandID uint, name string) (*entity.MasterPattern, error) {
	pattern := &entity.MasterPattern{BrandID: brandID, Name: name}
	err := r.db.Create(pattern).Error
	return pattern, err
}

func (r *masterRepository) UpdatePattern(id uint, name string) (*entity.MasterPattern, error) {
	var pattern entity.MasterPattern
	if err := r.db.First(&pattern, id).Error; err != nil {
		return nil, err
	}
	pattern.Name = name
	if err := r.db.Save(&pattern).Error; err != nil {
		return nil, err
	}
	return &pattern, nil
}

func (r *masterRepository) DeletePattern(id uint) error {
	return r.db.Delete(&entity.MasterPattern{}, id).Error
}

func (r *masterRepository) ListReasons() ([]*entity.MasterReason, error) {
	var reasons []*entity.MasterReason
	err := r.db.Where("status = 'active'").Order("name ASC").Find(&reasons).Error
	return reasons, err
}

func (r *masterRepository) GetReasonByID(id uint) (*entity.MasterReason, error) {
	var reason entity.MasterReason
	err := r.db.First(&reason, id).Error
	if err != nil {
		return nil, err
	}
	return &reason, nil
}

func (r *masterRepository) CreateReason(name string) (*entity.MasterReason, error) {
	reason := &entity.MasterReason{Name: name}
	err := r.db.Create(reason).Error
	return reason, err
}

func (r *masterRepository) UpdateReason(id uint, name string) (*entity.MasterReason, error) {
	var reason entity.MasterReason
	if err := r.db.First(&reason, id).Error; err != nil {
		return nil, err
	}
	reason.Name = name
	if err := r.db.Save(&reason).Error; err != nil {
		return nil, err
	}
	return &reason, nil
}

func (r *masterRepository) DeleteReason(id uint) error {
	return r.db.Delete(&entity.MasterReason{}, id).Error
}

func (r *masterRepository) GetActionByID(id uint) (*entity.MasterAction, error) {
	var action entity.MasterAction
	err := r.db.First(&action, id).Error
	if err != nil {
		return nil, err
	}
	return &action, nil
}

func (r *masterRepository) CreateAction(name string) (*entity.MasterAction, error) {
	action := &entity.MasterAction{Name: name}
	err := r.db.Create(action).Error
	return action, err
}

func (r *masterRepository) UpdateAction(id uint, name string) (*entity.MasterAction, error) {
	var action entity.MasterAction
	if err := r.db.First(&action, id).Error; err != nil {
		return nil, err
	}
	action.Name = name
	if err := r.db.Save(&action).Error; err != nil {
		return nil, err
	}
	return &action, nil
}

func (r *masterRepository) DeleteAction(id uint) error {
	return r.db.Delete(&entity.MasterAction{}, id).Error
}

func (r *masterRepository) GetRemarkByID(id uint) (*entity.MasterRemark, error) {
	var remark entity.MasterRemark
	err := r.db.First(&remark, id).Error
	if err != nil {
		return nil, err
	}
	return &remark, nil
}

func (r *masterRepository) CreateRemark(name string) (*entity.MasterRemark, error) {
	remark := &entity.MasterRemark{Name: name}
	err := r.db.Create(remark).Error
	return remark, err
}

func (r *masterRepository) UpdateRemark(id uint, name string) (*entity.MasterRemark, error) {
	var remark entity.MasterRemark
	if err := r.db.First(&remark, id).Error; err != nil {
		return nil, err
	}
	remark.Name = name
	if err := r.db.Save(&remark).Error; err != nil {
		return nil, err
	}
	return &remark, nil
}

func (r *masterRepository) DeleteRemark(id uint) error {
	return r.db.Delete(&entity.MasterRemark{}, id).Error
}

func (r *masterRepository) ListActions() ([]*entity.MasterAction, error) {
	var actions []*entity.MasterAction
	err := r.db.Where("status = 'active'").Order("name ASC").Find(&actions).Error
	return actions, err
}

func (r *masterRepository) ListRemarks() ([]*entity.MasterRemark, error) {
	var remarks []*entity.MasterRemark
	err := r.db.Where("status = 'active'").Order("name ASC").Find(&remarks).Error
	return remarks, err
}

func (r *masterRepository) GetUnitTypeConfig(unitType string) (*entity.UnitTypeConfig, error) {
	var config entity.UnitTypeConfig
	err := r.db.Where("unit_type = ?", unitType).First(&config).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func (r *masterRepository) ListUnitTypeConfigs() ([]*entity.UnitTypeConfig, error) {
	var configs []*entity.UnitTypeConfig
	err := r.db.Order("display_name ASC").Find(&configs).Error
	return configs, err
}

func (r *masterRepository) GetUnitTypeConfigByID(id uint) (*entity.UnitTypeConfig, error) {
	var config entity.UnitTypeConfig
	err := r.db.First(&config, id).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func (r *masterRepository) CreateUnitTypeConfig(unitType, displayName string, maxPosition int, positionConfigs entity.PositionConfigs) (*entity.UnitTypeConfig, error) {
	config := &entity.UnitTypeConfig{
		UnitType:       unitType,
		DisplayName:    displayName,
		MaxPosition:    maxPosition,
		PositionConfig: positionConfigs,
	}
	err := r.db.Create(config).Error
	return config, err
}

func (r *masterRepository) UpdateUnitTypeConfig(id uint, displayName string, maxPosition int, positionConfigs entity.PositionConfigs) (*entity.UnitTypeConfig, error) {
	var config entity.UnitTypeConfig
	if err := r.db.First(&config, id).Error; err != nil {
		return nil, err
	}
	config.DisplayName = displayName
	config.MaxPosition = maxPosition
	config.PositionConfig = positionConfigs
	if err := r.db.Save(&config).Error; err != nil {
		return nil, err
	}
	return &config, nil
}

func (r *masterRepository) DeleteUnitTypeConfig(id uint) error {
	return r.db.Delete(&entity.UnitTypeConfig{}, id).Error
}
