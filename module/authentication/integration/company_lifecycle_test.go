package integration

import (
	"practicev2/module/authentication/auth"
	"practicev2/module/authentication/invitation"
	"practicev2/module/authentication/models"
	"practicev2/module/authentication/organization"
	"practicev2/module/authentication/role"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// CompanyLifecycleTestSuite tests the complete lifecycle of a company, including e-commerce customer flows.
type CompanyLifecycleTestSuite struct {
	AuthLifecycleTestSuite
}

// TestCompanyLifecycleTestSuite runs the entire suite.
func TestCompanyLifecycleTestSuite(t *testing.T) {
	suite.Run(t, new(CompanyLifecycleTestSuite))
}

// TestCompleteCompanyLifecycle contains the full end-to-end test for a company.
func (suite *CompanyLifecycleTestSuite) TestCompleteCompanyLifecycle() {
	t := suite.T()
	require := require.New(t)

	// --- Test Data Definition ---
	companyOrgData := organization.OrganizationRegistrationDTO{
		Identity: auth.RegisterDTO{
			FirstName: "Ricardo",
			LastName:  "Mendoza",
			Email:     "ceo@techsolutions.com",
			Password:  "CEOPassword2025!",
		},
		Name: "TechSolutions S.A.S",
		Type: "company",
	}

	// --- Phase 1: Create the Company and its CEO ---
	t.Log("Fase 1: Creación de la Empresa y su CEO")
	companyOrg, err := suite.orgService.CreateOrganization(&companyOrgData)
	require.NoError(err)
	require.NotNil(companyOrg)
	ceo, err := suite.authRepo.FindIdentityByEmail(companyOrgData.Identity.Email)
	require.NoError(err)
	require.NotNil(ceo)
	t.Logf("  ✓ Empresa '%s' y CEO '%s' creados.", companyOrg.Name, ceo.Email)

	// --- Phase 2: Create Corporate Roles ---
	t.Log("Fase 2: Creación de Roles Corporativos")
	corporateRoles := map[string]role.RoleDTO{
		"department_manager": {
			Name: "department_manager", DisplayName: "Gerente de Departamento", HierarchyLevel: 90,
			Permissions: []role.PermissionDTO{{Resource: "employees", Actions: []string{"read", "update", "invite"}, Scope: "department"}},
		},
		"senior_developer": {
			Name: "senior_developer", DisplayName: "Desarrollador Senior", HierarchyLevel: 70,
			Permissions: []role.PermissionDTO{{Resource: "projects", Actions: []string{"read", "update"}, Scope: "department"}},
		},
		"accountant": {
			Name: "accountant", DisplayName: "Contador", HierarchyLevel: 60,
			Permissions: []role.PermissionDTO{{Resource: "invoices", Actions: []string{"read", "create"}, Scope: "department"}},
		},
	}
	createdRoles := make(map[string]*models.Role)
	for name, roleDTO := range corporateRoles {
		newRole, err := suite.roleService.CreateRole(companyOrg.ID, &roleDTO)
		require.NoError(err)
		createdRoles[name] = newRole
		t.Logf("  ✓ Rol '%s' creado.", newRole.DisplayName)
	}
	require.Len(createdRoles, 3)

	// --- Phase 3: Onboard Employees ---
	t.Log("Fase 3: Incorporación de Empleados")
	suite.onboardCompanyUser(
		"gerente.it@techsolutions.com", "ManagerIT2025!", "David", "Torres",
		companyOrg.ID, ceo.ID, createdRoles["department_manager"],
	)
	suite.onboardCompanyUser(
		"dev.senior@techsolutions.com", "Employee2025!", "Miguel", "Herrera",
		companyOrg.ID, ceo.ID, createdRoles["senior_developer"],
	)
	t.Log("  ✓ Empleados clave incorporados.")

	// --- Phase 4: E-commerce Customer Lifecycle ---
	t.Log("Fase 4: Ciclo de Vida del Cliente de E-commerce")
	suite.testECommerceCustomerFlow()
}

// onboardCompanyUser is a helper function to invite and accept a user into the company.
func (suite *CompanyLifecycleTestSuite) onboardCompanyUser(email, password, firstName, lastName, orgID, inviterID string, role *models.Role) *models.Identity {
	require := require.New(suite.T())
	inviteDTO := &invitation.InvitationDTO{Email: email, RoleID: role.ID}
	newInvite, err := suite.invitationService.CreateInvitation(orgID, inviterID, inviteDTO)
	require.NoError(err)
	acceptDTO := &invitation.AcceptInvitationDTO{
		Token: newInvite.Token, FirstName: firstName, LastName: lastName, Password: password,
	}
	err = suite.invitationService.AcceptInvitation(acceptDTO)
	require.NoError(err)
	identity, err := suite.authRepo.FindIdentityByEmail(email)
	require.NoError(err)
	return identity
}

func (suite *CompanyLifecycleTestSuite) testECommerceCustomerFlow() {
	t := suite.T()
	require := require.New(t)

	t.Log("  Iniciando flujo de cliente de e-commerce...")

	// 1. Visitante anónimo llega y crea una sesión de invitado
	guestSession, err := suite.guestService.CreateGuestSession()
	require.NoError(err)
	require.NotEmpty(guestSession.SessionToken)
	t.Logf("    ✓ Sesión de invitado creada: %s", guestSession.SessionToken)

	// 2. Invitado añade productos a su carrito
	cartData := `{"items":[{"sku":"TS-001","name":"Enterprise Firewall","price":2999.99,"quantity":1}]}`
	_, err = suite.guestService.UpdateCart(guestSession.SessionToken, cartData)
	require.NoError(err)
	t.Log("    ✓ Carrito de invitado actualizado.")

	// 3. Invitado decide registrarse para completar la compra
	convertDTO := &auth.ConvertGuestDTO{
		GuestSessionToken: guestSession.SessionToken,
		FirstName:         "John",
		LastName:          "Customer",
		Email:             "john.customer@example.com",
		Password:          "a-Good-Customer-Password123!",
	}
	newIdentity, err := suite.customerService.RegisterCustomerFromGuest(convertDTO)
	require.NoError(err)
	require.NotNil(newIdentity)
	t.Logf("    ✓ Invitado convertido a cliente registrado: %s", newIdentity.Email)

	// 4. Verificar que el perfil de cliente fue creado
	customerProfile, err := suite.customerRepo.FindProfileByIdentityID(newIdentity.ID)
	require.NoError(err)
	require.NotNil(customerProfile)
	t.Logf("    ✓ Perfil de cliente creado: %s", customerProfile.CustomerNumber)

	// 5. Verificar que la sesión de invitado fue eliminada
	_, err = suite.guestService.GetSession(guestSession.SessionToken)
	require.Error(err, "La sesión de invitado debería haber sido eliminada después de la conversión")
	t.Log("    ✓ Sesión de invitado eliminada.")
}
