package auth

import (
	"fmt"
	"practicev2/config"
	"practicev2/module/authentication/models"
	"practicev2/module/authentication/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type AuthTestSuite struct {
	suite.Suite
	db         *gorm.DB
	repo       AuthRepository
	service    AuthService
	jwtService *utils.JWTService
}

// setupTestDatabase initializes an in-memory SQLite database for testing purposes.
func (suite *AuthTestSuite) setupTestDatabase() *gorm.DB {
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
		&models.CustomerProfile{},
		&models.RefreshToken{},
		&models.PasswordResetToken{},
	)
	if err != nil {
		suite.T().Fatalf("Failed to run migrations: %v", err)
	}

	return db
}

// SetupSuite runs once before the entire test suite
func (suite *AuthTestSuite) SetupSuite() {
	suite.db = suite.setupTestDatabase()
	suite.repo = NewAuthRepository(suite.db)
	suite.jwtService = utils.NewJWTService()
	suite.service = NewAuthService(suite.repo, suite.jwtService)
}

// SetupTest runs before each individual test to ensure clean state
func (suite *AuthTestSuite) SetupTest() {
	// Clean all test data to ensure complete isolation (only existing tables)
	suite.db.Exec("DELETE FROM refresh_token")
	suite.db.Exec("DELETE FROM password_reset_token")
	suite.db.Exec("DELETE FROM customer_profile")
	suite.db.Exec("DELETE FROM organizational_membership")
	suite.db.Exec("DELETE FROM identity")

	// Small sleep to ensure different timestamps between tests
	time.Sleep(5 * time.Millisecond)
}

// generateUniqueEmail creates a unique email for each test to avoid conflicts
func (suite *AuthTestSuite) generateUniqueEmail(prefix string) string {
	// Use crypto/rand for better randomness in tests
	return fmt.Sprintf("%s.%d@example.com", prefix, time.Now().UnixNano()%100000)
}

func TestAuthTestSuite(t *testing.T) {
	suite.Run(t, new(AuthTestSuite))
}

func (suite *AuthTestSuite) TestRegisterAndLogin() {
	// 1. Register a new user
	registerDTO := &RegisterDTO{
		FirstName: "John",
		LastName:  "Doe",
		Email:     suite.generateUniqueEmail("john.doe"),
		Password:  "strong-password-123",
	}

	identity, err := suite.service.Register(registerDTO)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), identity)
	assert.Equal(suite.T(), registerDTO.Email, identity.Email)

	// 2. Attempt to register the same email again (should fail)
	_, err = suite.service.Register(registerDTO)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "email already in use", err.Error())

	// 3. Login with correct credentials
	loginDTO := &LoginDTO{
		Email:    registerDTO.Email,
		Password: "strong-password-123",
	}

	tokenResponse, err := suite.service.Login(loginDTO)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), tokenResponse)
	assert.NotEmpty(suite.T(), tokenResponse.AccessToken)
	assert.NotEmpty(suite.T(), tokenResponse.RefreshToken)
	assert.Equal(suite.T(), identity.ID, tokenResponse.Identity.ID)

	// 4. Validate the received access token
	claims, err := suite.jwtService.ValidateToken(tokenResponse.AccessToken)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), identity.ID, claims.IdentityID)
	assert.Equal(suite.T(), identity.Email, claims.Email)

	// 5. Login with incorrect credentials (should fail)
	invalidLoginDTO := &LoginDTO{
		Email:    "john.doe@example.com",
		Password: "wrong-password",
	}
	_, err = suite.service.Login(invalidLoginDTO)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "invalid credentials", err.Error())
}

func (suite *AuthTestSuite) TestRefreshToken() {
	// 1. Register and login a user to get a refresh token
	registerDTO := &RegisterDTO{
		FirstName: "Jane",
		LastName:  "Doe",
		Email:     suite.generateUniqueEmail("jane.doe"),
		Password:  "a-different-password",
	}
	_, err := suite.service.Register(registerDTO)
	assert.NoError(suite.T(), err)

	loginDTO := &LoginDTO{Email: registerDTO.Email, Password: registerDTO.Password}
	loginResponse, err := suite.service.Login(loginDTO)
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), loginResponse.RefreshToken)

	// Significant delay to ensure different timestamp for refresh token generation
	time.Sleep(50 * time.Millisecond)

	// 2. Use the refresh token to get a new access token
	refreshDTO := &RefreshTokenDTO{RefreshToken: loginResponse.RefreshToken}
	refreshResponse, err := suite.service.RefreshToken(refreshDTO)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), refreshResponse)
	assert.NotEmpty(suite.T(), refreshResponse.AccessToken)
	assert.NotEqual(suite.T(), loginResponse.AccessToken, refreshResponse.AccessToken)

	// 3. Try to use the old refresh token again (should fail as it's revoked)
	_, err = suite.service.RefreshToken(refreshDTO)
	assert.Error(suite.T(), err)
}

func (suite *AuthTestSuite) TestLogout() {
	// 1. Register and login
	registerDTO := &RegisterDTO{
		FirstName: "Logout",
		LastName:  "User",
		Email:     suite.generateUniqueEmail("logout.user"),
		Password:  "password-to-logout",
	}
	_, err := suite.service.Register(registerDTO)
	assert.NoError(suite.T(), err)
	loginDTO := &LoginDTO{Email: registerDTO.Email, Password: registerDTO.Password}
	loginResponse, err := suite.service.Login(loginDTO)
	assert.NoError(suite.T(), err)

	// Significant delay to ensure different timestamp
	time.Sleep(50 * time.Millisecond)

	// 2. Logout using the refresh token
	logoutDTO := &RefreshTokenDTO{RefreshToken: loginResponse.RefreshToken}
	err = suite.service.Logout(logoutDTO)
	assert.NoError(suite.T(), err)

	// 3. Verify the token is revoked by trying to use it again
	_, err = suite.service.RefreshToken(logoutDTO)
	assert.Error(suite.T(), err, "The refresh token should be revoked after logout")
}
