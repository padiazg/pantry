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

func TestProduct_ApplyMovement(t *testing.T) {
	p, _ := domain.NewProduct(validEAN, "Arroz", "", "kg", 0, "")
	p.CurrentStock = 10.0

	require.NoError(t, p.ApplyMovement(domain.MovementTypeIn, 5.0))
	assert.Equal(t, 15.0, p.CurrentStock)

	require.NoError(t, p.ApplyMovement(domain.MovementTypeOut, 3.0))
	assert.Equal(t, 12.0, p.CurrentStock)

	err := p.ApplyMovement(domain.MovementTypeOut, 100.0)
	require.Error(t, err, "should fail: insufficient stock")
	assert.Equal(t, 12.0, p.CurrentStock, "stock unchanged on failure")
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
