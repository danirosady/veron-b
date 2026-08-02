package entity

import (
	"time"
)

type ReplacementDetail struct {
	ID            uint        `gorm:"primaryKey" json:"id"`
	ReplacementID uint        `gorm:"column:replacement_id;not null" json:"replacement_id"`
	Position      string      `gorm:"column:position;type:varchar(20);not null" json:"position"`
	Action        string      `gorm:"column:action;size:50;not null" json:"action"`

	// Old Tyre
	OldTyreID        *uint          `gorm:"column:old_tyre_id" json:"old_tyre_id,omitempty"`
	OldTyreSerialNum string        `gorm:"column:old_tyre_serial_number;size:100" json:"old_tyre_serial_number,omitempty"`
	OldTyrePattern   string        `gorm:"column:old_tyre_pattern;size:100" json:"old_tyre_pattern,omitempty"`
	OldTyreSize      string        `gorm:"column:old_tyre_size;size:50" json:"old_tyre_size,omitempty"`
	OldTyreTread1    *float64      `gorm:"column:old_tyre_tread_1;type:decimal(5,2)" json:"old_tyre_tread_1,omitempty"`
	OldTyreTread2    *float64      `gorm:"column:old_tyre_tread_2;type:decimal(5,2)" json:"old_tyre_tread_2,omitempty"`
	OldTyreLifetime  *float64      `gorm:"column:old_tyre_lifetime;type:decimal(15,2)" json:"old_tyre_lifetime,omitempty"`
	OldTyreStatus    string        `gorm:"column:old_tyre_status;size:20" json:"old_tyre_status,omitempty"`
	FailureReasonID  *uint         `gorm:"column:failure_reason_id" json:"failure_reason_id,omitempty"`
	FailureReason    *MasterReason `gorm:"foreignKey:FailureReasonID" json:"failure_reason,omitempty"`
	FromUnitID       string        `gorm:"column:from_unit_id;size:50" json:"from_unit_id,omitempty"`

	// New Tyre
	NewTyreID          *uint    `gorm:"column:new_tyre_id" json:"new_tyre_id,omitempty"`
	NewTyreSerialNum  string   `gorm:"column:new_tyre_serial_number;size:100" json:"new_tyre_serial_number,omitempty"`
	NewTyrePattern    string   `gorm:"column:new_tyre_pattern;size:100" json:"new_tyre_pattern,omitempty"`
	NewTyreSize       string   `gorm:"column:new_tyre_size;size:50" json:"new_tyre_size,omitempty"`
	NewTyreTread1     *float64 `gorm:"column:new_tyre_tread_1;type:decimal(5,2)" json:"new_tyre_tread_1,omitempty"`
	NewTyreTread2     *float64 `gorm:"column:new_tyre_tread_2;type:decimal(5,2)" json:"new_tyre_tread_2,omitempty"`
	NewTyreCurrentLife float64 `gorm:"column:new_tyre_current_lifetime;type:decimal(15,2);not null;default:0" json:"new_tyre_current_lifetime"`
	NewTyreStatus     string   `gorm:"column:new_tyre_status;size:50" json:"new_tyre_status,omitempty"`

	Remark   string    `gorm:"column:remark;type:text" json:"remark,omitempty"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
}

func (ReplacementDetail) TableName() string { return "replacement_details" }

type ReplacementAction string

const (
	ActionMount    ReplacementAction = "mount"
	ActionDismount ReplacementAction = "dismount"
	ActionSwap     ReplacementAction = "swap"
)
