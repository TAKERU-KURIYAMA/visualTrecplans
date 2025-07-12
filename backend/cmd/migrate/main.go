package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"

	"github.com/visualtrecplans/backend/pkg/config"
	"github.com/visualtrecplans/backend/pkg/logger"
)

func main() {
	var (
		direction = flag.String("direction", "up", "Migration direction: up or down")
		steps     = flag.Int("steps", 0, "Number of migration steps (0 for all)")
		version   = flag.Int("version", -1, "Migrate to specific version")
		force     = flag.Int("force", -1, "Force migration to specific version")
		drop      = flag.Bool("drop", false, "Drop database (dangerous!)")
	)
	flag.Parse()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Connect to database
	dsn := cfg.GetDatabaseDSN()
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Test connection
	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	// Create migrate driver
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal("Failed to create migration driver:", err)
	}

	// Create migrate instance
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres",
		driver,
	)
	if err != nil {
		log.Fatal("Failed to create migrate instance:", err)
	}

	logger.Info("Migration tool started",
		logger.String("direction", *direction),
		logger.Int("steps", *steps),
		logger.Int("version", *version),
	)

	// Execute migration based on flags
	switch {
	case *drop:
		if err := confirmDrop(); err != nil {
			log.Fatal("Operation cancelled:", err)
		}
		if err := m.Drop(); err != nil {
			log.Fatal("Failed to drop database:", err)
		}
		logger.Info("Database dropped successfully")

	case *force != -1:
		if err := m.Force(*force); err != nil {
			log.Fatal("Failed to force migration:", err)
		}
		logger.Info("Migration forced successfully", logger.Int("version", *force))

	case *version != -1:
		if err := m.Migrate(uint(*version)); err != nil && err != migrate.ErrNoChange {
			log.Fatal("Failed to migrate to version:", err)
		}
		logger.Info("Migration completed successfully", logger.Int("version", *version))

	case *direction == "up":
		if *steps > 0 {
			if err := m.Steps(*steps); err != nil && err != migrate.ErrNoChange {
				log.Fatal("Failed to migrate up:", err)
			}
		} else {
			if err := m.Up(); err != nil && err != migrate.ErrNoChange {
				log.Fatal("Failed to migrate up:", err)
			}
		}
		logger.Info("Migration up completed successfully")

	case *direction == "down":
		if *steps > 0 {
			if err := m.Steps(-*steps); err != nil && err != migrate.ErrNoChange {
				log.Fatal("Failed to migrate down:", err)
			}
		} else {
			if err := m.Down(); err != nil && err != migrate.ErrNoChange {
				log.Fatal("Failed to migrate down:", err)
			}
		}
		logger.Info("Migration down completed successfully")

	default:
		flag.Usage()
		os.Exit(1)
	}

	// Get current version
	version_num, dirty, err := m.Version()
	if err != nil {
		logger.Warn("Failed to get migration version", logger.Error(err))
	} else {
		logger.Info("Current migration status",
			logger.Int("version", int(version_num)),
			logger.Bool("dirty", dirty),
		)
	}
}

func confirmDrop() error {
	fmt.Print("Are you sure you want to drop the entire database? This action cannot be undone. (yes/no): ")
	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		return err
	}

	if response != "yes" {
		return fmt.Errorf("operation cancelled by user")
	}

	return nil
}

func printUsage() {
	fmt.Println("Migration Tool Usage:")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  go run cmd/migrate/main.go -direction=up              # Migrate up (all)")
	fmt.Println("  go run cmd/migrate/main.go -direction=down            # Migrate down (all)")
	fmt.Println("  go run cmd/migrate/main.go -direction=up -steps=1     # Migrate up 1 step")
	fmt.Println("  go run cmd/migrate/main.go -direction=down -steps=1   # Migrate down 1 step")
	fmt.Println("  go run cmd/migrate/main.go -version=2                 # Migrate to version 2")
	fmt.Println("  go run cmd/migrate/main.go -force=1                   # Force to version 1")
	fmt.Println("  go run cmd/migrate/main.go -drop                      # Drop database")
	fmt.Println()
}