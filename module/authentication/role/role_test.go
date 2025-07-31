package role

import (
	"practicev2/module/authentication/auth"
	"practicev2/module/authentication/models"
	"practicev2/module/authentication/organization"
	"practicev2/module/authentication/test"
	"practicev2/module/authentication/user"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type RoleTestSuite struct {
	suite.Suite
	db          *gorm.DB
	roleService RoleService
	userService user.UserService
	testOrg     *models.Organization
	testUser    *models.Identity
}

func (suite *RoleTestSuite) SetupSuite() {
	suite.db = test.SetupTestDatabase(suite.T())
	roleRepo := NewRoleRepository(suite.db)
	userRepo := user.NewUserRepository(suite.db)
	suite.roleService = NewRoleService(roleRepo)
	suite.userService = user.NewUserService(userRepo)

	// We need an organization and a user to exist for these tests
	orgRepo := organization.NewOrganizationRepository(suite.db)
	authRepo := auth.NewAuthRepository(suite.db)
	orgService := organization.NewOrganizationService(orgRepo, authRepo)

	dto := &organization.OrganizationRegistrationDTO{
		Identity: auth.RegisterDTO{
			FirstName: "Role",
			LastName:  "Tester",
			Email:     "role.tester@example.com",
			Password:  "password123",
		},
		Name: "Role Test Corp",
		Type: "company",
	}
	org, err := orgService.CreateOrganization(dto)
	assert.NoError(suite.T(), err)
	suite.testOrg = org

	user, err := authRepo.FindIdentityByEmail(dto.Identity.Email)
	assert.NoError(suite.T(), err)
	suite.testUser = user
}

func TestRoleTestSuite(t *testing.T) {
	suite.Run(t, new(RoleTestSuite))
}

func (suite *RoleTestSuite) TestRoleLifecycle() {
	// 1. Create a new role
	roleDTO := &RoleDTO{
		Name:        "editor",
		DisplayName: "Content Editor",
		Permissions: []PermissionDTO{
			{Resource: "articles", Actions: []string{"create", "read", "update"}, Scope: "organization"},
		},
	}
	createdRole, err := suite.roleService.CreateRole(suite.testOrg.ID, roleDTO)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "editor", createdRole.Name)

	// 2. Assign the role to a user
	err = suite.userService.AssignRoleToUser(suite.testOrg.ID, suite.testUser.ID, createdRole.ID)
	assert.NoError(suite.T(), err)

	// 3. List users in that role
	users, err := suite.roleService.ListUsersInRole(suite.testOrg.ID, createdRole.ID)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), users, 1)
	assert.Equal(suite.T(), suite.testUser.ID, users[0].ID)

	// 4. Update the role
	updateDTO := &RoleDTO{
		Name:        "editor_v2",
		DisplayName: "Content Editor V2",
		Permissions: []PermissionDTO{
			{Resource: "articles", Actions: []string{"create", "read", "update", "delete"}, Scope: "organization"},
		},
	}
	updatedRole, err := suite.roleService.UpdateRole(suite.testOrg.ID, createdRole.ID, updateDTO)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "editor_v2", updatedRole.Name)

	// 5. Delete the role
	err = suite.roleService.DeleteRole(suite.testOrg.ID, createdRole.ID)
	assert.NoError(suite.T(), err)
}
