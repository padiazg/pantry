package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMovementFilter_Args(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name   string
		idx    int
		want   string
		want2  []any
		before func(*MovementFilter)
	}{
		{
			name:   "empty filter returns empty",
			idx:    1,
			want:   "",
			want2:  nil,
			before: func(f *MovementFilter) {},
		},
		{
			name:  "filter by product ean13",
			idx:   1,
			want:  " AND product_ean13 = $1",
			want2: []any{"1234567890128"},
			before: func(f *MovementFilter) {
				f.ProductEan13 = "1234567890128"
			},
		},
		{
			name:  "filter by type",
			idx:   1,
			want:  " AND type = $1",
			want2: []any{"in"},
			before: func(f *MovementFilter) {
				f.Type = "in"
			},
		},
		{
			name:  "filter by from date",
			idx:   1,
			want:  " AND created_at >= $1",
			want2: []any{now},
			before: func(f *MovementFilter) {
				f.From = now
			},
		},
		{
			name:  "filter by to date",
			idx:   1,
			want:  " AND created_at <= $1",
			want2: []any{now},
			before: func(f *MovementFilter) {
				f.To = now
			},
		},
		{
			name:  "filter by product and type",
			idx:   1,
			want:  " AND product_ean13 = $1 AND type = $2",
			want2: []any{"1234567890128", "out"},
			before: func(f *MovementFilter) {
				f.ProductEan13 = "1234567890128"
				f.Type = "out"
			},
		},
		{
			name:  "filter by date range",
			idx:   1,
			want:  " AND created_at >= $1 AND created_at <= $2",
			want2: []any{now.Add(-24 * time.Hour), now},
			before: func(f *MovementFilter) {
				f.From = now.Add(-24 * time.Hour)
				f.To = now
			},
		},
		{
			name:  "all filters combined",
			idx:   1,
			want:  " AND product_ean13 = $1 AND type = $2 AND created_at >= $3 AND created_at <= $4",
			want2: []any{"1234567890128", "in", now.Add(-24 * time.Hour), now},
			before: func(f *MovementFilter) {
				f.ProductEan13 = "1234567890128"
				f.Type = "in"
				f.From = now.Add(-24 * time.Hour)
				f.To = now
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			s := &MovementFilter{}
			if tt.before != nil {
				tt.before(s)
			}
			r, r2 := s.Args()
			assert.Equal(t, tt.want, r)
			assert.Equal(t, tt.want2, r2)
		})
	}
}
