package domain_test

import (
	"testing"

	"github.com/padiazg/pantry/internal/core/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewStockLevel(t *testing.T) {
	tests := []struct {
		name    string
		value   float64
		wantErr bool
	}{
		{name: "positive integer", value: 10, wantErr: false},
		{name: "positive decimal", value: 0.5, wantErr: false},
		{name: "very small positive", value: 0.001, wantErr: false},
		{name: "zero", value: 0, wantErr: true},
		{name: "negative", value: -1, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vo, err := domain.NewStockLevel(tt.value)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.value, vo.Value)
			}
		})
	}
}

func TestStockLevel_Equals(t *testing.T) {
	a := domain.StockLevel{Value: 5.0}
	b := domain.StockLevel{Value: 5.0}
	c := domain.StockLevel{Value: 3.0}

	assert.True(t, a.Equals(b))
	assert.False(t, a.Equals(c))
}

func TestStockLevel_Validate(t *testing.T) {
	tests := []struct {
		name    string
		vo      domain.StockLevel
		wantErr bool
	}{
		{name: "valid", vo: domain.StockLevel{Value: 1.0}, wantErr: false},
		{name: "zero", vo: domain.StockLevel{Value: 0}, wantErr: true},
		{name: "negative", vo: domain.StockLevel{Value: -5}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.vo.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
