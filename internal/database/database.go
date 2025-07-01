package database

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"venturo-core/configs"

	"github.com/golang-migrate/migrate/v4"
	mysqlMigrate "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// ConnectDB connects to the database using the provided configuration.
func ConnectDB(config *configs.Config) {
	var err error

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

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}

	slog.Info("Database connection successful.")
}

// newMigrate creates a new migrate instance.
func newMigrate() (*migrate.Migrate, error) {
	if DB == nil {
		return nil, errors.New("database connection is not initialized")
	}

	// Call the DB() method to get the underlying *sql.DB instance
	sqlDB, err := DB.DB()
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
func MigrateUp() {
	m, err := newMigrate()
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
func MigrateDown() {
	m, err := newMigrate()
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
func Drop() {
	m, err := newMigrate()
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
