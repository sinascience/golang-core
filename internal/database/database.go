package database

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"venturo-core/configs"

	"github.com/golang-migrate/migrate/v4"
	mysqlMigrate "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Removed global DB variable - now using dependency injection

// sanitizeDSN removes password from DSN for safe logging
func sanitizeDSN(dsn string) string {
	// Replace password in DSN with ***
	parts := strings.Split(dsn, "@")
	if len(parts) < 2 {
		return dsn
	}
	
	userParts := strings.Split(parts[0], ":")
	if len(userParts) >= 2 {
		userParts[1] = "***"
		parts[0] = strings.Join(userParts, ":")
	}
	
	return strings.Join(parts, "@")
}

// ConnectDB connects to the database using the provided configuration and returns the DB instance.
func ConnectDB(config *configs.Config) *gorm.DB {
	credentials := config.DBUser
	if config.DBPassword != "" {
		credentials = fmt.Sprintf("%s:%s", config.DBUser, config.DBPassword)
	}

	dsn := fmt.Sprintf("%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		credentials,
		config.DBHost,
		config.DBPort,
		config.DBName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		// Sanitize DSN for logging (remove password)
		sanitizedDSN := sanitizeDSN(dsn)
		slog.Error("Failed to connect to database", "error", err, "dsn", sanitizedDSN)
		os.Exit(1)
	}

	slog.Info("Database connection successful.")
	return db
}

// newMigrate creates a new migrate instance.
func newMigrate(db *gorm.DB) (*migrate.Migrate, error) {
	if db == nil {
		return nil, errors.New("database connection is not initialized")
	}

	// Call the DB() method to get the underlying *sql.DB instance
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	driver, err := mysqlMigrate.WithInstance(sqlDB, &mysqlMigrate.Config{})
	if err != nil {
		return nil, err
	}
	return migrate.NewWithDatabaseInstance("file://database/migrations", "mysql", driver)
}

// MigrateUp applies all available up migrations.
func MigrateUp(db *gorm.DB) {
	m, err := newMigrate(db)
	if err != nil {
		slog.Error("Migration failed", "error", err)
		os.Exit(1)
	}
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		slog.Error("An error occurred while migrating up", "error", err)
		os.Exit(1)
	}
	slog.Info("Database migrated up successfully.")
}

// MigrateDown rolls back the last applied migration.
func MigrateDown(db *gorm.DB) {
	m, err := newMigrate(db)
	if err != nil {
		slog.Error("Migration failed", "error", err)
		os.Exit(1)
	}
	if err := m.Down(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		slog.Error("An error occurred while migrating down", "error", err)
		os.Exit(1)
	}
	slog.Info("Database migrated down successfully.")
}

// Drop deletes everything in the database.
func Drop(db *gorm.DB) {
	m, err := newMigrate(db)
	if err != nil {
		slog.Error("Migration failed", "error", err)
		os.Exit(1)
	}
	if err := m.Drop(); err != nil {
		slog.Error("An error occurred while dropping database", "error", err)
		os.Exit(1)
	}
	slog.Info("Database dropped successfully.")
}
