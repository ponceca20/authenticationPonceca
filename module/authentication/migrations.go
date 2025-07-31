package authentication

import (
	"fmt"
	"practicev2/module/authentication/models"

	"gorm.io/gorm"
)

// RunMigrations executes the database migrations for the authentication module.
func RunMigrations(db *gorm.DB) error {
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
	// Create composite index for role name within organization (MySQL compatible)
	err := db.Exec(`
		CREATE INDEX idx_role_org_name 
		ON role (organization_id, name)
	`).Error
	if err != nil {
		fmt.Printf("⚠️  Role organization name index may already exist: %v\n", err)
	}

	// Index for audit logs performance
	err = db.Exec(`
		CREATE INDEX idx_audit_log_timestamp_org 
		ON audit_log (organization_id, timestamp)
	`).Error
	if err != nil {
		fmt.Printf("⚠️  Audit log index may already exist: %v\n", err)
	}

	// Index for guest session cleanup
	err = db.Exec(`
		CREATE INDEX idx_guest_session_expires 
		ON guest_session (expires_at)
	`).Error
	if err != nil {
		fmt.Printf("⚠️  Guest session expiry index may already exist: %v\n", err)
	}

	// Composite index for organizational membership (identity + organization)
	err = db.Exec(`
		CREATE INDEX idx_org_membership_identity_org 
		ON organizational_membership (identity_id, organization_id, is_active)
	`).Error
	if err != nil {
		fmt.Printf("⚠️  Organizational membership index may already exist: %v\n", err)
	}

	return nil
}
