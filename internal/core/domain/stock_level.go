package domain

import (
	"errors"
	"fmt"
)

// StockLevel is a value object representing a positive quantity of stock.
// The value must be strictly greater than zero.
type StockLevel struct {
	Value float64
}

// NewStockLevel creates a validated StockLevel.
func NewStockLevel(value float64) (StockLevel, error) {
	vo := StockLevel{Value: value}
	if err := vo.Validate(); err != nil {
		return StockLevel{}, err
	}
	return vo, nil
}

// Validate ensures the stock level is positive.
func (v StockLevel) Validate() error {
	if v.Value <= 0 {
		return errors.New("stock level must be greater than 0")
	}
	return nil
}

// Equals compares two StockLevel instances by value.
func (v StockLevel) Equals(other StockLevel) bool {
	return v.Value == other.Value
}

// String returns a formatted representation.
func (v StockLevel) String() string {
	return fmt.Sprintf("%.3f", v.Value)
}
