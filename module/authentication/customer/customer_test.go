package customer

import (
	"practicev2/module/authentication/auth"
	"practicev2/module/authentication/models"
	"practicev2/module/authentication/test"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type CustomerTestSuite struct {
	suite.Suite
	db             *gorm.DB
	customerRepo   CustomerRepository
	authRepo       auth.AuthRepository
	customerService CustomerService
	testIdentity   *models.Identity
}

func (suite *CustomerTestSuite) SetupSuite() {
	suite.db = test.SetupTestDatabase(suite.T())
	suite.customerRepo = NewCustomerRepository(suite.db)
	suite.authRepo = auth.NewAuthRepository(suite.db)
	suite.customerService = NewCustomerService(suite.customerRepo, suite.authRepo)

	// Create a test customer
	dto := &CustomerRegistrationDTO{
		FirstName: "Customer",
		LastName:  "Test",
		Email:     "customer.test@example.com",
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
		Theme:            "dark",
		Language:         "en",
		EmailNotifications: false,
	}
	updatedPrefs, err := suite.customerService.UpdatePreferences(suite.testIdentity.ID, updateDTO)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "dark", updatedPrefs.Theme)
	assert.False(suite.T(), updatedPrefs.EmailNotifications)
}
