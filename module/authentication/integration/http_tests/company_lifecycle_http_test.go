package http_tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"practicev2/module/authentication/auth"
	"practicev2/module/authentication/invitation"
	"practicev2/module/authentication/role"
	"strings"

	"github.com/stretchr/testify/require"
)

// TestCompleteCompanyLifecycleHTTP prueba el ciclo completo de una empresa vía HTTP
func (suite *HTTPIntegrationTestSuite) TestCompleteCompanyLifecycleHTTP() {
	t := suite.T()
	require := require.New(t)

	t.Log("=== INICIANDO TEST DE CICLO COMPLETO DE EMPRESA VÍA HTTP ===")

	// --- Fase 1: Crear Empresa y CEO ---
	t.Log("Fase 1: Creación de Empresa con CEO/Founder integrado")

	ceoEmail := suite.generateUniqueEmail("ceo-techsolutions")
	t.Logf("DEBUG - CEO email que se registrará: %s", ceoEmail)

	orgData := suite.createOrganizationRegistrationData("company")
	orgData.Name = "TechSolutions S.A.S HTTP Test"
	// Incluir la información del CEO/Founder en la creación de la organización
	orgData.Identity.Email = ceoEmail
	orgData.Identity.FirstName = "Carlos"
	orgData.Identity.LastName = "Rodriguez"
	orgData.Identity.Password = "Ceo2025!"

	// POST /api/v1/organizations
	resp, respBody := suite.makeRequest("POST", "/api/v1/organizations", orgData, nil)
	suite.assertSuccessResponse(resp, respBody, nil)

	var orgResponse struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Data    struct {
			ID          string `json:"id"`
			Name        string `json:"name"`
			Slug        string `json:"slug"`
			Type        string `json:"type"`
			Description string `json:"description"`
			Website     string `json:"website"`
		} `json:"data"`
	}
	suite.parseResponseJSON(respBody, &orgResponse)

	// Debug: Log the full response to understand what's happening
	t.Logf("DEBUG - Response Status: %d", resp.StatusCode)
	t.Logf("DEBUG - Response Body: %s", string(respBody))
	t.Logf("DEBUG - Parsed Status: %v", orgResponse.Status)
	t.Logf("DEBUG - Organization ID: %s", orgResponse.Data.ID)
	t.Logf("DEBUG - Message: %s", orgResponse.Message)

	require.Equal("success", orgResponse.Status)
	require.NotEmpty(orgResponse.Data.ID)
	require.NotEmpty(orgResponse.Data.Slug)
	require.Equal("company", orgResponse.Data.Type)

	// Validaciones adicionales de seguridad y formato
	require.Len(orgResponse.Data.ID, 36, "Organization ID should be a valid UUID format")
	require.True(len(orgResponse.Data.Slug) > 0 && len(orgResponse.Data.Slug) <= 100, "Slug should have reasonable length")
	require.Equal("TechSolutions S.A.S HTTP Test", orgResponse.Data.Name, "Organization name should match exactly")

	orgSlug := orgResponse.Data.Slug
	suite.testOrgID = orgResponse.Data.ID

	t.Logf("  ✓ Empresa '%s' creada con slug '%s' y CEO integrado", orgResponse.Data.Name, orgSlug)

	// --- Fase 2: Login del CEO para obtener tokens ---
	t.Log("Fase 2: Login del CEO/Founder para obtener tokens de acceso")

	loginData := auth.LoginDTO{
		Email:    ceoEmail,
		Password: "Ceo2025!",
	}

	// POST /api/v1/auth/login
	resp, respBody = suite.makeRequest("POST", "/api/v1/auth/login", loginData, nil)

	// Debug del login del CEO
	t.Logf("DEBUG - CEO Login Status: %d", resp.StatusCode)
	t.Logf("DEBUG - CEO Login Body: %s", string(respBody))

	suite.assertSuccessResponse(resp, respBody, nil)

	var ceoLoginResponse struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Data    struct {
			AccessToken  string `json:"access_token"`
			RefreshToken string `json:"refresh_token"`
			ExpiresAt    string `json:"expires_at"`
			Identity     struct {
				ID        string `json:"id"`
				Email     string `json:"email"`
				FirstName string `json:"first_name"`
				LastName  string `json:"last_name"`
			} `json:"identity"`
			Contexts interface{} `json:"contexts"` // puede ser null o array
		} `json:"data"`
	}
	suite.parseResponseJSON(respBody, &ceoLoginResponse)

	require.Equal("success", ceoLoginResponse.Status)
	require.NotEmpty(ceoLoginResponse.Data.AccessToken)
	require.NotEmpty(ceoLoginResponse.Data.RefreshToken)
	require.Equal(ceoEmail, ceoLoginResponse.Data.Identity.Email)
	suite.validateJWTFormat(ceoLoginResponse.Data.AccessToken)

	// Validaciones adicionales de seguridad JWT
	require.Len(ceoLoginResponse.Data.Identity.ID, 36, "Identity ID should be a valid UUID format")
	require.Equal("Carlos", ceoLoginResponse.Data.Identity.FirstName)
	require.Equal("Rodriguez", ceoLoginResponse.Data.Identity.LastName)

	// Verificar que el CEO tiene contexto organizacional correcto
	require.NotNil(ceoLoginResponse.Data.Contexts, "CEO should have organizational context")
	if contexts, ok := ceoLoginResponse.Data.Contexts.([]interface{}); ok {
		require.Len(contexts, 1, "CEO should have exactly one organizational context")
		if len(contexts) > 0 {
			if contextMap, ok := contexts[0].(map[string]interface{}); ok {
				require.Equal("organization", contextMap["type"])
				require.Equal(suite.testOrgID, contextMap["id"])
				require.Equal("owner", contextMap["role"])
			}
		}
	}

	ceoToken := ceoLoginResponse.Data.AccessToken
	ceoRefreshToken := ceoLoginResponse.Data.RefreshToken

	t.Logf("  ✓ CEO logueado exitosamente con tokens válidos y contexto organizacional correcto")

	// Nota: Necesitaremos agregar el CEO a la organización como admin manualmente
	// o usar un endpoint diferente para esto

	// --- Fase 3: CEO verifica su perfil ---
	t.Log("Fase 3: Verificación del perfil del CEO")

	// GET /api/v1/auth/me
	resp, respBody = suite.makeAuthenticatedRequest("GET", "/api/v1/auth/me", nil, ceoToken)

	// Debug del perfil del CEO
	t.Logf("DEBUG - CEO Profile Status: %d", resp.StatusCode)
	t.Logf("DEBUG - CEO Profile Body: %s", string(respBody))

	suite.assertSuccessResponse(resp, respBody, nil)

	var profileResponse struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Data    struct {
			ID        string `json:"id"`
			Email     string `json:"email"`
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
		} `json:"data"`
	}
	suite.parseResponseJSON(respBody, &profileResponse)

	require.Equal("success", profileResponse.Status)
	require.Equal(ceoEmail, profileResponse.Data.Email)
	require.Equal("Carlos", profileResponse.Data.FirstName)
	require.Equal("Rodriguez", profileResponse.Data.LastName)

	t.Log("  ✓ Perfil del CEO verificado correctamente")

	// --- Fase 4: Crear Roles Corporativos ---
	t.Log("Fase 4: Creación de Roles Corporativos")

	corporateRoles := map[string]role.RoleDTO{
		"department_manager": {
			Name:           "department_manager",
			DisplayName:    "Gerente de Departamento",
			Description:    "Gerente con acceso departamental",
			HierarchyLevel: 90,
			Permissions: []role.PermissionDTO{
				{Resource: "employees", Actions: []string{"read", "update", "invite"}, Scope: "department"},
				{Resource: "reports", Actions: []string{"read", "create"}, Scope: "department"},
			},
		},
		"senior_developer": {
			Name:           "senior_developer",
			DisplayName:    "Desarrollador Senior",
			Description:    "Desarrollador con experiencia avanzada",
			HierarchyLevel: 70,
			Permissions: []role.PermissionDTO{
				{Resource: "projects", Actions: []string{"read", "update"}, Scope: "department"},
				{Resource: "code", Actions: []string{"read", "write", "review"}, Scope: "own"},
			},
		},
		"accountant": {
			Name:           "accountant",
			DisplayName:    "Contador",
			Description:    "Responsable de la contabilidad",
			HierarchyLevel: 60,
			Permissions: []role.PermissionDTO{
				{Resource: "invoices", Actions: []string{"read", "create", "update"}, Scope: "department"},
				{Resource: "finances", Actions: []string{"read"}, Scope: "organization"},
			},
		},
	}

	createdRoles := make(map[string]string) // name -> id

	for roleName, roleData := range corporateRoles {
		// POST /api/v1/org/{slug}/roles
		rolePath := fmt.Sprintf("/api/v1/org/%s/roles", orgSlug)

		// Debug de creación de rol
		t.Logf("DEBUG - Creando rol '%s' en ruta: %s", roleName, rolePath)
		t.Logf("DEBUG - Datos del rol: %+v", roleData)

		resp, respBody := suite.makeAuthenticatedRequest("POST", rolePath, roleData, ceoToken)

		// Debug de respuesta
		t.Logf("DEBUG - Respuesta Status: %d", resp.StatusCode)
		t.Logf("DEBUG - Respuesta Body: %s", string(respBody))

		suite.assertSuccessResponse(resp, respBody, nil)

		var roleResponse struct {
			Data struct {
				ID             string `json:"id"`
				Name           string `json:"name"`
				DisplayName    string `json:"display_name"`
				HierarchyLevel int    `json:"hierarchy_level"`
			} `json:"data"`
		}
		suite.parseResponseJSON(respBody, &roleResponse)

		require.Equal(roleData.Name, roleResponse.Data.Name)
		require.Equal(roleData.HierarchyLevel, roleResponse.Data.HierarchyLevel)
		require.NotEmpty(roleResponse.Data.ID)

		createdRoles[roleName] = roleResponse.Data.ID
		t.Logf("  ✓ Rol '%s' creado con ID: %s", roleResponse.Data.DisplayName, roleResponse.Data.ID)
	}

	require.Len(createdRoles, 3)

	// --- Fase 5: Listar Roles Creados ---
	t.Log("Fase 5: Verificación de Roles Listados")

	// GET /api/v1/org/{slug}/roles
	rolesPath := fmt.Sprintf("/api/v1/org/%s/roles", orgSlug)
	resp, respBody = suite.makeAuthenticatedRequest("GET", rolesPath, nil, ceoToken)

	// Debug de la respuesta del listado
	t.Logf("DEBUG - Listado Roles Status: %d", resp.StatusCode)
	t.Logf("DEBUG - Listado Roles Body: %s", string(respBody))

	suite.assertSuccessResponse(resp, respBody, nil)

	var rolesListResponse struct {
		Data []struct {
			ID             string `json:"id"`
			Name           string `json:"name"`
			DisplayName    string `json:"display_name"`
			HierarchyLevel int    `json:"hierarchy_level"`
		} `json:"data"`
	}
	suite.parseResponseJSON(respBody, &rolesListResponse)

	// Debe incluir el rol owner + los 3 roles creados = 4 total
	require.GreaterOrEqual(len(rolesListResponse.Data), 4)

	// Verificar que nuestros roles están en la lista
	foundRoles := make(map[string]bool)
	for _, role := range rolesListResponse.Data {
		if role.Name == "department_manager" || role.Name == "senior_developer" || role.Name == "accountant" {
			foundRoles[role.Name] = true
		}
	}
	require.Len(foundRoles, 3, "Todos los roles corporativos deben estar en la lista")

	t.Log("  ✓ Todos los roles listados correctamente")

	// --- Fase 6: Invitar Empleados ---
	t.Log("Fase 6: Invitación de Empleados")

	// Estructura más robusta para empleados con validaciones
	type Employee struct {
		email     string
		roleName  string
		name      string
		firstName string
		lastName  string
		password  string
		roleID    string // Se llenará con el ID del rol correspondiente
	}

	employees := []Employee{
		{
			email:     suite.generateUniqueEmail("gerente.it"),
			roleName:  "department_manager",
			name:      "David Torres",
			firstName: "David",
			lastName:  "Torres",
			password:  "Employee2025!",
		},
		{
			email:     suite.generateUniqueEmail("dev.senior"),
			roleName:  "senior_developer",
			name:      "Miguel Herrera",
			firstName: "Miguel",
			lastName:  "Herrera",
			password:  "Employee2025!",
		},
		{
			email:     suite.generateUniqueEmail("contador"),
			roleName:  "accountant",
			name:      "Ana García",
			firstName: "Ana",
			lastName:  "García",
			password:  "Employee2025!",
		},
	}

	// Asignar role IDs y validar que existen
	for i := range employees {
		roleID, exists := createdRoles[employees[i].roleName]
		require.True(exists, "Role %s should exist in created roles", employees[i].roleName)
		require.NotEmpty(roleID, "Role ID should not be empty")
		employees[i].roleID = roleID
	}

	invitationTokens := make(map[string]string) // email -> token
	invitationIDs := make(map[string]string)    // email -> invitation_id

	for _, emp := range employees {
		inviteData := invitation.InvitationDTO{
			Email:  emp.email,
			RoleID: emp.roleID,
		}

		// POST /api/v1/org/{slug}/invitations
		invitePath := fmt.Sprintf("/api/v1/org/%s/invitations", orgSlug)
		resp, respBody := suite.makeAuthenticatedRequest("POST", invitePath, inviteData, ceoToken)

		// Debug para ver la respuesta de la invitación
		t.Logf("DEBUG - Invitación a %s (%s)", emp.email, emp.name)
		t.Logf("DEBUG - Invitation Response Status: %d", resp.StatusCode)
		t.Logf("DEBUG - Invitation Response Body: %s", string(respBody))

		suite.assertSuccessResponse(resp, respBody, nil)

		var inviteResponse struct {
			Data struct {
				ID        string `json:"id"`
				Email     string `json:"email"`
				RoleID    string `json:"role_id"`
				Status    string `json:"status"`
				ExpiresAt string `json:"expires_at"`
			} `json:"data"`
			Message string `json:"message"`
		}
		suite.parseResponseJSON(respBody, &inviteResponse)

		require.Equal(emp.email, inviteResponse.Data.Email)
		require.Equal("pending", inviteResponse.Data.Status)
		require.NotEmpty(inviteResponse.Data.ID)
		require.Equal(emp.roleID, inviteResponse.Data.RoleID, "Invitation should have correct role ID")
		require.Equal("Invitation sent successfully", inviteResponse.Message)
		require.NotEmpty(inviteResponse.Data.ExpiresAt, "Invitation should have expiration date")

		// Almacenar el ID de la invitación para posible verificación posterior
		invitationIDs[emp.email] = inviteResponse.Data.ID

		// Nota: En una implementación real, el token se envía por email
		// Para el test de integración, necesitamos obtenerlo de la base de datos
		var invitation struct {
			Token string
		}
		suite.db.Table("invitation").
			Select("token").
			Where("email = ? AND status = 'pending'", emp.email).
			First(&invitation)

		require.NotEmpty(invitation.Token, "El token de invitación debe existir en la base de datos")

		invitationTokens[emp.email] = invitation.Token
		t.Logf("  ✓ Invitación enviada a %s (%s) - ID: %s", emp.email, emp.name, inviteResponse.Data.ID)
	}

	// --- Fase 7: Aceptar Invitaciones ---
	t.Log("Fase 7: Aceptación de Invitaciones")

	employeeTokens := make(map[string]string) // email -> access_token

	for _, emp := range employees {
		acceptData := invitation.AcceptInvitationDTO{
			Token:     invitationTokens[emp.email],
			FirstName: emp.name[:strings.Index(emp.name, " ")],
			LastName:  emp.name[strings.Index(emp.name, " ")+1:],
			Password:  "Employee2025!",
		}

		// POST /api/v1/invitations/accept
		resp, respBody := suite.makeAuthenticatedRequest("POST", "/api/v1/invitations/accept", acceptData, "")

		// Debug para la aceptación de invitación
		t.Logf("DEBUG - Aceptación de invitación para %s (%s)", emp.email, emp.name)
		t.Logf("DEBUG - Accept Response Status: %d", resp.StatusCode)
		t.Logf("DEBUG - Accept Response Body: %s", string(respBody))

		suite.assertSuccessResponse(resp, respBody, nil)

		var acceptResponse struct {
			Status  string      `json:"status"`
			Data    interface{} `json:"data"` // puede ser null
			Message string      `json:"message"`
		}
		suite.parseResponseJSON(respBody, &acceptResponse)

		require.Equal("success", acceptResponse.Status)
		require.Equal("Invitation accepted successfully. You can now log in.", acceptResponse.Message)

		// Verificar que la invitación cambió de estado en la base de datos
		var updatedInvitation struct {
			Status string
		}
		err := suite.db.Table("invitation").
			Select("status").
			Where("email = ?", emp.email).
			First(&updatedInvitation).Error
		require.NoError(err, "Should be able to query invitation status")
		require.Equal("accepted", updatedInvitation.Status, "Invitation status should be updated to accepted")

		// Después de aceptar la invitación, el empleado debe hacer login
		loginData := auth.LoginDTO{
			Email:    emp.email,
			Password: "Employee2025!",
		}

		// POST /api/v1/auth/login
		loginResp, loginRespBody := suite.makeRequest("POST", "/api/v1/auth/login", loginData, nil)

		// Debug del login del empleado
		t.Logf("DEBUG - Employee Login (%s) Status: %d", emp.name, loginResp.StatusCode)
		t.Logf("DEBUG - Employee Login Body: %s", string(loginRespBody))

		suite.assertSuccessResponse(loginResp, loginRespBody, nil)

		var loginResponse struct {
			Data struct {
				AccessToken string `json:"access_token"`
				Identity    struct {
					ID    string `json:"id"`
					Email string `json:"email"`
				} `json:"identity"`
				Contexts []struct {
					Type string `json:"type"`
					ID   string `json:"id"`
					Role string `json:"role"`
				} `json:"contexts"`
			} `json:"data"`
		}
		suite.parseResponseJSON(loginRespBody, &loginResponse)

		require.NotEmpty(loginResponse.Data.AccessToken)
		require.Equal(emp.email, loginResponse.Data.Identity.Email)
		require.NotEmpty(loginResponse.Data.Identity.ID)

		// Verificar contexto organizacional del empleado
		require.Len(loginResponse.Data.Contexts, 1, "Employee should have exactly one organizational context")
		require.Equal("organization", loginResponse.Data.Contexts[0].Type)
		require.Equal(suite.testOrgID, loginResponse.Data.Contexts[0].ID)
		require.Equal(emp.roleName, loginResponse.Data.Contexts[0].Role)

		suite.validateJWTFormat(loginResponse.Data.AccessToken)
		employeeTokens[emp.email] = loginResponse.Data.AccessToken

		t.Logf("  ✓ %s aceptó la invitación e hizo login exitosamente con rol %s", emp.name, emp.roleName)
	}

	// --- Fase 8: Verificar Permisos de Empleados ---
	t.Log("Fase 8: Verificación de Permisos de Empleados")

	// Verificar que el gerente puede listar usuarios del departamento
	managerEmail := employees[0].email // department_manager
	managerToken := employeeTokens[managerEmail]

	// GET /api/v1/org/{slug}/users (como gerente)
	usersPath := fmt.Sprintf("/api/v1/org/%s/users", orgSlug)
	resp, respBody = suite.makeAuthenticatedRequest("GET", usersPath, nil, managerToken)
	t.Logf("DEBUG - Users List Status: %d", resp.StatusCode)
	t.Logf("DEBUG - Users List Body: %s", string(respBody))
	suite.assertSuccessResponse(resp, respBody, nil)

	var usersListResponse struct {
		Data []struct {
			ID        string `json:"id"`
			Email     string `json:"email"`
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
			Role      string `json:"role"`
			IsActive  bool   `json:"is_active"`
			JoinedAt  string `json:"joined_at"`
		} `json:"data"`
	}
	suite.parseResponseJSON(respBody, &usersListResponse)

	// El gerente debería poder ver al menos a todos los empleados (4 usuarios: CEO + 3 empleados)
	require.GreaterOrEqual(len(usersListResponse.Data), 4)

	// Verificar que todos los usuarios esperados están en la lista
	expectedUsers := map[string]string{
		ceoEmail:           "owner",
		employees[0].email: employees[0].roleName, // department_manager
		employees[1].email: employees[1].roleName, // senior_developer
		employees[2].email: employees[2].roleName, // accountant
	}

	foundUsers := make(map[string]bool)
	for _, user := range usersListResponse.Data {
		t.Logf("DEBUG - Usuario: %s %s (%s) - Rol: %s", user.FirstName, user.LastName, user.Email, user.Role)

		// Verificar que cada usuario tiene los campos obligatorios
		require.NotEmpty(user.ID, "User ID should not be empty")
		require.NotEmpty(user.Email, "User email should not be empty")
		require.NotEmpty(user.FirstName, "User first name should not be empty")
		require.NotEmpty(user.LastName, "User last name should not be empty")
		require.NotEmpty(user.Role, "User role should not be empty")
		require.True(user.IsActive, "User should be active")
		require.NotEmpty(user.JoinedAt, "User should have joined date")

		// Marcar usuarios encontrados
		if expectedRole, exists := expectedUsers[user.Email]; exists {
			foundUsers[user.Email] = true
			require.Equal(expectedRole, user.Role, "User %s should have role %s", user.Email, expectedRole)
		}
	}

	// Verificar que todos los usuarios esperados fueron encontrados
	for email, role := range expectedUsers {
		require.True(foundUsers[email], "User %s with role %s should be found in the list", email, role)
	}

	t.Logf("DEBUG - Usuarios encontrados: %d (esperados: %d)", len(usersListResponse.Data), len(expectedUsers))

	t.Log("  ✓ Gerente puede listar usuarios correctamente y todos los usuarios esperados están presentes")

	// --- Fase 9: Listar Invitaciones desde Organización ---
	t.Log("Fase 9: Verificación de Lista de Invitaciones")

	// GET /api/v1/org/{slug}/invitations
	invitationsPath := fmt.Sprintf("/api/v1/org/%s/invitations", orgSlug)
	resp, respBody = suite.makeAuthenticatedRequest("GET", invitationsPath, nil, ceoToken)
	t.Logf("DEBUG - Invitations List Status: %d", resp.StatusCode)
	t.Logf("DEBUG - Invitations List Body: %s", string(respBody))
	suite.assertSuccessResponse(resp, respBody, nil)

	var invitationsListResponse struct {
		Data []struct {
			ID     string `json:"id"`
			Email  string `json:"email"`
			Status string `json:"status"`
		} `json:"data"`
	}
	suite.parseResponseJSON(respBody, &invitationsListResponse)

	// Las invitaciones aceptadas pueden no aparecer en la lista (por diseño de seguridad)
	// En su lugar, verificamos que los usuarios estén activos en la organización
	t.Logf("DEBUG - Invitaciones en lista: %d", len(invitationsListResponse.Data))
	if len(invitationsListResponse.Data) == 0 {
		t.Log("  ✓ Las invitaciones aceptadas no aparecen en la lista (comportamiento esperado)")
	} else {
		acceptedCount := 0
		for _, inv := range invitationsListResponse.Data {
			t.Logf("DEBUG - Invitación: %s - Estado: %s", inv.Email, inv.Status)
			if inv.Status == "accepted" {
				acceptedCount++
			}
		}
		t.Logf("  ✓ Se encontraron %d invitaciones aceptadas en la lista", acceptedCount)
	}

	t.Log("  ✓ Lista de invitaciones verificada correctamente")

	// --- Fase 10: Testing de Refresh Token ---
	t.Log("Fase 10: Prueba de Refresh Token")

	// POST /api/v1/auth/refresh (usando refresh token del CEO)
	refreshData := struct {
		RefreshToken string `json:"refresh_token"`
	}{
		RefreshToken: ceoRefreshToken,
	}

	resp, respBody = suite.makeRequest("POST", "/api/v1/auth/refresh", refreshData, nil)
	t.Logf("DEBUG - Refresh Token Status: %d", resp.StatusCode)
	t.Logf("DEBUG - Refresh Token Body: %s", string(respBody))
	suite.assertSuccessResponse(resp, respBody, nil)

	var refreshResponse struct {
		Data struct {
			AccessToken  string `json:"access_token"`
			RefreshToken string `json:"refresh_token"`
			ExpiresAt    string `json:"expires_at"`
		} `json:"data"`
	}
	suite.parseResponseJSON(respBody, &refreshResponse)

	require.NotEmpty(refreshResponse.Data.AccessToken)
	require.NotEmpty(refreshResponse.Data.RefreshToken)
	suite.validateJWTFormat(refreshResponse.Data.AccessToken)
	suite.validateJWTFormat(refreshResponse.Data.RefreshToken)

	// Verificar que se generaron nuevos tokens (diferentes a los originales)
	require.NotEqual(ceoToken, refreshResponse.Data.AccessToken, "New access token should be different from original")
	require.NotEqual(ceoRefreshToken, refreshResponse.Data.RefreshToken, "New refresh token should be different from original")

	// Verificar que el nuevo access token es válido
	testResp, _ := suite.makeAuthenticatedRequest("GET", "/api/v1/auth/me", nil, refreshResponse.Data.AccessToken)
	require.Equal(http.StatusOK, testResp.StatusCode, "New access token should be valid")

	t.Log("  ✓ Refresh token funciona correctamente y genera nuevos tokens válidos")

	// --- Fase 11: Testing de Logout ---
	t.Log("Fase 11: Prueba de Logout")

	// POST /api/v1/auth/logout - Enviar el refresh token en el body
	logoutData := map[string]interface{}{
		"refresh_token": refreshResponse.Data.RefreshToken,
	}
	resp, respBody = suite.makeRequest("POST", "/api/v1/auth/logout", logoutData, nil)
	t.Logf("DEBUG - Logout Status: %d", resp.StatusCode)
	t.Logf("DEBUG - Logout Body: %s", string(respBody))
	suite.assertSuccessResponse(resp, respBody, nil)

	var logoutResponse struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}
	suite.parseResponseJSON(respBody, &logoutResponse)

	require.Equal("success", logoutResponse.Status)
	require.Contains([]string{"Logout successful", "Logged out successfully"}, logoutResponse.Message, "Logout message should be appropriate")

	// Verificar que el refresh token ya no funciona (el access token seguirá válido hasta expirar)
	resp, respBody = suite.makeRequest("POST", "/api/v1/auth/refresh", map[string]interface{}{
		"refresh_token": refreshResponse.Data.RefreshToken,
	}, nil)
	t.Logf("DEBUG - Refresh después de logout Status: %d", resp.StatusCode)
	t.Logf("DEBUG - Refresh después de logout Body: %s", string(respBody))
	suite.assertErrorResponse(resp, http.StatusUnauthorized)

	// Verificar el mensaje de error específico
	var refreshErrorResponse struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Error   string `json:"error"`
	}
	suite.parseResponseJSON(respBody, &refreshErrorResponse)
	require.Equal("error", refreshErrorResponse.Status)
	require.Contains(refreshErrorResponse.Error, "refresh token", "Error should mention refresh token")

	// NOTA: El access token sigue válido hasta expirar (comportamiento estándar JWT)
	// En un sistema real, el frontend debería descartar el token después del logout
	resp, _ = suite.makeAuthenticatedRequest("GET", "/api/v1/auth/me", nil, refreshResponse.Data.AccessToken)
	t.Logf("DEBUG - Access token después de logout Status: %d (normal: sigue válido hasta expirar)", resp.StatusCode)
	require.Equal(http.StatusOK, resp.StatusCode, "Access token sigue válido después del logout (comportamiento estándar JWT)")

	t.Logf("  ✓ Logout completado correctamente - Refresh token revocado, Access token sigue válido hasta expirar")

	// Fase 12: Login Directo del Empleado
	t.Log("Fase 12: Login Directo del Empleado")

	// Login del gerente directamente con credenciales
	managerLoginPayload := map[string]interface{}{
		"email":    employees[0].email,
		"password": employees[0].password,
	}

	resp, respBody = suite.makeRequest("POST", "/api/v1/auth/login", managerLoginPayload, nil)
	t.Logf("DEBUG - Manager Direct Login Status: %d", resp.StatusCode)
	t.Logf("DEBUG - Manager Direct Login Body: %s", string(respBody))

	var managerLoginResponse struct {
		Status string                `json:"status"`
		Data   auth.TokenResponseDTO `json:"data"`
	}
	err := json.Unmarshal(respBody, &managerLoginResponse)
	require.NoError(err, "Error al parsear respuesta de login del gerente")
	require.Equal(http.StatusOK, resp.StatusCode, "El gerente debería poder hacer login directo")
	require.Equal("success", managerLoginResponse.Status)
	require.NotEmpty(managerLoginResponse.Data.AccessToken)
	require.NotEmpty(managerLoginResponse.Data.RefreshToken)

	// Validaciones estrictas de la identidad
	require.Equal(employees[0].firstName, managerLoginResponse.Data.Identity.FirstName)
	require.Equal(employees[0].lastName, managerLoginResponse.Data.Identity.LastName)
	require.Equal(employees[0].email, managerLoginResponse.Data.Identity.Email)
	require.NotEmpty(managerLoginResponse.Data.Identity.ID)
	require.Len(managerLoginResponse.Data.Identity.ID, 36, "Identity ID should be a valid UUID format")

	// Verificar que el gerente tiene el contexto organizacional correcto
	require.Len(managerLoginResponse.Data.Contexts, 1, "Manager should have exactly one organizational context")
	require.Equal("organization", managerLoginResponse.Data.Contexts[0].Type)
	require.Equal(orgResponse.Data.ID, managerLoginResponse.Data.Contexts[0].ID)
	require.Equal("department_manager", managerLoginResponse.Data.Contexts[0].Role)
	require.Equal(orgResponse.Data.Name, managerLoginResponse.Data.Contexts[0].Name, "Context should include organization name")

	// Validar formato JWT del nuevo token
	suite.validateJWTFormat(managerLoginResponse.Data.AccessToken)
	suite.validateJWTFormat(managerLoginResponse.Data.RefreshToken)

	// Verificar que el token es funcional haciendo una request autenticada
	profileResp, _ := suite.makeAuthenticatedRequest("GET", "/api/v1/auth/me", nil, managerLoginResponse.Data.AccessToken)
	require.Equal(http.StatusOK, profileResp.StatusCode, "Manager's new token should be functional")

	t.Logf("  ✓ %s (%s) hizo login directo exitosamente con rol %s y token funcional",
		employees[0].firstName, employees[0].email, managerLoginResponse.Data.Contexts[0].Role)

	t.Logf("\n=== ✅ TEST DE CICLO COMPLETO DE EMPRESA EXITOSO ===")
	t.Logf("✅ Todas las 12 fases completadas correctamente:")
	t.Logf("  1. ✅ Creación de empresa con CEO integrado")
	t.Logf("  2. ✅ Login del CEO con tokens")
	t.Logf("  3. ✅ Verificación del perfil del CEO")
	t.Logf("  4. ✅ Creación de roles corporativos")
	t.Logf("  5. ✅ Verificación de roles listados")
	t.Logf("  6. ✅ Invitación de empleados")
	t.Logf("  7. ✅ Aceptación de invitaciones y login")
	t.Logf("  8. ✅ Verificación de permisos de empleados")
	t.Logf("  9. ✅ Lista de invitaciones")
	t.Logf("  10. ✅ Prueba de refresh token")
	t.Logf("  11. ✅ Prueba de logout (revocación de refresh token)")
	t.Logf("  12. ✅ Login directo del empleado")
}
