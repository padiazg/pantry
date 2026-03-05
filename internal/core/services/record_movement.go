package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/padiazg/pantry/internal/core/domain"
)

// RecordMovementInput holds the data needed to record a stock movement.
type RecordMovementInput struct {
	ProductEan13 string
	Type         string // "in" or "out"
	Quantity     float64
	Reason       string
	Notes        string
	CreatedBy    string
}

// RecordMovementOutput contains the persisted movement and updated product.
type RecordMovementOutput struct {
	Movement *domain.Movement
	Product  *domain.Product
}

// RecordMovementService records inbound or outbound stock movements.
// Business rules enforced:
//   - EAN-13 must be valid.
//   - Quantity must be > 0.
//   - Outbound movements cannot make CurrentStock negative.
//   - Movement creation and stock update happen in a single atomic transaction.
//   - Movement ID is a UUID v4 generated here, not delegated to the database.
type RecordMovementService struct {
	products  domain.ProductRepository
	movements domain.MovementRepository
}

// NewRecordMovementService creates a new RecordMovementService.
func NewRecordMovementService(products domain.ProductRepository, movements domain.MovementRepository) *RecordMovementService {
	return &RecordMovementService{products: products, movements: movements}
}

// Execute validates and records a stock movement, updating product stock atomically.
func (s *RecordMovementService) Execute(ctx context.Context, input RecordMovementInput) (*RecordMovementOutput, error) {
	// Validate movement type.
	movType, err := domain.NewMovementType(input.Type)
	if err != nil {
		return nil, fmt.Errorf("invalid movement type: %w", err)
	}

	// Validate quantity (must be > 0).
	if _, err := domain.NewStockLevel(input.Quantity); err != nil {
		return nil, fmt.Errorf("invalid quantity: %w", err)
	}

	// Fetch product (also validates EAN-13 exists in the system).
	product, err := s.products.FindByEAN13(ctx, input.ProductEan13)
	if err != nil {
		return nil, fmt.Errorf("finding product: %w", err)
	}

	// Apply movement to product — enforces non-negative stock rule for "out".
	if err := product.ApplyMovement(movType, input.Quantity); err != nil {
		return nil, err
	}

	// Build the movement entity with a fresh UUID v4.
	movement, err := domain.NewMovement(
		uuid.NewString(),
		input.ProductEan13,
		movType,
		input.Quantity,
		input.Reason,
		input.Notes,
		input.CreatedBy,
	)
	if err != nil {
		return nil, fmt.Errorf("creating movement entity: %w", err)
	}

	// Persist atomically: insert movement + update product stock.
	if err := s.movements.CreateWithStockUpdate(ctx, movement, product.CurrentStock); err != nil {
		return nil, fmt.Errorf("recording movement: %w", err)
	}

	return &RecordMovementOutput{Movement: movement, Product: product}, nil
}
