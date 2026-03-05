package services

import (
	"context"
	"fmt"
	"time"

	"github.com/padiazg/pantry/internal/core/domain"
)

// CreateProductInput holds the data needed to create a new product.
type CreateProductInput struct {
	EAN13       string
	Name        string
	Description string
	Unit        string
	MinStock    float64
	CategoryID  string
}

// UpdateProductInput holds the data needed to update an existing product.
type UpdateProductInput struct {
	EAN13       string
	Name        string
	Description string
	Unit        string
	MinStock    float64
	CategoryID  string
}

// ManageProductService handles product lifecycle: create, update, activate/deactivate.
type ManageProductService struct {
	products domain.ProductRepository
}

// NewManageProductService creates a new ManageProductService.
func NewManageProductService(products domain.ProductRepository) *ManageProductService {
	return &ManageProductService{products: products}
}

// Create validates and persists a new product.
func (s *ManageProductService) Create(ctx context.Context, input CreateProductInput) (*domain.Product, error) {
	product, err := domain.NewProduct(input.EAN13, input.Name, input.Description, input.Unit, input.MinStock, input.CategoryID)
	if err != nil {
		return nil, fmt.Errorf("invalid product: %w", err)
	}
	if err := s.products.Create(ctx, product); err != nil {
		return nil, fmt.Errorf("creating product: %w", err)
	}
	return product, nil
}

// Update replaces the mutable fields of an existing product.
func (s *ManageProductService) Update(ctx context.Context, input UpdateProductInput) (*domain.Product, error) {
	product, err := s.products.FindByEAN13(ctx, input.EAN13)
	if err != nil {
		return nil, fmt.Errorf("finding product: %w", err)
	}
	product.Name = input.Name
	product.Description = input.Description
	product.Unit = input.Unit
	product.MinStock = input.MinStock
	product.CategoryID = input.CategoryID
	product.UpdatedAt = time.Now()
	if err := product.Validate(); err != nil {
		return nil, fmt.Errorf("invalid product: %w", err)
	}
	if err := s.products.Update(ctx, product); err != nil {
		return nil, fmt.Errorf("updating product: %w", err)
	}
	return product, nil
}

// SetActive activates or deactivates a product (soft delete when active=false).
func (s *ManageProductService) SetActive(ctx context.Context, ean13 string, active bool) error {
	product, err := s.products.FindByEAN13(ctx, ean13)
	if err != nil {
		return fmt.Errorf("finding product: %w", err)
	}
	if active {
		product.Activate()
	} else {
		product.Deactivate()
	}
	if err := s.products.Update(ctx, product); err != nil {
		return fmt.Errorf("updating product active state: %w", err)
	}
	return nil
}
