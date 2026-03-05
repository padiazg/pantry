package domain_test

import (
	"testing"

	"github.com/padiazg/pantry/internal/core/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewMovementType(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{name: "in", value: "in", wantErr: false},
		{name: "out", value: "out", wantErr: false},
		{name: "empty", value: "", wantErr: true},
		{name: "invalid value", value: "transfer", wantErr: true},
		{name: "uppercase IN", value: "IN", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vo, err := domain.NewMovementType(tt.value)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.value, vo.String())
			}
		})
	}
}

func TestMovementType_Equals(t *testing.T) {
	assert.True(t, domain.MovementTypeIn.Equals(domain.MovementTypeIn))
	assert.True(t, domain.MovementTypeOut.Equals(domain.MovementTypeOut))
	assert.False(t, domain.MovementTypeIn.Equals(domain.MovementTypeOut))
}

func TestMovementType_Validate(t *testing.T) {
	tests := []struct {
		name    string
		vo      domain.MovementType
		wantErr bool
	}{
		{name: "in", vo: domain.MovementTypeIn, wantErr: false},
		{name: "out", vo: domain.MovementTypeOut, wantErr: false},
		{name: "invalid", vo: domain.MovementType("bad"), wantErr: true},
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
