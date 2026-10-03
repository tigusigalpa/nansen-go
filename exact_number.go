package nansen

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// ExactNumber preserves a JSON numeric lexeme and distinguishes absent, null,
// and zero values. It is intended for provider-supplied economic values.
type ExactNumber struct {
	Lexeme  string
	Present bool
	Null    bool
}

// UnmarshalJSON records the original JSON number without converting it to float64.
func (n *ExactNumber) UnmarshalJSON(data []byte) error {
	n.Present = true
	n.Null = bytes.Equal(bytes.TrimSpace(data), []byte("null"))
	n.Lexeme = ""
	if n.Null {
		return nil
	}

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	var value any
	if err := decoder.Decode(&value); err != nil {
		return fmt.Errorf("nansen: invalid exact number: %w", err)
	}
	if _, ok := value.(json.Number); !ok {
		return fmt.Errorf("nansen: exact number must be a JSON number or null")
	}
	if decoder.More() {
		return fmt.Errorf("nansen: invalid exact number")
	}
	n.Lexeme = string(bytes.TrimSpace(data))
	return nil
}

// MarshalJSON emits the preserved numeric lexeme.
func (n ExactNumber) MarshalJSON() ([]byte, error) {
	if !n.Present || n.Null {
		return []byte("null"), nil
	}
	return []byte(n.Lexeme), nil
}
