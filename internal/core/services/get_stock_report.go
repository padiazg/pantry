package services

import (
	"context"
	"fmt"

	"github.com/padiazg/pantry/internal/core/domain"
)

// StockReportItem summarises the stock situation for one product.
type StockReportItem struct {
	Product    *domain.Product
	IsLowStock bool
}

// StockReport is the full stock summary across all active products.
type StockReport struct {
	TotalProducts    int
	LowStockProducts int
	Items            []*StockReportItem
}

// GetStockReportService generates stock summaries and low-stock alerts.
type GetStockReportService struct {
	products domain.ProductRepository
}

// NewGetStockReportService creates a new GetStockReportService.
func NewGetStockReportService(products domain.ProductRepository) *GetStockReportService {
	return &GetStockReportService{products: products}
}

// Execute returns a full stock report for all active products.
func (s *GetStockReportService) Execute(ctx context.Context) (*StockReport, error) {
	active := true
	all, err := s.products.List(ctx, domain.ProductFilter{Active: &active})
	if err != nil {
		return nil, fmt.Errorf("listing products for report: %w", err)
	}

	report := &StockReport{
		TotalProducts: len(all),
		Items:         make([]*StockReportItem, 0, len(all)),
	}
	for _, p := range all {
		item := &StockReportItem{Product: p, IsLowStock: p.IsLowStock()}
		report.Items = append(report.Items, item)
		if item.IsLowStock {
			report.LowStockProducts++
		}
	}
	return report, nil
}

// GetLowStock returns only the active products whose stock is at or below the minimum.
func (s *GetStockReportService) GetLowStock(ctx context.Context) ([]*domain.Product, error) {
	products, err := s.products.List(ctx, domain.ProductFilter{LowStock: true})
	if err != nil {
		return nil, fmt.Errorf("listing low-stock products: %w", err)
	}
	return products, nil
}
