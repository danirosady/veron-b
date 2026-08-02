package entity

import (
	"time"
)

type Unit struct {
	ID               uint       `gorm:"primaryKey" json:"id"`
	CompanyID        uint       `gorm:"column:company_id;not null" json:"company_id"`
	Company          *Company   `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	ProjectID        uint       `gorm:"column:project_id;not null" json:"project_id"`
	Project          *Project   `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
	UnitID           string     `gorm:"column:unit_id;size:50;uniqueIndex;not null" json:"unit_id"`
	UnitModel       string     `gorm:"column:unit_model;size:255;not null" json:"unit_model"`
	PlateNumber     string     `gorm:"column:plate_number;size:50" json:"plate_number,omitempty"`
	TyreSizeDefault string     `gorm:"column:tyre_size_default;size:50;not null" json:"tyre_size_default"`
	UnitType        string     `gorm:"column:unit_type;size:50;not null;default:'ADT'" json:"unit_type"`
	MaxPosition     int        `gorm:"column:max_position;not null;default:6" json:"max_position"`
	CurrentHM       float64    `gorm:"column:current_hm;type:decimal(15,2);not null;default:0" json:"current_hm"`
	Status          string     `gorm:"column:status;size:20;not null;default:'active'" json:"status"`
	CreatedAt       time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt       time.Time  `gorm:"column:updated_at" json:"updated_at"`
}

func (Unit) TableName() string { return "units" }
