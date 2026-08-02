package entity

import (
	"time"
)

type MasterSize struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"column:name;size:50;uniqueIndex;not null" json:"name"`
	Status    string    `gorm:"column:status;size:20;not null;default:'active'" json:"status"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (MasterSize) TableName() string { return "master_sizes" }
