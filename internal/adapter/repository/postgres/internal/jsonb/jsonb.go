package jsonb

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
)

type JSONB string

func NewNull() JSONB {
	return ""
}

func NewString(s string) JSONB {
	return JSONB(s)
}

func NewFromMarshaler(s any) (JSONB, error) {
	b, err := json.Marshal(s)
	if err != nil {
		return "", fmt.Errorf("json marshal: %w", err)
	}

	return JSONB(b), nil
}

func (j JSONB) IsEmpty() bool {
	return j == ""
}

func (j JSONB) IsEmptyObject() bool {
	return strings.TrimSpace(string(j)) == "{}"
}

func (j JSONB) Value() (driver.Value, error) {
	if j == "" {
		return "{}", nil
	}

	return string(j), nil
}

func ConvertTo[T any](j JSONB, v *T) error {
	if err := json.Unmarshal([]byte(j), v); err != nil {
		return fmt.Errorf("json unmarshal: %w", err)
	}

	return nil
}
