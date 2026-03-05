package domain

import (
	"errors"
	"time"
)

// Category groups products in the pantry.
type Category struct {
	Id          string
	Name        string
	Description string
	CreatedAt   time.Time
}

// NewCategory creates a validated Category. The id must be a pre-generated UUID.
func NewCategory(id, name, description string) (*Category, error) {
	entity := &Category{
		Id:          id,
		Name:        name,
		Description: description,
		CreatedAt:   time.Now(),
	}
	if err := entity.Validate(); err != nil {
		return nil, err
	}
	return entity, nil
}

// Validate ensures the category has required fields.
func (e *Category) Validate() error {
	if e.Id == "" {
		return errors.New("category id cannot be empty")
	}
	if e.Name == "" {
		return errors.New("category name cannot be empty")
	}
	return nil
}
