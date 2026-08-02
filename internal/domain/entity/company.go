package entity

import (
	"time"
)

type Company struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	Name          string    `gorm:"column:name;size:255;not null" json:"name"`
	Address       string    `gorm:"column:address;type:text" json:"address,omitempty"`
	ContactPerson string    `gorm:"column:contact_person;size:255" json:"contact_person,omitempty"`
	Phone         string    `gorm:"column:phone;size:50" json:"phone,omitempty"`
	Email         string    `gorm:"column:email;size:255" json:"email,omitempty"`
	Status        string    `gorm:"column:status;size:20;not null;default:'active'" json:"status"`
	CreatedAt     time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (Company) TableName() string { return "companies" }
