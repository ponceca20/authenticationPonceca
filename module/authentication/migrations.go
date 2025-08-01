package authentication

import (
	"fmt"
	"practicev2/module/authentication/models"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// RunMigrations executes the database migrations for the authentication module.
func RunMigrations(db *gorm.DB) error {
	// Configurar logger silencioso para las migraciones
	originalLogger := db.Logger
	defer func() {
		db.Logger = originalLogger // Restaurar logger original
	}()

	// Usar logger silencioso solo para migraciones
	db.Logger = db.Logger.LogMode(logger.Silent)

	// GORM's AutoMigrate will create tables, missing foreign keys, constraints, columns, and indexes.
	// It will NOT delete unused columns to protect your data.
	err := db.AutoMigrate(
		// Core Identity & Organization
		&models.Identity{},
		&models.Organization{},
		&models.Permission{},
		&models.Role{},
		&models.Department{},
		&models.OrganizationalMembership{},

		// Invitations
		&models.Invitation{},

		// E-commerce
		&models.CustomerProfile{},
		&models.ShippingAddress{},
		&models.CustomerPreferences{},
		&models.GuestSession{},

		// User Profile
		&models.UserProfile{},

		// Security & Auditing
		&models.RefreshToken{},
		&models.PasswordResetToken{},
		&models.AuditLog{},
	)

	if err != nil {
		return fmt.Errorf("authentication module migration failed: %w", err)
	}

	// Create additional composite indexes for better performance and data integrity
	err = createAdditionalIndexes(db)
	if err != nil {
		return fmt.Errorf("failed to create additional indexes: %w", err)
	}

	fmt.Println("✅ Authentication module migrated successfully.")
	return nil
}

// createAdditionalIndexes creates composite indexes for better performance and data integrity
func createAdditionalIndexes(db *gorm.DB) error {
	indexes := []indexDefinition{
		{
			Name:        "idx_role_org_name",
			Table:       "role",
			Query:       "CREATE INDEX idx_role_org_name ON role (organization_id, name)",
			Description: "Role organization name index",
		},
		{
			Name:        "idx_audit_log_timestamp_org",
			Table:       "audit_log",
			Query:       "CREATE INDEX idx_audit_log_timestamp_org ON audit_log (organization_id, timestamp)",
			Description: "Audit log timestamp organization index",
		},
		{
			Name:        "idx_guest_session_expires",
			Table:       "guest_session",
			Query:       "CREATE INDEX idx_guest_session_expires ON guest_session (expires_at)",
			Description: "Guest session expiry index",
		},
		{
			Name:        "idx_org_membership_identity_org",
			Table:       "organizational_membership",
			Query:       "CREATE INDEX idx_org_membership_identity_org ON organizational_membership (identity_id, organization_id, is_active)",
			Description: "Organizational membership composite index",
		},
	}

	// Crear índices de forma silenciosa
	createdCount := 0
	skippedCount := 0

	for _, idx := range indexes {
		if err := createIndexIfNotExists(db, idx, &createdCount, &skippedCount); err != nil {
			return fmt.Errorf("failed to create %s: %w", idx.Description, err)
		}
	}

	// Solo mostrar un resumen
	if createdCount > 0 || skippedCount > 0 {
		fmt.Printf("📊 Índices: %d creados, %d existentes\n", createdCount, skippedCount)
	}

	return nil
}

// indexDefinition represents an index to be created
type indexDefinition struct {
	Name        string
	Table       string
	Query       string
	Description string
}

// createIndexIfNotExists creates an index only if it doesn't already exist
func createIndexIfNotExists(db *gorm.DB, idx indexDefinition, createdCount, skippedCount *int) error {
	// Check if index already exists
	exists, err := indexExists(db, idx.Table, idx.Name)
	if err != nil {
		return fmt.Errorf("error checking if index %s exists: %w", idx.Name, err)
	}

	if exists {
		*skippedCount++
		return nil
	}

	// Create the index with silent logger
	originalLogger := db.Logger
	db.Logger = db.Logger.LogMode(logger.Silent)
	err = db.Exec(idx.Query).Error
	db.Logger = originalLogger

	if err != nil {
		// Check if the error is due to duplicate key name (index already exists)
		if isDuplicateIndexError(err) {
			*skippedCount++
			return nil
		}
		return fmt.Errorf("error creating index %s: %w", idx.Name, err)
	}

	*createdCount++
	return nil
}

// indexExists checks if an index exists on a table
func indexExists(db *gorm.DB, tableName, indexName string) (bool, error) {
	var count int64

	// Query to check if index exists in MySQL with silent logging
	query := `
		SELECT COUNT(*) 
		FROM information_schema.statistics 
		WHERE table_schema = DATABASE() 
		AND table_name = ? 
		AND index_name = ?
	`

	// Ejecutar con logger silencioso
	originalLogger := db.Logger
	db.Logger = db.Logger.LogMode(logger.Silent)
	err := db.Raw(query, tableName, indexName).Scan(&count).Error
	db.Logger = originalLogger

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// isDuplicateIndexError checks if the error is due to a duplicate index
func isDuplicateIndexError(err error) bool {
	if err == nil {
		return false
	}

	errStr := err.Error()
	// MySQL error codes for duplicate index/key
	duplicateErrors := []string{
		"Error 1061", // Duplicate key name
		"Error 1062", // Duplicate entry
		"Duplicate key name",
		"duplicate key name",
	}

	for _, dupErr := range duplicateErrors {
		if strings.Contains(errStr, dupErr) {
			return true
		}
	}

	return false
}
