package gastos

import (
	"practicev2/registry"
)

// The init function is run automatically when the package is imported.
// It's used here to register the module's database migrations.
func init() {
	// Register the migration for the gastos module with priority 2.
	// This ensures that the gastos tables are created after the authentication
	// module tables (which have priority 1).
	registry.RegisterMigration(2, migrateExpenseModels)
}
