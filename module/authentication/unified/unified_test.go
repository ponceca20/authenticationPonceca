package unified

import (
	"practicev2/config"
	"practicev2/module/authentication/auth"
	"practicev2/module/authentication/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type UnifiedTestSuite struct {
	suite.Suite
	db             *gorm.DB
	unifiedService UnifiedService
	unifiedRepo    UnifiedRepository
	testUser       *models.Identity
}

// setupTestDatabase initializes an in-memory SQLite database for testing purposes.
func (suite *UnifiedTestSuite) setupTestDatabase() *gorm.DB {
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
		&models.RefreshToken{},
	)
	if err != nil {
		suite.T().Fatalf("Failed to run migrations: %v", err)
	}

	return db
}

func (suite *UnifiedTestSuite) SetupSuite() {
	suite.db = suite.setupTestDatabase()
	suite.unifiedRepo = NewUnifiedRepository(suite.db)
	suite.unifiedService = NewUnifiedService(suite.unifiedRepo)

	// Create a test user
	authRepo := auth.NewAuthRepository(suite.db)
	authService := auth.NewAuthService(authRepo, nil)
	registerDTO := &auth.RegisterDTO{
		FirstName: "Unified",
		LastName:  "User",
		Email:     "unified.user@example.com",
		Password:  "unified-password",
	}
	identity, err := authService.Register(registerDTO)
	assert.NoError(suite.T(), err)
	suite.testUser = identity
}

func TestUnifiedTestSuite(t *testing.T) {
	suite.Run(t, new(UnifiedTestSuite))
}

func (suite *UnifiedTestSuite) TestGetDashboardData() {
	// 1. Get dashboard data for the test user
	dashboardData, err := suite.unifiedService.GetDashboardData(suite.testUser)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), dashboardData)

	// 2. Verify welcome message contains user's first name
	expectedWelcome := "Welcome back, Unified!"
	assert.Equal(suite.T(), expectedWelcome, dashboardData.WelcomeMessage)

	// 3. Verify recent activities are present
	assert.NotEmpty(suite.T(), dashboardData.RecentActivities)
	assert.Len(suite.T(), dashboardData.RecentActivities, 1)
	assert.Equal(suite.T(), "You logged in.", dashboardData.RecentActivities[0].Description)
	assert.Equal(suite.T(), "Just now", dashboardData.RecentActivities[0].Timestamp)

	// 4. Verify pending tasks are present
	assert.NotEmpty(suite.T(), dashboardData.PendingTasks)
	assert.Len(suite.T(), dashboardData.PendingTasks, 1)
	assert.Equal(suite.T(), "Complete your profile", dashboardData.PendingTasks[0].Title)
	assert.Equal(suite.T(), "/api/v1/profiles/me", dashboardData.PendingTasks[0].Link)
}

func (suite *UnifiedTestSuite) TestGetDashboardDataMultipleUsers() {
	// 1. Create another test user
	authRepo := auth.NewAuthRepository(suite.db)
	authService := auth.NewAuthService(authRepo, nil)
	registerDTO := &auth.RegisterDTO{
		FirstName: "Second",
		LastName:  "User",
		Email:     "second.unified@example.com",
		Password:  "second-password",
	}
	secondUser, err := authService.Register(registerDTO)
	assert.NoError(suite.T(), err)

	// 2. Get dashboard data for both users
	dashboard1, err := suite.unifiedService.GetDashboardData(suite.testUser)
	assert.NoError(suite.T(), err)

	dashboard2, err := suite.unifiedService.GetDashboardData(secondUser)
	assert.NoError(suite.T(), err)

	// 3. Verify each user gets personalized welcome message
	assert.Equal(suite.T(), "Welcome back, Unified!", dashboard1.WelcomeMessage)
	assert.Equal(suite.T(), "Welcome back, Second!", dashboard2.WelcomeMessage)

	// 4. Verify both have the same structure but different personalization
	assert.Len(suite.T(), dashboard1.RecentActivities, 1)
	assert.Len(suite.T(), dashboard2.RecentActivities, 1)
	assert.Len(suite.T(), dashboard1.PendingTasks, 1)
	assert.Len(suite.T(), dashboard2.PendingTasks, 1)
}

func (suite *UnifiedTestSuite) TestGetDashboardDataStructure() {
	// Test the structure and types of the dashboard data
	dashboardData, err := suite.unifiedService.GetDashboardData(suite.testUser)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), dashboardData)

	// 1. Verify all fields are present
	assert.NotEmpty(suite.T(), dashboardData.WelcomeMessage)
	assert.NotNil(suite.T(), dashboardData.RecentActivities)
	assert.NotNil(suite.T(), dashboardData.PendingTasks)

	// 2. Verify activity structure
	if len(dashboardData.RecentActivities) > 0 {
		activity := dashboardData.RecentActivities[0]
		assert.NotEmpty(suite.T(), activity.Description)
		assert.NotEmpty(suite.T(), activity.Timestamp)
	}

	// 3. Verify task structure
	if len(dashboardData.PendingTasks) > 0 {
		task := dashboardData.PendingTasks[0]
		assert.NotEmpty(suite.T(), task.Title)
		assert.NotEmpty(suite.T(), task.Link)
	}
}

func (suite *UnifiedTestSuite) TestGetDashboardDataWithDifferentUserNames() {
	// Test with users with different name patterns
	authRepo := auth.NewAuthRepository(suite.db)
	authService := auth.NewAuthService(authRepo, nil)

	testCases := []struct {
		firstName       string
		email           string
		expectedWelcome string
	}{
		{"John", "john@example.com", "Welcome back, John!"},
		{"María", "maria@example.com", "Welcome back, María!"},
		{"李", "li@example.com", "Welcome back, 李!"},
		{"", "empty@example.com", "Welcome back, !"},
	}

	for i, tc := range testCases {
		registerDTO := &auth.RegisterDTO{
			FirstName: tc.firstName,
			LastName:  "TestUser",
			Email:     tc.email,
			Password:  "test-password",
		}
		user, err := authService.Register(registerDTO)
		assert.NoError(suite.T(), err, "Failed to create user %d", i)

		dashboardData, err := suite.unifiedService.GetDashboardData(user)
		assert.NoError(suite.T(), err, "Failed to get dashboard for user %d", i)
		assert.Equal(suite.T(), tc.expectedWelcome, dashboardData.WelcomeMessage, "Wrong welcome message for user %d", i)
	}
}
