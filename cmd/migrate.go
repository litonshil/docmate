package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"docmate/config"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Database migration commands",
}

var migrateUpCmd = &cobra.Command{
	Use:   "up",
	Short: "Apply all up migrations",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runMigration("up")
	},
}

var migrateDownCmd = &cobra.Command{
	Use:   "down",
	Short: "Apply all down migrations",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runMigration("down")
	},
}

func RegisterMigrateCommand() {
	migrateCmd.AddCommand(migrateUpCmd)
	migrateCmd.AddCommand(migrateDownCmd)
	RegisterSubCommand(migrateCmd)
}

func runMigration(direction string) error {
	//config.Load()
	dbCfg := config.DB().Db
	if dbCfg == nil {
		return fmt.Errorf("DB config not loaded")
	}

	fmt.Printf("DSN: postgres://%s:****@%s:%d/%s\n", dbCfg.Username, dbCfg.Host, dbCfg.Port, dbCfg.Name)
	// Build PostgreSQL DSN
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		dbCfg.Username, dbCfg.Password, dbCfg.Host, dbCfg.Port, dbCfg.Name)

	// Dynamically determine migrations directory
	var migrationsDir string
	candidateDirs := []string{
		filepath.Join("..", "migrations"),
		filepath.Join(".", "migrations"),
	}
	found := false
	for _, dir := range candidateDirs {
		if fi, err := os.Stat(dir); err == nil && fi.IsDir() {
			migrationsDir = "file://" + dir
			found = true
			fmt.Printf("Migrations directory: %s\n", migrationsDir)
			break
		}
	}
	if !found {
		return fmt.Errorf("no migrations directory found in any of: %v", candidateDirs)
	}

	m, err := migrate.New(migrationsDir, dsn)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	defer func(m *migrate.Migrate) {
		err, _ := m.Close()
		if err != nil {
			fmt.Printf("Failed to close migrate instance: %v\n", err)
		}
	}(m)

	switch direction {
	case "up":
		fmt.Println("Applying all up migrations...")
		err = m.Up()
		if err != nil && !errors.Is(err, migrate.ErrNoChange) {
			return fmt.Errorf("migration up failed: %w", err)
		}
		fmt.Println("Migration up completed successfully.")
	case "down":
		fmt.Println("Applying all down migrations...")
		err = m.Down()
		if err != nil && !errors.Is(err, migrate.ErrNoChange) {
			return fmt.Errorf("migration down failed: %w", err)
		}
		fmt.Println("Migration down completed successfully.")
	default:
		return fmt.Errorf("invalid direction: %s", direction)
	}

	return nil
}
