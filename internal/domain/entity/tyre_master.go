package entity

import (
	"time"
)

type TyreMaster struct {
	ID       uint `gorm:"primaryKey" json:"id"`
	CompanyID uint `gorm:"not null" json:"company_id"`
	Company   *Company  `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	UnitID    *uint     `json:"unit_id,omitempty"`
	Unit      *Unit     `gorm:"foreignKey:UnitID" json:"unit,omitempty"`
	MountedPosition *int `json:"mounted_position,omitempty"`

	Barcode      string `gorm:"size:100;uniqueIndex;not null" json:"barcode"`
	SerialNumber string `gorm:"size:100;uniqueIndex;not null" json:"serial_number"`
	DOTCode     string `gorm:"size:100" json:"dot_code,omitempty"`
	Type        string `gorm:"size:50;not null;default:'Radial'" json:"type"`

	SizeID    uint          `gorm:"not null" json:"size_id"`
	Size      *MasterSize   `gorm:"foreignKey:SizeID" json:"size,omitempty"`
	BrandID   uint          `gorm:"not null" json:"brand_id"`
	Brand     *MasterBrand  `gorm:"foreignKey:BrandID" json:"brand,omitempty"`
	PatternID uint          `gorm:"not null" json:"pattern_id"`
	Pattern   *MasterPattern `gorm:"foreignKey:PatternID" json:"pattern,omitempty"`

	OTD      float64  `gorm:"type:decimal(5,2);not null;default:0" json:"otd"`
	RTD      float64  `gorm:"type:decimal(5,2);not null;default:0" json:"rtd"`
	RTD1     *float64 `gorm:"type:decimal(5,2)" json:"rtd_1,omitempty"`
	RTD2     *float64 `gorm:"type:decimal(5,2)" json:"rtd_2,omitempty"`
	Lifetime float64  `gorm:"type:decimal(15,2);not null;default:0" json:"lifetime"`
	PSI      *float64 `gorm:"type:decimal(5,2)" json:"psi,omitempty"`

	Status  string `gorm:"size:20;not null;default:'spare'" json:"status"`
	Remarks string `gorm:"type:text" json:"remarks,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (TyreMaster) TableName() string { return "tyre_master" }

type TyreStatus string

const (
	TyreStatusSpare      TyreStatus = "spare"
	TyreStatusMounted     TyreStatus = "mounted"
	TyreStatusDismounted TyreStatus = "dismounted"
	TyreStatusScrap      TyreStatus = "scrap"
)
