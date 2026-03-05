/*
Copyright © 2026
*/
// Package httpserver provides the HTTP server implementation using Chi.
//
// This is a PRIMARY adapter (inbound).
// It receives HTTP requests and delegates to core services.
package httpserver

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/padiazg/pantry/internal/config"
	"github.com/padiazg/pantry/internal/observability"
	"github.com/padiazg/pantry/pkg/logger"
	srv "github.com/padiazg/pantry/pkg/server"
)

var _ srv.Server = (*server)(nil)

// server is the Chi-backed HTTP server.
type server struct {
	router     chi.Router
	httpSrv    *http.Server
	cfg        *config.Config
	log        logger.Logger
	metrics    *observability.PrometheusMetrics
	products   *ProductHandler
	categories *CategoryHandler
	movements  *MovementHandler
}

type ServerConfig struct {
	Config     *config.Config
	Logger     logger.Logger
	Metrics    *observability.PrometheusMetrics
	Products   *ProductHandler
	Categories *CategoryHandler
	Movements  *MovementHandler
}

// New creates and configures a new Chi HTTP server with all handlers wired in.
func New(cfg *ServerConfig) srv.Server {
	router := chi.NewRouter()
	router.Use(chimiddleware.Logger)
	router.Use(chimiddleware.Recoverer)

	s := &server{
		router: router,
		httpSrv: &http.Server{
			Addr:         fmt.Sprintf(":%d", cfg.Config.Server.Port),
			Handler:      router,
			ReadTimeout:  cfg.Config.Server.ReadTimeout,
			WriteTimeout: cfg.Config.Server.WriteTimeout,
			IdleTimeout:  cfg.Config.Server.IdleTimeout,
		},
		cfg:        cfg.Config,
		log:        cfg.Logger,
		metrics:    cfg.Metrics,
		products:   cfg.Products,
		categories: cfg.Categories,
		movements:  cfg.Movements,
	}

	s.setupRoutes()
	return s
}

// Run starts the HTTP server in a non-blocking goroutine.
func (s *server) Run(errChan chan<- error) {
	go func() {
		s.log.Info("Server listening on port %d", s.cfg.Server.Port)
		if err := s.httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()
}

// Stop gracefully shuts down the server.
func (s *server) Stop(ctx context.Context) error {
	return s.httpSrv.Shutdown(ctx)
}

// setupRoutes registers all API routes.
func (s *server) setupRoutes() {
	// Health check
	s.router.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"healthy","service":"pantry"}`)
	})

	s.router.Route("/api/v1", func(r chi.Router) {
		// Categories
		r.Route("/categories", func(r chi.Router) {
			r.Get("/", s.categories.List)
			r.Post("/", s.categories.Create)
			r.Get("/{id}", s.categories.GetByID)
			r.Put("/{id}", s.categories.Update)
		})

		// Products
		r.Route("/products", func(r chi.Router) {
			r.Get("/", s.products.List)
			r.Post("/{ean13}", s.products.Create)
			r.Get("/{ean13}", s.products.GetByEAN13)
			r.Put("/{ean13}", s.products.Update)
			r.Delete("/{ean13}", s.products.Deactivate)
			r.Get("/{ean13}/stock", s.products.GetStock)
			r.Get("/{ean13}/movements", s.products.GetMovements)
		})

		// Movements
		r.Route("/movements", func(r chi.Router) {
			r.Get("/", s.movements.List)
			r.Post("/", s.movements.Create)
			r.Get("/{id}", s.movements.GetByID)
		})

		// Reports
		r.Route("/reports", func(r chi.Router) {
			r.Get("/stock", s.products.GetStockReport)
			r.Get("/low-stock", s.products.GetLowStock)
		})
	})
}
