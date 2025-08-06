package organization_config

import (
	"practicev2/module/authentication/middleware"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

//=============================================================================
// RUTAS PARA CONFIGURACIÓN DINÁMICA DE MÓDULOS ORGANIZACIONALES
//=============================================================================

// SetupOrganizationConfigRoutes configura las rutas para la gestión de módulos organizacionales
func SetupOrganizationConfigRoutes(api fiber.Router, db *gorm.DB) {
	handler := NewOrganizationConfigHandler(db)

	// Inicializar middleware unificado si no está disponible globalmente
	if middleware.GlobalAuth == nil {
		middleware.InitGlobalAuth(middleware.DefaultUnifiedAuthConfig())
	}

	// Grupo de rutas para configuración de módulos organizacionales
	orgConfig := api.Group("/org/:slug/config",
		middleware.RequireAuth()) // Requiere autenticación

	//=========================================================================
	// RUTAS PRINCIPALES DE GESTIÓN DE MÓDULOS
	//=========================================================================

	// Obtener módulos configurados de la organización
	orgConfig.Get("/modules",
		middleware.Protect("org_modules:read:organization"),
		handler.GetOrganizationModules)

	// Instalar un módulo predefinido
	orgConfig.Post("/modules/install",
		middleware.Protect("org_modules:install:organization"),
		handler.InstallModule)

	// Habilitar/deshabilitar módulo
	orgConfig.Put("/modules/:module/enable",
		middleware.Protect("org_modules:enable:organization"),
		handler.EnableModule)

	orgConfig.Put("/modules/:module/disable",
		middleware.Protect("org_modules:disable:organization"),
		handler.DisableModule)

	// Obtener detalles de un módulo específico
	orgConfig.Get("/modules/:module",
		middleware.Protect("org_modules:read:organization"),
		handler.GetModuleDetails)

	//=========================================================================
	// RUTAS DE INFORMACIÓN GENERAL (SIN CONTEXTO ORGANIZACIONAL)
	//=========================================================================

	// Obtener módulos disponibles para instalar (público para admins)
	api.Get("/modules/available",
		middleware.RequireAuth(),
		middleware.Protect("org_modules:read:all"),
		handler.GetAvailableModules)
}

//=============================================================================
// CONFIGURACIÓN PARA USO EN EL MAIN ROUTER
//=============================================================================

// RegisterOrganizationConfigModule registra el módulo completo en el router principal
func RegisterOrganizationConfigModule(app *fiber.App, db *gorm.DB) {
	// Crear grupo API v1
	v1 := app.Group("/api/v1")

	// Configurar rutas
	SetupOrganizationConfigRoutes(v1, db)

	// Registrar el módulo en el motor RBAC dinámico
	registerModuleInRBAC()
}

// registerModuleInRBAC registra este módulo en el sistema RBAC dinámico
func registerModuleInRBAC() {
	// Verificar que GlobalAuth esté inicializado
	if middleware.GlobalAuth == nil {
		middleware.InitGlobalAuth(middleware.DefaultUnifiedAuthConfig())
	}

	// Recursos para configuración de módulos organizacionales
	moduleResources := map[string][]string{
		"module_config": {"read", "create", "update", "delete", "install", "enable", "disable"},
		"org_settings":  {"read", "update", "manage"},
		"org_modules":   {"read", "configure", "install", "uninstall", "enable", "disable"},
	}

	// Registrar el módulo usando el motor RBAC dinámico del GlobalAuth
	if middleware.GlobalAuth != nil && middleware.GlobalAuth.GetRBACEngine() != nil {
		rbacEngine := middleware.GlobalAuth.GetRBACEngine()

		// Registrar usando QuickRegisterModule
		err := rbacEngine.QuickRegisterModule("organization_config", moduleResources)
		if err != nil {
			// Log error pero no fallar el startup
			// logger.Error("Failed to register organization_config module in RBAC:", err)
		}
	}
}

//=============================================================================
// EJEMPLO DE USO EN MAIN.GO
//=============================================================================

/*
func main() {
	app := fiber.New()
	db := database.DBconn // Tu conexión a la base de datos

	// Registrar el módulo completo
	organization_config.RegisterOrganizationConfigModule(app, db)

	app.Listen(":8080")
}

// O si quieres más control, usar SetupOrganizationConfigRoutes directamente:

func setupAPI() {
	app := fiber.New()
	db := database.DBconn

	v1 := app.Group("/api/v1")

	// Otros módulos...
	auth.SetupAuthRoutes(v1, db)
	organization.SetupOrganizationRoutes(v1, db)

	// Tu nuevo módulo de configuración
	organization_config.SetupOrganizationConfigRoutes(v1, db)

	app.Listen(":8080")
}
*/
