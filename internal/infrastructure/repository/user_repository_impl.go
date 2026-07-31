package repository

import (
	"errors"

	"github.com/tms/tyre/internal/domain/entity"
	"github.com/tms/tyre/internal/domain/repository"
	"gorm.io/gorm"
)

type userRepository struct {
	*BaseRepository
}

func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &userRepository{NewBaseRepository(db)}
}

func (r *userRepository) Create(user *entity.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) GetByID(id uint) (*entity.User, error) {
	var user entity.User
	err := r.db.Preload("Company").First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(email string) (*entity.User, error) {
	var user entity.User
	err := r.db.Preload("Company").Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(user *entity.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) Delete(id uint) error {
	return r.db.Delete(&entity.User{}, id).Error
}

func (r *userRepository) List(page, perPage int) ([]*entity.User, int64, error) {
	var users []*entity.User
	var total int64

	offset := (page - 1) * perPage

	if err := r.db.Model(&entity.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Preload("Company").Offset(offset).Limit(perPage).Order("id DESC").Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *userRepository) GetByCompanyID(companyID uint) ([]*entity.User, error) {
	var users []*entity.User
	err := r.db.Where("company_id = ?", companyID).Find(&users).Error
	return users, err
}
