package organization

import (
	"practicev2/module/authentication/auth"
	"practicev2/module/authentication/models"
	"practicev2/module/authentication/test"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type OrganizationTestSuite struct {
	suite.Suite
	db      *gorm.DB
	orgRepo OrganizationRepository
	authRepo auth.AuthRepository
	service OrganizationService
}

func (suite *OrganizationTestSuite) SetupSuite() {
	suite.db = test.SetupTestDatabase(suite.T())
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
