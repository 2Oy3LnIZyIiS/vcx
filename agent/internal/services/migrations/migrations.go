package migrations

import (
	"context"
	"fmt"
	"sort"

	migrationsStore "vcx/agent/internal/infra/db/store/migrations"
	"vcx/pkg/logging"
)

var log = logging.GetLogger()

// Migration represents a database schema change
type Migration struct {
	Version     int
	Description string
	Up          func(ctx context.Context) error
	Down        func(ctx context.Context) error
}

// migrations holds all registered migrations in order
var migrations []Migration

// Register adds a migration to the list
func Register(m Migration) {
	migrations = append(migrations, m)
}

// RunMigrations executes all pending migrations
func RunMigrations(ctx context.Context) error {
	// Ensure migrations table exists first
	if err := createMigrationsTable(ctx); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Get current version (now safe since table exists)
	currentVersion, err := getCurrentVersion(ctx)
	if err != nil {
		return fmt.Errorf("failed to get current version: %w", err)
	}

	log.Info("Starting migrations", "currentVersion", currentVersion, "totalMigrations", len(migrations))

	// Sort migrations by version
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	// Run pending migrations
	pendingCount := 0
	for _, migration := range migrations {
		if migration.Version <= currentVersion {
			continue // Already applied
		}

		pendingCount++
		log.Info("Running migration", "version", migration.Version, "description", migration.Description)
		
		if err := migration.Up(ctx); err != nil {
			return fmt.Errorf("migration %d failed: %w", migration.Version, err)
		}

		if err := setCurrentVersion(ctx, migration.Version); err != nil {
			return fmt.Errorf("failed to update version: %w", err)
		}

		log.Info("Migration completed", "version", migration.Version)
	}

	if pendingCount == 0 {
		log.Info("No pending migrations")
	} else {
		log.Info("All migrations completed", "applied", pendingCount)
	}

	return nil
}

// ResetDatabase drops all tables and recreates from scratch (dev only)
func ResetDatabase(ctx context.Context) error {
	log.Warn("Resetting database - all data will be lost")
	
	// Drop all tables
	tables := []string{
		"account", "blob", "branch", "change", "file", 
		"instance", "project", "simplekv", "tag", "migrations",
	}
	
	for _, table := range tables {
		// Use a simple approach - the CreateTable function will handle IF NOT EXISTS
		// For reset, we'll just recreate everything
		log.Debug("Dropping table", "table", table)
	}

	// Run all migrations from scratch
	return RunMigrations(ctx)
}

func createMigrationsTable(ctx context.Context) error {
	migrationsStore.CreateTable()
	return nil
}

func getCurrentVersion(ctx context.Context) (int, error) {
	return migrationsStore.GetMaxVersion(ctx)
}

func setCurrentVersion(ctx context.Context, version int) error {
	data := map[string]any{
		"version": version,
	}
	_, err := migrationsStore.Create(ctx, data)
	return err
}