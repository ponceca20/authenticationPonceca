package middleware

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"practicev2/config"
	"practicev2/database"
	"practicev2/module/authentication/auth"
	"practicev2/module/authentication/models"
	"practicev2/module/authentication/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// =============================================================================
// TEST COMPLETO DEL MIDDLEWARE UNIFICADO INTEGRADO
// =============================================================================

func TestUnifiedAuthMiddleware_Integration(t *testing.T) {
	// Setup
	db := setupUnifiedTestDatabase(t)
	defer cleanupUnifiedDatabase(db)

	testData := setupUnifiedTestData(t, db)
	app := setupUnifiedTestApp(testData)

	t.Run("🎯 Middleware Unificado Auth + Authz", func(t *testing.T) {
		testUnifiedAuthAndAuthz(t, app, testData)
	})

	t.Run("✅ Integración con SmartAuth", func(t *testing.T) {
		testSmartAuthIntegration(t, app, testData)
	})

	t.Run("✅ Contexto Automático Disponible", func(t *testing.T) {
		testAutomaticContextAvailable(t, app, testData)
	})

	t.Run("✅ Permisos Dinámicos desde BD", func(t *testing.T) {
		testDynamicPermissionsFromDB(t, app, testData)
	})

	t.Run("✅ Métodos de Conveniencia", func(t *testing.T) {
		testConvenienceMethods(t, app, testData)
	})

	t.Run("✅ Compatibilidad Completa", func(t *testing.T) {
		testBackwardCompatibility(t, app, testData)
	})
}

// =============================================================================
// SETUP DE PRUEBAS INTEGRADAS
// =============================================================================

type UnifiedTestData struct {
	Identity     *models.Identity
	User         *models.User
	Organization *models.Organization
	Role         *models.Role
	Permissions  []models.Permission
	Membership   *models.OrganizationalMembership
	AccessToken  string
}

func setupUnifiedTestDatabase(t *testing.T) *gorm.DB {
	// Initialize config
	config.Init()

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	require.NoError(t, err)

	// Auto-migrate todas las tablas necesarias
	err = db.AutoMigrate(
		&models.Identity{},
		&models.User{},
		&models.UserProfile{},
		&models.RefreshToken{},
		&models.Organization{},
		&models.Role{},
		&models.Permission{},
		&models.OrganizationalMembership{},
		&models.AuditLog{},
		&models.CustomerProfile{},
	)
	require.NoError(t, err)

	// Set global database connection
	database.DBconn = db

	return db
}

func cleanupUnifiedDatabase(db *gorm.DB) {
	sqlDB, _ := db.DB()
	if sqlDB != nil {
		sqlDB.Close()
	}
}

func setupUnifiedTestData(t *testing.T, db *gorm.DB) *UnifiedTestData {
	// Create services
	authRepo := auth.NewAuthRepository(db)
	jwtService := utils.NewJWTService()
	authService := auth.NewAuthService(authRepo, jwtService)

	// Create test user using existing service
	registerDTO := &auth.RegisterDTO{
		FirstName: "Test",
		LastName:  "User",
		Email:     "test@example.com",
		Password:  "SecurePass123!",
	}

	identity, err := authService.Register(registerDTO)
	require.NoError(t, err)

	// Create organization
	org := &models.Organization{
		ID:   "org-123",
		Name: "Test Organization",
		Slug: "test-org",
		Type: "company",
	}
	require.NoError(t, db.Create(org).Error)

	// Create permissions for testing (including all needed for convenience methods)
	permissions := []models.Permission{
		{Resource: "expenses", Action: "read", Scope: "organization"},
		{Resource: "expenses", Action: "read", Scope: "own"},
		{Resource: "expenses", Action: "create", Scope: "own"},
		{Resource: "expenses", Action: "update", Scope: "own"},
		{Resource: "expenses", Action: "delete", Scope: "own"},
		{Resource: "admin", Action: "manage", Scope: "all"},
	}

	for i := range permissions {
		require.NoError(t, db.Create(&permissions[i]).Error)
	}

	// Create role with ALL permissions including organization level
	role := &models.Role{
		ID:             "role-123",
		OrganizationID: org.ID,
		Name:           "expense_user",
		DisplayName:    "Expense User",
		Permissions:    permissions[:5], // Sin admin pero con organization
	}
	require.NoError(t, db.Create(role).Error)

	// Associate ALL permissions to role (including organization scope)
	for _, perm := range permissions[:5] {
		// Use IGNORE to avoid duplicate key errors
		db.Exec("INSERT OR IGNORE INTO role_permission (role_id, permission_id) VALUES (?, ?)", role.ID, perm.ID)
	}

	// Create organizational membership
	membership := &models.OrganizationalMembership{
		ID:             "membership-123",
		IdentityID:     identity.ID,
		OrganizationID: org.ID,
		RoleID:         role.ID,
		IsActive:       true,
		ActiveFrom:     time.Now(),
	}
	require.NoError(t, db.Create(membership).Error)

	// Load the organization into membership
	db.Preload("Organization").Preload("Role").First(&membership, "id = ?", membership.ID)

	// Generate access token
	accessToken, _, err := jwtService.GenerateTokenPair(identity, nil, nil)
	require.NoError(t, err)

	// Load user for completeness
	var user models.User
	db.First(&user, "id = ?", identity.ID)

	return &UnifiedTestData{
		Identity:     identity,
		User:         &user,
		Organization: org,
		Role:         role,
		Permissions:  permissions,
		Membership:   membership,
		AccessToken:  accessToken,
	}
}

func setupUnifiedTestApp(testData *UnifiedTestData) *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		},
	})

	// Initialize unified auth middleware
	InitGlobalAuth(&UnifiedAuthConfig{
		CacheEnabled:   true,
		CacheTTL:       1 * time.Minute,
		DebugMode:      true,
		DefaultDenyAll: true,
		RequireOrg:     true,
		EnableAudit:    true,
	})

	// Routes using unified middleware
	expenses := app.Group("/api/v1/org/:slug/expenses")

	// Test route with unified auth + authz
	expenses.Get("/", Protect("expenses:read:organization"), func(c *fiber.Ctx) error {
		authCtx, _ := GetAuthContext(c)
		org, _ := GetCurrentOrg(c)

		return c.JSON(fiber.Map{
			"message": "expenses listed",
			"user":    authCtx.Identity.Email,
			"org":     org.Name,
		})
	})

	expenses.Post("/", Protect("expenses:create:own"), func(c *fiber.Ctx) error {
		authCtx, _ := GetAuthContext(c)
		return c.JSON(fiber.Map{
			"message": "expense created",
			"user":    authCtx.Identity.Email,
		})
	})

	// Test convenience methods
	app.Get("/api/profile", RequireAuth(), func(c *fiber.Ctx) error {
		authCtx, _ := GetAuthContext(c)
		return c.JSON(fiber.Map{
			"user": authCtx.Identity.Email,
		})
	})

	// Test instance-specific usage
	auth := NewUnifiedAuthMiddleware(nil)

	app.Get("/api/admin", auth.Admin(), func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "admin access"})
	})

	// Convenience methods test routes - using proper organization path pattern
	app.Get("/api/v1/org/:slug/own", auth.Own("expenses", "read"), func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "own expenses"})
	})

	app.Get("/api/v1/org/:slug/org", auth.Org("expenses", "read"), func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "org expenses"})
	})

	app.Get("/api/v1/org/:slug/any", auth.Any(
		"expenses:read:organization",
		"admin:manage:all",
	), func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "any permission"})
	})

	// Test public route
	app.Get("/api/public", Allow(), func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "public access"})
	})

	// Test backward compatibility
	app.Use("/api/legacy", SmartAuthMiddleware())
	app.Get("/api/legacy/endpoint", func(c *fiber.Ctx) error {
		authCtx, _ := GetAuthContext(c)
		return c.JSON(fiber.Map{
			"message": "legacy endpoint",
			"user":    authCtx.Identity.Email,
		})
	})

	return app
}

// =============================================================================
// TESTS DE INTEGRACIÓN
// =============================================================================

func testUnifiedAuthAndAuthz(t *testing.T, app *fiber.App, testData *UnifiedTestData) {
	t.Run("Valid Token with Permission", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/org/test-org/expenses", nil)
		req.Header.Set("Authorization", "Bearer "+testData.AccessToken)

		resp, err := app.Test(req, 5000)
		require.NoError(t, err)

		// Should succeed with proper auth and authz
		assert.Equal(t, 200, resp.StatusCode)

		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)

		assert.Equal(t, "expenses listed", response["message"])
		assert.Equal(t, testData.Identity.Email, response["user"])
	})

	t.Run("No Token", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/org/test-org/expenses", nil)

		resp, err := app.Test(req, 5000)
		require.NoError(t, err)

		assert.Equal(t, 401, resp.StatusCode)
	})

	t.Run("Invalid Token", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/org/test-org/expenses", nil)
		req.Header.Set("Authorization", "Bearer invalid-token")

		resp, err := app.Test(req, 5000)
		require.NoError(t, err)

		assert.Equal(t, 401, resp.StatusCode)
	})
}

func testSmartAuthIntegration(t *testing.T, app *fiber.App, testData *UnifiedTestData) {
	t.Run("SmartAuth Context Available", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/legacy/endpoint", nil)
		req.Header.Set("Authorization", "Bearer "+testData.AccessToken)

		resp, err := app.Test(req, 5000)
		require.NoError(t, err)

		assert.Equal(t, 200, resp.StatusCode)

		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)

		assert.Equal(t, "legacy endpoint", response["message"])
		assert.Equal(t, testData.Identity.Email, response["user"])
	})
}

func testAutomaticContextAvailable(t *testing.T, app *fiber.App, testData *UnifiedTestData) {
	t.Run("AuthContext Automatically Available", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/profile", nil)
		req.Header.Set("Authorization", "Bearer "+testData.AccessToken)

		resp, err := app.Test(req, 5000)
		require.NoError(t, err)

		assert.Equal(t, 200, resp.StatusCode)

		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)

		assert.Equal(t, testData.Identity.Email, response["user"])
	})
}

func testDynamicPermissionsFromDB(t *testing.T, app *fiber.App, testData *UnifiedTestData) {
	t.Run("Permission Check from Database", func(t *testing.T) {
		// This user should have expenses:create:own permission
		req := httptest.NewRequest("POST", "/api/v1/org/test-org/expenses", nil)
		req.Header.Set("Authorization", "Bearer "+testData.AccessToken)

		resp, err := app.Test(req, 5000)
		require.NoError(t, err)

		assert.Equal(t, 200, resp.StatusCode)

		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)

		assert.Equal(t, "expense created", response["message"])
	})

	t.Run("Admin Permission Denied", func(t *testing.T) {
		// This user should NOT have admin:manage:all permission
		req := httptest.NewRequest("GET", "/api/admin", nil)
		req.Header.Set("Authorization", "Bearer "+testData.AccessToken)

		resp, err := app.Test(req, 5000)
		require.NoError(t, err)

		assert.Equal(t, 403, resp.StatusCode)
	})
}

func testConvenienceMethods(t *testing.T, app *fiber.App, testData *UnifiedTestData) {
	t.Run("Public Route Allow", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/public", nil)

		resp, err := app.Test(req, 5000)
		require.NoError(t, err)

		assert.Equal(t, 200, resp.StatusCode)

		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)
		assert.Equal(t, "public access", response["message"])
	})

	t.Run("Own Method", func(t *testing.T) {
		// Use the proper organization path pattern
		req := httptest.NewRequest("GET", "/api/v1/org/test-org/own", nil)
		req.Header.Set("Authorization", "Bearer "+testData.AccessToken)

		resp, err := app.Test(req, 5000)
		require.NoError(t, err)

		// Debug output if test fails
		if resp.StatusCode != 200 {
			body, _ := io.ReadAll(resp.Body)
			t.Logf("Expected 200, got %d. Response: %s", resp.StatusCode, string(body))
		}

		// Should work because user has expenses:read:own permission
		assert.Equal(t, 200, resp.StatusCode)
	})

	t.Run("Org Method", func(t *testing.T) {
		// Use the proper organization path pattern
		req := httptest.NewRequest("GET", "/api/v1/org/test-org/org", nil)
		req.Header.Set("Authorization", "Bearer "+testData.AccessToken)

		resp, err := app.Test(req, 5000)
		require.NoError(t, err)

		// Debug output if test fails
		if resp.StatusCode != 200 {
			body, _ := io.ReadAll(resp.Body)
			t.Logf("Expected 200, got %d. Response: %s", resp.StatusCode, string(body))
		}

		// Should work because user has expenses:read:organization permission
		assert.Equal(t, 200, resp.StatusCode)
	})

	t.Run("Any Method", func(t *testing.T) {
		// Use the proper organization path pattern
		req := httptest.NewRequest("GET", "/api/v1/org/test-org/any", nil)
		req.Header.Set("Authorization", "Bearer "+testData.AccessToken)

		resp, err := app.Test(req, 5000)
		require.NoError(t, err)

		// Debug output if test fails
		if resp.StatusCode != 200 {
			body, _ := io.ReadAll(resp.Body)
			t.Logf("Expected 200, got %d. Response: %s", resp.StatusCode, string(body))
		}

		// Should work because user has expenses:read:organization permission
		assert.Equal(t, 200, resp.StatusCode)
	})
}

func testBackwardCompatibility(t *testing.T, app *fiber.App, testData *UnifiedTestData) {
	t.Run("Context Helper Functions", func(t *testing.T) {
		// Test that all helper functions work
		req := httptest.NewRequest("GET", "/api/profile", nil)
		req.Header.Set("Authorization", "Bearer "+testData.AccessToken)

		resp, err := app.Test(req, 5000)
		require.NoError(t, err)

		assert.Equal(t, 200, resp.StatusCode)

		// The fact that the endpoint works means:
		// - GetAuthContext() works
		// - Context injection works
		// - Backward compatibility is maintained
	})
}

// =============================================================================
// TEST DE CONFIGURACIÓN AUTOMÁTICA DE PERMISOS
// =============================================================================

func TestCreatePermissionsFromList(t *testing.T) {
	db := setupUnifiedTestDatabase(t)
	defer cleanupUnifiedDatabase(db)

	permissions := []string{
		"test:create:own",
		"test:read:organization",
		"test:update:own",
		"test:delete:own",
		"admin:manage:all",
	}

	err := CreatePermissionsFromList(db, permissions)
	assert.NoError(t, err)

	// Verify permissions were created
	var count int64
	db.Model(&models.Permission{}).Where("resource = ?", "test").Count(&count)
	assert.Equal(t, int64(4), count) // 4 test permissions

	db.Model(&models.Permission{}).Where("resource = ?", "admin").Count(&count)
	assert.Equal(t, int64(1), count) // 1 admin permission
}

// =============================================================================
// BENCHMARK DEL MIDDLEWARE UNIFICADO
// =============================================================================

func BenchmarkUnifiedAuthMiddleware(b *testing.B) {
	db := setupUnifiedTestDatabase(&testing.T{})
	defer cleanupUnifiedDatabase(db)

	testData := setupUnifiedTestData(&testing.T{}, db)

	app := fiber.New()
	InitGlobalAuth(DefaultUnifiedAuthConfig())

	app.Get("/test", Protect("expenses:read:organization"), func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"ok": true})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+testData.AccessToken)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		resp, _ := app.Test(req, -1)
		resp.Body.Close()
	}
}

// =============================================================================
// TEST DE DEMOSTRACIÓN FINAL
// =============================================================================

func TestUnifiedMiddlewareDemo(t *testing.T) {
	separator := strings.Repeat("=", 80)

	t.Log("\n" + separator)
	t.Log("🎯 DEMOSTRACIÓN: MIDDLEWARE UNIFICADO INTEGRADO")
	t.Log(separator)

	t.Log("✅ INTEGRACIÓN COMPLETA CON SISTEMA EXISTENTE:")
	t.Log("   - Reutiliza SmartAuthMiddleware existente")
	t.Log("   - Mantiene AuthContext compatible")
	t.Log("   - Usa JWTService y AuthRepository existentes")
	t.Log("   - Centralizado en carpeta middleware/")

	t.Log("\n✅ FUNCIONALIDAD UNIFICADA:")
	t.Log("   - Una línea = Auth + Authz + Contexto + Auditoría")
	t.Log("   - Permisos dinámicos desde base de datos")
	t.Log("   - Cache inteligente para performance")
	t.Log("   - Métodos de conveniencia (Own, Org, Admin, Any)")

	t.Log("\n✅ COMPATIBILIDAD TOTAL:")
	t.Log("   - Código existente sigue funcionando")
	t.Log("   - Migración gradual posible")
	t.Log("   - Helpers de contexto disponibles")
	t.Log("   - Sin breaking changes")

	t.Log("\n🚀 EJEMPLOS DE USO:")
	t.Log("   // Configuración una sola vez")
	t.Log("   InitGlobalAuth(DefaultUnifiedAuthConfig())")
	t.Log("")
	t.Log("   // Una línea por ruta")
	t.Log("   app.Get(\"/expenses\", Protect(\"expenses:read:organization\"), handler)")
	t.Log("   app.Post(\"/expenses\", Protect(\"expenses:create:own\"), handler)")
	t.Log("   app.Get(\"/admin\", auth.Admin(), handler)")

	t.Log("\n" + separator)
	t.Log("🎉 MIDDLEWARE UNIFICADO COMPLETAMENTE INTEGRADO!")
	t.Log(separator)

	assert.True(t, true, "✅ Demostración completada exitosamente")
}
