package domain

import (
	"errors"
	"fmt"
	"time"
)

// Product represents an article stored in the pantry.
// The EAN-13 barcode is the natural key.
type Product struct {
	Ean13        string
	Name         string
	Description  string
	Unit         string
	MinStock     float64
	CurrentStock float64
	CategoryID   string
	Active       bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// NewProduct creates a new active Product with zero initial stock.
func NewProduct(ean13, name, description, unit string, minStock float64, categoryID string) (*Product, error) {
	now := time.Now()
	entity := &Product{
		Ean13:        ean13,
		Name:         name,
		Description:  description,
		Unit:         unit,
		MinStock:     minStock,
		CurrentStock: 0,
		CategoryID:   categoryID,
		Active:       true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if err := entity.Validate(); err != nil {
		return nil, err
	}
	return entity, nil
}

// Validate ensures the product is in a consistent state.
func (e *Product) Validate() error {
	if _, err := NewEAN13(e.Ean13); err != nil {
		return fmt.Errorf("invalid ean13: %w", err)
	}
	if e.Name == "" {
		return errors.New("product name cannot be empty")
	}
	if e.Unit == "" {
		return errors.New("product unit cannot be empty")
	}
	if e.MinStock < 0 {
		return errors.New("min stock cannot be negative")
	}
	return nil
}

// IsLowStock returns true when current stock is at or below the minimum threshold.
func (e *Product) IsLowStock() bool {
	return e.CurrentStock <= e.MinStock
}

// Activate marks the product as active.
func (e *Product) Activate() {
	e.Active = true
	e.UpdatedAt = time.Now()
}

// Deactivate performs a soft-delete by marking the product inactive.
func (e *Product) Deactivate() {
	e.Active = false
	e.UpdatedAt = time.Now()
}

// ApplyMovement adjusts the current stock according to movement type and quantity.
// Returns an error if an outbound movement would result in negative stock.
func (e *Product) ApplyMovement(movType MovementType, quantity float64) error {
	switch movType {
	case MovementTypeIn:
		e.CurrentStock += quantity
	case MovementTypeOut:
		if e.CurrentStock-quantity < 0 {
			return fmt.Errorf("insufficient stock: available %.3f, requested %.3f", e.CurrentStock, quantity)
		}
		e.CurrentStock -= quantity
	default:
		return fmt.Errorf("unknown movement type: %s", movType)
	}
	e.UpdatedAt = time.Now()
	return nil
}
