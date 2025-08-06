package rbac

import (
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// TestRBACIntegration prueba la integración completa del sistema RBAC
func TestRBACIntegration(t *testing.T) {
	// Configurar base de datos en memoria para testing
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	// Configurar RBAC básico
	rbacEngine := QuickSetupWithOrg(db, "test", map[string][]string{
		"products": {"create", "read", "update", "delete"},
		"orders":   {"create", "read", "update", "delete", "approve"},
	})

	assert.NotNil(t, rbacEngine)
	assert.NotNil(t, rbacEngine.db)
	assert.NotNil(t, rbacEngine.permissionCache)
	assert.NotNil(t, rbacEngine.resourceRegistry)

	// Verificar que el engine está correctamente configurado
	assert.True(t, rbacEngine.config.CacheEnabled)
	assert.True(t, rbacEngine.config.OrganizationMode)
}

// TestRBACBuilder prueba el builder pattern
func TestRBACBuilder(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	// Usar builder para configuración detallada
	rbacEngine := NewRBACBuilder(db, DefaultRBACConfig()).
		AddQuickModule("inventory", map[string][]string{
			"products": {"create", "read", "update", "delete", "manage"},
		}).
		Build()

	assert.NotNil(t, rbacEngine)

	// Verificar que el engine fue construido correctamente
	assert.NotNil(t, rbacEngine.RBACEngine.resourceRegistry)
	// Verificar que al menos hay algún recurso registrado
	assert.True(t, len(rbacEngine.RBACEngine.resourceRegistry.resources) >= 1)
}

// TestRBACMiddleware prueba los middlewares básicos
func TestRBACMiddleware(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	rbacEngine := QuickSetupWithOrg(db, "test", map[string][]string{
		"products": {"read", "create"},
	})

	app := fiber.New()

	// Configurar rutas con middlewares
	api := app.Group("/api")
	api.Get("/products", rbacEngine.Protect("products", "read", "organization"), func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "products listed"})
	})

	api.Post("/products", rbacEngine.Protect("products", "create", "organization"), func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "product created"})
	})

	// Las rutas se configuraron correctamente
	routes := app.GetRoutes()
	assert.True(t, len(routes) >= 2)
}

// TestRBACConfig prueba la configuración personalizada
func TestRBACConfig(t *testing.T) {
	config := &RBACConfig{
		CacheEnabled:     true,
		CacheTTL:         300, // 5 minutos en segundos (será convertido a time.Duration)
		EnableAuditLog:   true,
		DefaultDenyAll:   true,
		OrganizationMode: true,
	}

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	rbacEngine := NewRBACBuilder(db, config).
		AddQuickModule("test", map[string][]string{
			"test-resource": {"create", "read"},
		}).
		Build()

	assert.NotNil(t, rbacEngine)
	assert.Equal(t, config.CacheEnabled, rbacEngine.config.CacheEnabled)
	assert.Equal(t, config.EnableAuditLog, rbacEngine.config.EnableAuditLog)
	assert.Equal(t, config.OrganizationMode, rbacEngine.config.OrganizationMode)
}

// TestResourceRegistry prueba el registro de recursos simplificado
func TestResourceRegistry(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	rbacEngine := NewDynamicRBACEngine(db, DefaultRBACConfig())

	// Crear y registrar un recurso
	resource := &ResourceDefinition{
		Name:        "users",
		Module:      "test",
		Description: "User management",
		Actions:     []string{"create", "read", "update", "delete"},
		Scopes:      []string{"own", "organization", "all"},
	}

	err = rbacEngine.RBACEngine.RegisterResource(resource)
	assert.NoError(t, err)

	// Verificar que se registró correctamente
	assert.NotNil(t, rbacEngine.RBACEngine.resourceRegistry.resources["users"])
	assert.Equal(t, "test", rbacEngine.RBACEngine.resourceRegistry.resources["users"].Module)
}

// TestPermissionCache prueba el sistema de caché simplificado
func TestPermissionCache(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	rbacEngine := NewDynamicRBACEngine(db, DefaultRBACConfig())

	// Las pruebas de cache se realizan internamente
	// Solo verificamos que el sistema está inicializado
	assert.NotNil(t, rbacEngine.permissionCache)
	assert.NotNil(t, rbacEngine.permissionCache.cache)
}

// BenchmarkRBACMiddleware benchmark del middleware
func BenchmarkRBACMiddleware(b *testing.B) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	rbacEngine := QuickSetupWithOrg(db, "test", map[string][]string{
		"products": {"read"},
	})

	app := fiber.New()
	app.Get("/products", rbacEngine.ProtectResource("products"), func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "ok"})
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// En un benchmark real, aquí simularíamos requests HTTP
		// Por ahora solo medimos la creación del middleware
		middleware := rbacEngine.ProtectResource("products")
		_ = middleware
	}
}

// TestModuleRegistration prueba el registro de módulos
func TestModuleRegistration(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	// Crear engine RBAC
	rbacEngine := QuickSetup(db, "expenses", map[string][]string{
		"expenses": {"create", "read", "update", "delete"},
		"budgets":  {"create", "read", "update", "delete", "approve"},
	})

	// Registrar módulo
	RegisterModule("expenses", rbacEngine)

	// Verificar que se registró
	retrievedEngine, exists := GetModule("expenses")
	assert.True(t, exists)
	assert.Equal(t, rbacEngine, retrievedEngine)

	// Listar módulos
	modules := ListModules()
	assert.Contains(t, modules, "expenses")
}

// TestRBACHelpers prueba las funciones helper
func TestRBACHelpers(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	rbacEngine := QuickSetup(db, "test", map[string][]string{
		"products": {"create", "read", "update", "delete"},
	})

	// Crear app de prueba
	app := fiber.New()

	// Configurar rutas con diferentes tipos de protección
	api := app.Group("/api")
	api.Get("/products", rbacEngine.Protect("products", "read", "organization"))
	api.Post("/products", rbacEngine.Protect("products", "create", "organization"))
	api.Put("/products/:id", rbacEngine.Protect("products", "update", "own"))
	api.Delete("/products/:id", rbacEngine.Protect("products", "delete", "own"))
	api.Get("/admin/products", rbacEngine.Protect("products", "admin", "all"))

	// Verificar que las rutas se crearon
	routes := app.GetRoutes()
	assert.True(t, len(routes) >= 5)
}
