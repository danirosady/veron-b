package entity

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type PositionConfigs []PositionConfig

func (p *PositionConfigs) Scan(value interface{}) error {
	if value == nil {
		*p = nil
		return nil
	}

	var data []byte

	switch v := value.(type) {
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		return fmt.Errorf("cannot scan %T into PositionConfigs", value)
	}

	return json.Unmarshal(data, p)
}

func (p PositionConfigs) Value() (driver.Value, error) {
	return json.Marshal(p)
}
