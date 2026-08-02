package entity

type PositionConfig struct {
	Position string  `json:"position"`
	Label    string  `json:"label"`
	Side     string  `json:"side"`
	Axle     string  `json:"axle"`
	X        float64 `json:"x"`
	Y        float64 `json:"y"`
}
