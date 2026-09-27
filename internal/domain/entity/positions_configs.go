package entity

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type PositionConfigs []PositionConfig

// Scan implements sql.Scanner — reads JSONB from PostgreSQL.
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

// UnmarshalJSON delegates to the custom UnmarshalJSON on each element so that
// mixed number/string Position/MirrorOf values are handled correctly.
func (p *PositionConfigs) UnmarshalJSON(data []byte) error {
	var raw []json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	*p = make(PositionConfigs, len(raw))
	for i, item := range raw {
		if err := json.Unmarshal(item, &(*p)[i]); err != nil {
			return err
		}
	}
	return nil
}
