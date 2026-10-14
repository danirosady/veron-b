package entity

import (
	"time"
)

type EmployeeRole string

const (
	RoleDriver     EmployeeRole = "driver"
	RoleSupervisor EmployeeRole = "supervisor"
	RoleTyreman   EmployeeRole = "tyreman"
)

type Employee struct {
	ID         uint         `gorm:"primaryKey" json:"id"`
	CompanyID  uint         `gorm:"column:company_id;not null;uniqueIndex:idx_employees_company_role_emp,priority:1" json:"company_id"`
	Company    *Company    `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Name       string       `gorm:"column:name;size:255;not null" json:"name"`
	EmployeeID string       `gorm:"column:employee_id;size:50;not null;uniqueIndex:idx_employees_company_role_emp,priority:2" json:"employee_id"`
	Phone      string       `gorm:"column:phone;size:50" json:"phone,omitempty"`
	Role       EmployeeRole `gorm:"column:role;size:20;not null;uniqueIndex:idx_employees_company_role_emp,priority:3" json:"role"`
	Status     string       `gorm:"column:status;size:20;not null;default:'active'" json:"status"`
	CreatedAt  time.Time   `gorm:"column:created_at" json:"created_at"`
	UpdatedAt  time.Time   `gorm:"column:updated_at" json:"updated_at"`
}

func (Employee) TableName() string { return "employees" }
