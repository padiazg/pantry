/*
Copyright © 2026
*/
package cmd

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	infra "github.com/padiazg/pantry/internal/infrastructure/database"
	"github.com/padiazg/pantry/pkg/logger"
	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Database migration management",
	Long:  `Run, rollback or inspect database migrations.`,
}

var migrateUpCmd = &cobra.Command{
	Use:   "up",
	Short: "Apply all pending migrations",
	RunE: func(cmd *cobra.Command, args []string) error {
		return withMigrator(func(m *infra.Migrator) error {
			return m.Up()
		})
	},
}

var migrateDownCmd = &cobra.Command{
	Use:   "down",
	Short: "Roll back the last applied migration",
	RunE: func(cmd *cobra.Command, args []string) error {
		return withMigrator(func(m *infra.Migrator) error {
			return m.Down()
		})
	},
}

var migrateVersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show the current migration version",
	RunE: func(cmd *cobra.Command, args []string) error {
		return withMigrator(func(m *infra.Migrator) error {
			version, dirty, err := m.Version()
			if err != nil {
				return err
			}
			dirtyFlag := ""
			if dirty {
				dirtyFlag = " (dirty)"
			}
			fmt.Printf("Current migration version: %d%s\n", version, dirtyFlag)
			return nil
		})
	},
}

// withMigrator opens the DB, creates a Migrator and calls fn, then closes the DB.
func withMigrator(fn func(*infra.Migrator) error) error {
	cfg := GetConfig()
	log := logger.New(&logger.Config{
		Level:  cfg.LogLevel,
		Format: cfg.LogFormat,
	})

	db, err := sql.Open("postgres", cfg.Database.URL)
	if err != nil {
		return fmt.Errorf("opening database: %w", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		return fmt.Errorf("connecting to database: %w", err)
	}

	return fn(infra.NewMigrator(db, log))
}

func init() {
	migrateCmd.AddCommand(migrateUpCmd)
	migrateCmd.AddCommand(migrateDownCmd)
	migrateCmd.AddCommand(migrateVersionCmd)
	rootCmd.AddCommand(migrateCmd)
}
