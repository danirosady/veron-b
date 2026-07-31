package database

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) error {
	// Run GORM AutoMigrate for any model-level changes
	if err := db.AutoMigrate(); err != nil {
		return fmt.Errorf("auto migrate failed: %w", err)
	}

	// Ensure migrations table exists
	db.Exec(`CREATE TABLE IF NOT EXISTS schema_migrations (
		version VARCHAR(255) PRIMARY KEY,
		applied_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
	)`)

	// Get applied migrations
	var applied []string
	db.Raw(`SELECT version FROM schema_migrations ORDER BY applied_at`).Scan(&applied)
	appliedSet := make(map[string]bool)
	for _, v := range applied {
		appliedSet[v] = true
	}

	// Find migration files
	migrationDir := "migrations"
	entries, err := os.ReadDir(migrationDir)
	if err != nil {
		// If migrations dir doesn't exist, skip
		fmt.Printf("[MIGRATION] migrations/ directory not found, skipping SQL migrations\n")
		return nil
	}

	var migrationFiles []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".up.sql") {
			migrationFiles = append(migrationFiles, entry.Name())
		}
	}
	sort.Strings(migrationFiles)

	for _, filename := range migrationFiles {
		version := strings.TrimSuffix(filename, ".up.sql")

		if appliedSet[version] {
			continue
		}

		content, err := os.ReadFile(filepath.Join(migrationDir, filename))
		if err != nil {
			return fmt.Errorf("failed to read migration %s: %w", filename, err)
		}

		statements := splitStatements(string(content))
		for _, stmt := range statements {
			stmt = strings.TrimSpace(stmt)
			if stmt == "" {
				continue
			}
			if err := db.Exec(stmt).Error; err != nil {
				return fmt.Errorf("failed to execute migration %s: %w\nStatement: %s",
					version, err, truncate(stmt, 200))
			}
		}

		db.Exec(`INSERT INTO schema_migrations (version) VALUES (?)`, version)
		fmt.Printf("[MIGRATION] Applied: %s\n", version)
	}

	return nil
}

func splitStatements(sql string) []string {
	var stmts []string
	var current strings.Builder
	inString := false
	stringChar := byte(0)

	for i := 0; i < len(sql); i++ {
		c := sql[i]

		if !inString && (c == '\'' || c == '"') {
			inString = true
			stringChar = c
			current.WriteByte(c)
			continue
		}

		if inString && c == stringChar {
			if i+1 < len(sql) && sql[i+1] == stringChar {
				current.WriteByte(c)
				current.WriteByte(sql[i+1])
				i++
				continue
			}
			inString = false
			current.WriteByte(c)
			continue
		}

		if inString {
			current.WriteByte(c)
			continue
		}

		// Skip line comments
		if c == '-' && i+1 < len(sql) && sql[i+1] == '-' {
			for i < len(sql) && sql[i] != '\n' {
				i++
			}
			continue
		}

		// Skip block comments
		if c == '/' && i+1 < len(sql) && sql[i+1] == '*' {
			i += 2
			for i < len(sql)-1 {
				if sql[i] == '*' && sql[i+1] == '/' {
					i += 2
					break
				}
				i++
			}
			continue
		}

		if c == ';' {
			stmts = append(stmts, current.String())
			current.Reset()
			continue
		}

		current.WriteByte(c)
	}

	if s := strings.TrimSpace(current.String()); s != "" {
		stmts = append(stmts, s)
	}

	return stmts
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
}

func RunSeed(db *gorm.DB) error {
	return nil
}

func RunDownMigrations(db *gorm.DB) error {
	migrationDir := "migrations"
	entries, err := os.ReadDir(migrationDir)
	if err != nil {
		return err
	}

	var migrationFiles []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".down.sql") {
			migrationFiles = append(migrationFiles, entry.Name())
		}
	}
	sort.Sort(sort.Reverse(sort.StringSlice(migrationFiles)))

	for _, filename := range migrationFiles {
		version := strings.TrimSuffix(filename, ".down.sql")

		var count int64
		db.Raw(`SELECT COUNT(*) FROM schema_migrations WHERE version = ?`, version).Scan(&count)
		if count == 0 {
			continue
		}

		content, err := os.ReadFile(filepath.Join(migrationDir, filename))
		if err != nil {
			return fmt.Errorf("failed to read down migration %s: %w", filename, err)
		}

		statements := splitStatements(string(content))
		for _, stmt := range statements {
			stmt = strings.TrimSpace(stmt)
			if stmt == "" {
				continue
			}
			if err := db.Exec(stmt).Error; err != nil {
				fmt.Printf("[MIGRATION] Warning: down migration %s failed: %v\n", version, err)
			}
		}

		db.Exec(`DELETE FROM schema_migrations WHERE version = ?`, version)
		fmt.Printf("[MIGRATION] Rolled back: %s\n", version)
	}

	return nil
}
