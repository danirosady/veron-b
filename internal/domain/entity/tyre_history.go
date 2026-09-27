package entity

import (
	"time"
)

type TyreHistoryItem struct {
	ID              uint      `json:"id"`
	ReplacementDate time.Time `json:"replacement_date"`
	HM              float64   `json:"hm"`
	Position        string    `json:"position"`
	Action          string    `json:"action"`
	RTD             *float64 `json:"rtd,omitempty"`
	Remarks         string    `json:"remarks,omitempty"`
	Unit            *Unit     `json:"unit,omitempty"`
}
