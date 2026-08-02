package entity

import (
	"time"
)

type Project struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	CompanyID uint       `gorm:"column:company_id;not null" json:"company_id"`
	Company   *Company   `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Name      string     `gorm:"column:name;size:255;not null" json:"name"`
	Location  string     `gorm:"column:location;type:text" json:"location,omitempty"`
	StartDate *time.Time `gorm:"column:start_date" json:"start_date,omitempty"`
	EndDate   *time.Time `gorm:"column:end_date" json:"end_date,omitempty"`
	Status    string     `gorm:"column:status;size:20;not null;default:'active'" json:"status"`
	CreatedAt time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at" json:"updated_at"`
}

func (Project) TableName() string { return "projects" }
