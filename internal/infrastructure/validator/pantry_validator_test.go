package validator_test

import (
	"testing"

	"github.com/padiazg/pantry/internal/infrastructure/validator"
	"github.com/stretchr/testify/assert"
)

func TestPantryValidator_Required(t *testing.T) {
	tests := []struct {
		name    string
		field   string
		value   string
		wantErr bool
	}{
		{
			name:    "valid value",
			field:   "name",
			value:   "John",
			wantErr: false,
		},
		{
			name:    "empty value",
			field:   "name",
			value:   "",
			wantErr: true,
		},
		{
			name:    "whitespace only",
			field:   "name",
			value:   "   ",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := validator.NewPantryValidator()
			v.Required(tt.field, tt.value)

			if tt.wantErr {
				assert.True(t, v.HasErrors())
			} else {
				assert.False(t, v.HasErrors())
			}
		})
	}
}

func TestPantryValidator_Email(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		{
			name:    "valid email",
			email:   "user@example.com",
			wantErr: false,
		},
		{
			name:    "invalid email - no @",
			email:   "userexample.com",
			wantErr: true,
		},
		{
			name:    "invalid email - no domain",
			email:   "user@",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := validator.NewPantryValidator()
			v.Email("email", tt.email)

			if tt.wantErr {
				assert.True(t, v.HasErrors())
			} else {
				assert.False(t, v.HasErrors())
			}
		})
	}
}

func TestPantryValidator_MinLength(t *testing.T) {
	v := validator.NewPantryValidator()

	v.MinLength("password", "short", 8)
	assert.True(t, v.HasErrors())

	v = validator.NewPantryValidator()
	v.MinLength("password", "longenough", 8)
	assert.False(t, v.HasErrors())
}
