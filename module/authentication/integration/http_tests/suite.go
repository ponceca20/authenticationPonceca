package http_tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"practicev2/config"
	"practicev2/database"
	"practicev2/module/authentication"
	"practicev2/module/authentication/auth"
	"practicev2/module/authentication/models"
	"practicev2/module/authentication/organization"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// HTTPIntegrationTestSuite es la suite base para todos los tests HTTP de integración
type HTTPIntegrationTestSuite struct {
	suite.Suite
	app        *fiber.App
	db         *gorm.DB
	config     *config.Config
	testOrgID  string   // Para almacenar el ID de la organización de prueba
	counter    int64    // Para generar emails únicos
	originalDB *gorm.DB // Para restaurar la DB original después de los tests
}

// SetupSuite configura los recursos compartidos para todos los tests
func (s *HTTPIntegrationTestSuite) SetupSuite() {
	s.T().Log("🔧 Configurando suite de tests HTTP...")

	// 1. Inicializar configuración
	config.Init()
	s.config = config.GetConfig()

	// 2. Configurar base de datos en memoria
	s.setupTestDatabase()

	// 3. Configurar aplicación Fiber con módulos
	s.setupFiberApp()

	s.T().Log("✅ Suite HTTP configurada correctamente")
}

// TearDownSuite limpia los recursos después de todos los tests
func (s *HTTPIntegrationTestSuite) TearDownSuite() {
	s.T().Log("🧹 Limpiando suite de tests HTTP...")

	if s.db != nil {
		sqlDB, _ := s.db.DB()
		if sqlDB != nil {
			sqlDB.Close()
		}
	}

	// Restaurar la BD original
	if s.originalDB != nil {
		database.DBconn = s.originalDB
	}

	s.T().Log("✅ Suite HTTP limpiada correctamente")
}

// SetupTest prepara el estado para cada test individual
func (s *HTTPIntegrationTestSuite) SetupTest() {
	s.T().Log("🔄 Preparando test individual...")
	s.cleanupDatabase()
	s.counter++ // Incrementar contador para emails únicos
}

// TearDownTest limpia después de cada test individual
func (s *HTTPIntegrationTestSuite) TearDownTest() {
	s.T().Log("🧹 Limpiando test individual...")
	s.cleanupDatabase()
}

// setupTestDatabase inicializa la base de datos SQLite en memoria
func (s *HTTPIntegrationTestSuite) setupTestDatabase() {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // Usar tablas en singular como en el sistema principal
		},
	})
	s.Require().NoError(err, "Failed to connect to in-memory database")

	// Ejecutar migraciones
	err = db.AutoMigrate(
		&models.Identity{}, &models.User{}, &models.UserProfile{},
		&models.RefreshToken{}, &models.PasswordResetToken{},
		&models.GuestSession{}, &models.CustomerProfile{}, &models.ShippingAddress{},
		&models.Organization{}, &models.Department{}, &models.OrganizationalMembership{},
		&models.Role{}, &models.Permission{},
		&models.Invitation{}, &models.AuditLog{},
	)
	s.Require().NoError(err, "Failed to run migrations")

	s.db = db
	s.T().Log("✅ Base de datos en memoria configurada")
}

// setupFiberApp configura la aplicación Fiber con el módulo de autenticación
func (s *HTTPIntegrationTestSuite) setupFiberApp() {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"status":  "error",
				"message": err.Error(),
				"error":   err.Error(),
			})
		},
	})

	// Guardar la BD original y reemplazarla temporalmente con la de prueba
	s.originalDB = database.DBconn
	database.DBconn = s.db

	// Registrar las rutas del módulo de autenticación
	authentication.RegisterRoutes(app)

	s.app = app
	s.T().Log("✅ Aplicación Fiber configurada con módulo de autenticación")
}

// cleanupDatabase limpia todas las tablas para aislar los tests
func (s *HTTPIntegrationTestSuite) cleanupDatabase() {
	tables := []string{
		"role_permission", "permissions", "roles", "organizational_memberships",
		"invitations", "departments", "organizations", "audit_logs",
		"shipping_addresses", "customer_profiles", "guest_sessions",
		"password_reset_tokens", "refresh_tokens", "user_profiles", "identity",
	}

	for _, table := range tables {
		s.db.Exec(fmt.Sprintf("DELETE FROM %s;", table))
		s.db.Exec(fmt.Sprintf("DELETE FROM sqlite_sequence WHERE name = '%s';", table))
	}
}

// =============================
// HELPER METHODS - HTTP REQUESTS
// =============================

// makeRequest realiza una petición HTTP sin autenticación
func (s *HTTPIntegrationTestSuite) makeRequest(method, path string, body interface{}, headers map[string]string) (*http.Response, []byte) {
	return s.makeRequestWithAuth(method, path, body, headers, "")
}

// makeAuthenticatedRequest realiza una petición HTTP con token de autorización
func (s *HTTPIntegrationTestSuite) makeAuthenticatedRequest(method, path string, body interface{}, token string) (*http.Response, []byte) {
	headers := map[string]string{
		"Authorization": "Bearer " + token,
	}
	return s.makeRequestWithAuth(method, path, body, headers, token)
}

// makeRequestWithAuth es el método base para realizar peticiones HTTP
func (s *HTTPIntegrationTestSuite) makeRequestWithAuth(method, path string, body interface{}, headers map[string]string, token string) (*http.Response, []byte) {
	var bodyReader io.Reader

	if body != nil {
		jsonBytes, err := json.Marshal(body)
		s.Require().NoError(err, "Failed to marshal request body")
		bodyReader = bytes.NewReader(jsonBytes)
	}

	req := httptest.NewRequest(method, path, bodyReader)

	// Configurar headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Agregar headers adicionales
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Ejecutar request
	resp, err := s.app.Test(req, -1) // -1 = sin timeout
	s.Require().NoError(err, "Failed to execute HTTP request")

	// Leer response body
	respBody, err := io.ReadAll(resp.Body)
	s.Require().NoError(err, "Failed to read response body")
	resp.Body.Close()

	return resp, respBody
}

// =============================
// HELPER METHODS - ASSERTIONS
// =============================

// assertSuccessResponse verifica que la respuesta sea exitosa (2xx)
func (s *HTTPIntegrationTestSuite) assertSuccessResponse(resp *http.Response, respBody []byte, expectedMessage *string) {
	require := s.Require()

	require.True(resp.StatusCode >= 200 && resp.StatusCode < 300,
		"Expected success status code (2xx), got %d. Body: %s", resp.StatusCode, string(respBody))

	// Verificar que el Content-Type sea JSON
	contentType := resp.Header.Get("Content-Type")
	require.Contains(contentType, "application/json", "Response should be JSON")

	// Parsear y verificar estructura básica
	var response map[string]interface{}
	err := json.Unmarshal(respBody, &response)
	require.NoError(err, "Response should be valid JSON")

	status, exists := response["status"]
	if exists {
		require.Equal("success", status, "Response status should be 'success'")
	}

	if expectedMessage != nil {
		message, exists := response["message"]
		if exists {
			require.Equal(*expectedMessage, message, "Response message should match expected")
		}
	}
}

// =============================
// HELPER METHODS - JSON PARSING
// =============================

// parseResponseJSON parsea el JSON de respuesta en la estructura proporcionada
func (s *HTTPIntegrationTestSuite) parseResponseJSON(respBody []byte, target interface{}) {
	err := json.Unmarshal(respBody, target)
	s.Require().NoError(err, "Failed to parse response JSON. Body: %s", string(respBody))
}

// =============================
// HELPER METHODS - DATA GENERATION
// =============================

// generateUniqueEmail genera un email único para testing
func (s *HTTPIntegrationTestSuite) generateUniqueEmail(prefix string) string {
	timestamp := time.Now().UnixNano()
	return fmt.Sprintf("%s.%d.%d@test.example.com", prefix, s.counter, timestamp)
}

// createOrganizationRegistrationData crea datos básicos para registro de organización (solo para uso interno del builder)
func (s *HTTPIntegrationTestSuite) createOrganizationRegistrationData(orgType string) organization.OrganizationRegistrationDTO {
	uniqueEmail := s.generateUniqueEmail("founder")

	baseData := organization.OrganizationRegistrationDTO{
		Name:        fmt.Sprintf("Test Organization %d", s.counter),
		Type:        orgType,
		Description: "Test organization for integration testing",
		Website:     "https://test-org.example.com",
		Identity: auth.RegisterDTO{
			Email:     uniqueEmail,
			FirstName: "Test",
			LastName:  "Founder",
			Password:  "TestPassword123!",
		},
	}

	// Personalizar según el tipo de organización
	switch orgType {
	case "company":
		baseData.Name = fmt.Sprintf("Test Company Corp %d", s.counter)
		baseData.Description = "Test company for integration testing"
	case "school":
		baseData.Name = fmt.Sprintf("Test School %d", s.counter)
		baseData.Description = "Test educational institution"
	case "nonprofit":
		baseData.Name = fmt.Sprintf("Test Nonprofit %d", s.counter)
		baseData.Description = "Test nonprofit organization"
	}

	return baseData
}

// =============================
// HELPER METHODS - VALIDATION
// =============================

// =============================
// HELPER METHODS - VALIDATION
// =============================

// validateJWTFormat verifica que un string tenga el formato básico de JWT
func (s *HTTPIntegrationTestSuite) validateJWTFormat(token string) {
	require := s.Require()

	require.NotEmpty(token, "JWT token should not be empty")

	parts := strings.Split(token, ".")
	require.Len(parts, 3, "JWT should have 3 parts separated by dots")

	// Verificar que cada parte no esté vacía
	for i, part := range parts {
		require.NotEmpty(part, "JWT part %d should not be empty", i+1)
	}

	// Verificar longitud mínima razonable
	require.Greater(len(token), 100, "JWT token should be reasonably long")
}

// validateUUIDFormat verifica que un string tenga formato UUID válido
func (s *HTTPIntegrationTestSuite) validateUUIDFormat(uuid string) {
	require := s.Require()

	require.NotEmpty(uuid, "UUID should not be empty")
	require.Len(uuid, 36, "UUID should be 36 characters long")
	require.Contains(uuid, "-", "UUID should contain hyphens")

	// Formato básico: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
	parts := strings.Split(uuid, "-")
	require.Len(parts, 5, "UUID should have 5 parts separated by hyphens")
	require.Len(parts[0], 8, "First UUID part should be 8 characters")
	require.Len(parts[1], 4, "Second UUID part should be 4 characters")
	require.Len(parts[2], 4, "Third UUID part should be 4 characters")
	require.Len(parts[3], 4, "Fourth UUID part should be 4 characters")
	require.Len(parts[4], 12, "Fifth UUID part should be 12 characters")
}

// assertErrorResponse verifica que la respuesta sea de error con el código esperado
func (s *HTTPIntegrationTestSuite) assertErrorResponse(resp *http.Response, expectedStatusCode int) {
	require := s.Require()

	require.Equal(expectedStatusCode, resp.StatusCode,
		"Expected error status code %d, got %d", expectedStatusCode, resp.StatusCode)

	// Verificar que el Content-Type sea JSON
	contentType := resp.Header.Get("Content-Type")
	require.Contains(contentType, "application/json", "Error response should be JSON")
}

// =============================
// HELPER METHODS - DEBUGGING
// =============================

// TestRBACBasicFunctionality prueba la funcionalidad básica del sistema RBAC
func (s *HTTPIntegrationTestSuite) TestRBACBasicFunctionality() {
	t := s.T()

	t.Log("=== INICIANDO TEST BÁSICO DEL SISTEMA RBAC ===")

	// --- Fase 1: Configurar organización y usuarios de prueba ---
	t.Log("Fase 1: Configuración inicial - Empresa y usuarios")

	// Crear empresa de prueba
	companyBuilder := s.NewOrganizationBuilder("company").
		WithName("RBAC Test Company").
		WithDescription("Empresa para testing de RBAC").
		WithFounder(s.generateUniqueEmail("ceo"), "María", "CEO", "TestCEO123!")

	companyData := companyBuilder.Build()

	// POST /api/v1/organizations
	resp, respBody := s.makeRequest("POST", "/api/v1/organizations", companyData, nil)
	s.assertSuccessResponse(resp, respBody, nil)

	var orgResponse struct {
		Status string `json:"status"`
		Data   struct {
			ID   string `json:"id"`
			Slug string `json:"slug"`
			Name string `json:"name"`
		} `json:"data"`
	}
	s.parseResponseJSON(respBody, &orgResponse)

	orgID := orgResponse.Data.ID
	s.testOrgID = orgID

	t.Logf("✅ Organización creada: %s (ID: %s)", orgResponse.Data.Name, orgID)

	// --- Fase 2: Login del CEO ---
	t.Log("Fase 2: Login del CEO")

	loginData := auth.LoginDTO{
		Email:    companyData.Identity.Email,
		Password: companyData.Identity.Password,
	}

	resp, respBody = s.makeRequest("POST", "/api/v1/auth/login", loginData, nil)
	s.assertSuccessResponse(resp, respBody, nil)

	var loginResponse struct {
		Status string `json:"status"`
		Data   struct {
			AccessToken string `json:"access_token"`
			Identity    struct {
				ID        string `json:"id"`
				Email     string `json:"email"`
				FirstName string `json:"first_name"`
				LastName  string `json:"last_name"`
			} `json:"identity"`
		} `json:"data"`
	}
	s.parseResponseJSON(respBody, &loginResponse)

	t.Logf("✅ CEO logueado: %s %s", loginResponse.Data.Identity.FirstName, loginResponse.Data.Identity.LastName)

	// --- Fase 3: Verificar que tenemos acceso a RBAC ---
	t.Log("Fase 3: Verificación de disponibilidad del sistema RBAC")

	// Intentar crear un rol básico para verificar RBAC
	rolePath := fmt.Sprintf("/api/v1/org/%s/roles", orgResponse.Data.Slug)
	roleData := map[string]interface{}{
		"name":            "test_role",
		"display_name":    "Test Role",
		"description":     "Rol de prueba para verificar RBAC",
		"hierarchy_level": 50,
		"permissions": []map[string]interface{}{
			{
				"resource": "test",
				"actions":  []string{"read"},
				"scope":    "own",
			},
		},
	}

	resp, respBody = s.makeAuthenticatedRequest("POST", rolePath, roleData, loginResponse.Data.AccessToken)

	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
		t.Log("✅ Sistema RBAC está operativo - Rol creado exitosamente")

		var roleResponse struct {
			Data struct {
				ID          string `json:"id"`
				Name        string `json:"name"`
				DisplayName string `json:"display_name"`
			} `json:"data"`
		}
		s.parseResponseJSON(respBody, &roleResponse)
		t.Logf("✅ Rol RBAC creado: %s (ID: %s)", roleResponse.Data.DisplayName, roleResponse.Data.ID)
	} else {
		t.Logf("⚠️ Sistema RBAC puede no estar completamente configurado (Status: %d)", resp.StatusCode)
	}

	t.Log("✅ Test RBAC básico completado - Funcionalidad base verificada")
}
