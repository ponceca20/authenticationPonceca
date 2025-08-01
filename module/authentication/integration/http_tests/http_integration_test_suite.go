package http_tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"practicev2/config"
	"practicev2/database"
	"practicev2/module/authentication"
	"practicev2/module/authentication/auth"
	"practicev2/module/authentication/invitation"
	"practicev2/module/authentication/organization"
	"practicev2/module/authentication/role"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// HTTPIntegrationTestSuite prueba el sistema completo de autenticación a través de HTTP
type HTTPIntegrationTestSuite struct {
	suite.Suite
	app      *fiber.App
	db       *gorm.DB
	baseURL  string
	testPort string
	cleanup  []func()

	// Tokens para diferentes contextos
	adminToken     string
	customerToken  string
	guestSessionID string

	// IDs de entidades creadas
	testOrgID      string
	testUserID     string
	testRoleID     string
	testCustomerID string
}

// SetupSuite configura el entorno de testing HTTP completo
func (suite *HTTPIntegrationTestSuite) SetupSuite() {
	t := suite.T()

	// Configurar variables de entorno para testing
	os.Setenv("ENV", "testing")
	os.Setenv("DB_NAME", "gastos_ia_test")

	// Leer configuración desde .env para obtener puerto y host
	config.Init()

	// Usar puerto de test diferente para no interferir con desarrollo
	suite.testPort = "3031"
	if testPort := os.Getenv("TEST_PORT"); testPort != "" {
		suite.testPort = testPort
	}

	suite.baseURL = fmt.Sprintf("http://127.0.0.1:%s", suite.testPort)

	// Configurar base de datos de testing
	suite.setupTestDatabase()

	// Crear aplicación Fiber con todas las rutas
	suite.app = fiber.New(fiber.Config{
		DisableStartupMessage: true,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
				"code":  code,
			})
		},
	})

	// Registrar todas las rutas de autenticación
	authentication.RegisterRoutes(suite.app)

	// Inicializar datos de prueba base
	suite.seedTestData()

	t.Log("🚀 HTTP Integration Test Suite configurado correctamente")
}

// setupTestDatabase configura la base de datos para testing
func (suite *HTTPIntegrationTestSuite) setupTestDatabase() {
	// Crear la base de datos de testing si no existe
	suite.createTestDatabaseIfNotExists()

	// Conectar a la base de datos de testing
	database.ConnectDatabase()
	suite.db = database.DBconn
	require.NotNil(suite.T(), suite.db)

	// Ejecutar migraciones para crear las tablas
	err := authentication.RunMigrations(suite.db)
	require.NoError(suite.T(), err, "Error ejecutando migraciones")

	suite.T().Log("✅ Migraciones ejecutadas correctamente")

	// Limpiar tablas si existen (orden importante por foreign keys)
	// Deshabilitar las restricciones de foreign key temporalmente
	suite.db.Exec("SET FOREIGN_KEY_CHECKS = 0")

	tables := []string{
		"role_permission",
		"organizational_membership",
		"password_reset_token",
		"refresh_token",
		"audit_log",
		"invitation",
		"shipping_address",
		"customer_preference",
		"customer_profile",
		"guest_session",
		"user_profile",
		"role",
		"department",
		"organization",
		"identity",
	}

	for _, table := range tables {
		result := suite.db.Exec(fmt.Sprintf("DELETE FROM %s", table))
		if result.Error != nil {
			suite.T().Logf("Warning: Error cleaning table %s: %v", table, result.Error)
		}
	}

	// Rehabilitar las restricciones de foreign key
	suite.db.Exec("SET FOREIGN_KEY_CHECKS = 1")

	suite.T().Log("✅ Base de datos de testing limpia")
}

// createTestDatabaseIfNotExists crea la base de datos de testing si no existe
func (suite *HTTPIntegrationTestSuite) createTestDatabaseIfNotExists() {
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	testDBName := os.Getenv("DB_NAME") // should be gastos_ia_test for testing

	// Conectar a MySQL sin especificar base de datos
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				LogLevel: logger.Error, // Solo errores
			},
		),
	})
	if err != nil {
		suite.T().Fatalf("Error conectando a MySQL para crear base de datos: %v", err)
	}

	// Crear la base de datos si no existe
	createDBSQL := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", testDBName)
	result := db.Exec(createDBSQL)
	if result.Error != nil {
		suite.T().Fatalf("Error creando base de datos de testing %s: %v", testDBName, result.Error)
	}

	suite.T().Logf("✅ Base de datos de testing %s verificada/creada", testDBName)

	// Cerrar la conexión temporal
	sqlDB, _ := db.DB()
	if sqlDB != nil {
		sqlDB.Close()
	}
}

// seedTestData crea datos base necesarios para los tests
func (suite *HTTPIntegrationTestSuite) seedTestData() {
	// Los datos se crearán dinámicamente en cada test para aislamiento
	suite.T().Log("✅ Datos de prueba base preparados")
}

// TearDownSuite limpia el entorno después de todos los tests
func (suite *HTTPIntegrationTestSuite) TearDownSuite() {
	// Ejecutar funciones de limpieza
	for _, cleanup := range suite.cleanup {
		cleanup()
	}

	// Limpiar base de datos
	if suite.db != nil {
		suite.setupTestDatabase() // Reutilizar lógica de limpieza
	}

	suite.T().Log("🧹 HTTP Integration Test Suite limpiado")
}

// TearDownTest limpia después de cada test individual
func (suite *HTTPIntegrationTestSuite) TearDownTest() {
	// Limpiar tokens para aislamiento entre tests
	suite.adminToken = ""
	suite.customerToken = ""
	suite.guestSessionID = ""

	// Limpiar IDs de entidades
	suite.testOrgID = ""
	suite.testUserID = ""
	suite.testRoleID = ""
	suite.testCustomerID = ""
}

// makeRequest realiza una petición HTTP y devuelve la respuesta
func (suite *HTTPIntegrationTestSuite) makeRequest(method, path string, body interface{}, headers map[string]string) (*http.Response, []byte) {
	var reqBody []byte
	var err error

	if body != nil {
		reqBody, err = json.Marshal(body)
		require.NoError(suite.T(), err)
	}

	req := httptest.NewRequest(method, path, bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// Agregar headers adicionales
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := suite.app.Test(req, -1) // -1 = sin timeout
	require.NoError(suite.T(), err)

	respBody := make([]byte, resp.ContentLength)
	if resp.ContentLength > 0 {
		_, err = resp.Body.Read(respBody)
		if err != nil && err.Error() != "EOF" {
			require.NoError(suite.T(), err)
		}
	}

	return resp, respBody
}

// makeAuthenticatedRequest realiza una petición con token de autorización
func (suite *HTTPIntegrationTestSuite) makeAuthenticatedRequest(method, path string, body interface{}, token string) (*http.Response, []byte) {
	headers := map[string]string{
		"Authorization": "Bearer " + token,
	}
	return suite.makeRequest(method, path, body, headers)
}

// parseResponseJSON parsea el cuerpo de respuesta como JSON
func (suite *HTTPIntegrationTestSuite) parseResponseJSON(respBody []byte, target interface{}) {
	err := json.Unmarshal(respBody, target)
	require.NoError(suite.T(), err)
}

// assertSuccessResponse verifica que la respuesta sea exitosa y parsea el JSON
func (suite *HTTPIntegrationTestSuite) assertSuccessResponse(resp *http.Response, respBody []byte, target interface{}) {
	require.True(suite.T(), resp.StatusCode >= 200 && resp.StatusCode < 300,
		"Expected success status, got %d. Body: %s", resp.StatusCode, string(respBody))

	if target != nil {
		suite.parseResponseJSON(respBody, target)
	}
}

// assertErrorResponse verifica que la respuesta sea de error
func (suite *HTTPIntegrationTestSuite) assertErrorResponse(resp *http.Response, expectedStatus int) {
	require.Equal(suite.T(), expectedStatus, resp.StatusCode,
		"Expected status %d, got %d", expectedStatus, resp.StatusCode)
}

// Helper para extraer token de la respuesta de login
func (suite *HTTPIntegrationTestSuite) extractTokenFromLoginResponse(respBody []byte) string {
	var loginResp struct {
		AccessToken string `json:"access_token"`
	}
	suite.parseResponseJSON(respBody, &loginResp)
	require.NotEmpty(suite.T(), loginResp.AccessToken)
	return loginResp.AccessToken
}

// Helper para crear datos de registro de organización
func (suite *HTTPIntegrationTestSuite) createOrganizationRegistrationData(orgType string) organization.OrganizationRegistrationDTO {
	timestamp := time.Now().Unix()

	baseData := organization.OrganizationRegistrationDTO{
		Identity: auth.RegisterDTO{
			FirstName: "Test",
			LastName:  "Admin",
			Email:     fmt.Sprintf("admin-%d@test-%s.com", timestamp, orgType),
			Password:  "TestPassword123!",
		},
		Name:    fmt.Sprintf("Test %s %d", strings.Title(orgType), timestamp),
		Type:    orgType,
		Address: "Test Address 123",
		Phone:   "+57 300 123 4567",
	}

	// Configuración específica por tipo
	switch orgType {
	case "company":
		baseData.Description = "Technology company focused on innovation"
		baseData.Website = "https://techcompany.com"
	case "educational_institution":
		baseData.Description = "Educational institution providing quality education"
		baseData.Website = "https://school.edu.co"
	}

	return baseData
}

// Helper para crear datos de rol
func (suite *HTTPIntegrationTestSuite) createRoleData(roleName string, hierarchyLevel int) role.RoleDTO {
	return role.RoleDTO{
		Name:           roleName,
		DisplayName:    strings.Title(strings.ReplaceAll(roleName, "_", " ")),
		Description:    fmt.Sprintf("Test role: %s", roleName),
		HierarchyLevel: hierarchyLevel,
		Permissions: []role.PermissionDTO{
			{
				Resource: "users",
				Actions:  []string{"read", "create"},
				Scope:    "department",
			},
		},
	}
}

// Helper para crear datos de invitación
func (suite *HTTPIntegrationTestSuite) createInvitationData(email, roleID string) invitation.InvitationDTO {
	return invitation.InvitationDTO{
		Email:  email,
		RoleID: roleID,
	}
}

// Helper para generar email único
func (suite *HTTPIntegrationTestSuite) generateUniqueEmail(prefix string) string {
	timestamp := time.Now().UnixNano() // Usar nanosegundos para mayor precisión
	return fmt.Sprintf("%s-%d@test.com", prefix, timestamp)
}

// Helper para verificar estructura de respuesta estándar
func (suite *HTTPIntegrationTestSuite) verifyStandardResponse(respBody []byte) {
	var response map[string]interface{}
	suite.parseResponseJSON(respBody, &response)

	// Verificar que tiene campos estándar esperados
	require.Contains(suite.T(), response, "success", "Response should have success field")
}

// Helper para obtener slug de organización desde respuesta
func (suite *HTTPIntegrationTestSuite) extractOrgSlugFromResponse(respBody []byte) string {
	var orgResp struct {
		Organization struct {
			Slug string `json:"slug"`
		} `json:"organization"`
	}
	suite.parseResponseJSON(respBody, &orgResp)
	require.NotEmpty(suite.T(), orgResp.Organization.Slug)
	return orgResp.Organization.Slug
}

// Helper para validar formato de token JWT
func (suite *HTTPIntegrationTestSuite) validateJWTFormat(token string) {
	require.NotEmpty(suite.T(), token)
	parts := strings.Split(token, ".")
	require.Equal(suite.T(), 3, len(parts), "JWT should have 3 parts separated by dots")

	// Verificar que cada parte no esté vacía
	for i, part := range parts {
		require.NotEmpty(suite.T(), part, "JWT part %d should not be empty", i)
	}
}

// Helper para convertir puerto de string a int de manera segura
func (suite *HTTPIntegrationTestSuite) getTestPortInt() int {
	port, err := strconv.Atoi(suite.testPort)
	require.NoError(suite.T(), err)
	return port
}

// TestHTTPIntegrationTestSuite ejecuta toda la suite
func TestHTTPIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(HTTPIntegrationTestSuite))
}
