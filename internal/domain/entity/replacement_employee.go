package entity

import (
	"time"
)

// ReplacementEmployee is a junction table linking a replacement detail to employees (tyremen).
type ReplacementEmployee struct {
	ID                  uint       `gorm:"primaryKey" json:"id"`
	ReplacementDetailID uint       `gorm:"column:replacement_detail_id;not null;uniqueIndex:uq_re_detail_employee" json:"replacement_detail_id"`
	EmployeeID          uint       `gorm:"column:employee_id;not null;uniqueIndex:uq_re_detail_employee" json:"employee_id"`
	Employee            *Employee  `gorm:"foreignKey:EmployeeID" json:"employee,omitempty"`
	CreatedAt          time.Time  `gorm:"column:created_at" json:"created_at"`
}

func (ReplacementEmployee) TableName() string { return "replacement_employees" }
