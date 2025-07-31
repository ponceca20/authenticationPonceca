package main

import (
	"fmt"
	"log"
	"os"
	"practicev2/config"
	"practicev2/database"
	"practicev2/module/authentication/models"
	"practicev2/registry"

	// Import authentication module to register migrations
	_ "practicev2/module/authentication"

	"gorm.io/gorm"
)

func main() {
	fmt.Println("🧪 Testing Database Migrations...")

	// Initialize configuration
	config.Init()

	// Set default database environment variables if not set
	setDefaultEnvVars()

	// Initialize database connection
	database.ConnectDatabase()
	db := database.DBconn

	if db == nil {
		log.Fatal("Failed to initialize database connection")
	}

	// Test migration
	fmt.Println("📋 Running migrations...")
	registry.RunMigrations(db)

	// Verify all tables exist
	verifyTables(db)

	// Test basic model operations
	testModelOperations(db)

	fmt.Println("✅ All migration tests passed!")
}

func setDefaultEnvVars() {
	// Set default database environment variables for testing if not already set
	envVars := map[string]string{
		"DB_USER":     "root",
		"DB_PASSWORD": "password",
		"DB_HOST":     "localhost",
		"DB_PORT":     "3306",
		"DB_NAME":     "baselogin_test",
	}

	for key, defaultValue := range envVars {
		if os.Getenv(key) == "" {
			os.Setenv(key, defaultValue)
			fmt.Printf("⚠️  Set default %s = %s\n", key, defaultValue)
		}
	}
}

func verifyTables(db *gorm.DB) {
	fmt.Println("🔍 Verifying table existence...")

	tables := []string{
		"identity", // Singular form due to SingularTable: true
		"organization",
		"permission",
		"role",
		"role_permission", // Many-to-many table
		"department",
		"organizational_membership",
		"invitation",
		"customer_profile",
		"shipping_address",
		"customer_preference",
		"guest_session",
		"user_profile",
		"refresh_token",
		"password_reset_token",
		"audit_log",
	}

	for _, table := range tables {
		if db.Migrator().HasTable(table) {
			fmt.Printf("  ✅ Table '%s' exists\n", table)
		} else {
			fmt.Printf("  ❌ Table '%s' missing\n", table)
		}
	}
}

func testModelOperations(db *gorm.DB) {
	fmt.Println("🧪 Testing basic model operations...")

	// Test creating a basic identity
	identity := models.Identity{
		ID:           "test-identity-123",
		Email:        "test@example.com",
		PasswordHash: "hashed_password",
		FirstName:    "Test",
		LastName:     "User",
	}

	// Test validation without actually saving
	if err := db.Create(&identity).Error; err != nil {
		fmt.Printf("  ❌ Identity creation failed: %v\n", err)
	} else {
		fmt.Println("  ✅ Identity creation works")
		// Clean up
		db.Delete(&identity)
	}

	// Test organization creation
	org := models.Organization{
		ID:   "test-org-123",
		Name: "Test Organization",
		Slug: "test-org",
		Type: "company",
	}

	if err := db.Create(&org).Error; err != nil {
		fmt.Printf("  ❌ Organization creation failed: %v\n", err)
	} else {
		fmt.Println("  ✅ Organization creation works")
		// Clean up
		db.Delete(&org)
	}

	// Test permission creation
	permission := models.Permission{
		Resource: "users",
		Action:   "create",
		Scope:    "organization",
	}

	if err := db.Create(&permission).Error; err != nil {
		fmt.Printf("  ❌ Permission creation failed: %v\n", err)
	} else {
		fmt.Println("  ✅ Permission creation works")
		// Clean up
		db.Delete(&permission)
	}

	// Test customer preferences creation (with proper foreign key)
	// First create an identity and customer profile
	testIdentity := models.Identity{
		ID:           "test-identity-456",
		Email:        "customer@example.com",
		PasswordHash: "hashed_password",
		FirstName:    "Customer",
		LastName:     "Test",
	}

	testCustomerProfile := models.CustomerProfile{
		ID:             "test-customer-123",
		IdentityID:     "test-identity-456",
		CustomerNumber: "CUST-12345",
	}

	customerPref := models.CustomerPreferences{
		ID:                 "test-pref-123",
		CustomerProfileID:  "test-customer-123",
		Theme:              "dark",
		Language:           "en",
		TimeZone:           "America/New_York",
		EmailNotifications: true,
		SmsNotifications:   false,
	}

	// Create in order: Identity -> CustomerProfile -> CustomerPreferences
	if err := db.Create(&testIdentity).Error; err == nil {
		if err := db.Create(&testCustomerProfile).Error; err == nil {
			if err := db.Create(&customerPref).Error; err != nil {
				fmt.Printf("  ❌ CustomerPreferences creation failed: %v\n", err)
			} else {
				fmt.Println("  ✅ CustomerPreferences creation works")
				// Clean up in reverse order
				db.Delete(&customerPref)
				db.Delete(&testCustomerProfile)
			}
		} else {
			fmt.Printf("  ❌ CustomerProfile creation failed: %v\n", err)
		}
		db.Delete(&testIdentity)
	} else {
		fmt.Printf("  ❌ Test Identity creation failed: %v\n", err)
	}
}
