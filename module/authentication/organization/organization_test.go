package organization

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

type OrganizationTestSuite struct {
	suite.Suite
	db       *gorm.DB
	orgRepo  OrganizationRepository
	authRepo auth.AuthRepository
	service  OrganizationService
}

// setupTestDatabase initializes an in-memory SQLite database for testing purposes.
func (suite *OrganizationTestSuite) setupTestDatabase() *gorm.DB {
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

func (suite *OrganizationTestSuite) SetupSuite() {
	suite.db = suite.setupTestDatabase()
	suite.orgRepo = NewOrganizationRepository(suite.db)
	suite.authRepo = auth.NewAuthRepository(suite.db)
	suite.service = NewOrganizationService(suite.orgRepo, suite.authRepo)
}

func TestOrganizationTestSuite(t *testing.T) {
	suite.Run(t, new(OrganizationTestSuite))
}

func (suite *OrganizationTestSuite) TestCreateOrganization() {
	dto := &OrganizationRegistrationDTO{
		Identity: auth.RegisterDTO{
			FirstName: "Admin",
			LastName:  "User",
			Email:     "org.admin@example.com",
			Password:  "secure-org-password",
		},
		Name: "Test Organization Inc.",
		Type: "company",
	}

	org, err := suite.service.CreateOrganization(dto)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), org)
	assert.Equal(suite.T(), "Test Organization Inc.", org.Name)

	// Verify that the founder was created
	founder, err := suite.authRepo.FindIdentityByEmail(dto.Identity.Email)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), founder)
	assert.Equal(suite.T(), "Admin", founder.FirstName)

	// Verify that the membership was created
	var membership models.OrganizationalMembership
	err = suite.db.Where("identity_id = ? AND organization_id = ?", founder.ID, org.ID).First(&membership).Error
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), membership)
}
