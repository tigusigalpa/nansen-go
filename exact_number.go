package nansen

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		return fmt.Errorf("nansen: invalid exact number")
	}
	n.Lexeme = string(bytes.TrimSpace(data))
	return nil
}

func validateExactNumber(n ExactNumber) error {
	if !n.Present {
		return fmt.Errorf("nansen: exact number must be present")
	}
	if n.Null {
		return fmt.Errorf("nansen: null exact range bounds are not supported")
	}
	var parsed ExactNumber
	if err := parsed.UnmarshalJSON([]byte(n.Lexeme)); err != nil {
		return err
	}
	if parsed.Null || parsed.Lexeme != n.Lexeme {
		return fmt.Errorf("nansen: exact number must be a JSON number")
	}
	return nil
}

// ExactNumericRangeFilter serializes optional numeric bounds directly from
// validated JSON numeric lexemes, without a float64 conversion. A zero-value
// bound is omitted; null bounds are rejected because Nansen range filters do
// not document nullability.
type ExactNumericRangeFilter struct {
	Min ExactNumber
	Max ExactNumber
}

// MarshalJSON serializes present exact bounds as JSON numbers.
func (f ExactNumericRangeFilter) MarshalJSON() ([]byte, error) {
	type wire struct {
		Min *json.RawMessage `json:"min,omitempty"`
		Max *json.RawMessage `json:"max,omitempty"`
	}
	var out wire
	if f.Min.Present {
		if err := validateExactNumber(f.Min); err != nil {
			return nil, fmt.Errorf("nansen: invalid exact min bound: %w", err)
		}
		raw := json.RawMessage(f.Min.Lexeme)
		out.Min = &raw
	}
	if f.Max.Present {
		if err := validateExactNumber(f.Max); err != nil {
			return nil, fmt.Errorf("nansen: invalid exact max bound: %w", err)
		}
		raw := json.RawMessage(f.Max.Lexeme)
		out.Max = &raw
	}
	return json.Marshal(out)
}

// ExactNumberFromLexeme constructs a present exact number after validating its JSON number grammar.
func ExactNumberFromLexeme(lexeme string) (ExactNumber, error) {
	n := ExactNumber{Lexeme: lexeme, Present: true}
	if err := validateExactNumber(n); err != nil {
		return ExactNumber{}, err
	}
	return n, nil
}

// MarshalJSON emits the preserved numeric lexeme.
func (n ExactNumber) MarshalJSON() ([]byte, error) {
	if !n.Present || n.Null {
		return []byte("null"), nil
	}
	return []byte(n.Lexeme), nil
}
