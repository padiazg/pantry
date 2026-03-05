package domain_test

import (
	"testing"

	"github.com/padiazg/pantry/internal/core/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewMovement(t *testing.T) {
	tests := []struct {
		name     string
		id       string
		ean13    string
		movType  domain.MovementType
		quantity float64
		wantErr  bool
	}{
		{name: "valid in", id: "uuid-1", ean13: validEAN, movType: domain.MovementTypeIn, quantity: 5.0, wantErr: false},
		{name: "valid out", id: "uuid-1", ean13: validEAN, movType: domain.MovementTypeOut, quantity: 2.5, wantErr: false},
		{name: "empty id", id: "", ean13: validEAN, movType: domain.MovementTypeIn, quantity: 1, wantErr: true},
		{name: "invalid ean13", id: "uuid-1", ean13: "bad", movType: domain.MovementTypeIn, quantity: 1, wantErr: true},
		{name: "zero quantity", id: "uuid-1", ean13: validEAN, movType: domain.MovementTypeIn, quantity: 0, wantErr: true},
		{name: "negative quantity", id: "uuid-1", ean13: validEAN, movType: domain.MovementTypeIn, quantity: -1, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entity, err := domain.NewMovement(tt.id, tt.ean13, tt.movType, tt.quantity, "", "", "user")
			if tt.wantErr {
				require.Error(t, err)
				assert.Nil(t, entity)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, entity)
				assert.Equal(t, string(tt.movType), entity.Type)
				assert.Equal(t, tt.movType, entity.MovementTypeVO())
			}
		})
	}
}

func TestMovement_Validate(t *testing.T) {
	tests := []struct {
		name    string
		entity  *domain.Movement
		wantErr bool
	}{
		{name: "valid", entity: &domain.Movement{Id: "1", ProductEan13: validEAN, Type: "in", Quantity: 1}, wantErr: false},
		{name: "empty id", entity: &domain.Movement{Id: "", ProductEan13: validEAN, Type: "in", Quantity: 1}, wantErr: true},
		{name: "zero qty", entity: &domain.Movement{Id: "1", ProductEan13: validEAN, Type: "in", Quantity: 0}, wantErr: true},
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
