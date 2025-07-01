package database

import (
	"errors"
	"fmt"
	"log"
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

	// Conditionally build the credentials part of the DSN
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
		log.Fatalf("🔥 Failed to connect to database: %v", err)
	}

	log.Println("✅ Database connection successful.")
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
		log.Fatalf("migration failed: %v", err)
	}
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("an error occurred while migrating up: %v", err)
	}
	log.Println("✅ Database migrated up successfully.")
}

// MigrateDown rolls back the last applied migration.
func MigrateDown() {
	m, err := newMigrate()
	if err != nil {
		log.Fatalf("migration failed: %v", err)
	}
	if err := m.Down(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("an error occurred while migrating down: %v", err)
	}
	log.Println("✅ Database migrated down successfully.")
}

// Drop deletes everything in the database.
func Drop() {
	m, err := newMigrate()
	if err != nil {
		log.Fatalf("migration failed: %v", err)
	}
	if err := m.Drop(); err != nil {
		log.Fatalf("an error occurred while dropping database: %v", err)
	}
	log.Println("✅ Database dropped successfully.")
}
