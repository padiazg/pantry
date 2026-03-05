package domain

import "fmt"

// MovementType represents the direction of a stock movement.
// Implemented as a typed string constant to ensure type safety in the domain.
type MovementType string

const (
	MovementTypeIn  MovementType = "in"
	MovementTypeOut MovementType = "out"
)

// NewMovementType creates a validated MovementType from a raw string.
func NewMovementType(value string) (MovementType, error) {
	mt := MovementType(value)
	if err := mt.Validate(); err != nil {
		return "", err
	}
	return mt, nil
}

// Validate ensures the movement type is either "in" or "out".
func (v MovementType) Validate() error {
	switch v {
	case MovementTypeIn, MovementTypeOut:
		return nil
	}
	return fmt.Errorf("invalid movement type %q: must be %q or %q", string(v), MovementTypeIn, MovementTypeOut)
}

// Equals compares two MovementType values.
func (v MovementType) Equals(other MovementType) bool {
	return v == other
}

// IsIn returns true when the movement is an inbound stock entry.
func (v MovementType) IsIn() bool { return v == MovementTypeIn }

// IsOut returns true when the movement is an outbound stock withdrawal.
func (v MovementType) IsOut() bool { return v == MovementTypeOut }

// String returns the string representation.
func (v MovementType) String() string { return string(v) }
