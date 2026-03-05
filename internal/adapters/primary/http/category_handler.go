package httpserver

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/padiazg/pantry/internal/core/domain"
	"github.com/padiazg/pantry/internal/core/services"
)

// CategoryHandler handles HTTP requests for categories.
type CategoryHandler struct {
	manage *services.ManageCategoryService
}

// NewCategoryHandler creates a new CategoryHandler.
func NewCategoryHandler(manage *services.ManageCategoryService) *CategoryHandler {
	return &CategoryHandler{manage: manage}
}

// --- DTOs ---

type createCategoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type updateCategoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type categoryResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

func toCategoryResponse(c *domain.Category) categoryResponse {
	return categoryResponse{
		ID:          c.Id,
		Name:        c.Name,
		Description: c.Description,
		CreatedAt:   c.CreatedAt,
	}
}

// --- Handlers ---

// List handles GET /api/v1/categories
func (h *CategoryHandler) List(w http.ResponseWriter, r *http.Request) {
	cats, err := h.manage.List(r.Context())
	if err != nil {
		respondHTTPError(w, err)
		return
	}
	resp := make([]categoryResponse, 0, len(cats))
	for _, c := range cats {
		resp = append(resp, toCategoryResponse(c))
	}
	respondJSON(w, http.StatusOK, resp)
}

// Create handles POST /api/v1/categories
func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	cat, err := h.manage.Create(r.Context(), services.CreateCategoryInput{
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		respondError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, toCategoryResponse(cat))
}

// GetByID handles GET /api/v1/categories/{id}
func (h *CategoryHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	cat, err := h.manage.GetByID(r.Context(), id)
	if err != nil {
		respondHTTPError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, toCategoryResponse(cat))
}

// Update handles PUT /api/v1/categories/{id}
func (h *CategoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req updateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	cat, err := h.manage.Update(r.Context(), services.UpdateCategoryInput{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		respondHTTPError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, toCategoryResponse(cat))
}
