package entity

import (
	"time"
)

type UnitTypeConfig struct {
	ID             uint            `gorm:"primaryKey" json:"id"`
	UnitType       string          `gorm:"size:50;uniqueIndex;not null" json:"unit_type"`
	DisplayName    string          `gorm:"size:100;not null" json:"display_name"`
	MaxPosition    int             `gorm:"not null" json:"max_position"`
	PositionConfig PositionConfigs `gorm:"type:jsonb" json:"position_config"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
}

func (UnitTypeConfig) TableName() string { return "unit_type_configs" }
