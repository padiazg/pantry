package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/padiazg/pantry/internal/core/domain"
)

// dbTX is the subset of database/sql.DB used by this repository.
type dbTX interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

// ProductRepository implements domain.ProductRepository using PostgreSQL.
type ProductRepository struct {
	db dbTX
}

// NewProductRepository creates a new ProductRepository.
func NewProductRepository(db dbTX) *ProductRepository {
	return &ProductRepository{db: db}
}

// compile-time check that ProductRepository satisfies the port.
var _ domain.ProductRepository = (*ProductRepository)(nil)

// Create inserts a new product into the database.
func (r *ProductRepository) Create(ctx context.Context, p *domain.Product) error {
	q := `INSERT INTO products
		(ean13, name, description, unit, min_stock, current_stock, category_id, active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NULLIF($7,''), $8, $9, $10)`
	_, err := r.db.ExecContext(ctx, q,
		p.Ean13, p.Name, p.Description, p.Unit,
		p.MinStock, p.CurrentStock, p.CategoryID,
		p.Active, p.CreatedAt, p.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert product: %w", err)
	}
	return nil
}

// FindByEAN13 retrieves a product by its EAN-13 code.
func (r *ProductRepository) FindByEAN13(ctx context.Context, ean13 string) (*domain.Product, error) {
	q := `SELECT ean13, name, description, unit, min_stock, current_stock,
		COALESCE(category_id,''), active, created_at, updated_at
		FROM products WHERE ean13 = $1`
	row := r.db.QueryRowContext(ctx, q, ean13)
	p := &domain.Product{}
	err := row.Scan(
		&p.Ean13, &p.Name, &p.Description, &p.Unit,
		&p.MinStock, &p.CurrentStock, &p.CategoryID,
		&p.Active, &p.CreatedAt, &p.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("scan product: %w", err)
	}
	return p, nil
}

// Update saves updated product fields (everything except ean13 and created_at).
func (r *ProductRepository) Update(ctx context.Context, p *domain.Product) error {
	q := `UPDATE products
		SET name=$1, description=$2, unit=$3, min_stock=$4, current_stock=$5,
		    category_id=NULLIF($6,''), active=$7, updated_at=$8
		WHERE ean13=$9`
	res, err := r.db.ExecContext(ctx, q,
		p.Name, p.Description, p.Unit, p.MinStock, p.CurrentStock,
		p.CategoryID, p.Active, p.UpdatedAt, p.Ean13,
	)
	if err != nil {
		return fmt.Errorf("update product: %w", err)
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// List returns products matching the provided filter, ordered by name.
func (r *ProductRepository) List(ctx context.Context, filter domain.ProductFilter) ([]*domain.Product, error) {
	q := `SELECT ean13, name, description, unit, min_stock, current_stock,
		COALESCE(category_id,''), active, created_at, updated_at
		FROM products WHERE 1=1`

	q0, args := filter.Args()
	q += q0 + " ORDER BY name"

	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("list products: %w", err)
	}
	defer rows.Close()

	var products []*domain.Product
	for rows.Next() {
		p := &domain.Product{}
		if err := rows.Scan(
			&p.Ean13, &p.Name, &p.Description, &p.Unit,
			&p.MinStock, &p.CurrentStock, &p.CategoryID,
			&p.Active, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan product row: %w", err)
		}
		products = append(products, p)
	}
	return products, rows.Err()
}
