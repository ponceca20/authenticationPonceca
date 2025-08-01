package admin

import (
	"practicev2/config"
	"practicev2/module/authentication/auth"
	"practicev2/module/authentication/models"
	"practicev2/module/authentication/organization"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type AdminTestSuite struct {
	suite.Suite
	db           *gorm.DB
	adminService AdminService
	authRepo     auth.AuthRepository
	orgService   organization.OrganizationService
}

// setupTestDatabase initializes an in-memory SQLite database for testing purposes.
func (suite *AdminTestSuite) setupTestDatabase() *gorm.DB {
	// Initialize config with default values for tests
	config.Init()

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		suite.T().Fatalf("Failed to connect to in-memory database: %v", err)
	}

	// Auto-migrate the models we need for testing
	err = db.AutoMigrate(
		&models.Identity{},
		&models.User{},
		&models.UserProfile{},
		&models.Organization{},
		&models.OrganizationalMembership{},
		&models.Role{},
		&models.RefreshToken{},
	)
	if err != nil {
		suite.T().Fatalf("Failed to run migrations: %v", err)
	}

	return db
}

func (suite *AdminTestSuite) SetupSuite() {
	suite.db = suite.setupTestDatabase()
	adminRepo := NewAdminRepository(suite.db)
	suite.adminService = NewAdminService(adminRepo)
	suite.authRepo = auth.NewAuthRepository(suite.db)

	// Setup organization service for creating test data
	orgRepo := organization.NewOrganizationRepository(suite.db)
	suite.orgService = organization.NewOrganizationService(orgRepo, suite.authRepo)
}

func TestAdminTestSuite(t *testing.T) {
	suite.Run(t, new(AdminTestSuite))
}

func (suite *AdminTestSuite) TestGetSystemStats() {
	// 1. Get initial stats (should be empty)
	stats, err := suite.adminService.GetSystemStats()
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), stats)
	assert.Equal(suite.T(), 0, stats.TotalUsers)
	assert.Equal(suite.T(), 0, stats.TotalOrganizations)
	assert.Equal(suite.T(), 0, stats.ActiveSessions)

	// 2. Create some test data
	orgDTO := &organization.OrganizationRegistrationDTO{
		Identity: auth.RegisterDTO{
			FirstName: "Admin",
			LastName:  "User",
			Email:     "admin.user@example.com",
			Password:  "Secure-Admin-Password123", // Fixed: Added uppercase, lowercase, and number
		},
		Name: "Admin Test Corp",
		Type: "company",
	}
	_, err = suite.orgService.CreateOrganization(orgDTO)
	assert.NoError(suite.T(), err)

	// Create another user without organization
	registerDTO := &auth.RegisterDTO{
		FirstName: "Standalone",
		LastName:  "User",
		Email:     "standalone.user@example.com",
		Password:  "Standalone-Password123", // Fixed: Added uppercase, lowercase, and number
	}
	authService := auth.NewAuthService(suite.authRepo, nil)
	_, err = authService.Register(registerDTO)
	assert.NoError(suite.T(), err)

	// 3. Get updated stats
	updatedStats, err := suite.adminService.GetSystemStats()
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), updatedStats)
	assert.Equal(suite.T(), 2, updatedStats.TotalUsers)         // 2 users created
	assert.Equal(suite.T(), 1, updatedStats.TotalOrganizations) // 1 organization created
	assert.Equal(suite.T(), 0, updatedStats.ActiveSessions)     // Always 0 for now
}
