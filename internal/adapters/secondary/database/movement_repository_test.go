package database_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/padiazg/pantry/internal/adapters/secondary/database"
	"github.com/padiazg/pantry/internal/core/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMovementRepository_Create(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		before  func(t *testing.T, sm sqlmock.Sqlmock)
		wantErr bool
		wantSub string
	}{
		{
			name: "success",
			before: func(t *testing.T, sm sqlmock.Sqlmock) {
				sm.ExpectExec(`INSERT INTO movements`).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
		},
		{
			name: "insert error",
			before: func(t *testing.T, sm sqlmock.Sqlmock) {
				sm.ExpectExec(`INSERT INTO movements`).
					WillReturnError(assert.AnError)
			},
			wantErr: true,
			wantSub: "insert movement",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			m, _, sm := newMockDB(t)
			tt.before(t, sm)

			movement := &domain.Movement{
				Id:           uuid.NewString(),
				ProductEan13: "1234567890128",
				Type:         "in",
				Quantity:     10.0,
				Reason:       "restock",
				Notes:        "initial stock",
				CreatedBy:    "user1",
				CreatedAt:    time.Now(),
			}

			s := database.NewMovementRepository(m)
			err := s.Create(context.Background(), movement)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantSub)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestMovementRepository_FindByID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		before  func(t *testing.T, sm sqlmock.Sqlmock, id string)
		wantErr bool
		wantSub string
	}{
		{
			name: "success",
			before: func(t *testing.T, sm sqlmock.Sqlmock, id string) {
				sm.ExpectQuery(`.*FROM movements.*`).
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "product_ean13", "type", "quantity",
						"reason", "notes", "created_by", "created_at",
					}).AddRow(
						id,
						"1234567890128",
						"in",
						10.0,
						"restock",
						"notes",
						"user1",
						time.Now(),
					))
			},
		},
		{
			name: "not found",
			before: func(t *testing.T, sm sqlmock.Sqlmock, id string) {
				sm.ExpectQuery(`.*FROM movements.*`).
					WillReturnError(sql.ErrNoRows)
			},
			wantErr: true,
		},
		{
			name: "query error",
			before: func(t *testing.T, sm sqlmock.Sqlmock, id string) {
				sm.ExpectQuery(`.*FROM movements.*`).
					WillReturnError(assert.AnError)
			},
			wantErr: true,
			wantSub: "scan movement",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			m, _, sm := newMockDB(t)
			id := uuid.NewString()
			tt.before(t, sm, id)

			s := database.NewMovementRepository(m)
			movement, err := s.FindByID(context.Background(), id)

			if tt.wantErr {
				require.Error(t, err)
				if tt.wantSub != "" {
					assert.Contains(t, err.Error(), tt.wantSub)
				}
				if tt.name == "not found" {
					assert.True(t, errors.Is(err, domain.ErrNotFound), "expected ErrNotFound")
				}
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, id, movement.Id)
			assert.Equal(t, "1234567890128", movement.ProductEan13)
			assert.Equal(t, "in", movement.Type)
			assert.Equal(t, 10.0, movement.Quantity)
			assert.Equal(t, "restock", movement.Reason)
			assert.Equal(t, "notes", movement.Notes)
			assert.Equal(t, "user1", movement.CreatedBy)
		})
	}
}

func TestMovementRepository_List(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		filter  domain.MovementFilter
		before  func(t *testing.T, sm sqlmock.Sqlmock)
		wantErr bool
		wantSub string
		wantLen int
	}{
		{
			name:    "empty filter returns all",
			filter:  domain.MovementFilter{},
			wantLen: 2,
			before: func(t *testing.T, sm sqlmock.Sqlmock) {
				sm.ExpectQuery(`.*FROM movements.*`).
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "product_ean13", "type", "quantity",
						"reason", "notes", "created_by", "created_at",
					}).AddRow(
						uuid.NewString(),
						"1234567890128",
						"in",
						10.0,
						"restock",
						"notes",
						"user1",
						time.Now(),
					).AddRow(
						uuid.NewString(),
						"1234567890128",
						"out",
						2.0,
						"sale",
						"",
						"user2",
						time.Now().Add(-time.Hour),
					))
			},
		},
		{
			name: "filter by product ean13",
			filter: domain.MovementFilter{
				ProductEan13: "1234567890128",
			},
			wantLen: 1,
			before: func(t *testing.T, sm sqlmock.Sqlmock) {
				sm.ExpectQuery(`.*FROM movements.*`).
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "product_ean13", "type", "quantity",
						"reason", "notes", "created_by", "created_at",
					}).AddRow(
						uuid.NewString(),
						"1234567890128",
						"in",
						10.0,
						"restock",
						"",
						"user1",
						time.Now(),
					))
			},
		},
		{
			name: "filter by type",
			filter: domain.MovementFilter{
				Type: "out",
			},
			wantLen: 1,
			before: func(t *testing.T, sm sqlmock.Sqlmock) {
				sm.ExpectQuery(`.*FROM movements.*`).
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "product_ean13", "type", "quantity",
						"reason", "notes", "created_by", "created_at",
					}).AddRow(
						uuid.NewString(),
						"1234567890128",
						"out",
						2.0,
						"sale",
						"",
						"user1",
						time.Now(),
					))
			},
		},
		{
			name: "filter by date range",
			filter: domain.MovementFilter{
				From: time.Now().Add(-24 * time.Hour),
				To:   time.Now(),
			},
			wantLen: 1,
			before: func(t *testing.T, sm sqlmock.Sqlmock) {
				sm.ExpectQuery(`.*FROM movements.*`).
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "product_ean13", "type", "quantity",
						"reason", "notes", "created_by", "created_at",
					}).AddRow(
						uuid.NewString(),
						"1234567890128",
						"in",
						10.0,
						"restock",
						"",
						"user1",
						time.Now(),
					))
			},
		},
		{
			name: "filter by product and type",
			filter: domain.MovementFilter{
				ProductEan13: "1234567890128",
				Type:         "in",
			},
			wantLen: 1,
			before: func(t *testing.T, sm sqlmock.Sqlmock) {
				sm.ExpectQuery(`.*FROM movements.*`).
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "product_ean13", "type", "quantity",
						"reason", "notes", "created_by", "created_at",
					}).AddRow(
						uuid.NewString(),
						"1234567890128",
						"in",
						10.0,
						"restock",
						"",
						"user1",
						time.Now(),
					))
			},
		},
		{
			name: "query error",
			filter: domain.MovementFilter{
				ProductEan13: "1234567890128",
			},
			before: func(t *testing.T, sm sqlmock.Sqlmock) {
				sm.ExpectQuery(`.*FROM movements.*`).
					WillReturnError(assert.AnError)
			},
			wantErr: true,
			wantSub: "list movements",
		},
		{
			name:    "no rows returns empty slice",
			filter:  domain.MovementFilter{},
			wantLen: 0,
			before: func(t *testing.T, sm sqlmock.Sqlmock) {
				sm.ExpectQuery(`.*FROM movements.*`).
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "product_ean13", "type", "quantity",
						"reason", "notes", "created_by", "created_at",
					}))
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			m, _, sm := newMockDB(t)
			tt.before(t, sm)

			s := database.NewMovementRepository(m)
			movements, err := s.List(context.Background(), tt.filter)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantSub)
				return
			}
			assert.NoError(t, err)
			assert.Len(t, movements, tt.wantLen)
		})
	}
}

func TestMovementRepository_CreateWithStockUpdate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		newStock float64
		before   func(t *testing.T, sm sqlmock.Sqlmock)
		wantErr  bool
		wantSub  string
	}{
		{
			name:     "success",
			newStock: 100.0,
			before: func(t *testing.T, sm sqlmock.Sqlmock) {
				sm.ExpectBegin()
				sm.ExpectExec(`INSERT INTO movements`).
					WillReturnResult(sqlmock.NewResult(0, 0))
				sm.ExpectExec(`UPDATE products SET current_stock`).
					WillReturnResult(sqlmock.NewResult(0, 1))
				sm.ExpectCommit()
			},
		},
		{
			name:     "insert fails rolls back",
			newStock: 100.0,
			before: func(t *testing.T, sm sqlmock.Sqlmock) {
				sm.ExpectBegin()
				sm.ExpectExec(`INSERT INTO movements`).
					WillReturnError(assert.AnError)
				sm.ExpectRollback()
			},
			wantErr: true,
			wantSub: "insert movement in tx",
		},
		{
			name:     "update stock fails rolls back",
			newStock: 100.0,
			before: func(t *testing.T, sm sqlmock.Sqlmock) {
				sm.ExpectBegin()
				sm.ExpectExec(`INSERT INTO movements`).
					WillReturnResult(sqlmock.NewResult(0, 0))
				sm.ExpectExec(`UPDATE products SET current_stock`).
					WillReturnError(assert.AnError)
				sm.ExpectRollback()
			},
			wantErr: true,
			wantSub: "update stock in tx",
		},
		{
			name:     "begin tx fails",
			newStock: 100.0,
			before: func(t *testing.T, sm sqlmock.Sqlmock) {
				sm.ExpectBegin().WillReturnError(assert.AnError)
			},
			wantErr: true,
			wantSub: "begin transaction",
		},
	}

	for _, tt := range tests {
		ttt := tt
		t.Run(ttt.name, func(t *testing.T) {
			t.Parallel()
			m, _, sm := newMockDB(t)
			ttt.before(t, sm)

			s := database.NewMovementRepository(m)

			movement := &domain.Movement{
				Id:           uuid.NewString(),
				ProductEan13: "1234567890128",
				Type:         "in",
				Quantity:     10.0,
				Reason:       "restock",
				Notes:        "notes",
				CreatedBy:    "user1",
				CreatedAt:    time.Now(),
			}

			err := s.CreateWithStockUpdate(context.Background(), movement, ttt.newStock)
			if ttt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), ttt.wantSub)
				return
			}
			assert.NoError(t, err)
		})
	}
}
