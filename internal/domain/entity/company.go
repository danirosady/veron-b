package entity

import (
	"time"
)

type Company struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	Name          string    `gorm:"size:255;not null" json:"name"`
	Address       string    `gorm:"type:text" json:"address,omitempty"`
	ContactPerson string    `gorm:"size:255" json:"contact_person,omitempty"`
	Phone         string    `gorm:"size:50" json:"phone,omitempty"`
	Email         string    `gorm:"size:255" json:"email,omitempty"`
	Status        string    `gorm:"size:20;not null;default:'active'" json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (Company) TableName() string { return "companies" }
