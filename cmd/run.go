/*
Copyright © 2026
*/
package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"
	httpserver "github.com/padiazg/pantry/internal/adapters/primary/http"
	"github.com/padiazg/pantry/internal/adapters/secondary/database"
	"github.com/padiazg/pantry/internal/core/services"
	"github.com/padiazg/pantry/internal/observability"
	"github.com/padiazg/pantry/pkg/logger"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Start the pantry HTTP server",
	Long: `Start the pantry HTTP API server with graceful shutdown support.

The server will listen for SIGINT (Ctrl+C) and SIGTERM signals
and perform a graceful shutdown with a configurable timeout.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := GetConfig()

		log := logger.New(&logger.Config{
			Level:  cfg.LogLevel,
			Format: cfg.LogFormat,
		})

		log.Info("Starting pantry HTTP server...")

		// ── Observability ─────────────────────────────────────────────────
		enableObservability, _ := cmd.Flags().GetBool("observability")
		var shutdownObservability func()
		var metrics *observability.PrometheusMetrics

		if enableObservability {
			healthChecker := observability.NewHealthChecker()
			metrics = observability.NewPrometheusMetrics()

			healthChecker.RegisterCheck("api", func(ctx context.Context) observability.ComponentHealth {
				return observability.ComponentHealth{
					Status:  observability.StatusHealthy,
					Message: "API server running",
				}
			})

			obsAddr, _ := cmd.Flags().GetString("observability-addr")
			obsServer := observability.NewServer(obsAddr, healthChecker, metrics)
			shutdownObservability = obsServer.StartInBackground(log)
			defer shutdownObservability()
		}

		// ── Database ──────────────────────────────────────────────────────
		db, err := sql.Open("postgres", cfg.Database.URL)
		if err != nil {
			return fmt.Errorf("opening database: %w", err)
		}
		defer db.Close()

		if err := db.Ping(); err != nil {
			log.Error("Database unreachable: %v — continuing without DB (set PANTRY_DATABASE_URL)", err)
		}

		// ── Repositories (secondary adapters) ────────────────────────────
		productRepo := database.NewProductRepository(db)
		categoryRepo := database.NewCategoryRepository(db)
		movementRepo := database.NewMovementRepository(db)

		// ── Services (core) ───────────────────────────────────────────────
		manageProduct := services.NewManageProductService(productRepo)
		getProduct := services.NewGetProductService(productRepo)
		manageCategory := services.NewManageCategoryService(categoryRepo)
		recordMovement := services.NewRecordMovementService(productRepo, movementRepo)
		getMovements := services.NewGetMovementsService(movementRepo)
		getStockReport := services.NewGetStockReportService(productRepo)

		// ── Handlers (primary adapters) ───────────────────────────────────
		productHandler := httpserver.NewProductHandler(manageProduct, getProduct, getMovements, getStockReport)
		categoryHandler := httpserver.NewCategoryHandler(manageCategory)
		movementHandler := httpserver.NewMovementHandler(recordMovement, getMovements)

		// ── HTTP server ───────────────────────────────────────────────────
		srv := httpserver.New(&httpserver.ServerConfig{
			Config:     cfg,
			Logger:     log,
			Metrics:    metrics,
			Products:   productHandler,
			Categories: categoryHandler,
			Movements:  movementHandler,
		})

		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
		errChan := make(chan error, 1)

		srv.Run(errChan)

		select {
		case sig := <-sigChan:
			log.Info("Received signal: %v, initiating graceful shutdown...", sig)

			shutdownCtx, shutdownCancel := context.WithTimeout(
				context.Background(),
				cfg.Server.ShutdownTimeout,
			)
			defer shutdownCancel()

			if err := srv.Stop(shutdownCtx); err != nil {
				log.Error("Server shutdown error: %v", err)
				return err
			}

			log.Info("Server stopped gracefully")
			return nil

		case err := <-errChan:
			log.Error("Server error: %v", err)
			return err
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().Bool("observability", false, "Enable observability server (health + metrics)")
	runCmd.Flags().String("observability-addr", ":8081", "Observability server address")
}
