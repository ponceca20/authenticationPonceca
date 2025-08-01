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

// SchoolLifecycleTestSuite tests the complete lifecycle of an educational institution.
// It inherits the setup from AuthLifecycleTestSuite.
type SchoolLifecycleTestSuite struct {
	AuthLifecycleTestSuite
}

// TestSchoolLifecycle runs the entire suite.
func TestSchoolLifecycleTestSuite(t *testing.T) {
	suite.Run(t, new(SchoolLifecycleTestSuite))
}

// TestCompleteSchoolLifecycle contains the full end-to-end test.
func (suite *SchoolLifecycleTestSuite) TestCompleteSchoolLifecycle() {
	t := suite.T()
	require := require.New(t)

	// --- Test Data Definition ---
	schoolOrgData := organization.OrganizationRegistrationDTO{
		Identity: auth.RegisterDTO{
			FirstName: "María",
			LastName:  "González",
			Email:     "director@colegiosanmartin.edu.co",
			Password:  "Director2025!",
		},
		Name: "Colegio San Martín",
		Type: "educational_institution",
	}

	// --- Phase 1: Create the School and its Director ---
	t.Log("Fase 1: Creación del Colegio y su Director")
	schoolOrg, err := suite.orgService.CreateOrganization(&schoolOrgData)
	require.NoError(err)
	require.NotNil(schoolOrg)
	director, err := suite.authRepo.FindIdentityByEmail(schoolOrgData.Identity.Email)
	require.NoError(err)
	require.NotNil(director)
	t.Logf("  ✓ Colegio '%s' y Director '%s' creados.", schoolOrg.Name, director.Email)

	// --- Phase 2: Create Hierarchical Roles ---
	t.Log("Fase 2: Creación de Roles Jerárquicos")
	schoolRoles := map[string]role.RoleDTO{
		"academic_coordinator": {
			Name: "academic_coordinator", DisplayName: "Coordinador Académico", HierarchyLevel: 90,
			Permissions: []role.PermissionDTO{{Resource: "students", Actions: []string{"read", "update"}, Scope: "organization"}},
		},
		"teacher": {
			Name: "teacher", DisplayName: "Profesor", HierarchyLevel: 60,
			Permissions: []role.PermissionDTO{{Resource: "grades", Actions: []string{"read", "create", "update"}, Scope: "own"}},
		},
		"student": {
			Name: "student", DisplayName: "Estudiante", HierarchyLevel: 20,
			Permissions: []role.PermissionDTO{{Resource: "grades", Actions: []string{"read"}, Scope: "own"}},
		},
		"parent": {
			Name: "parent", DisplayName: "Padre de Familia", HierarchyLevel: 30,
			Permissions: []role.PermissionDTO{{Resource: "student_grades", Actions: []string{"read"}, Scope: "own"}},
		},
	}
	createdRoles := make(map[string]*models.Role)
	for name, roleDTO := range schoolRoles {
		newRole, err := suite.roleService.CreateRole(schoolOrg.ID, &roleDTO)
		require.NoError(err)
		createdRoles[name] = newRole
		t.Logf("  ✓ Rol '%s' creado.", newRole.DisplayName)
	}
	require.Len(createdRoles, 4)

	// --- Phase 3: Onboard Staff, Teachers, Students, and Parents ---
	t.Log("Fase 3: Incorporación de Personal, Profesores, Estudiantes y Padres")
	teacher := suite.onboardSchoolUser(
		"patricia.lopez@colegiosanmartin.edu.co", "ProfeMates2025!", "Patricia", "López",
		schoolOrg.ID, director.ID, createdRoles["teacher"],
	)
	student := suite.onboardSchoolUser(
		"juan.perez@estudiante.colegiosanmartin.edu.co", "Estudiante2025!", "Juan", "Pérez",
		schoolOrg.ID, director.ID, createdRoles["student"],
	)
	parent := suite.onboardSchoolUser(
		"pedro.perez@familia.colegiosanmartin.edu.co", "PadreJuan2025!", "Pedro", "Pérez",
		schoolOrg.ID, director.ID, createdRoles["parent"],
	)
	require.NotNil(teacher)
	require.NotNil(student)
	require.NotNil(parent)
	t.Log("  ✓ Todos los tipos de usuarios incorporados exitosamente.")

	// --- Phase 4: Verify Permissions and Security ---
	t.Log("Fase 4: Verificación de Permisos y Seguridad")
	// For this service-level test, we verify that the membership linkage is correct.
	// A full permission check would require HTTP-level tests to validate middleware.
	var teacherMembership models.OrganizationalMembership
	err = suite.db.Where("identity_id = ? AND role_id = ?", teacher.ID, createdRoles["teacher"].ID).First(&teacherMembership).Error
	require.NoError(err, "Teacher should have a membership with the teacher role")

	var studentMembership models.OrganizationalMembership
	err = suite.db.Where("identity_id = ? AND role_id = ?", student.ID, createdRoles["student"].ID).First(&studentMembership).Error
	require.NoError(err, "Student should have a membership with the student role")

	var parentMembership models.OrganizationalMembership
	err = suite.db.Where("identity_id = ? AND role_id = ?", parent.ID, createdRoles["parent"].ID).First(&parentMembership).Error
	require.NoError(err, "Parent should have a membership with the parent role")
	t.Log("  ✓ Membresías y roles de usuarios verificados en la base de datos.")

	t.Log("=== CICLO DE VIDA DEL COLEGIO COMPLETADO (Implementación en progreso) ===")
}

// onboardSchoolUser is a helper function to invite and accept a user into the school organization.
func (suite *SchoolLifecycleTestSuite) onboardSchoolUser(email, password, firstName, lastName, orgID, inviterID string, role *models.Role) *models.Identity {
	require := require.New(suite.T())

	// 1. Director invites the user
	inviteDTO := &invitation.InvitationDTO{Email: email, RoleID: role.ID}
	newInvite, err := suite.invitationService.CreateInvitation(orgID, inviterID, inviteDTO)
	require.NoError(err)

	// 2. The new user accepts the invitation
	acceptDTO := &invitation.AcceptInvitationDTO{
		Token:     newInvite.Token,
		FirstName: firstName,
		LastName:  lastName,
		Password:  password,
	}
	err = suite.invitationService.AcceptInvitation(acceptDTO)
	require.NoError(err)

	// 3. Get the newly created user identity
	identity, err := suite.authRepo.FindIdentityByEmail(email)
	require.NoError(err)
	return identity
}
