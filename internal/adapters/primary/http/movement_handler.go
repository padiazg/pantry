package httpserver

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/padiazg/pantry/internal/core/domain"
	"github.com/padiazg/pantry/internal/core/services"
)

// MovementHandler handles HTTP requests for stock movements.
type MovementHandler struct {
	record *services.RecordMovementService
	get    *services.GetMovementsService
}

// NewMovementHandler creates a new MovementHandler.
func NewMovementHandler(record *services.RecordMovementService, get *services.GetMovementsService) *MovementHandler {
	return &MovementHandler{record: record, get: get}
}

// --- DTOs ---

type createMovementRequest struct {
	ProductEan13 string  `json:"product_ean13"`
	Type         string  `json:"type"`
	Quantity     float64 `json:"quantity"`
	Reason       string  `json:"reason"`
	Notes        string  `json:"notes"`
	CreatedBy    string  `json:"created_by"`
}

type movementResponse struct {
	ID           string    `json:"id"`
	ProductEan13 string    `json:"product_ean13"`
	Type         string    `json:"type"`
	Quantity     float64   `json:"quantity"`
	Reason       string    `json:"reason,omitempty"`
	Notes        string    `json:"notes,omitempty"`
	CreatedBy    string    `json:"created_by,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

func toMovementResponse(m *domain.Movement) movementResponse {
	return movementResponse{
		ID:           m.Id,
		ProductEan13: m.ProductEan13,
		Type:         m.Type,
		Quantity:     m.Quantity,
		Reason:       m.Reason,
		Notes:        m.Notes,
		CreatedBy:    m.CreatedBy,
		CreatedAt:    m.CreatedAt,
	}
}

// --- Handlers ---

// List handles GET /api/v1/movements
func (h *MovementHandler) List(w http.ResponseWriter, r *http.Request) {
	filter := domain.MovementFilter{
		ProductEan13: r.URL.Query().Get("ean13"),
		Type:         r.URL.Query().Get("type"),
	}
	if v := r.URL.Query().Get("from"); v != "" {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			filter.From = t
		}
	}
	if v := r.URL.Query().Get("to"); v != "" {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			filter.To = t
		}
	}

	mvs, err := h.get.Execute(r.Context(), filter)
	if err != nil {
		respondHTTPError(w, err)
		return
	}
	resp := make([]movementResponse, 0, len(mvs))
	for _, m := range mvs {
		resp = append(resp, toMovementResponse(m))
	}
	respondJSON(w, http.StatusOK, resp)
}

// Create handles POST /api/v1/movements
func (h *MovementHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createMovementRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	out, err := h.record.Execute(r.Context(), services.RecordMovementInput{
		ProductEan13: req.ProductEan13,
		Type:         req.Type,
		Quantity:     req.Quantity,
		Reason:       req.Reason,
		Notes:        req.Notes,
		CreatedBy:    req.CreatedBy,
	})
	if err != nil {
		respondError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, toMovementResponse(out.Movement))
}

// GetByID handles GET /api/v1/movements/{id}
func (h *MovementHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	mv, err := h.get.GetByID(r.Context(), id)
	if err != nil {
		respondHTTPError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, toMovementResponse(mv))
}
