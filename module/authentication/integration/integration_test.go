package integration

import (
	"fmt"
	"practicev2/config"
	"practicev2/module/authentication/auth"
	"practicev2/module/authentication/models"
	"practicev2/module/authentication/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type IntegrationTestSuite struct {
	suite.Suite
	db                 *gorm.DB
	integrationService IntegrationService
	jwtService         *utils.JWTService
	authService        auth.AuthService
	testUser           *models.Identity
}

// setupTestDatabase initializes an in-memory SQLite database for testing purposes.
func (suite *IntegrationTestSuite) setupTestDatabase() *gorm.DB {
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

func (suite *IntegrationTestSuite) SetupSuite() {
	suite.db = suite.setupTestDatabase()
	suite.jwtService = utils.NewJWTService()
	suite.integrationService = NewIntegrationService(suite.jwtService)

	// Setup auth service to create test users and tokens
	authRepo := auth.NewAuthRepository(suite.db)
	suite.authService = auth.NewAuthService(authRepo, suite.jwtService)
}

// SetupTest runs before each individual test to ensure clean state
func (suite *IntegrationTestSuite) SetupTest() {
	// Clean all test data to ensure complete isolation
	suite.db.Exec("DELETE FROM refresh_token")
	suite.db.Exec("DELETE FROM password_reset_token")
	suite.db.Exec("DELETE FROM customer_profile")
	suite.db.Exec("DELETE FROM organizational_membership")
	suite.db.Exec("DELETE FROM identity")

	// Small sleep to ensure different timestamps between tests
	time.Sleep(5 * time.Millisecond)

	// Create a fresh test user for each test
	uniqueEmail := fmt.Sprintf("integration.user.%d@example.com", time.Now().UnixNano()%100000)
	registerDTO := &auth.RegisterDTO{
		FirstName: "Integration",
		LastName:  "User",
		Email:     uniqueEmail,
		Password:  "integration-password",
	}
	identity, err := suite.authService.Register(registerDTO)
	if err != nil {
		suite.T().Fatalf("Failed to create test user: %v", err)
	}
	suite.testUser = identity
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func (suite *IntegrationTestSuite) TestValidateToken() {
	// 1. Generate a valid token
	accessToken, _, err := suite.jwtService.GenerateTokenPair(suite.testUser, nil, nil)
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), accessToken)

	// 2. Validate the token through integration service
	response, err := suite.integrationService.ValidateToken(accessToken)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), response)
	assert.True(suite.T(), response.IsValid)
	assert.Equal(suite.T(), suite.testUser.ID, response.UserID)
	assert.Equal(suite.T(), suite.testUser.Email, response.Email)
}

func (suite *IntegrationTestSuite) TestValidateInvalidToken() {
	// 1. Test with completely invalid token
	response, err := suite.integrationService.ValidateToken("invalid-token")
	assert.NoError(suite.T(), err) // Service should not return error, but mark as invalid
	assert.NotNil(suite.T(), response)
	assert.False(suite.T(), response.IsValid)
	assert.Empty(suite.T(), response.UserID)
	assert.Empty(suite.T(), response.Email)

	// 2. Test with empty token
	response, err = suite.integrationService.ValidateToken("")
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), response)
	assert.False(suite.T(), response.IsValid)

	// 3. Test with malformed JWT
	response, err = suite.integrationService.ValidateToken("Bearer invalid.jwt.token")
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), response)
	assert.False(suite.T(), response.IsValid)
}

func (suite *IntegrationTestSuite) TestValidateExpiredToken() {
	// Test with completely malformed tokens that should definitely be invalid
	testCases := []struct {
		name  string
		token string
	}{
		{"corrupted JWT", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.INVALID_PAYLOAD.INVALID_SIGNATURE"},
		{"malformed structure", "invalid.jwt.structure.here"},
		{"empty segments", ".."},
		{"wrong signature", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.wrong_signature_here"},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			response, err := suite.integrationService.ValidateToken(tc.token)
			assert.NoError(t, err, "Service should not return error for invalid tokens")
			assert.NotNil(t, response, "Response should not be nil")
			assert.False(t, response.IsValid, "Token should be marked as invalid for case: %s", tc.name)
			assert.Empty(t, response.UserID, "UserID should be empty for invalid token")
			assert.Empty(t, response.Email, "Email should be empty for invalid token")
		})
	}
}

func (suite *IntegrationTestSuite) TestValidateTokenMultipleUsers() {
	// 1. Create another test user with unique email
	uniqueEmail := fmt.Sprintf("second.user.%d@example.com", time.Now().UnixNano()%100000)
	registerDTO := &auth.RegisterDTO{
		FirstName: "Second",
		LastName:  "User",
		Email:     uniqueEmail,
		Password:  "second-password",
	}
	secondUser, err := suite.authService.Register(registerDTO)
	assert.NoError(suite.T(), err)

	// 2. Generate tokens for both users
	token1, _, err := suite.jwtService.GenerateTokenPair(suite.testUser, nil, nil)
	assert.NoError(suite.T(), err)

	token2, _, err := suite.jwtService.GenerateTokenPair(secondUser, nil, nil)
	assert.NoError(suite.T(), err)

	// 3. Validate both tokens
	response1, err := suite.integrationService.ValidateToken(token1)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), response1.IsValid)
	assert.Equal(suite.T(), suite.testUser.ID, response1.UserID)
	assert.Equal(suite.T(), suite.testUser.Email, response1.Email)

	response2, err := suite.integrationService.ValidateToken(token2)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), response2.IsValid)
	assert.Equal(suite.T(), secondUser.ID, response2.UserID)
	assert.Equal(suite.T(), secondUser.Email, response2.Email)

	// 4. Verify tokens are not interchangeable
	assert.NotEqual(suite.T(), response1.UserID, response2.UserID)
	assert.NotEqual(suite.T(), response1.Email, response2.Email)
}
