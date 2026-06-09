package domain

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// ErrNotFound is returned by repositories when the requested resource does not exist.
var ErrNotFound = errors.New("not found")

// ProductFilter holds optional criteria for listing products.
type ProductFilter struct {
	Active     *bool
	CategoryID string
	LowStock   bool
}

// MovementFilter holds optional criteria for listing movements.
type MovementFilter struct {
	From         time.Time
	To           time.Time
	ProductEan13 string
	Type         string
}

func (filter *MovementFilter) Args() (string, []any) {
	var (
		q    string
		args []any
	)

	idx := 1

	if filter.ProductEan13 != "" {
		q += fmt.Sprintf(" AND product_ean13 = $%d", idx)
		args = append(args, filter.ProductEan13)
		idx++
	}
	if filter.Type != "" {
		q += fmt.Sprintf(" AND type = $%d", idx)
		args = append(args, filter.Type)
		idx++
	}
	if !filter.From.IsZero() {
		q += fmt.Sprintf(" AND created_at >= $%d", idx)
		args = append(args, filter.From)
		idx++
	}
	if !filter.To.IsZero() {
		q += fmt.Sprintf(" AND created_at <= $%d", idx)
		args = append(args, filter.To)
	}

	return q, args
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
