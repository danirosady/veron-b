package entity

import (
	"time"
)

type Replacement struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	CompanyID    uint       `gorm:"not null" json:"company_id"`
	Company      *Company   `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	ProjectID    uint       `gorm:"not null" json:"project_id"`
	Project      *Project   `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
	UnitID       uint       `gorm:"not null" json:"unit_id"`
	Unit         *Unit      `gorm:"foreignKey:UnitID" json:"unit,omitempty"`
	DriverID     uint       `gorm:"not null" json:"driver_id"`
	Driver       *Driver    `gorm:"foreignKey:DriverID" json:"driver,omitempty"`
	Date         time.Time  `gorm:"type:date;not null" json:"date"`
	HMUpdate     float64    `gorm:"type:decimal(15,2);not null;default:0" json:"hm_update"`
	CurrentLifeHM float64   `gorm:"type:decimal(15,2);not null" json:"current_life_hm"`
	HMPlan       float64    `gorm:"type:decimal(15,2);not null" json:"hm_plan"`
	Remarks      string     `gorm:"type:text" json:"remarks,omitempty"`
	CreatedBy    uint       `gorm:"not null" json:"created_by"`
	Creator      *User      `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`

	Details []ReplacementDetail `gorm:"foreignKey:ReplacementID" json:"details,omitempty"`
}

func (Replacement) TableName() string { return "replacements" }
