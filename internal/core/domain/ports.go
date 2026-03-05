package domain

import (
	"context"
	"errors"
	"time"
)

// ErrNotFound is returned by repositories when the requested resource does not exist.
var ErrNotFound = errors.New("not found")

// ProductFilter holds optional criteria for listing products.
type ProductFilter struct {
	CategoryID string
	Active     *bool
	LowStock   bool
}

// MovementFilter holds optional criteria for listing movements.
type MovementFilter struct {
	ProductEan13 string
	Type         string
	From         time.Time
	To           time.Time
}

// ProductRepository defines the secondary port for product persistence.
type ProductRepository interface {
	Create(ctx context.Context, product *Product) error
	FindByEAN13(ctx context.Context, ean13 string) (*Product, error)
	Update(ctx context.Context, product *Product) error
	List(ctx context.Context, filter ProductFilter) ([]*Product, error)
}

// CategoryRepository defines the secondary port for category persistence.
type CategoryRepository interface {
	Create(ctx context.Context, category *Category) error
	FindByID(ctx context.Context, id string) (*Category, error)
	Update(ctx context.Context, category *Category) error
	List(ctx context.Context) ([]*Category, error)
}

// MovementRepository defines the secondary port for movement persistence.
type MovementRepository interface {
	Create(ctx context.Context, movement *Movement) error
	FindByID(ctx context.Context, id string) (*Movement, error)
	List(ctx context.Context, filter MovementFilter) ([]*Movement, error)
	// CreateWithStockUpdate saves the movement and updates product stock atomically.
	CreateWithStockUpdate(ctx context.Context, movement *Movement, newStock float64) error
}
