package entity

import (
	"time"
)

type TyreHistoryItem struct {
	ID              uint      `json:"id"`
	ReplacementDate time.Time `json:"replacement_date"`
	HM              float64   `json:"hm"`
	CurrentLifeHM   float64   `json:"current_life_hm"`
	Position        string    `json:"position"`
	Action          string    `json:"action"`
	RTD             *float64 `json:"rtd,omitempty"`
	RTDATMount      *float64 `json:"rtd_at_mount,omitempty"`
	RTDAtDismount  *float64 `json:"rtd_at_dismount,omitempty"`
	LifeHM         *float64 `json:"life_hm,omitempty"`
	Remarks         string    `json:"remarks,omitempty"`
	Unit            *Unit     `json:"unit,omitempty"`
}
