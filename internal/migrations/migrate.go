package migrations

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"gorm.io/gorm"
)

// RunMigrations executes .sql files in database_schema folder in lexical order.
// It tracks applied migrations in a schema_migrations table to avoid reapplying.
func RunMigrations(db *gorm.DB) error {
	if err := ensureTrackingTable(db); err != nil {
		return err
	}
	files, err := collectSQLFiles("database_schema")
	if err != nil {
		return err
	}
	for _, f := range files {
		applied, err := isApplied(db, f)
		if err != nil {
			return err
		}
		if applied {
			continue
		}
		if err := applyMigrationFile(db, f); err != nil {
			return err
		}
	}
	return nil
}

func ensureTrackingTable(db *gorm.DB) error {
	return db.Exec(`
        CREATE TABLE IF NOT EXISTS schema_migrations (
            filename TEXT PRIMARY KEY,
            applied_at TIMESTAMP NOT NULL DEFAULT NOW()
        );
    `).Error
}

func collectSQLFiles(baseDir string) ([]string, error) {
	var files []string
	if err := filepath.WalkDir(baseDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if filepath.Ext(path) == ".sql" {
			files = append(files, path)
		}
		return nil
	}); err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, fmt.Errorf("walk %s: %w", baseDir, err)
	}
	sort.Strings(files)
	return files, nil
}

func applyMigrationFile(db *gorm.DB, filename string) error {
	sqlBytes, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("read %s: %w", filename, err)
	}

	content := string(sqlBytes)

	// Ambil hanya bagian Up
	if strings.Contains(content, "-- +goose Up") {
		parts := strings.Split(content, "-- +goose Down")
		upPart := strings.Split(parts[0], "-- +goose Up")
		if len(upPart) > 1 {
			content = upPart[1]
		} else {
			content = upPart[0]
		}
	}

	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec(content).Error; err != nil {
			return fmt.Errorf("exec %s: %w", filename, err)
		}
		if err := tx.Exec("INSERT INTO schema_migrations (filename, applied_at) VALUES (?, ?)", filename, time.Now()).Error; err != nil {
			return fmt.Errorf("record %s: %w", filename, err)
		}
		return nil
	})
}

func isApplied(db *gorm.DB, filename string) (bool, error) {
	var count int64
	if err := db.Table("schema_migrations").Where("filename = ?", filename).Count(&count).Error; err != nil {
		return false, fmt.Errorf("check applied %s: %w", filename, err)
	}
	return count > 0, nil
}
