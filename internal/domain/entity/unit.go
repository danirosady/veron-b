package entity

import (
	"time"
)

type Unit struct {
	ID              uint       `gorm:"primaryKey" json:"id"`
	CompanyID       uint       `gorm:"not null" json:"company_id"`
	Company         *Company   `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	ProjectID       uint       `gorm:"not null" json:"project_id"`
	Project         *Project   `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
	UnitID          string     `gorm:"size:50;uniqueIndex;not null" json:"unit_id"`
	UnitModel       string     `gorm:"size:255;not null" json:"unit_model"`
	PlateNumber     string     `gorm:"size:50" json:"plate_number,omitempty"`
	TyreSizeDefault string     `gorm:"size:50;not null" json:"tyre_size_default"`
	UnitType        string     `gorm:"size:50;not null;default:'ADT'" json:"unit_type"`
	MaxPosition     int        `gorm:"not null;default:6" json:"max_position"`
	CurrentHM       float64    `gorm:"type:decimal(15,2);not null;default:0" json:"current_hm"`
	Status          string     `gorm:"size:20;not null;default:'active'" json:"status"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

func (Unit) TableName() string { return "units" }
