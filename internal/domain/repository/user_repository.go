package repository

import "github.com/tms/tyre/internal/domain/entity"

type UserRepository interface {
	Create(user *entity.User) error
	GetByID(id uint) (*entity.User, error)
	GetByEmail(email string) (*entity.User, error)
	Update(user *entity.User) error
	Delete(id uint) error
	List(page, perPage int) ([]*entity.User, int64, error)
	GetByCompanyID(companyID uint) ([]*entity.User, error)
}
