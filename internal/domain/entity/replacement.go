package entity

import (
	"time"
)

type Replacement struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	CompanyID   uint       `gorm:"column:company_id;not null" json:"company_id"`
	Company     *Company   `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	ProjectID   uint       `gorm:"column:project_id;not null" json:"project_id"`
	Project     *Project   `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
	UnitID      uint       `gorm:"column:unit_id;not null" json:"unit_id"`
	Unit        *Unit      `gorm:"foreignKey:UnitID;references:ID" json:"unit,omitempty"`
	DriverID    uint       `gorm:"column:driver_id;not null" json:"driver_id"`
	Driver      *Driver    `gorm:"foreignKey:DriverID" json:"driver,omitempty"`
	Date        time.Time  `gorm:"column:date;type:date;not null" json:"date"`
	HMUpdate    float64    `gorm:"column:hm_update;type:decimal(15,2);not null;default:0" json:"hm_update"`
	CurrentLifeHM float64 `gorm:"column:current_life_hm;type:decimal(15,2);not null" json:"current_life_hm"`
	HMPlan      float64    `gorm:"column:hm_plan;type:decimal(15,2);not null" json:"hm_plan"`
	Remarks     string     `gorm:"column:remarks;type:text" json:"remarks,omitempty"`
	CreatedBy   uint       `gorm:"column:created_by;not null" json:"created_by"`
	Creator     *User      `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
	CreatedAt   time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"column:updated_at" json:"updated_at"`

	Details []ReplacementDetail `gorm:"foreignKey:ReplacementID" json:"details,omitempty"`
}

func (Replacement) TableName() string { return "replacements" }
