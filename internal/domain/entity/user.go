package entity

import (
	"time"
)

type User struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	Name      string     `gorm:"column:name;size:255;not null" json:"name"`
	Email     string     `gorm:"column:email;size:255;uniqueIndex;not null" json:"email"`
	Password  string     `gorm:"column:password;size:255;not null" json:"-"`
	Role      string     `gorm:"column:role;size:50;not null;default:'admin_company'" json:"role"`
	CompanyID *uint      `gorm:"column:company_id" json:"company_id,omitempty"`
	Company   *Company   `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Status    string     `gorm:"column:status;size:20;not null;default:'active'" json:"status"`
	CreatedAt time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at" json:"updated_at"`
}

func (User) TableName() string { return "users" }

type UserRole string

const (
	RoleSuperadmin    UserRole = "superadmin"
	RoleAdminCompany UserRole = "admin_company"
)

type UserStatus string

const (
	UserStatusActive   UserStatus = "active"
	UserStatusInactive UserStatus = "inactive"
)
