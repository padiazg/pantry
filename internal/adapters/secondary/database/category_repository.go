package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/padiazg/pantry/internal/core/domain"
)

// CategoryRepository implements domain.CategoryRepository using PostgreSQL.
type CategoryRepository struct {
	db *sql.DB
}

// NewCategoryRepository creates a new CategoryRepository.
func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

// compile-time check that CategoryRepository satisfies the port.
var _ domain.CategoryRepository = (*CategoryRepository)(nil)

// Create inserts a new category.
func (r *CategoryRepository) Create(ctx context.Context, c *domain.Category) error {
	q := `INSERT INTO categories (id, name, description, created_at) VALUES ($1, $2, $3, $4)`
	_, err := r.db.ExecContext(ctx, q, c.Id, c.Name, c.Description, c.CreatedAt)
	if err != nil {
		return fmt.Errorf("insert category: %w", err)
	}
	return nil
}

// FindByID retrieves a category by its UUID.
func (r *CategoryRepository) FindByID(ctx context.Context, id string) (*domain.Category, error) {
	q := `SELECT id, name, description, created_at FROM categories WHERE id = $1`
	row := r.db.QueryRowContext(ctx, q, id)
	c := &domain.Category{}
	err := row.Scan(&c.Id, &c.Name, &c.Description, &c.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("scan category: %w", err)
	}
	return c, nil
}

// Update saves updated category fields.
func (r *CategoryRepository) Update(ctx context.Context, c *domain.Category) error {
	q := `UPDATE categories SET name=$1, description=$2 WHERE id=$3`
	res, err := r.db.ExecContext(ctx, q, c.Name, c.Description, c.Id)
	if err != nil {
		return fmt.Errorf("update category: %w", err)
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// List returns all categories ordered by name.
func (r *CategoryRepository) List(ctx context.Context) ([]*domain.Category, error) {
	q := `SELECT id, name, description, created_at FROM categories ORDER BY name`
	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("list categories: %w", err)
	}
	defer rows.Close()

	var cats []*domain.Category
	for rows.Next() {
		c := &domain.Category{}
		if err := rows.Scan(&c.Id, &c.Name, &c.Description, &c.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan category row: %w", err)
		}
		cats = append(cats, c)
	}
	return cats, rows.Err()
}
