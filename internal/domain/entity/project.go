package entity

import (
	"time"
)

type Project struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	CompanyID uint       `gorm:"not null" json:"company_id"`
	Company   *Company   `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Name      string     `gorm:"size:255;not null" json:"name"`
	Location  string     `gorm:"type:text" json:"location,omitempty"`
	StartDate *time.Time `json:"start_date,omitempty"`
	EndDate   *time.Time `json:"end_date,omitempty"`
	Status    string     `gorm:"size:20;not null;default:'active'" json:"status"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func (Project) TableName() string { return "projects" }
