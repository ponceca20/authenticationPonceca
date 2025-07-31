package profile

import (
	"fmt"
	"practicev2/config"
	"practicev2/module/authentication/auth"
	"practicev2/module/authentication/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type ProfileTestSuite struct {
	suite.Suite
	db             *gorm.DB
	profileService ProfileService
	profileRepo    ProfileRepository
	authRepo       auth.AuthRepository
	testUser       *models.Identity
}

func (suite *ProfileTestSuite) setupTestDatabase() *gorm.DB {
	config.Init()

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		suite.T().Fatalf("Failed to connect to in-memory database: %v", err)
	}

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

func (suite *ProfileTestSuite) SetupSuite() {
	suite.db = suite.setupTestDatabase()
	suite.profileRepo = NewProfileRepository(suite.db)
	suite.profileService = NewProfileService(suite.profileRepo, suite.db)
	suite.authRepo = auth.NewAuthRepository(suite.db)
}

// SetupTest runs before each individual test to ensure clean state
func (suite *ProfileTestSuite) SetupTest() {
	// Clean all test data to ensure complete isolation
	suite.db.Exec("DELETE FROM user_profile")
	suite.db.Exec("DELETE FROM refresh_token")
	suite.db.Exec("DELETE FROM identity")

	// Small sleep to ensure different timestamps between tests
	time.Sleep(5 * time.Millisecond)

	// Create a fresh test user for each test
	authService := auth.NewAuthService(suite.authRepo, nil)
	uniqueEmail := fmt.Sprintf("profile.user.%d@example.com", time.Now().UnixNano()%100000)

	registerDTO := &auth.RegisterDTO{
		FirstName: "Profile",
		LastName:  "User",
		Email:     uniqueEmail,
		Password:  "profile-password",
	}
	identity, err := authService.Register(registerDTO)
	if err != nil {
		suite.T().Fatalf("Failed to create test user: %v", err)
	}
	suite.testUser = identity
}

func TestProfileTestSuite(t *testing.T) {
	suite.Run(t, new(ProfileTestSuite))
}

func (suite *ProfileTestSuite) TestGetProfile() {
	profile, err := suite.profileService.GetProfile(suite.testUser.ID)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), profile)
	assert.Equal(suite.T(), suite.testUser.ID, profile.ID)
	assert.Equal(suite.T(), suite.testUser.Email, profile.Email)
	assert.Equal(suite.T(), suite.testUser.FirstName, profile.FirstName)
	assert.Equal(suite.T(), suite.testUser.LastName, profile.LastName)

	_, err = suite.profileService.GetProfile("non-existent-id")
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "user not found")
}

func (suite *ProfileTestSuite) TestUpdateProfile() {
	updateDTO := &ProfileDTO{
		FirstName: "Updated",
		LastName:  "Name",
		Avatar:    "https://example.com/avatar.jpg",
		Phone:     "+1234567890",
		Bio:       "This is my updated bio",
		Location:  "New York, USA",
		Website:   "https://example.com",
	}

	updatedProfile, err := suite.profileService.UpdateProfile(suite.testUser.ID, updateDTO)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), updatedProfile)
	assert.Equal(suite.T(), "Updated", updatedProfile.FirstName)
	assert.Equal(suite.T(), "Name", updatedProfile.LastName)
}

func (suite *ProfileTestSuite) TestUpdateNonExistentUser() {
	updateDTO := &ProfileDTO{
		FirstName: "Should",
		LastName:  "Fail",
	}

	_, err := suite.profileService.UpdateProfile("non-existent-id", updateDTO)
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "user not found")
}
