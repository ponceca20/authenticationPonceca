package invitation

import (
	"practicev2/module/authentication/auth"
	"practicev2/module/authentication/models"
	"practicev2/module/authentication/organization"
	"practicev2/module/authentication/test"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type InvitationTestSuite struct {
	suite.Suite
	db              *gorm.DB
	invitationService InvitationService
	authRepo        auth.AuthRepository
	testOrg         *models.Organization
	inviter         *models.Identity
	testRole        *models.Role
}

func (suite *InvitationTestSuite) SetupSuite() {
	suite.db = test.SetupTestDatabase(suite.T())
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
		Password:  "new-password-strong",
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
