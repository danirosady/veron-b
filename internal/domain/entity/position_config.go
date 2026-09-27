package entity

import (
	"encoding/json"
	"strconv"
)

// PositionConfig represents one tyre position within a unit type template.
// Coordinates x/y are normalised [0,1] relative to the chassis image bounds.
type PositionConfig struct {
	Position string  `json:"position"`   // "1", "2", ...
	Label    string  `json:"label"`     // "Tyre 1", ...
	Side     string  `json:"side"`      // "left" | "right"
	Axle     string  `json:"axle"`      // "poros_1" | "poros_2" | ...
	X        float64 `json:"x"`         // normalised x [0,1]
	Y        float64 `json:"y"`         // normalised y [0,1]
	Col      int     `json:"col"`       // column index (0, 1, ...)
	MirrorOf string  `json:"mirror_of"` // position string of mirrored tyre, e.g. "2"
}

// UnmarshalJSON accepts position and mirror_of as either a JSON string ("1") or a
// JSON number (1). The canvas and seed data may send either form.
func (p *PositionConfig) UnmarshalJSON(data []byte) error {
	var raw struct {
		Position any    `json:"position"`
		Label    string `json:"label"`
		Side     string `json:"side"`
		Axle     string `json:"axle"`
		X        float64 `json:"x"`
		Y        float64 `json:"y"`
		Col      int     `json:"col"`
		MirrorOf any    `json:"mirror_of"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	p.Label = raw.Label
	p.Side  = raw.Side
	p.Axle  = raw.Axle
	p.X     = raw.X
	p.Y     = raw.Y
	p.Col   = raw.Col

	switch v := raw.Position.(type) {
	case string:
		p.Position = v
	case float64:
		p.Position = strconv.FormatFloat(v, 'f', -1, 64)
	}

	switch v := raw.MirrorOf.(type) {
	case string:
		p.MirrorOf = v
	case float64:
		p.MirrorOf = strconv.FormatFloat(v, 'f', -1, 64)
	}

	return nil
}
