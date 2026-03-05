package httpserver

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/padiazg/pantry/internal/core/domain"
	"github.com/padiazg/pantry/internal/core/services"
)

// ProductHandler handles HTTP requests for products and stock reports.
type ProductHandler struct {
	manage    *services.ManageProductService
	get       *services.GetProductService
	movements *services.GetMovementsService
	report    *services.GetStockReportService
}

// NewProductHandler creates a new ProductHandler.
func NewProductHandler(
	manage *services.ManageProductService,
	get *services.GetProductService,
	movements *services.GetMovementsService,
	report *services.GetStockReportService,
) *ProductHandler {
	return &ProductHandler{manage: manage, get: get, movements: movements, report: report}
}

// --- DTOs ---

type createProductRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Unit        string  `json:"unit"`
	MinStock    float64 `json:"min_stock"`
	CategoryID  string  `json:"category_id"`
}

type updateProductRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Unit        string  `json:"unit"`
	MinStock    float64 `json:"min_stock"`
	CategoryID  string  `json:"category_id"`
}

type productResponse struct {
	EAN13        string    `json:"ean13"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Unit         string    `json:"unit"`
	MinStock     float64   `json:"min_stock"`
	CurrentStock float64   `json:"current_stock"`
	CategoryID   string    `json:"category_id,omitempty"`
	Active       bool      `json:"active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func toProductResponse(p *domain.Product) productResponse {
	return productResponse{
		EAN13:        p.Ean13,
		Name:         p.Name,
		Description:  p.Description,
		Unit:         p.Unit,
		MinStock:     p.MinStock,
		CurrentStock: p.CurrentStock,
		CategoryID:   p.CategoryID,
		Active:       p.Active,
		CreatedAt:    p.CreatedAt,
		UpdatedAt:    p.UpdatedAt,
	}
}

// --- Handlers ---

// List handles GET /api/v1/products
func (h *ProductHandler) List(w http.ResponseWriter, r *http.Request) {
	filter := domain.ProductFilter{
		CategoryID: r.URL.Query().Get("category"),
		Active:     parseBoolQuery(r, "active"),
	}
	if r.URL.Query().Get("low_stock") == "true" {
		filter.LowStock = true
	}

	products, err := h.get.List(r.Context(), filter)
	if err != nil {
		respondHTTPError(w, err)
		return
	}

	resp := make([]productResponse, 0, len(products))
	for _, p := range products {
		resp = append(resp, toProductResponse(p))
	}
	respondJSON(w, http.StatusOK, resp)
}

// Create handles POST /api/v1/products
func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	ean13 := chi.URLParam(r, "ean13")

	product, err := h.manage.Create(r.Context(), services.CreateProductInput{
		EAN13:       ean13,
		Name:        req.Name,
		Description: req.Description,
		Unit:        req.Unit,
		MinStock:    req.MinStock,
		CategoryID:  req.CategoryID,
	})
	if err != nil {
		respondError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, toProductResponse(product))
}

// GetByEAN13 handles GET /api/v1/products/{ean13}
func (h *ProductHandler) GetByEAN13(w http.ResponseWriter, r *http.Request) {
	ean13 := chi.URLParam(r, "ean13")
	product, err := h.get.GetByEAN13(r.Context(), ean13)
	if err != nil {
		respondHTTPError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, toProductResponse(product))
}

// Update handles PUT /api/v1/products/{ean13}
func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	ean13 := chi.URLParam(r, "ean13")
	var req updateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	product, err := h.manage.Update(r.Context(), services.UpdateProductInput{
		EAN13:       ean13,
		Name:        req.Name,
		Description: req.Description,
		Unit:        req.Unit,
		MinStock:    req.MinStock,
		CategoryID:  req.CategoryID,
	})
	if err != nil {
		respondHTTPError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, toProductResponse(product))
}

// Deactivate handles DELETE /api/v1/products/{ean13} (soft delete)
func (h *ProductHandler) Deactivate(w http.ResponseWriter, r *http.Request) {
	ean13 := chi.URLParam(r, "ean13")
	if err := h.manage.SetActive(r.Context(), ean13, false); err != nil {
		respondHTTPError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// GetStock handles GET /api/v1/products/{ean13}/stock
func (h *ProductHandler) GetStock(w http.ResponseWriter, r *http.Request) {
	ean13 := chi.URLParam(r, "ean13")
	product, err := h.get.GetByEAN13(r.Context(), ean13)
	if err != nil {
		respondHTTPError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, map[string]any{
		"ean13":         product.Ean13,
		"current_stock": product.CurrentStock,
		"min_stock":     product.MinStock,
		"is_low_stock":  product.IsLowStock(),
	})
}

// GetMovements handles GET /api/v1/products/{ean13}/movements
func (h *ProductHandler) GetMovements(w http.ResponseWriter, r *http.Request) {
	ean13 := chi.URLParam(r, "ean13")
	filter := domain.MovementFilter{ProductEan13: ean13}
	mvs, err := h.movements.Execute(r.Context(), filter)
	if err != nil {
		respondHTTPError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, mvs)
}

// GetStockReport handles GET /api/v1/reports/stock
func (h *ProductHandler) GetStockReport(w http.ResponseWriter, r *http.Request) {
	report, err := h.report.Execute(r.Context())
	if err != nil {
		respondHTTPError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, report)
}

// GetLowStock handles GET /api/v1/reports/low-stock
func (h *ProductHandler) GetLowStock(w http.ResponseWriter, r *http.Request) {
	products, err := h.report.GetLowStock(r.Context())
	if err != nil {
		respondHTTPError(w, err)
		return
	}
	resp := make([]productResponse, 0, len(products))
	for _, p := range products {
		resp = append(resp, toProductResponse(p))
	}
	respondJSON(w, http.StatusOK, resp)
}
