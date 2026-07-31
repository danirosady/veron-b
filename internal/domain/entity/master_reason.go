package entity

import (
	"time"
)

type MasterReason struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:255;not null" json:"name"`
	Status    string    `gorm:"size:20;not null;default:'active'" json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (MasterReason) TableName() string { return "master_reasons" }
