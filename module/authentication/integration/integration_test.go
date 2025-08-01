package integration

import (
	"fmt"
	"practicev2/config"
	"practicev2/module/authentication/auth"
	"practicev2/module/authentication/customer"
	"practicev2/module/authentication/guest"
	"practicev2/module/authentication/invitation"
	"practicev2/module/authentication/models"
	"practicev2/module/authentication/organization"
	"practicev2/module/authentication/role"
	"practicev2/module/authentication/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type AuthLifecycleTestSuite struct {
	suite.Suite
	db                *gorm.DB
	jwtService        *utils.JWTService
	authService       auth.AuthService
	orgService        organization.OrganizationService
	roleService       role.RoleService
	invitationService invitation.InvitationService
	guestService      guest.GuestService
	customerService   customer.CustomerService
	authRepo          auth.AuthRepository
	customerRepo      customer.CustomerRepository

	// Test-specific data
	adminUser *models.Identity
}

// setupTestDatabase initializes an in-memory SQLite database for testing.
func (suite *AuthLifecycleTestSuite) setupTestDatabase() {
	config.Init() // Initialize config with default values for tests
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		suite.T().Fatalf("Failed to connect to in-memory database: %v", err)
	}

	err = db.AutoMigrate(
		&models.Identity{}, &models.User{}, &models.UserProfile{},
		&models.RefreshToken{}, &models.PasswordResetToken{},
		&models.GuestSession{}, &models.CustomerProfile{}, &models.ShippingAddress{},
		&models.Organization{}, &models.Department{}, &models.OrganizationalMembership{},
		&models.Role{}, &models.Permission{},
		&models.Invitation{}, &models.AuditLog{},
	)
	if err != nil {
		suite.T().Fatalf("Failed to run migrations: %v", err)
	}
	suite.db = db
}

func (suite *AuthLifecycleTestSuite) SetupSuite() {
	suite.setupTestDatabase()
	suite.jwtService = utils.NewJWTService()

	// Repositories
	suite.authRepo = auth.NewAuthRepository(suite.db)
	orgRepo := organization.NewOrganizationRepository(suite.db)
	roleRepo := role.NewRoleRepository(suite.db)
	invitationRepo := invitation.NewInvitationRepository(suite.db)
	guestRepo := guest.NewGuestRepository(suite.db)
	suite.customerRepo = customer.NewCustomerRepository(suite.db)

	// Services
	suite.guestService = guest.NewGuestService(guestRepo)
	suite.customerService = customer.NewCustomerService(suite.customerRepo, suite.authRepo, suite.guestService)
	suite.authService = auth.NewAuthService(suite.authRepo, suite.jwtService)
	suite.orgService = organization.NewOrganizationService(orgRepo, suite.authRepo)
	suite.roleService = role.NewRoleService(roleRepo)
	suite.invitationService = invitation.NewInvitationService(invitationRepo)
}

func (suite *AuthLifecycleTestSuite) SetupTest() {
	tables := []string{
		"role_permission", "permissions", "roles", "organizational_memberships",
		"invitations", "departments", "organizations", "audit_logs",
		"shipping_addresses", "customer_profiles", "guest_sessions",
		"password_reset_tokens", "refresh_tokens", "user_profiles", "identity",
	}
	for _, table := range tables {
		suite.db.Exec(fmt.Sprintf("DELETE FROM %s;", table))
		suite.db.Exec(fmt.Sprintf("DELETE FROM sqlite_sequence WHERE name = '%s';", table))
	}
	uniqueEmail := fmt.Sprintf("admin.%d@example.com", time.Now().UnixNano())
	registerDTO := &auth.RegisterDTO{
		FirstName: "Admin", LastName: "User", Email: uniqueEmail, Password: "a-Strong-Password123!",
	}
	identity, err := suite.authService.Register(registerDTO)
	if err != nil {
		suite.T().Fatalf("Failed to create admin user for test: %v", err)
	}
	suite.adminUser = identity
}

func TestAuthLifecycleTestSuite(t *testing.T) {
	suite.Run(t, new(AuthLifecycleTestSuite))
}

func (suite *AuthLifecycleTestSuite) TestPhase1_IdentityLifecycle() {
	require := require.New(suite.T())
	assert := assert.New(suite.T())
	suite.T().Run("should create user successfully on registration", func(t *testing.T) {
		require.NotZero(suite.adminUser.ID)
		require.False(suite.adminUser.EmailVerified)
		require.NotEmpty(suite.adminUser.PasswordHash)
	})
	suite.T().Run("should reject weak password during registration", func(t *testing.T) {
		registerDTO := &auth.RegisterDTO{
			FirstName: "Weak", LastName: "Password", Email: "weak.password@example.com", Password: "12345",
		}
		_, err := suite.authService.Register(registerDTO)
		assert.Error(err, "Expected an error for weak password")
	})
	suite.T().Run("should allow login and token refresh with correct credentials", func(t *testing.T) {
		loginDTO := &auth.LoginDTO{Email: suite.adminUser.Email, Password: "a-Strong-Password123!"}
		loginResponse, err := suite.authService.Login(loginDTO)
		require.NoError(err)
		require.NotEmpty(loginResponse.AccessToken)
		require.NotEmpty(loginResponse.RefreshToken)
		refreshDTO := &auth.RefreshTokenDTO{RefreshToken: loginResponse.RefreshToken}
		refreshResponse, err := suite.authService.RefreshToken(refreshDTO)
		require.NoError(err)
		require.NotEmpty(refreshResponse.AccessToken)
		assert.NotEqual(loginResponse.AccessToken, refreshResponse.AccessToken, "New access token should be different")
	})
	suite.T().Run("should deny login with incorrect credentials", func(t *testing.T) {
		loginDTO := &auth.LoginDTO{Email: suite.adminUser.Email, Password: "wrong-password"}
		_, err := suite.authService.Login(loginDTO)
		assert.Error(err, "Expected an error for incorrect credentials")
	})
	suite.T().Run("should lock account after multiple failed login attempts", func(t *testing.T) {
		loginDTO := &auth.LoginDTO{Email: suite.adminUser.Email, Password: "another-wrong-password"}
		for i := 0; i < 6; i++ {
			_, _ = suite.authService.Login(loginDTO)
		}
		correctLoginDTO := &auth.LoginDTO{Email: suite.adminUser.Email, Password: "a-Strong-Password123!"}
		_, err := suite.authService.Login(correctLoginDTO)
		require.Error(err, "Expected login to fail for locked account")
		assert.Contains(err.Error(), "account is locked", "Error message should indicate account lock")
		var identity models.Identity
		result := suite.db.Where("email = ?", suite.adminUser.Email).First(&identity)
		require.NoError(result.Error)
		require.NotNil(identity.LockedUntil, "LockedUntil should be set in the database")
		assert.True(identity.LockedUntil.After(time.Now()), "Lockout time should be in the future")
	})
}

func (suite *AuthLifecycleTestSuite) TestPhase2_GuestLifecycle() {
	require := require.New(suite.T())
	assert := assert.New(suite.T())
	suite.T().Run("should create and update a guest session", func(t *testing.T) {
		session, err := suite.guestService.CreateGuestSession()
		require.NoError(err)
		require.NotNil(session)
		require.NotEmpty(session.SessionToken)
		retrievedSession, err := suite.guestService.GetSession(session.SessionToken)
		require.NoError(err)
		require.Equal(session.ID, retrievedSession.ID)
		assert.Empty(retrievedSession.CartData)
		cartJSON := `{"items":[{"id":"prod_123","quantity":2}],"total":99.98}`
		updatedSession, err := suite.guestService.UpdateCart(session.SessionToken, cartJSON)
		require.NoError(err)
		assert.Equal(cartJSON, updatedSession.CartData)
		retrievedAfterUpdate, err := suite.guestService.GetSession(session.SessionToken)
		require.NoError(err)
		assert.Equal(cartJSON, retrievedAfterUpdate.CartData)
	})
	suite.T().Run("should convert guest to registered user", func(t *testing.T) {
		t.Skip("Skipping test: Guest-to-user conversion feature is not yet implemented in any service.")
	})
}

func (suite *AuthLifecycleTestSuite) TestPhase3_OrganizationLifecycle() {
	require := require.New(suite.T())
	assert := assert.New(suite.T())
	suite.T().Run("should create an organization and assign owner", func(t *testing.T) {
		orgDTO := &organization.OrganizationRegistrationDTO{
			Identity: auth.RegisterDTO{
				FirstName: "CEO", LastName: "TestCorp", Email: "ceo.testcorp@example.com", Password: "a-Very-Strong-Password123!",
			},
			Name: "TestCorp Inc.", Type: "company",
		}
		newOrg, err := suite.orgService.CreateOrganization(orgDTO)
		require.NoError(err, "Failed to create organization")
		require.NotNil(newOrg)
		assert.Equal("TestCorp Inc.", newOrg.Name)
		assert.Equal("testcorp-inc", newOrg.Slug)
		founder, err := suite.authRepo.FindIdentityByEmail(orgDTO.Identity.Email)
		require.NoError(err)
		require.NotNil(founder)
		assert.Equal("CEO", founder.FirstName)
		var membership models.OrganizationalMembership
		err = suite.db.Where("organization_id = ? AND identity_id = ?", newOrg.ID, founder.ID).First(&membership).Error
		require.NoError(err, "Failed to find organizational membership for the founder")
		assert.True(membership.IsActive)
		var role models.Role
		err = suite.db.Where("id = ?", membership.RoleID).First(&role).Error
		require.NoError(err, "Failed to find the role assigned to the founder")
		assert.Equal("owner", role.Name)
		assert.True(role.IsSystemRole, "The owner role should be a system role")
		assert.Equal(newOrg.ID, role.OrganizationID, "Role should be linked to the new organization")
	})
}

func (suite *AuthLifecycleTestSuite) TestPhase4_RBACAndInvitationLifecycle() {
	require := require.New(suite.T())
	assert := assert.New(suite.T())
	orgDTO := &organization.OrganizationRegistrationDTO{
		Identity: auth.RegisterDTO{
			FirstName: "CEO", LastName: "RBAC-Test", Email: "ceo.rbac@example.com", Password: "a-Very-Strong-Password123!",
		},
		Name: "RBAC Test Corp", Type: "company",
	}
	org, err := suite.orgService.CreateOrganization(orgDTO)
	require.NoError(err)
	founder, err := suite.authRepo.FindIdentityByEmail(orgDTO.Identity.Email)
	require.NoError(err)
	var managerRole *models.Role
	var invite *models.Invitation
	suite.T().Run("should allow owner to create a custom role", func(t *testing.T) {
		roleDTO := &role.RoleDTO{
			Name:        "manager",
			DisplayName: "Manager",
			Permissions: []role.PermissionDTO{
				{Resource: "users", Actions: []string{"read", "invite"}, Scope: "organization"},
				{Resource: "orders", Actions: []string{"read", "update"}, Scope: "department"},
			},
		}
		newRole, err := suite.roleService.CreateRole(org.ID, roleDTO)
		require.NoError(err)
		require.NotNil(newRole)
		assert.Equal("manager", newRole.Name)
		managerRole = newRole
		var permissions []models.Permission
		err = suite.db.Model(newRole).Association("Permissions").Find(&permissions)
		require.NoError(err)
		assert.Len(permissions, 4, "Should have 4 permissions (2 resources * 2 actions)")
	})
	suite.T().Run("should allow owner to invite a new member", func(t *testing.T) {
		require.NotNil(managerRole, "Manager role must be created first")
		inviteDTO := &invitation.InvitationDTO{
			Email:  "new.manager@example.com",
			RoleID: managerRole.ID,
		}
		newInvite, err := suite.invitationService.CreateInvitation(org.ID, founder.ID, inviteDTO)
		require.NoError(err)
		require.NotNil(newInvite)
		assert.Equal("pending", newInvite.Status)
		assert.Equal(org.ID, newInvite.OrganizationID)
		assert.Equal(managerRole.ID, newInvite.RoleID)
		invite = newInvite
	})
	suite.T().Run("new user should be able to accept invitation", func(t *testing.T) {
		require.NotNil(invite, "Invitation must be created first")
		acceptDTO := &invitation.AcceptInvitationDTO{
			Token:     invite.Token,
			FirstName: "New",
			LastName:  "Manager",
			Password:  "a-Secure-Manager-Password123!",
		}
		err := suite.invitationService.AcceptInvitation(acceptDTO)
		require.NoError(err)
		newMember, err := suite.authRepo.FindIdentityByEmail("new.manager@example.com")
		require.NoError(err)
		require.NotNil(newMember)
		assert.True(newMember.EmailVerified)
		var acceptedInvite models.Invitation
		suite.db.First(&acceptedInvite, "id = ?", invite.ID)
		assert.Equal("accepted", acceptedInvite.Status)
		var newMembership models.OrganizationalMembership
		err = suite.db.Where("identity_id = ? AND organization_id = ?", newMember.ID, org.ID).First(&newMembership).Error
		require.NoError(err)
		assert.Equal(managerRole.ID, newMembership.RoleID)
	})
}

func (suite *AuthLifecycleTestSuite) TestPhase5_SecurityAndEdgeCases() {
	require := require.New(suite.T())
	assert := assert.New(suite.T())
	suite.T().Run("should prevent cross-tenant access", func(t *testing.T) {
		t.Skip("Skipping test: Tenant isolation is enforced by middleware and requires HTTP-level testing, which is not yet set up.")
	})
	suite.T().Run("should create audit logs for key events", func(t *testing.T) {
		t.Skip("Skipping test: Audit log verification requires integration with the AuditService, which is not yet set up.")
	})
	suite.T().Run("should allow password recovery flow", func(t *testing.T) {
		err := suite.authService.ForgotPassword(&auth.ForgotPasswordDTO{Email: suite.adminUser.Email})
		require.NoError(err)
		var resetToken models.PasswordResetToken
		err = suite.db.Where("identity_id = ?", suite.adminUser.ID).First(&resetToken).Error
		require.NoError(err, "Password reset token was not created in the database")
		require.NotEmpty(resetToken.Token)
		newPassword := "a-Brand-New-Password123!"
		err = suite.authService.ResetPassword(&auth.ResetPasswordDTO{
			Token:       resetToken.Token,
			NewPassword: newPassword,
		})
		require.NoError(err)
		var tokenAfterUse models.PasswordResetToken
		err = suite.db.Where("token = ?", resetToken.Token).First(&tokenAfterUse).Error
		assert.ErrorIs(err, gorm.ErrRecordNotFound, "Reset token should be deleted after use")
		loginDTO := &auth.LoginDTO{Email: suite.adminUser.Email, Password: newPassword}
		_, err = suite.authService.Login(loginDTO)
		assert.NoError(err, "Login with new password should succeed")
		oldLoginDTO := &auth.LoginDTO{Email: suite.adminUser.Email, Password: "a-Strong-Password123!"}
		_, err = suite.authService.Login(oldLoginDTO)
		assert.Error(err, "Login with old password should fail")
	})
}
