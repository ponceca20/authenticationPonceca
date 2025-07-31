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
		&models.Role{},
		&models.Permission{},
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

	fmt.Println("✅ Authentication module migrated successfully.")
	return nil
}
