package customer

import (
	"fmt"
	"math/rand"
	"practicev2/config"
	"practicev2/module/authentication/auth"
	"practicev2/module/authentication/guest"
	"practicev2/module/authentication/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type CustomerTestSuite struct {
	suite.Suite
	db              *gorm.DB
	customerRepo    CustomerRepository
	authRepo        auth.AuthRepository
	guestRepo       guest.GuestRepository
	guestService    guest.GuestService
	customerService CustomerService
	testIdentity    *models.Identity
}

// setupTestDatabase initializes an in-memory SQLite database for testing purposes.
func (suite *CustomerTestSuite) setupTestDatabase() *gorm.DB {
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
		&models.CustomerProfile{},
		&models.CustomerPreferences{},
		&models.ShippingAddress{},
		&models.RefreshToken{},
		&models.GuestSession{},
	)
	if err != nil {
		suite.T().Fatalf("Failed to run migrations: %v", err)
	}

	return db
}

func (suite *CustomerTestSuite) SetupSuite() {
	suite.db = suite.setupTestDatabase()
	suite.customerRepo = NewCustomerRepository(suite.db)
	suite.authRepo = auth.NewAuthRepository(suite.db)
	suite.guestRepo = guest.NewGuestRepository(suite.db)
	suite.guestService = guest.NewGuestService(suite.guestRepo)
	suite.customerService = NewCustomerService(suite.customerRepo, suite.authRepo, suite.guestService)

	// Create a test customer with unique email
	rand.Seed(time.Now().UnixNano())
	uniqueEmail := fmt.Sprintf("customer.test.%d@example.com", rand.Intn(100000))

	dto := &CustomerRegistrationDTO{
		FirstName: "Customer",
		LastName:  "Test",
		Email:     uniqueEmail,
		Password:  "password123",
	}
	profile, err := suite.customerService.RegisterCustomer(dto)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), profile)

	identity, err := suite.authRepo.FindIdentityByEmail(dto.Email)
	assert.NoError(suite.T(), err)
	suite.testIdentity = identity
}

func TestCustomerTestSuite(t *testing.T) {
	suite.Run(t, new(CustomerTestSuite))
}

func (suite *CustomerTestSuite) TestAddressManagement() {
	// 1. Add an address
	addrDTO := &AddressDTO{
		AddressLine1: "123 Main St",
		City:         "Anytown",
		State:        "CA",
		PostalCode:   "12345",
		Country:      "USA",
	}
	addedAddr, err := suite.customerService.AddAddress(suite.testIdentity.ID, addrDTO)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "123 Main St", addedAddr.AddressLine1)

	// 2. List addresses
	addresses, err := suite.customerService.ListAddresses(suite.testIdentity.ID)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), addresses, 1)

	// 3. Update address
	updateDTO := &AddressDTO{AddressLine1: "456 Oak Ave"}
	updatedAddr, err := suite.customerService.UpdateAddress(suite.testIdentity.ID, addedAddr.ID, updateDTO)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "456 Oak Ave", updatedAddr.AddressLine1)

	// 4. Delete address
	err = suite.customerService.DeleteAddress(suite.testIdentity.ID, addedAddr.ID)
	assert.NoError(suite.T(), err)

	// 5. Verify deletion
	_, err = suite.customerRepo.FindAddressByID(addedAddr.ID)
	assert.Error(suite.T(), err) // Should be gorm.ErrRecordNotFound
}

func (suite *CustomerTestSuite) TestPreferencesManagement() {
	// 1. Get initial preferences (should be created on the fly)
	prefs, err := suite.customerService.GetPreferences(suite.testIdentity.ID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "light", prefs.Theme) // Default value

	// 2. Update preferences
	updateDTO := &PreferencesDTO{
		Theme:              "dark",
		Language:           "en",
		EmailNotifications: false,
	}
	updatedPrefs, err := suite.customerService.UpdatePreferences(suite.testIdentity.ID, updateDTO)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "dark", updatedPrefs.Theme)
	assert.False(suite.T(), updatedPrefs.EmailNotifications)
}
