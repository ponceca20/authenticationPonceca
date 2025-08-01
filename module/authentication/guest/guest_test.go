package guest

import (
	"practicev2/config"
	"practicev2/module/authentication/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type GuestTestSuite struct {
	suite.Suite
	db           *gorm.DB
	guestService GuestService
	guestRepo    GuestRepository
}

// setupTestDatabase initializes an in-memory SQLite database for testing purposes.
func (suite *GuestTestSuite) setupTestDatabase() *gorm.DB {
	// Initialize config with default values for tests
	config.Init()

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		suite.T().Fatalf("Failed to connect to in-memory database: %v", err)
	}

	// Auto-migrate the models we need for testing
	err = db.AutoMigrate(
		&models.GuestSession{},
	)
	if err != nil {
		suite.T().Fatalf("Failed to run migrations: %v", err)
	}

	return db
}

func (suite *GuestTestSuite) SetupSuite() {
	suite.db = suite.setupTestDatabase()
	suite.guestRepo = NewGuestRepository(suite.db)
	suite.guestService = NewGuestService(suite.guestRepo)
}

func TestGuestTestSuite(t *testing.T) {
	suite.Run(t, new(GuestTestSuite))
}

func (suite *GuestTestSuite) TestCreateGuestSession() {
	// 1. Create a new guest session
	session, err := suite.guestService.CreateGuestSession()
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), session)
	assert.NotEmpty(suite.T(), session.ID)
	assert.NotEmpty(suite.T(), session.SessionToken)
	assert.True(suite.T(), session.ExpiresAt.After(time.Now()))
	assert.WithinDuration(suite.T(), time.Now(), session.LastActivity, time.Second)

	// 2. Verify session can be retrieved
	retrievedSession, err := suite.guestService.GetSession(session.SessionToken)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), retrievedSession)
	assert.Equal(suite.T(), session.ID, retrievedSession.ID)
	assert.Equal(suite.T(), session.SessionToken, retrievedSession.SessionToken)
}

func (suite *GuestTestSuite) TestGetSession() {
	// 1. Create a session first
	session, err := suite.guestService.CreateGuestSession()
	assert.NoError(suite.T(), err)

	// 2. Retrieve the session by token
	retrievedSession, err := suite.guestService.GetSession(session.SessionToken)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), retrievedSession)
	assert.Equal(suite.T(), session.ID, retrievedSession.ID)

	// 3. Try to get a non-existent session
	_, err = suite.guestService.GetSession("invalid-token")
	assert.Error(suite.T(), err)
}

func (suite *GuestTestSuite) TestUpdateCart() {
	// 1. Create a session first
	session, err := suite.guestService.CreateGuestSession()
	assert.NoError(suite.T(), err)

	// Small delay to ensure different timestamps
	time.Sleep(10 * time.Millisecond)

	// 2. Update cart data
	cartData := `{"items":[{"id":"item1","quantity":2},{"id":"item2","quantity":1}]}`
	updatedSession, err := suite.guestService.UpdateCart(session.SessionToken, cartData)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), updatedSession)
	assert.NotNil(suite.T(), updatedSession.CartData)
	assert.Equal(suite.T(), cartData, *updatedSession.CartData)
	assert.True(suite.T(), updatedSession.LastActivity.After(session.LastActivity) || updatedSession.LastActivity.Equal(session.LastActivity))

	// 3. Verify the cart data persisted
	retrievedSession, err := suite.guestService.GetSession(session.SessionToken)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), cartData, retrievedSession.CartData)

	// 4. Try to update cart with invalid token
	_, err = suite.guestService.UpdateCart("invalid-token", "some data")
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "invalid or expired session token")
}

func (suite *GuestTestSuite) TestSessionLifecycle() {
	// 1. Create multiple sessions
	session1, err := suite.guestService.CreateGuestSession()
	assert.NoError(suite.T(), err)

	session2, err := suite.guestService.CreateGuestSession()
	assert.NoError(suite.T(), err)

	// 2. Verify they have different tokens
	assert.NotEqual(suite.T(), session1.SessionToken, session2.SessionToken)
	assert.NotEqual(suite.T(), session1.ID, session2.ID)

	// 3. Update cart for each session independently
	cartData1 := `{"items":[{"id":"item1","quantity":1}]}`
	cartData2 := `{"items":[{"id":"item2","quantity":3}]}`

	_, err = suite.guestService.UpdateCart(session1.SessionToken, cartData1)
	assert.NoError(suite.T(), err)

	_, err = suite.guestService.UpdateCart(session2.SessionToken, cartData2)
	assert.NoError(suite.T(), err)

	// 4. Verify each session has correct cart data
	retrieved1, err := suite.guestService.GetSession(session1.SessionToken)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), cartData1, retrieved1.CartData)

	retrieved2, err := suite.guestService.GetSession(session2.SessionToken)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), cartData2, retrieved2.CartData)
}
