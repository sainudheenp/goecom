package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// JSONStringSlice is a custom type for []string stored as JSON
type JSONStringSlice []string

// Scan implements sql.Scanner interface
func (j *JSONStringSlice) Scan(value interface{}) error {
	if value == nil {
		*j = []string{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan JSONStringSlice")
	}

	return json.Unmarshal(bytes, j)
}

// Value implements driver.Valuer interface
func (j JSONStringSlice) Value() (driver.Value, error) {
	if len(j) == 0 {
		return json.Marshal([]string{})
	}
	return json.Marshal(j)
}

// JSONMap is a custom type for map[string]interface{} stored as JSON
type JSONMap map[string]interface{}

// Scan implements sql.Scanner interface
func (j *JSONMap) Scan(value interface{}) error {
	if value == nil {
		*j = make(map[string]interface{})
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan JSONMap")
	}

	return json.Unmarshal(bytes, j)
}

// Value implements driver.Valuer interface
func (j JSONMap) Value() (driver.Value, error) {
	if len(j) == 0 {
		return json.Marshal(map[string]interface{}{})
	}
	return json.Marshal(j)
}
