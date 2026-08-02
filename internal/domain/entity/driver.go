package entity

import (
	"time"
)

type Driver struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	CompanyID  uint      `gorm:"column:company_id;not null" json:"company_id"`
	Company    *Company  `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Name       string    `gorm:"column:name;size:255;not null" json:"name"`
	EmployeeID string    `gorm:"column:employee_id;size:50;not null" json:"employee_id"`
	Phone      string    `gorm:"column:phone;size:50" json:"phone,omitempty"`
	Status     string    `gorm:"column:status;size:20;not null;default:'active'" json:"status"`
	CreatedAt  time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (Driver) TableName() string { return "drivers" }
