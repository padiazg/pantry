package services

import (
	"context"
	"fmt"

	"github.com/padiazg/pantry/internal/core/domain"
)

// GetMovementsService handles movement history queries.
type GetMovementsService struct {
	movements domain.MovementRepository
}

// NewGetMovementsService creates a new GetMovementsService.
func NewGetMovementsService(movements domain.MovementRepository) *GetMovementsService {
	return &GetMovementsService{movements: movements}
}

// Execute returns a list of movements matching the given filter.
func (s *GetMovementsService) Execute(ctx context.Context, filter domain.MovementFilter) ([]*domain.Movement, error) {
	result, err := s.movements.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("listing movements: %w", err)
	}
	return result, nil
}

// GetByID returns a single movement by its UUID.
func (s *GetMovementsService) GetByID(ctx context.Context, id string) (*domain.Movement, error) {
	movement, err := s.movements.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("finding movement %s: %w", id, err)
	}
	return movement, nil
}
