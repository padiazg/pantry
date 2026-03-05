package domain_test

import (
	"testing"

	"github.com/padiazg/pantry/internal/core/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewEAN13(t *testing.T) {
	tests := []struct {
		name    string
		code    string
		wantErr bool
	}{
		{name: "all zeros", code: "0000000000000", wantErr: false},
		{name: "valid 1234567890128", code: "1234567890128", wantErr: false},
		{name: "valid 4006381333931", code: "4006381333931", wantErr: false},
		{name: "too short", code: "123456789012", wantErr: true},
		{name: "too long", code: "12345678901234", wantErr: true},
		{name: "wrong check digit", code: "1234567890120", wantErr: true},
		{name: "non-digit chars", code: "12345678901AB", wantErr: true},
		{name: "empty string", code: "", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vo, err := domain.NewEAN13(tt.code)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.code, vo.Value)
				assert.Equal(t, tt.code, vo.String())
			}
		})
	}
}

func TestEAN13_Equals(t *testing.T) {
	a, _ := domain.NewEAN13("1234567890128")
	b, _ := domain.NewEAN13("1234567890128")
	c, _ := domain.NewEAN13("0000000000000")

	assert.True(t, a.Equals(b))
	assert.False(t, a.Equals(c))
}

func TestEAN13_Validate(t *testing.T) {
	tests := []struct {
		name    string
		vo      domain.EAN13
		wantErr bool
	}{
		{name: "valid", vo: domain.EAN13{Value: "1234567890128"}, wantErr: false},
		{name: "invalid check digit", vo: domain.EAN13{Value: "1234567890120"}, wantErr: true},
		{name: "empty", vo: domain.EAN13{Value: ""}, wantErr: true},
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
