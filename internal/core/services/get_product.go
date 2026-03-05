package services

import (
	"context"
	"fmt"

	"github.com/padiazg/pantry/internal/core/domain"
)

// GetProductService handles product retrieval.
type GetProductService struct {
	products domain.ProductRepository
}

// NewGetProductService creates a new GetProductService.
func NewGetProductService(products domain.ProductRepository) *GetProductService {
	return &GetProductService{products: products}
}

// GetByEAN13 returns a single product by its EAN-13 code.
func (s *GetProductService) GetByEAN13(ctx context.Context, ean13 string) (*domain.Product, error) {
	product, err := s.products.FindByEAN13(ctx, ean13)
	if err != nil {
		return nil, fmt.Errorf("finding product %s: %w", ean13, err)
	}
	return product, nil
}

// List returns products matching the given filter criteria.
func (s *GetProductService) List(ctx context.Context, filter domain.ProductFilter) ([]*domain.Product, error) {
	products, err := s.products.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("listing products: %w", err)
	}
	return products, nil
}
