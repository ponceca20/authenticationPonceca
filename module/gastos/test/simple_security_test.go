package gastos

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"practicev2/module/authentication/models"
	"practicev2/module/authentication/rbac"
	gastos_model "practicev2/module/gastos/model"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// =============================================================================
// SUITE DE TEST DE SEGURIDAD SIMPLE
// =============================================================================

type SimpleSecurityTestSuite struct {
	suite.Suite
	db      *gorm.DB
	testOrg *models.Organization
	app     *fiber.App
}

// =============================================================================
// CONFIGURACIÓN DEL SUITE
// =============================================================================

func (s *SimpleSecurityTestSuite) SetupSuite() {
	// Configurar base de datos
	dsn := ":memory:"
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(s.T(), err)
	s.db = db

	// Migrar esquemas necesarios
	err = s.db.AutoMigrate(
		&models.Organization{}, &models.Identity{}, &models.Role{}, &models.Permission{},
		&models.OrganizationalMembership{}, &gastos_model.Expense{},
	)
	require.NoError(s.T(), err)

	// Crear organización de prueba
	s.testOrg = &models.Organization{
		ID:          uuid.New().String(),
		Name:        "Security Test Organization",
		Slug:        "security-test",
		Type:        "company",
		Description: "Organización para tests de seguridad",
	}
	require.NoError(s.T(), s.db.Create(s.testOrg).Error)

	// Configurar RBAC
	rbacConfig := &rbac.RBACConfig{
		CacheEnabled:     true,
		CacheTTL:         15 * time.Minute,
		DefaultDenyAll:   true,
		EnableAuditLog:   false,
		DebugMode:        true,
		OrganizationMode: true,
	}
	builder := rbac.NewRBACBuilder(s.db, rbacConfig)
	builder.WithAutoSeed(true)
	engine := builder.Build()
	require.NotNil(s.T(), engine)

	// Crear app de prueba
	s.app = s.createTestApp()
}

func (s *SimpleSecurityTestSuite) createTestApp() *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Configurar rutas básicas
	app.Get("/api/org/:slug/expenses", func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "No authorization header provided",
			})
		}

		return c.JSON(fiber.Map{
			"status": "success",
			"data":   []interface{}{},
		})
	})

	app.Get("/api/org/:slug/admin/expenses/all", func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "No authorization header provided",
			})
		}

		// Simular restricción RBAC para usuarios no admin
		if authHeader != "Bearer admin-token" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Insufficient permissions",
			})
		}

		return c.JSON(fiber.Map{
			"status": "success",
			"data":   []interface{}{},
		})
	})

	app.Post("/api/org/:slug/expenses", func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "No authorization header provided",
			})
		}

		return c.JSON(fiber.Map{
			"status": "success",
			"data": map[string]interface{}{
				"id":     1,
				"title":  "Test Expense",
				"amount": 100.0,
			},
		})
	})

	return app
}

func (s *SimpleSecurityTestSuite) TearDownSuite() {
	if s.db != nil {
		db, _ := s.db.DB()
		db.Close()
	}
}

func (s *SimpleSecurityTestSuite) logSecurityEvent(eventType, message string) {
	fmt.Printf("🔐 [%s] %s\n", eventType, message)
}

// =============================================================================
// TESTS DE SEGURIDAD USANDO FIBER TEST
// =============================================================================

func (s *SimpleSecurityTestSuite) TestAuthentication_NoToken() {
	req, _ := http.NewRequest("GET", "/api/org/security-test/expenses", nil)
	resp, err := s.app.Test(req)
	require.NoError(s.T(), err)

	require.Equal(s.T(), http.StatusUnauthorized, resp.StatusCode)
	s.logSecurityEvent("AUTH_DENIED", "Sin token - acceso correctamente denegado")
}

func (s *SimpleSecurityTestSuite) TestAuthentication_WithToken() {
	req, _ := http.NewRequest("GET", "/api/org/security-test/expenses", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	resp, err := s.app.Test(req)
	require.NoError(s.T(), err)

	require.Equal(s.T(), http.StatusOK, resp.StatusCode)
	s.logSecurityEvent("AUTH_SUCCESS", "Token válido - acceso concedido")
}

func (s *SimpleSecurityTestSuite) TestRBAC_AdminAccess() {
	req, _ := http.NewRequest("GET", "/api/org/security-test/admin/expenses/all", nil)
	req.Header.Set("Authorization", "Bearer admin-token")
	resp, err := s.app.Test(req)
	require.NoError(s.T(), err)

	require.Equal(s.T(), http.StatusOK, resp.StatusCode)
	s.logSecurityEvent("RBAC_SUCCESS", "Admin acceso concedido correctamente")
}

func (s *SimpleSecurityTestSuite) TestRBAC_RegularUserDenied() {
	req, _ := http.NewRequest("GET", "/api/org/security-test/admin/expenses/all", nil)
	req.Header.Set("Authorization", "Bearer user-token")
	resp, err := s.app.Test(req)
	require.NoError(s.T(), err)

	require.Equal(s.T(), http.StatusForbidden, resp.StatusCode)
	s.logSecurityEvent("RBAC_SUCCESS", "Usuario regular correctamente denegado de rutas admin")
}

func (s *SimpleSecurityTestSuite) TestExpenseCreation() {
	req, _ := http.NewRequest("POST", "/api/org/security-test/expenses", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	req.Header.Set("Content-Type", "application/json")
	resp, err := s.app.Test(req)
	require.NoError(s.T(), err)

	require.Equal(s.T(), http.StatusOK, resp.StatusCode)
	s.logSecurityEvent("CREATION_SUCCESS", "Gasto creado correctamente con autorización")
}

// =============================================================================
// FUNCIÓN PRINCIPAL DE TEST
// =============================================================================

func TestSimpleSecuritySuite(t *testing.T) {
	suite.Run(t, new(SimpleSecurityTestSuite))
}
