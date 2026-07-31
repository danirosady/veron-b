package entity

import (
	"time"
)

type Driver struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	CompanyID  uint      `gorm:"not null" json:"company_id"`
	Company    *Company  `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Name       string    `gorm:"size:255;not null" json:"name"`
	EmployeeID string    `gorm:"size:50;not null" json:"employee_id"`
	Phone      string    `gorm:"size:50" json:"phone,omitempty"`
	Status     string    `gorm:"size:20;not null;default:'active'" json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (Driver) TableName() string { return "drivers" }
