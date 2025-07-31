package authentication

import (
	"practicev2/registry"
)

// The init function is run automatically when the package is imported.
// It's used here to register the module's database migrations.
func init() {
	// Register the migration for the authentication module with priority 1.
	// This ensures that the core tables for users, organizations, etc.,
	// are created before other modules that might depend on them.
	registry.RegisterMigration(1, RunMigrations)
}
