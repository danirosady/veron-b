package entity

import (
	"time"
)

type MasterRemark struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	StatusName string         `gorm:"column:status_name;size:100;not null" json:"status_name"`
	Status     string         `gorm:"size:20;not null;default:'active'" json:"status"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
}

func (MasterRemark) TableName() string { return "master_remarks" }
