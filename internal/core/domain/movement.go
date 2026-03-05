package domain

import (
	"errors"
	"fmt"
	"time"
)

// Movement records a single inbound or outbound stock event.
type Movement struct {
	Id           string
	ProductEan13 string
	Type         string // "in" or "out"
	Quantity     float64
	Reason       string
	Notes        string
	CreatedBy    string
	CreatedAt    time.Time
}

// NewMovement creates a validated Movement. The id must be a pre-generated UUID v4.
func NewMovement(id, productEan13 string, movType MovementType, quantity float64, reason, notes, createdBy string) (*Movement, error) {
	entity := &Movement{
		Id:           id,
		ProductEan13: productEan13,
		Type:         string(movType),
		Quantity:     quantity,
		Reason:       reason,
		Notes:        notes,
		CreatedBy:    createdBy,
		CreatedAt:    time.Now(),
	}
	if err := entity.Validate(); err != nil {
		return nil, err
	}
	return entity, nil
}

// Validate ensures the movement is consistent before persisting.
func (e *Movement) Validate() error {
	if e.Id == "" {
		return errors.New("movement id cannot be empty")
	}
	if _, err := NewEAN13(e.ProductEan13); err != nil {
		return fmt.Errorf("invalid product ean13: %w", err)
	}
	if _, err := NewMovementType(e.Type); err != nil {
		return fmt.Errorf("invalid movement type: %w", err)
	}
	if e.Quantity <= 0 {
		return errors.New("movement quantity must be greater than 0")
	}
	return nil
}

// MovementTypeVO returns the type as a typed MovementType value object.
func (e *Movement) MovementTypeVO() MovementType {
	return MovementType(e.Type)
}
