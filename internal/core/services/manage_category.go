package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/padiazg/pantry/internal/core/domain"
)

// CreateCategoryInput holds the data needed to create a new category.
type CreateCategoryInput struct {
	Name        string
	Description string
}

// UpdateCategoryInput holds the data needed to update an existing category.
type UpdateCategoryInput struct {
	ID          string
	Name        string
	Description string
}

// ManageCategoryService handles category lifecycle: create, update, and retrieval.
type ManageCategoryService struct {
	categories domain.CategoryRepository
}

// NewManageCategoryService creates a new ManageCategoryService.
func NewManageCategoryService(categories domain.CategoryRepository) *ManageCategoryService {
	return &ManageCategoryService{categories: categories}
}

// Create validates and persists a new category with an auto-generated UUID.
func (s *ManageCategoryService) Create(ctx context.Context, input CreateCategoryInput) (*domain.Category, error) {
	category, err := domain.NewCategory(uuid.NewString(), input.Name, input.Description)
	if err != nil {
		return nil, fmt.Errorf("invalid category: %w", err)
	}
	if err := s.categories.Create(ctx, category); err != nil {
		return nil, fmt.Errorf("creating category: %w", err)
	}
	return category, nil
}

// Update replaces the mutable fields of an existing category.
func (s *ManageCategoryService) Update(ctx context.Context, input UpdateCategoryInput) (*domain.Category, error) {
	category, err := s.categories.FindByID(ctx, input.ID)
	if err != nil {
		return nil, fmt.Errorf("finding category: %w", err)
	}
	category.Name = input.Name
	category.Description = input.Description
	if err := category.Validate(); err != nil {
		return nil, fmt.Errorf("invalid category: %w", err)
	}
	if err := s.categories.Update(ctx, category); err != nil {
		return nil, fmt.Errorf("updating category: %w", err)
	}
	return category, nil
}

// GetByID returns a single category by its ID.
func (s *ManageCategoryService) GetByID(ctx context.Context, id string) (*domain.Category, error) {
	return s.categories.FindByID(ctx, id)
}

// List returns all categories ordered by name.
func (s *ManageCategoryService) List(ctx context.Context) ([]*domain.Category, error) {
	return s.categories.List(ctx)
}
