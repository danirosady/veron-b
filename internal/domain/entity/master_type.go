package entity

import (
	"time"
)

type MasterType struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:50;uniqueIndex;not null" json:"name"`
	Status    string    `gorm:"size:20;not null;default:'active'" json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (MasterType) TableName() string { return "master_types" }
