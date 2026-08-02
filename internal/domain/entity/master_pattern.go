package entity

import (
	"time"
)

type MasterPattern struct {
	ID        uint        `gorm:"primaryKey" json:"id"`
	BrandID   uint        `gorm:"column:brand_id;not null" json:"brand_id"`
	Brand     *MasterBrand `gorm:"foreignKey:BrandID" json:"brand,omitempty"`
	Name      string      `gorm:"column:name;size:100;not null" json:"name"`
	Status    string      `gorm:"column:status;size:20;not null;default:'active'" json:"status"`
	CreatedAt time.Time   `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time   `gorm:"column:updated_at" json:"updated_at"`
}

func (MasterPattern) TableName() string { return "master_patterns" }
