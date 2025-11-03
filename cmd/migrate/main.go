package main

import (
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	"github.com/zjoart/eunoia/internal/config"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	cfg := config.LoadConfig()

	databaseURL := fmt.Sprintf("mysql://%s:%s@tcp(%s:%s)/%s",
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.Name,
	)

	m, err := migrate.New(
		"file://migrations",
		databaseURL,
	)
	if err != nil {
		log.Fatal("Failed to create migration instance:", err)
	}
	defer m.Close()

	if len(os.Args) < 2 {
		printUsage()
		return
	}

	command := os.Args[1]

	switch command {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatal("Failed to run migrations:", err)
		}
		fmt.Println("Migrations applied successfully")

	case "down":
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatal("Failed to rollback migrations:", err)
		}
		fmt.Println("Migrations rolled back successfully")

	case "version":
		version, dirty, err := m.Version()
		if err != nil {
			log.Fatal("Failed to get version:", err)
		}
		fmt.Printf("Current version: %d, Dirty: %v\n", version, dirty)

	case "force":
		if len(os.Args) < 3 {
			log.Fatal("Please provide version number: migrate force <version>")
		}
		var version int
		fmt.Sscanf(os.Args[2], "%d", &version)
		if err := m.Force(version); err != nil {
			log.Fatal("Failed to force version:", err)
		}
		fmt.Printf("Forced version to %d\n", version)

	case "steps":
		if len(os.Args) < 3 {
			log.Fatal("Please provide number of steps: migrate steps <n>")
		}
		var steps int
		fmt.Sscanf(os.Args[2], "%d", &steps)
		if err := m.Steps(steps); err != nil && err != migrate.ErrNoChange {
			log.Fatal("Failed to migrate steps:", err)
		}
		fmt.Printf("Migrated %d steps\n", steps)

	default:
		printUsage()
	}
}

func printUsage() {
	fmt.Println("Usage: migrate <command>")
	fmt.Println("\nCommands:")
	fmt.Println("  up              Apply all available migrations")
	fmt.Println("  down            Rollback all migrations")
	fmt.Println("  version         Show current migration version")
	fmt.Println("  steps <n>       Apply n migrations (use negative to rollback)")
	fmt.Println("  force <version> Force set migration version (use with caution)")
	fmt.Println("\nExamples:")
	fmt.Println("  migrate up")
	fmt.Println("  migrate down")
	fmt.Println("  migrate version")
	fmt.Println("  migrate steps 1")
	fmt.Println("  migrate steps -1")
}
