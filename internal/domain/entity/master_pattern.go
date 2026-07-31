package entity

import (
	"time"
)

type MasterPattern struct {
	ID        uint         `gorm:"primaryKey" json:"id"`
	BrandID   uint         `gorm:"not null" json:"brand_id"`
	Brand     *MasterBrand `gorm:"foreignKey:BrandID" json:"brand,omitempty"`
	Name      string       `gorm:"size:100;not null" json:"name"`
	Status    string       `gorm:"size:20;not null;default:'active'" json:"status"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

func (MasterPattern) TableName() string { return "master_patterns" }
