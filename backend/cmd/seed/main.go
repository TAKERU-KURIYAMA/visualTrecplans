package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	"github.com/visualtrecplans/backend/internal/database"
	"github.com/visualtrecplans/backend/pkg/config"
	"github.com/visualtrecplans/backend/pkg/logger"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Could not load .env file: %v", err)
	}

	// Initialize logger
	logger.Init("debug")

	// Load configuration
	cfg := config.Load()

	// Connect to database
	if err := database.Connect(cfg); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Run database migrations first
	if err := database.Migrate(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	fmt.Println("🌱 Starting database seeding...")

	// Get database instance
	db := database.GetDB()

	// Find and execute all seed files
	seedsDir := "seeds"
	if _, err := os.Stat(seedsDir); os.IsNotExist(err) {
		log.Fatalf("Seeds directory not found: %s", seedsDir)
	}

	// Read all SQL files in seeds directory
	files, err := filepath.Glob(filepath.Join(seedsDir, "*.sql"))
	if err != nil {
		log.Fatalf("Failed to read seed files: %v", err)
	}

	if len(files) == 0 {
		fmt.Println("No seed files found in seeds directory")
		return
	}

	// Sort files to ensure correct execution order
	// Files should be named like 001_*, 002_*, etc.
	for _, file := range files {
		fmt.Printf("📥 Executing seed file: %s\n", filepath.Base(file))
		
		if err := executeSeedFile(db, file); err != nil {
			log.Fatalf("Failed to execute seed file %s: %v", file, err)
		}
		
		fmt.Printf("✅ Successfully executed: %s\n", filepath.Base(file))
	}

	fmt.Println("🎉 Database seeding completed successfully!")
}

// executeSeedFile reads and executes a SQL seed file
func executeSeedFile(db interface{ Exec(sql string, values ...interface{}) interface{ Error() error } }, filename string) error {
	// Read the SQL file
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Split content by semicolons to handle multiple statements
	sqlContent := string(content)
	statements := splitSQLStatements(sqlContent)

	// Execute each statement
	for i, statement := range statements {
		statement = strings.TrimSpace(statement)
		if statement == "" || strings.HasPrefix(statement, "--") {
			continue // Skip empty lines and comments
		}

		result := db.Exec(statement)
		if result.Error() != nil {
			return fmt.Errorf("failed to execute statement %d: %w", i+1, result.Error())
		}
	}

	return nil
}

// splitSQLStatements splits SQL content into individual statements
func splitSQLStatements(content string) []string {
	var statements []string
	var currentStatement strings.Builder
	
	scanner := bufio.NewScanner(strings.NewReader(content))
	
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		
		// Skip comments and empty lines
		if line == "" || strings.HasPrefix(line, "--") {
			continue
		}
		
		currentStatement.WriteString(line)
		currentStatement.WriteString(" ")
		
		// If line ends with semicolon, it's the end of a statement
		if strings.HasSuffix(line, ";") {
			stmt := strings.TrimSuffix(strings.TrimSpace(currentStatement.String()), ";")
			if stmt != "" {
				statements = append(statements, stmt)
			}
			currentStatement.Reset()
		}
	}
	
	// Add any remaining statement
	if currentStatement.Len() > 0 {
		stmt := strings.TrimSpace(currentStatement.String())
		if stmt != "" {
			statements = append(statements, stmt)
		}
	}
	
	return statements
}