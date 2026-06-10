package database_test

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/padiazg/pantry/internal/adapters/secondary/database"
	"github.com/padiazg/pantry/internal/core/domain"
	"github.com/stretchr/testify/assert"
)

func TestProductRepository(t *testing.T) {
	t.Skip("requires integration test with a real PostgreSQL database")
}

type checkProductRepositoryListFn func(*testing.T, []*domain.Product, error)

var checkProductRepositoryList = func(fns ...checkProductRepositoryListFn) []checkProductRepositoryListFn { return fns }

func checkListError(want string) checkProductRepositoryListFn {
	return func(t *testing.T, _ []*domain.Product, err error) {
		t.Helper()
		if want == "" {
			assert.NoErrorf(t, err, "checkListError: expected no error, got %v", err)
			return
		}
		if assert.Errorf(t, err, "checkListError: expected error %q", want) {
			assert.Containsf(t, err.Error(), want, "checkListError mismatch")
		}
	}
}

func boolPtr(b bool) *bool {
	return &b
}

func TestProductRepository_List(t *testing.T) {
	columns := []string{
		"ean13", "name", "description", "unit", "min_stock", "current_stock",
		"category_id", "active", "created_at", "updated_at",
	}

	tests := []struct {
		name   string
		filter domain.ProductFilter
		expect []*domain.Product
		checks []checkProductRepositoryListFn
	}{
		{
			name:   "returns all products when no filter",
			filter: domain.ProductFilter{},
			expect: []*domain.Product{
				{Ean13: "1234567890128", Name: "Arroz", Unit: "kg", CurrentStock: 5.0, Active: true},
				{Ean13: "1234567890135", Name: "Leche", Unit: "l", CurrentStock: 2.0, Active: true},
			},
			checks: checkProductRepositoryList(
				checkListError(""),
				func(t *testing.T, products []*domain.Product, err error) {
					t.Helper()
					assert.Len(t, products, 2)
				},
			),
		},
		{
			name:   "filters by category ID",
			filter: domain.ProductFilter{CategoryID: "1234567890123"},
			expect: []*domain.Product{
				{Ean13: "1234567890128", Name: "Arroz", Unit: "kg", CurrentStock: 5.0, Active: true, CategoryID: "1234567890123"},
			},
			checks: checkProductRepositoryList(
				checkListError(""),
				func(t *testing.T, products []*domain.Product, err error) {
					t.Helper()
					assert.Len(t, products, 1)
					assert.Equal(t, "1234567890123", products[0].CategoryID)
				},
			),
		},
		{
			name:   "filters by active status",
			filter: domain.ProductFilter{Active: boolPtr(true)},
			expect: []*domain.Product{
				{Ean13: "1234567890128", Name: "Arroz", Unit: "kg", CurrentStock: 5.0, Active: true},
			},
			checks: checkProductRepositoryList(
				checkListError(""),
				func(t *testing.T, products []*domain.Product, err error) {
					t.Helper()
					assert.Len(t, products, 1)
					assert.True(t, products[0].Active)
				},
			),
		},
		{
			name:   "filters by low stock",
			filter: domain.ProductFilter{LowStock: true},
			expect: []*domain.Product{
				{Ean13: "1234567890135", Name: "Leche", Unit: "l", CurrentStock: 1.0, MinStock: 2.0, Active: true},
			},
			checks: checkProductRepositoryList(
				checkListError(""),
				func(t *testing.T, products []*domain.Product, err error) {
					t.Helper()
					assert.Len(t, products, 1)
					assert.True(t, products[0].CurrentStock <= products[0].MinStock)
				},
			),
		},
		{
			name:   "returns empty slice when no match",
			filter: domain.ProductFilter{CategoryID: "nonexistent"},
			expect: []*domain.Product{},
			checks: checkProductRepositoryList(
				checkListError(""),
				func(t *testing.T, products []*domain.Product, err error) {
					t.Helper()
					assert.Empty(t, products)
				},
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			q0, _ := tt.filter.Args()
			expectedQuery := `SELECT\s+ean13,\s+name,\s+description,\s+unit,\s+min_stock,\s+current_stock,\s+COALESCE\(category_id,''\),\s+active,\s+created_at,\s+updated_at\s+FROM\s+products\s+WHERE\s+1=1` + regexp.QuoteMeta(q0) + `\s+ORDER\s+BY\s+name`

			rows := sqlmock.NewRows(columns)
			for _, p := range tt.expect {
				rows.AddRow(p.Ean13, p.Name, p.Description, p.Unit, p.MinStock, p.CurrentStock, p.CategoryID, p.Active, p.CreatedAt, p.UpdatedAt)
			}

			mock.ExpectQuery(expectedQuery).WillReturnRows(rows)

			s := database.NewProductRepository(db)
			r, err := s.List(context.Background(), tt.filter)
			for _, c := range tt.checks {
				c(t, r, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
