package domain

import (
	"errors"
	"fmt"
)

// EAN13 is a value object representing a GS1 EAN-13 barcode.
// It validates both the 13-digit format and the GS1 check digit algorithm.
type EAN13 struct {
	Value string
}

// NewEAN13 creates a validated EAN13 value object.
func NewEAN13(code string) (EAN13, error) {
	vo := EAN13{Value: code}
	if err := vo.Validate(); err != nil {
		return EAN13{}, err
	}
	return vo, nil
}

// Validate ensures the code has exactly 13 digits and a valid GS1 check digit.
func (v EAN13) Validate() error {
	if len(v.Value) != 13 {
		return fmt.Errorf("EAN-13 must have exactly 13 digits, got %d", len(v.Value))
	}
	for _, r := range v.Value {
		if r < '0' || r > '9' {
			return errors.New("EAN-13 must contain only digits")
		}
	}
	// GS1 check digit: alternating weights 1 and 3 on the first 12 digits.
	sum := 0
	for i := 0; i < 12; i++ {
		d := int(v.Value[i] - '0')
		if i%2 == 0 {
			sum += d
		} else {
			sum += d * 3
		}
	}
	expected := (10 - (sum % 10)) % 10
	actual := int(v.Value[12] - '0')
	if expected != actual {
		return fmt.Errorf("invalid EAN-13 check digit: expected %d, got %d", expected, actual)
	}
	return nil
}

// Equals compares two EAN13 instances by value.
func (v EAN13) Equals(other EAN13) bool {
	return v.Value == other.Value
}

// String returns the EAN-13 code.
func (v EAN13) String() string {
	return v.Value
}
