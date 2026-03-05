package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/padiazg/pantry/internal/core/domain"
)

// MovementRepository implements domain.MovementRepository using PostgreSQL.
type MovementRepository struct {
	db *sql.DB
}

// NewMovementRepository creates a new MovementRepository.
func NewMovementRepository(db *sql.DB) *MovementRepository {
	return &MovementRepository{db: db}
}

// compile-time check that MovementRepository satisfies the port.
var _ domain.MovementRepository = (*MovementRepository)(nil)

// Create inserts a new movement record.
func (r *MovementRepository) Create(ctx context.Context, m *domain.Movement) error {
	q := `INSERT INTO movements (id, product_ean13, type, quantity, reason, notes, created_by, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.db.ExecContext(ctx, q,
		m.Id, m.ProductEan13, m.Type, m.Quantity,
		m.Reason, m.Notes, m.CreatedBy, m.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert movement: %w", err)
	}
	return nil
}

// FindByID retrieves a movement by its UUID.
func (r *MovementRepository) FindByID(ctx context.Context, id string) (*domain.Movement, error) {
	q := `SELECT id, product_ean13, type, quantity, reason, notes, created_by, created_at
		FROM movements WHERE id = $1`
	row := r.db.QueryRowContext(ctx, q, id)
	m := &domain.Movement{}
	err := row.Scan(&m.Id, &m.ProductEan13, &m.Type, &m.Quantity,
		&m.Reason, &m.Notes, &m.CreatedBy, &m.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("scan movement: %w", err)
	}
	return m, nil
}

// List returns movements matching the filter, ordered by created_at descending.
func (r *MovementRepository) List(ctx context.Context, filter domain.MovementFilter) ([]*domain.Movement, error) {
	q := `SELECT id, product_ean13, type, quantity, reason, notes, created_by, created_at
		FROM movements WHERE 1=1`
	args := []any{}
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
		idx++
	}
	q += " ORDER BY created_at DESC"

	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("list movements: %w", err)
	}
	defer rows.Close()

	var movements []*domain.Movement
	for rows.Next() {
		m := &domain.Movement{}
		if err := rows.Scan(&m.Id, &m.ProductEan13, &m.Type, &m.Quantity,
			&m.Reason, &m.Notes, &m.CreatedBy, &m.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan movement row: %w", err)
		}
		movements = append(movements, m)
	}
	return movements, rows.Err()
}

// CreateWithStockUpdate saves a movement and updates the product's current stock atomically.
// This is the key operation for business rule: atomic stock update.
func (r *MovementRepository) CreateWithStockUpdate(ctx context.Context, m *domain.Movement, newStock float64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback() //nolint:errcheck

	insertQ := `INSERT INTO movements (id, product_ean13, type, quantity, reason, notes, created_by, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	if _, err := tx.ExecContext(ctx, insertQ,
		m.Id, m.ProductEan13, m.Type, m.Quantity,
		m.Reason, m.Notes, m.CreatedBy, m.CreatedAt,
	); err != nil {
		return fmt.Errorf("insert movement in tx: %w", err)
	}

	updateQ := `UPDATE products SET current_stock=$1, updated_at=$2 WHERE ean13=$3`
	if _, err := tx.ExecContext(ctx, updateQ, newStock, time.Now(), m.ProductEan13); err != nil {
		return fmt.Errorf("update stock in tx: %w", err)
	}

	return tx.Commit()
}
