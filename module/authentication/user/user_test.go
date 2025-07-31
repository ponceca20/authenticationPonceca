package user

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

type UserTestSuite struct {
	suite.Suite
	db          *gorm.DB
	userRepo    UserRepository
	userService UserService
	orgService  organization.OrganizationService
	testOrg     *models.Organization
	testUser    *models.Identity
}

// setupTestDatabase initializes an in-memory SQLite database for testing purposes.
func (suite *UserTestSuite) setupTestDatabase() *gorm.DB {
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

func (suite *UserTestSuite) SetupSuite() {
	suite.db = suite.setupTestDatabase()
	suite.userRepo = NewUserRepository(suite.db)
	suite.userService = NewUserService(suite.userRepo)

	// We need an organization and a user to exist for these tests
	orgRepo := organization.NewOrganizationRepository(suite.db)
	authRepo := auth.NewAuthRepository(suite.db)
	suite.orgService = organization.NewOrganizationService(orgRepo, authRepo)

	// Create a test organization and user
	dto := &organization.OrganizationRegistrationDTO{
		Identity: auth.RegisterDTO{
			FirstName: "Test",
			LastName:  "User",
			Email:     "user.test@example.com",
			Password:  "password123",
		},
		Name: "User Test Corp",
		Type: "company",
	}
	org, err := suite.orgService.CreateOrganization(dto)
	assert.NoError(suite.T(), err)
	suite.testOrg = org

	user, err := authRepo.FindIdentityByEmail(dto.Identity.Email)
	assert.NoError(suite.T(), err)
	suite.testUser = user
}

func TestUserTestSuite(t *testing.T) {
	suite.Run(t, new(UserTestSuite))
}

func (suite *UserTestSuite) TestListAndUpdateUsers() {
	// 1. List users, should be at least one
	users, err := suite.userService.ListOrgUsers(suite.testOrg.ID)
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), users)
	assert.Equal(suite.T(), suite.testUser.ID, users[0].ID)

	// 2. Update the user's department
	updateDTO := &UpdateUserDTO{Department: "Engineering"}
	updatedUser, err := suite.userService.UpdateUser(suite.testOrg.ID, suite.testUser.ID, updateDTO)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Engineering", updatedUser.Department)

	// 3. Deactivate the user
	statusDTO := &UpdateUserStatusDTO{IsActive: false}
	err = suite.userService.UpdateUserStatus(suite.testOrg.ID, suite.testUser.ID, statusDTO)
	assert.NoError(suite.T(), err)

	// 4. Verify the user is inactive
	inactiveUser, err := suite.userService.GetOrgUser(suite.testOrg.ID, suite.testUser.ID)
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), inactiveUser.IsActive)
}

func (suite *UserTestSuite) TestBulkCreateUsers() {
	// For this test, we need a role to assign
	var role models.Role
	err := suite.db.Where("organization_id = ?", suite.testOrg.ID).First(&role).Error
	assert.NoError(suite.T(), err)

	bulkDTO := &BulkCreateUserDTO{
		Users: []BulkCreateUserItemDTO{
			{FirstName: "Bulk1", LastName: "User", Email: "bulk1@example.com", Password: "password", RoleID: role.ID},
			{FirstName: "Bulk2", LastName: "User", Email: "bulk2@example.com", Password: "password", RoleID: role.ID},
		},
	}

	err = suite.userService.BulkCreateUsers(suite.testOrg.ID, bulkDTO)
	assert.NoError(suite.T(), err)

	// Verify users were created
	users, err := suite.userService.ListOrgUsers(suite.testOrg.ID)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), users, 3) // Original user + 2 bulk created
}
