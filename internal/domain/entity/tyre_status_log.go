package entity

import (
	"time"
)

// TyreStatusLog records every status/remarks change on a tyre.
type TyreStatusLog struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	TyreID       uint      `gorm:"column:tyre_id;not null;index" json:"tyre_id"`
	Tyre         *TyreMaster `gorm:"foreignKey:TyreID" json:"tyre,omitempty"`

	// Old values (null on create)
	OldStatus      string `gorm:"column:old_status;size:20" json:"old_status,omitempty"`
	OldRemarksID   *uint  `gorm:"column:old_remarks_id" json:"old_remarks_id,omitempty"`
	OldRemarks     *MasterRemark `gorm:"foreignKey:OldRemarksID" json:"old_remarks,omitempty"`

	// New values
	NewStatus      string `gorm:"column:new_status;size:20;not null" json:"new_status"`
	NewRemarksID   *uint  `gorm:"column:new_remarks_id" json:"new_remarks_id,omitempty"`
	NewRemarks     *MasterRemark `gorm:"foreignKey:NewRemarksID" json:"new_remarks,omitempty"`

	// Who made the change (nil = system)
	ChangedBy *uint  `gorm:"column:changed_by" json:"changed_by,omitempty"`
	Reason    string `gorm:"column:reason;type:text" json:"reason,omitempty"`
	ChangedAt time.Time `gorm:"column:changed_at;not null;default:now()" json:"changed_at"`
}

func (TyreStatusLog) TableName() string { return "tyre_status_logs" }
