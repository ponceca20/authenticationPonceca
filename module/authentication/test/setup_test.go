package test

import (
	"practicev2/config"
	"practicev2/module/authentication"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// SetupTestDatabase initializes an in-memory SQLite database for testing purposes.
// It also initializes the application config and runs migrations.
func SetupTestDatabase(t *testing.T) *gorm.DB {
	// Initialize config with default values for tests
	config.Init()
	// Using "file::memory:?cache=shared" allows the in-memory database to be shared
	// between connections, which can be useful for complex tests.
	// For simple tests, "file::memory:" is sufficient.
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to in-memory database: %v", err)
	}

	// Run migrations for the authentication module
	err = authentication.RunMigrations(db)
	if err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	return db
}
