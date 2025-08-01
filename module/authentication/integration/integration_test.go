package integration

import (
	"fmt"
	"practicev2/config"
	"practicev2/module/authentication/audit"
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
	auditService      audit.AuditService
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
	auditRepo := audit.NewAuditRepository(suite.db)

	// Services
	suite.guestService = guest.NewGuestService(guestRepo)
	suite.customerService = customer.NewCustomerService(suite.customerRepo, suite.authRepo, suite.guestService)
	suite.authService = auth.NewAuthService(suite.authRepo, suite.jwtService)
	suite.orgService = organization.NewOrganizationService(orgRepo, suite.authRepo)
	suite.roleService = role.NewRoleService(roleRepo)
	suite.invitationService = invitation.NewInvitationService(invitationRepo)
	suite.auditService = audit.NewAuditService(auditRepo)
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
		assert.Nil(retrievedSession.CartData)
		cartJSON := `{"items":[{"id":"prod_123","quantity":2}],"total":99.98}`
		updatedSession, err := suite.guestService.UpdateCart(session.SessionToken, cartJSON)
		require.NoError(err)
		assert.NotNil(updatedSession.CartData)
		assert.Equal(cartJSON, *updatedSession.CartData)
		retrievedAfterUpdate, err := suite.guestService.GetSession(session.SessionToken)
		require.NoError(err)
		assert.NotNil(retrievedAfterUpdate.CartData)
		assert.Equal(cartJSON, *retrievedAfterUpdate.CartData)
	})
	suite.T().Run("should convert guest to registered user", func(t *testing.T) {
		// First, create a guest session with some cart data
		session, err := suite.guestService.CreateGuestSession()
		require.NoError(err)
		require.NotNil(session)

		cartJSON := `{"items":[{"id":"prod_456","quantity":1}],"total":49.99}`
		_, err = suite.guestService.UpdateCart(session.SessionToken, cartJSON)
		require.NoError(err)

		// Now convert the guest to a registered user
		convertDTO := &auth.ConvertGuestDTO{
			GuestSessionToken: session.SessionToken,
			FirstName:         "Former",
			LastName:          "Guest",
			Email:             "former.guest@example.com",
			Password:          "a-Strong-Guest-Password123!",
		}

		newIdentity, err := suite.customerService.RegisterCustomerFromGuest(convertDTO)
		require.NoError(err)
		require.NotNil(newIdentity)
		assert.Equal("Former", newIdentity.FirstName)
		assert.Equal("Guest", newIdentity.LastName)
		assert.Equal("former.guest@example.com", newIdentity.Email)

		// Verify that a customer profile was created
		customerProfile, err := suite.customerService.GetCustomerProfile(newIdentity.ID)
		require.NoError(err)
		require.NotNil(customerProfile)
		assert.NotEmpty(customerProfile.CustomerNumber)

		// Verify that the guest session was deleted
		_, err = suite.guestService.GetSession(session.SessionToken)
		assert.Error(err, "Guest session should be deleted after conversion")

		// Verify that the new user can log in
		loginDTO := &auth.LoginDTO{
			Email:    "former.guest@example.com",
			Password: "a-Strong-Guest-Password123!",
		}
		loginResponse, err := suite.authService.Login(loginDTO)
		require.NoError(err)
		require.NotEmpty(loginResponse.AccessToken)
		assert.Equal(newIdentity.ID, loginResponse.Identity.ID)
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
		// Test Strategy: Create two separate organizations and verify that users from one
		// organization cannot access resources from another organization at the service level

		// Step 1: Create first organization (TenantA)
		orgADTO := &organization.OrganizationRegistrationDTO{
			Identity: auth.RegisterDTO{
				FirstName: "CEO", LastName: "TenantA", Email: "ceo.tenanta@example.com", Password: "a-Very-Strong-Password123!",
			},
			Name: "Tenant A Corp", Type: "company",
		}
		orgA, err := suite.orgService.CreateOrganization(orgADTO)
		require.NoError(err, "Should create organization A")

		// Step 2: Create second organization (TenantB)
		orgBDTO := &organization.OrganizationRegistrationDTO{
			Identity: auth.RegisterDTO{
				FirstName: "CEO", LastName: "TenantB", Email: "ceo.tenantb@example.com", Password: "a-Very-Strong-Password123!",
			},
			Name: "Tenant B Corp", Type: "company",
		}
		orgB, err := suite.orgService.CreateOrganization(orgBDTO)
		require.NoError(err, "Should create organization B")

		// Step 3: Get the founders (they should only access their own organization)
		founderA, err := suite.authRepo.FindIdentityByEmail(orgADTO.Identity.Email)
		require.NoError(err)
		founderB, err := suite.authRepo.FindIdentityByEmail(orgBDTO.Identity.Email)
		require.NoError(err)

		// Step 4: Create roles in each organization
		roleADTO := &role.RoleDTO{
			Name:        "manager",
			DisplayName: "Manager A",
			Permissions: []role.PermissionDTO{
				{Resource: "users", Actions: []string{"read"}, Scope: "organization"},
			},
		}
		roleA, err := suite.roleService.CreateRole(orgA.ID, roleADTO)
		require.NoError(err, "Should create role in organization A")

		roleBDTO := &role.RoleDTO{
			Name:        "manager",
			DisplayName: "Manager B",
			Permissions: []role.PermissionDTO{
				{Resource: "users", Actions: []string{"read"}, Scope: "organization"},
			},
		}
		roleB, err := suite.roleService.CreateRole(orgB.ID, roleBDTO)
		require.NoError(err, "Should create role in organization B")

		// Step 5: Verify cross-tenant isolation at the service level

		// Test 5a: Verify that roles from org A cannot be accessed when querying for org B
		rolesInOrgA, err := suite.roleService.ListRoles(orgA.ID)
		require.NoError(err)
		rolesInOrgB, err := suite.roleService.ListRoles(orgB.ID)
		require.NoError(err)

		// Verify that each organization only sees its own roles
		foundRoleAInOrgA := false
		foundRoleBInOrgA := false
		for _, roleDto := range rolesInOrgA {
			if roleDto.ID == roleA.ID {
				foundRoleAInOrgA = true
			}
			if roleDto.ID == roleB.ID {
				foundRoleBInOrgA = true
			}
		}

		foundRoleAInOrgB := false
		foundRoleBInOrgB := false
		for _, roleDto := range rolesInOrgB {
			if roleDto.ID == roleA.ID {
				foundRoleAInOrgB = true
			}
			if roleDto.ID == roleB.ID {
				foundRoleBInOrgB = true
			}
		}

		// Assertions for role isolation
		assert.True(foundRoleAInOrgA, "Organization A should see its own role")
		assert.False(foundRoleBInOrgA, "Organization A should NOT see organization B's role")
		assert.True(foundRoleBInOrgB, "Organization B should see its own role")
		assert.False(foundRoleAInOrgB, "Organization B should NOT see organization A's role")

		// Test 5b: Verify organizational membership isolation
		// Get memberships for each organization
		var membershipsOrgA []models.OrganizationalMembership
		err = suite.db.Where("organization_id = ?", orgA.ID).Find(&membershipsOrgA).Error
		require.NoError(err)

		var membershipsOrgB []models.OrganizationalMembership
		err = suite.db.Where("organization_id = ?", orgB.ID).Find(&membershipsOrgB).Error
		require.NoError(err)

		// Verify that each organization only has its own members
		foundFounderAInOrgA := false
		foundFounderBInOrgA := false
		for _, membership := range membershipsOrgA {
			if membership.IdentityID == founderA.ID {
				foundFounderAInOrgA = true
			}
			if membership.IdentityID == founderB.ID {
				foundFounderBInOrgA = true
			}
		}

		foundFounderAInOrgB := false
		foundFounderBInOrgB := false
		for _, membership := range membershipsOrgB {
			if membership.IdentityID == founderA.ID {
				foundFounderAInOrgB = true
			}
			if membership.IdentityID == founderB.ID {
				foundFounderBInOrgB = true
			}
		}

		// Assertions for membership isolation
		assert.True(foundFounderAInOrgA, "Founder A should be member of organization A")
		assert.False(foundFounderBInOrgA, "Founder B should NOT be member of organization A")
		assert.True(foundFounderBInOrgB, "Founder B should be member of organization B")
		assert.False(foundFounderAInOrgB, "Founder A should NOT be member of organization B")

		// Test 5c: Create invitations and verify they're scoped to organizations
		inviteADTO := &invitation.InvitationDTO{
			Email:  "employee.a@example.com",
			RoleID: roleA.ID,
		}
		inviteA, err := suite.invitationService.CreateInvitation(orgA.ID, founderA.ID, inviteADTO)
		require.NoError(err, "Should create invitation for organization A")

		inviteBDTO := &invitation.InvitationDTO{
			Email:  "employee.b@example.com",
			RoleID: roleB.ID,
		}
		inviteB, err := suite.invitationService.CreateInvitation(orgB.ID, founderB.ID, inviteBDTO)
		require.NoError(err, "Should create invitation for organization B")

		// Verify invitation isolation
		var invitationsOrgA []models.Invitation
		err = suite.db.Where("organization_id = ?", orgA.ID).Find(&invitationsOrgA).Error
		require.NoError(err)

		var invitationsOrgB []models.Invitation
		err = suite.db.Where("organization_id = ?", orgB.ID).Find(&invitationsOrgB).Error
		require.NoError(err)

		foundInviteAInOrgA := false
		foundInviteBInOrgA := false
		for _, inv := range invitationsOrgA {
			if inv.ID == inviteA.ID {
				foundInviteAInOrgA = true
			}
			if inv.ID == inviteB.ID {
				foundInviteBInOrgA = true
			}
		}

		foundInviteAInOrgB := false
		foundInviteBInOrgB := false
		for _, inv := range invitationsOrgB {
			if inv.ID == inviteA.ID {
				foundInviteAInOrgB = true
			}
			if inv.ID == inviteB.ID {
				foundInviteBInOrgB = true
			}
		}

		// Assertions for invitation isolation
		assert.True(foundInviteAInOrgA, "Organization A should see its own invitation")
		assert.False(foundInviteBInOrgA, "Organization A should NOT see organization B's invitation")
		assert.True(foundInviteBInOrgB, "Organization B should see its own invitation")
		assert.False(foundInviteAInOrgB, "Organization B should NOT see organization A's invitation")

		// Test 5d: Verify audit logs are properly scoped (if we have organization-scoped logs)
		if suite.auditService != nil {
			// Log some organization-specific events
			err = suite.auditService.LogAction(
				founderA.ID,
				"role.create",
				"role",
				roleA.ID,
				"success",
				&orgA.ID,
				`{"role_name":"manager","organization":"Tenant A Corp"}`,
				"192.168.1.100",
				"test-browser",
			)
			require.NoError(err)

			err = suite.auditService.LogAction(
				founderB.ID,
				"role.create",
				"role",
				roleB.ID,
				"success",
				&orgB.ID,
				`{"role_name":"manager","organization":"Tenant B Corp"}`,
				"192.168.1.101",
				"test-browser",
			)
			require.NoError(err)

			// Query audit logs for each organization and verify isolation
			var auditLogsOrgA []models.AuditLog
			err = suite.db.Where("organization_id = ?", orgA.ID).Find(&auditLogsOrgA).Error
			require.NoError(err)

			var auditLogsOrgB []models.AuditLog
			err = suite.db.Where("organization_id = ?", orgB.ID).Find(&auditLogsOrgB).Error
			require.NoError(err)

			// Verify audit log isolation
			foundOrgALogInOrgA := false
			foundOrgBLogInOrgA := false
			for _, log := range auditLogsOrgA {
				if log.ResourceID == roleA.ID && log.Action == "role.create" {
					foundOrgALogInOrgA = true
				}
				if log.ResourceID == roleB.ID && log.Action == "role.create" {
					foundOrgBLogInOrgA = true
				}
			}

			foundOrgALogInOrgB := false
			foundOrgBLogInOrgB := false
			for _, log := range auditLogsOrgB {
				if log.ResourceID == roleA.ID && log.Action == "role.create" {
					foundOrgALogInOrgB = true
				}
				if log.ResourceID == roleB.ID && log.Action == "role.create" {
					foundOrgBLogInOrgB = true
				}
			}

			assert.True(foundOrgALogInOrgA, "Organization A should see its own audit logs")
			assert.False(foundOrgBLogInOrgA, "Organization A should NOT see organization B's audit logs")
			assert.True(foundOrgBLogInOrgB, "Organization B should see its own audit logs")
			assert.False(foundOrgALogInOrgB, "Organization B should NOT see organization A's audit logs")
		}

		// Test 5e: Verify that attempting to use a role from another organization fails
		// This simulates what would happen if someone tried to bypass tenant isolation

		// Try to create an invitation in Org A using a role from Org B (should fail)
		invalidInviteDTO := &invitation.InvitationDTO{
			Email:  "malicious.user@example.com",
			RoleID: roleB.ID, // This role belongs to Org B, not Org A
		}

		// This should fail because roleB.ID doesn't belong to orgA.ID
		_, err = suite.invitationService.CreateInvitation(orgA.ID, founderA.ID, invalidInviteDTO)
		assert.Error(err, "Should not be able to create invitation with role from different organization")
		if err != nil {
			assert.Contains(err.Error(), "role not found", "Error should indicate role validation failure")
		}
	})
	suite.T().Run("should create audit logs for key events", func(t *testing.T) {
		// Test 1: Log a login event
		err := suite.auditService.LogAction(
			suite.adminUser.ID,
			"user.login",
			"identity",
			suite.adminUser.ID,
			"success",
			nil, // system-level event
			`{"ip":"192.168.1.100","user_agent":"test-browser"}`,
			"192.168.1.100",
			"test-browser",
		)
		require.NoError(err, "Should be able to log login action")

		// Test 2: Log a password change event
		err = suite.auditService.LogAction(
			suite.adminUser.ID,
			"user.password_change",
			"identity",
			suite.adminUser.ID,
			"success",
			nil,
			`{"method":"manual"}`,
			"192.168.1.100",
			"test-browser",
		)
		require.NoError(err, "Should be able to log password change action")

		// Test 3: Log a failed login attempt
		err = suite.auditService.LogAction(
			suite.adminUser.ID,
			"user.login",
			"identity",
			suite.adminUser.ID,
			"failure",
			nil,
			`{"reason":"invalid_credentials","attempts":1}`,
			"192.168.1.101",
			"test-browser",
		)
		require.NoError(err, "Should be able to log failed login action")

		// Test 4: Verify audit logs were created
		auditLogs, err := suite.auditService.GetAuditLogsByIdentityID(suite.adminUser.ID)
		require.NoError(err, "Should be able to retrieve audit logs")
		assert.GreaterOrEqual(len(auditLogs), 3, "Should have at least 3 audit log entries")

		// Test 5: Verify the content of specific audit logs
		loginSuccessFound := false
		passwordChangeFound := false
		loginFailureFound := false

		for _, log := range auditLogs {
			assert.Equal(suite.adminUser.ID, log.IdentityID, "All logs should belong to the admin user")
			assert.NotEmpty(log.ID, "Audit log should have an ID")
			assert.NotZero(log.Timestamp, "Audit log should have a timestamp")

			switch log.Action {
			case "user.login":
				if log.Status == "success" {
					loginSuccessFound = true
					assert.Equal("identity", log.Resource)
					assert.Equal("192.168.1.100", log.IPAddress)
					assert.Equal("test-browser", log.UserAgent)
				} else if log.Status == "failure" {
					loginFailureFound = true
					assert.Equal("192.168.1.101", log.IPAddress)
					assert.Contains(log.Details, "invalid_credentials")
				}
			case "user.password_change":
				passwordChangeFound = true
				assert.Equal("success", log.Status)
				assert.Contains(log.Details, "manual")
			}
		}

		assert.True(loginSuccessFound, "Should find successful login audit log")
		assert.True(passwordChangeFound, "Should find password change audit log")
		assert.True(loginFailureFound, "Should find failed login audit log")

		// Test 6: Create an organization-scoped audit log
		// First create an organization for testing
		orgDTO := &organization.OrganizationRegistrationDTO{
			Identity: auth.RegisterDTO{
				FirstName: "Audit", LastName: "TestOrg", Email: "audit.testorg@example.com", Password: "a-Very-Strong-Password123!",
			},
			Name: "Audit Test Corp", Type: "company",
		}
		org, err := suite.orgService.CreateOrganization(orgDTO)
		require.NoError(err)

		founder, err := suite.authRepo.FindIdentityByEmail(orgDTO.Identity.Email)
		require.NoError(err)

		// Log an organization-scoped event
		err = suite.auditService.LogAction(
			founder.ID,
			"organization.create",
			"organization",
			org.ID,
			"success",
			&org.ID, // organization-scoped event
			`{"organization_name":"Audit Test Corp","type":"company"}`,
			"192.168.1.102",
			"test-browser",
		)
		require.NoError(err, "Should be able to log organization creation")

		// Verify organization-scoped audit log
		founderAuditLogs, err := suite.auditService.GetAuditLogsByIdentityID(founder.ID)
		require.NoError(err)

		orgCreationFound := false
		for _, log := range founderAuditLogs {
			if log.Action == "organization.create" && log.Status == "success" {
				orgCreationFound = true
				assert.Equal(org.ID, *log.OrganizationID, "Should have correct organization ID")
				assert.Equal("organization", log.Resource)
				assert.Equal(org.ID, log.ResourceID)
				assert.Contains(log.Details, "Audit Test Corp")
			}
		}
		assert.True(orgCreationFound, "Should find organization creation audit log")
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
