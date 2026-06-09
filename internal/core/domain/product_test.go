package domain_test

import (
	"testing"

	"github.com/padiazg/pantry/internal/core/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const validEAN = "1234567890128"

func TestNewProduct(t *testing.T) {
	tests := []struct {
		name        string
		ean13       string
		productName string
		unit        string
		minStock    float64
		wantErr     bool
	}{
		{name: "valid product", ean13: validEAN, productName: "Arroz", unit: "kg", minStock: 2.0, wantErr: false},
		{name: "invalid ean13", ean13: "1234567890120", productName: "Arroz", unit: "kg", minStock: 0, wantErr: true},
		{name: "empty name", ean13: validEAN, productName: "", unit: "kg", minStock: 0, wantErr: true},
		{name: "empty unit", ean13: validEAN, productName: "Arroz", unit: "", minStock: 0, wantErr: true},
		{name: "negative min stock", ean13: validEAN, productName: "Arroz", unit: "kg", minStock: -1, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entity, err := domain.NewProduct(tt.ean13, tt.productName, "", tt.unit, tt.minStock, "")
			if tt.wantErr {
				require.Error(t, err)
				assert.Nil(t, entity)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, entity)
				assert.True(t, entity.Active)
				assert.Equal(t, float64(0), entity.CurrentStock)
			}
		})
	}
}

func TestProduct_IsLowStock(t *testing.T) {
	p, _ := domain.NewProduct(validEAN, "Arroz", "", "kg", 5.0, "")
	assert.True(t, p.IsLowStock(), "currentStock=0 <= minStock=5 → low stock")

	p.CurrentStock = 5.0
	assert.True(t, p.IsLowStock(), "currentStock=5 == minStock=5 → low stock")

	p.CurrentStock = 6.0
	assert.False(t, p.IsLowStock(), "currentStock=6 > minStock=5 → not low stock")
}

func TestProduct_Activate_Deactivate(t *testing.T) {
	p, _ := domain.NewProduct(validEAN, "Arroz", "", "kg", 0, "")
	require.True(t, p.Active)

	p.Deactivate()
	assert.False(t, p.Active)

	p.Activate()
	assert.True(t, p.Active)
}

func TestProduct_Validate(t *testing.T) {
	tests := []struct {
		name    string
		entity  *domain.Product
		wantErr bool
	}{
		{name: "valid", entity: &domain.Product{Ean13: validEAN, Name: "X", Unit: "u"}, wantErr: false},
		{name: "bad ean13", entity: &domain.Product{Ean13: "bad", Name: "X", Unit: "u"}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.entity.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

type checkProductApplyMovementFn func(*testing.T, *domain.Product, error)

var checkProductApplyMovement = func(fns ...checkProductApplyMovementFn) []checkProductApplyMovementFn { return fns }

func checkApplyMovementError(want string) checkProductApplyMovementFn {
	return func(t *testing.T, p *domain.Product, err error) {
		t.Helper()
		if want == "" {
			assert.NoErrorf(t, err, "checkApplyMovementError: expected no error, got %v", err)
			return
		}
		if assert.Errorf(t, err, "checkApplyMovementError: expected error %q", want) {
			assert.Containsf(t, err.Error(), want, "checkApplyMovementError mismatch")
		}
	}
}

func TestProduct_ApplyMovement(t *testing.T) {
	tests := []struct {
		name     string
		movType  domain.MovementType
		quantity float64
		checks   []checkProductApplyMovementFn
		before   func(*domain.Product)
	}{
		{
			name:     "inbound movement increases stock",
			movType:  domain.MovementTypeIn,
			quantity: 5.0,
			checks: checkProductApplyMovement(
				checkApplyMovementError(""),
				func(t *testing.T, p *domain.Product, err error) {
					t.Helper()
					assert.Equal(t, 5.0, p.CurrentStock)
				},
			),
		},
		{
			name:     "outbound movement decreases stock",
			movType:  domain.MovementTypeOut,
			quantity: 3.0,
			before: func(p *domain.Product) {
				p.CurrentStock = 10.0
			},
			checks: checkProductApplyMovement(
				checkApplyMovementError(""),
				func(t *testing.T, p *domain.Product, err error) {
					t.Helper()
					assert.Equal(t, 7.0, p.CurrentStock)
				},
			),
		},
		{
			name:     "outbound movement from zero stock fails",
			movType:  domain.MovementTypeOut,
			quantity: 1.0,
			checks: checkProductApplyMovement(
				checkApplyMovementError("insufficient stock"),
			),
		},
		{
			name:     "outbound exceeding available stock fails",
			movType:  domain.MovementTypeOut,
			quantity: 100.0,
			before: func(p *domain.Product) {
				p.CurrentStock = 12.0
			},
			checks: checkProductApplyMovement(
				checkApplyMovementError("insufficient stock"),
				func(t *testing.T, p *domain.Product, err error) {
					t.Helper()
					assert.Equal(t, 12.0, p.CurrentStock, "stock unchanged on failure")
				},
			),
		},
		{
			name:     "unknown movement type fails",
			movType:  "unknown",
			quantity: 1.0,
			checks: checkProductApplyMovement(
				checkApplyMovementError("unknown movement type"),
			),
		},
		{
			name:     "inbound with zero quantity succeeds",
			movType:  domain.MovementTypeIn,
			quantity: 0.0,
			before: func(p *domain.Product) {
				p.CurrentStock = 5.0
			},
			checks: checkProductApplyMovement(
				checkApplyMovementError(""),
				func(t *testing.T, p *domain.Product, err error) {
					t.Helper()
					assert.Equal(t, 5.0, p.CurrentStock)
				},
			),
		},
		{
			name:     "quantity leaves zero stock",
			movType:  domain.MovementTypeOut,
			quantity: 1.0,
			before: func(p *domain.Product) {
				p.CurrentStock = 1.0
			},
			checks: checkProductApplyMovement(
				checkApplyMovementError(""),
				func(t *testing.T, p *domain.Product, err error) {
					t.Helper()
					assert.Equal(t, 0.0, p.CurrentStock)
				},
			),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			s, e := domain.NewProduct(validEAN, "Rice", "", "kg", 0.0, "")
			if e != nil {
				assert.Fail(t, e.Error())
				return
			}

			if tt.before != nil {
				tt.before(s)
			}

			err := s.ApplyMovement(tt.movType, tt.quantity)
			for _, c := range tt.checks {
				c(t, s, err)
			}
		})
	}
}
