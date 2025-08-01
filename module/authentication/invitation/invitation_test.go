package invitation

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

type InvitationTestSuite struct {
	suite.Suite
	db                *gorm.DB
	invitationService InvitationService
	authRepo          auth.AuthRepository
	testOrg           *models.Organization
	inviter           *models.Identity
	testRole          *models.Role
}

// setupTestDatabase initializes an in-memory SQLite database for testing purposes.
func (suite *InvitationTestSuite) setupTestDatabase() *gorm.DB {
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
		&models.Invitation{},
		&models.Role{},
		&models.RefreshToken{},
	)
	if err != nil {
		suite.T().Fatalf("Failed to run migrations: %v", err)
	}

	return db
}

func (suite *InvitationTestSuite) SetupSuite() {
	suite.db = suite.setupTestDatabase()
	invitationRepo := NewInvitationRepository(suite.db)
	suite.invitationService = NewInvitationService(invitationRepo)
	suite.authRepo = auth.NewAuthRepository(suite.db)

	// We need an organization, an inviter user, and a role to exist
	orgRepo := organization.NewOrganizationRepository(suite.db)
	orgService := organization.NewOrganizationService(orgRepo, suite.authRepo)

	dto := &organization.OrganizationRegistrationDTO{
		Identity: auth.RegisterDTO{
			FirstName: "Inviter",
			LastName:  "User",
			Email:     "inviter.user@example.com",
			Password:  "password123",
		},
		Name: "Invitation Test Corp",
		Type: "company",
	}
	org, err := orgService.CreateOrganization(dto)
	assert.NoError(suite.T(), err)
	suite.testOrg = org

	user, err := suite.authRepo.FindIdentityByEmail(dto.Identity.Email)
	assert.NoError(suite.T(), err)
	suite.inviter = user

	var role models.Role
	err = suite.db.Where("organization_id = ?", org.ID).First(&role).Error
	assert.NoError(suite.T(), err)
	suite.testRole = &role
}

func TestInvitationTestSuite(t *testing.T) {
	suite.Run(t, new(InvitationTestSuite))
}

func (suite *InvitationTestSuite) TestInvitationLifecycle() {
	// 1. Create an invitation
	inviteDTO := &InvitationDTO{
		Email:  "new.invitee@example.com",
		RoleID: suite.testRole.ID,
	}
	createdInvite, err := suite.invitationService.CreateInvitation(suite.testOrg.ID, suite.inviter.ID, inviteDTO)
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), createdInvite.Token)

	// 2. Accept the invitation
	acceptDTO := &AcceptInvitationDTO{
		Token:     createdInvite.Token,
		FirstName: "New",
		LastName:  "Invitee",
		Password:  "New-Password-Strong123", // Fixed: Added uppercase, lowercase, and number
	}
	err = suite.invitationService.AcceptInvitation(acceptDTO)
	assert.NoError(suite.T(), err)

	// 3. Verify the new user and membership were created
	newUser, err := suite.authRepo.FindIdentityByEmail("new.invitee@example.com")
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "New", newUser.FirstName)

	var membership models.OrganizationalMembership
	err = suite.db.Where("identity_id = ?", newUser.ID).First(&membership).Error
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), suite.testOrg.ID, membership.OrganizationID)

	// 4. Test canceling another invitation
	cancelInviteDTO := &InvitationDTO{
		Email:  "another.invitee@example.com",
		RoleID: suite.testRole.ID,
	}
	toCancel, err := suite.invitationService.CreateInvitation(suite.testOrg.ID, suite.inviter.ID, cancelInviteDTO)
	assert.NoError(suite.T(), err)

	err = suite.invitationService.CancelInvitation(suite.testOrg.ID, toCancel.ID)
	assert.NoError(suite.T(), err)
}
