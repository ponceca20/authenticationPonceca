package gastos

import (
	"time"

	"practicev2/module/authentication/middleware"
	"practicev2/module/gastos/expense"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// =============================================================================
// CONFIGURACIÓN DE RUTAS CON MIDDLEWARE UNIFICADO
// =============================================================================

// SetupExpenseRoutes configura todas las rutas del módulo de gastos usando el middleware unificado
func SetupExpenseRoutes(app *fiber.App, db *gorm.DB) {
	// Crear handler de gastos
	handler := expense.NewExpenseHandler(db)

	// =============================================================================
	// INICIALIZACIÓN DEL MIDDLEWARE CON LA DB ESPECÍFICA
	// =============================================================================

	// Configurar middleware unificado con la DB específica
	config := &middleware.UnifiedAuthConfig{
		CacheEnabled:   true,
		CacheTTL:       15 * time.Minute,
		DebugMode:      false,
		DefaultDenyAll: true,
		RequireOrg:     true,
		AllowAnonymous: false,
		EnableAudit:    true,
	}

	// Crear middleware específico para esta DB (para tests)
	var authMiddleware *middleware.UnifiedAuthMiddleware
	if middleware.GlobalAuth == nil {
		// En tests o cuando GlobalAuth no existe, crear uno específico
		authMiddleware = middleware.NewUnifiedAuthMiddlewareWithDB(config, db)
		middleware.InitGlobalAuth(config) // Para compatibilidad
	} else {
		// Usar el global existente
		authMiddleware = middleware.GlobalAuth
	}

	// =============================================================================
	// CREAR PERMISOS NECESARIOS PARA EL MÓDULO DE GASTOS
	// =============================================================================

	expensePermissions := []string{
		"expenses:read:own",          // Leer gastos propios
		"expenses:read:organization", // Leer gastos organizacionales
		"expenses:create:own",        // Crear gastos propios
		"expenses:update:own",        // Actualizar gastos propios
		"expenses:delete:own",        // Eliminar gastos propios
		"admin:manage:all",           // Permisos administrativos
	}

	// Crear permisos automáticamente
	err := middleware.CreatePermissionsFromList(db, expensePermissions)
	if err != nil {
		// Log del error pero no fallar - en producción usar logger
		// log.Printf("Warning: Could not create permissions: %v", err)
	}

	// =============================================================================
	// RUTAS API ORGANIZACIONALES (/api/org/:slug/expenses)
	// =============================================================================

	api := app.Group("/api/org/:slug/expenses")

	// Aplicar autenticación a todas las rutas del grupo usando el middleware específico
	api.Use(authMiddleware.RequireAuth())

	// CRUD básico con permisos específicos
	api.Get("/", authMiddleware.Protect("expenses:read:organization"), handler.GetAllExpenses)
	api.Post("/", authMiddleware.Protect("expenses:create:own"), handler.CreateExpense)
	api.Get("/:id", authMiddleware.Protect("expenses:read:organization"), handler.GetExpense)
	api.Put("/:id", authMiddleware.Protect("expenses:update:own"), handler.UpdateExpense)
	api.Delete("/:id", authMiddleware.Protect("expenses:delete:own"), handler.DeleteExpense)

	// Búsquedas con permisos organizacionales
	api.Get("/search/title", authMiddleware.Protect("expenses:read:organization"), handler.SearchByTitle)
	api.Get("/search/amount", authMiddleware.Protect("expenses:read:organization"), handler.GetByAmountRange)

	// =============================================================================
	// RUTAS DE USUARIO (/api/org/:slug/user/*)
	// =============================================================================

	user := app.Group("/api/org/:slug/user")
	user.Use(authMiddleware.RequireAuth()) // Autenticación requerida para todas las rutas de usuario

	// Dashboard y gastos personales
	user.Get("/expenses", authMiddleware.Protect("expenses:read:own"), handler.GetUserExpenses)
	user.Get("/dashboard", authMiddleware.Protect("expenses:read:own"), handler.GetUserDashboard)

	// =============================================================================
	// RUTAS ADMINISTRATIVAS (/api/org/:slug/admin/*)
	// =============================================================================

	admin := app.Group("/api/org/:slug/admin")
	admin.Use(authMiddleware.RequireAuth()) // Autenticación requerida

	// Solo administradores pueden acceder
	admin.Get("/expenses/all", authMiddleware.Protect("admin:manage:all"), handler.GetAllExpenses)
}
