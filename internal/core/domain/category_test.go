package domain_test

import (
	"testing"

	"github.com/padiazg/pantry/internal/core/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCategory(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		catName string
		wantErr bool
	}{
		{name: "valid", id: "uuid-1", catName: "Lácteos", wantErr: false},
		{name: "empty id", id: "", catName: "Lácteos", wantErr: true},
		{name: "empty name", id: "uuid-1", catName: "", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entity, err := domain.NewCategory(tt.id, tt.catName, "desc")
			if tt.wantErr {
				require.Error(t, err)
				assert.Nil(t, entity)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, entity)
				assert.Equal(t, tt.id, entity.Id)
				assert.False(t, entity.CreatedAt.IsZero())
			}
		})
	}
}

func TestCategory_Validate(t *testing.T) {
	tests := []struct {
		name    string
		entity  *domain.Category
		wantErr bool
	}{
		{name: "valid", entity: &domain.Category{Id: "1", Name: "X"}, wantErr: false},
		{name: "empty id", entity: &domain.Category{Id: "", Name: "X"}, wantErr: true},
		{name: "empty name", entity: &domain.Category{Id: "1", Name: ""}, wantErr: true},
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
