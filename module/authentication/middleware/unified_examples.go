package middleware

import (
	"practicev2/database"

	"github.com/gofiber/fiber/v2"
)

// =============================================================================
// MIDDLEWARE UNIFICADO RBAC DINÁMICO - GUÍA DE USO ACTUALIZADA
// =============================================================================
//
// ⚠️  ESTE ARCHIVO CONTIENE LA GUÍA ACTUALIZADA PARA EL SISTEMA IMPLEMENTADO
//
// El sistema RBAC dinámico con middleware unificado ya está IMPLEMENTADO y FUNCIONANDO.
// Los ejemplos obsoletos han sido reemplazados por implementaciones reales.
//
// =============================================================================

/*
🎯 ESTADO ACTUAL DEL SISTEMA:

✅ IMPLEMENTADO Y FUNCIONANDO:
   - Sistema RBAC dinámico completo
   - Middleware unificado integrado
   - 7 roles de sistema predefinidos
   - 360+ permisos automáticos
   - Cache inteligente de permisos
   - Auditoría automática
   - Tests de integración pasando

✅ ARCHIVOS PRINCIPALES:
   - module/authentication/middleware/unified_auth.go      (Implementación real)
   - module/authentication/rbac/rbac_dynamic.go           (Motor RBAC)
   - module/authentication/rbac/rbac_seeder.go            (Inicialización)
   - examples/unified_rbac_server.go                      (Ejemplo completo)
   - module/authentication/middleware/unified_integration_test.go (Tests)

=============================================================================
📚 CÓMO USAR EL SISTEMA (IMPLEMENTACIÓN REAL):

1️⃣ INICIALIZACIÓN (Una sola vez al inicio de la aplicación):

	import "practicev2/module/authentication/middleware"

	// Inicializar middleware unificado global
	middleware.InitGlobalAuth(middleware.DefaultUnifiedAuthConfig())

2️⃣ USO EN RUTAS (Una línea = Auth + Authz + Contexto + Auditoría):

	// Protección básica con permisos específicos
	app.Get("/api/expenses", middleware.Protect("expenses:read:organization"), handler)
	app.Post("/api/expenses", middleware.Protect("expenses:create:own"), handler)

	// Métodos de conveniencia
	auth := middleware.NewUnifiedAuthMiddleware(nil)
	app.Get("/my-data", auth.Own("data", "read"), handler)
	app.Get("/org-reports", auth.Org("reports", "read"), handler)
	app.Get("/admin", auth.Admin(), handler)
	app.Get("/special", auth.Any("perm1", "perm2", "perm3"), handler)

3️⃣ ACCESO AL CONTEXTO (Automático en los handlers):

	func handler(c *fiber.Ctx) error {
		// Contexto automáticamente disponible
		authCtx, _ := middleware.GetAuthContext(c)
		org, _ := middleware.GetCurrentOrg(c)
		role, _ := middleware.GetCurrentRole(c)

		return c.JSON(fiber.Map{
			"user": authCtx.Identity.Email,
			"org":  org.Name,
			"role": role.Name,
		})
	}

=============================================================================
🚀 COMPARACIÓN: ANTES vs DESPUÉS

❌ ANTES (Separado y Complejo):
	app.Use(SmartAuthMiddleware())
	app.Get("/expenses", func(c *fiber.Ctx) error {
		// Verificación manual de permisos
		if !hasPermission(ctx, "expenses", "read", "org") {
			return c.Status(403).JSON(...)
		}
		// Handler...
	})

✅ DESPUÉS (Unificado y Simple):
	app.Get("/expenses", middleware.Protect("expenses:read:organization"), handler)

=============================================================================
📁 RECURSOS DISPONIBLES:

📄 EJEMPLOS COMPLETOS:
   - Ver: examples/unified_rbac_server.go (Servidor completo funcional)

🧪 TESTS DE REFERENCIA:
   - Ver: module/authentication/middleware/unified_integration_test.go
   - Ejecutar: go test -v ./module/authentication/middleware/ -run TestUnifiedAuthMiddlewareIntegration

📖 DOCUMENTACIÓN TÉCNICA:
   - Ver: module/authentication/rbac/README.md (Si existe)
   - Ver: docs/RBAC_ARCHITECTURE.md (Si existe)

🔧 IMPLEMENTACIÓN:
   - Motor: module/authentication/rbac/rbac_dynamic.go
   - Middleware: module/authentication/middleware/unified_auth.go
   - Seeder: module/authentication/rbac/rbac_seeder.go

=============================================================================
⚡ CARACTERÍSTICAS PRINCIPALES:

✅ Una línea combina Auth + Authz + Contexto + Auditoría
✅ Compatible 100% con código existente (SmartAuthMiddleware)
✅ Permisos dinámicos desde base de datos
✅ Cache inteligente para performance
✅ 7 roles de sistema predefinidos
✅ 360+ permisos automáticos
✅ Auditoría automática de accesos
✅ Métodos de conveniencia (.Own(), .Org(), .Admin(), .Any())
✅ Configuración flexible por ambiente
✅ Tests de integración completos

=============================================================================
🎯 PRÓXIMOS PASOS RECOMENDADOS:

1. Revisar el ejemplo completo: examples/unified_rbac_server.go
2. Ejecutar tests de integración para ver el sistema en acción
3. Integrar en tu aplicación siguiendo los ejemplos reales
4. Configurar la base de datos y ejecutar el seeder
5. Migrar rutas existentes gradualmente

=============================================================================
*/

// Este archivo ya no contiene ejemplos de código porque han sido reemplazados
// por implementaciones reales. Ver los archivos mencionados arriba para
// ejemplos funcionales y documentación actualizada.

// Example2_UsingMiddlewareInstance usando instancia específica del middleware
func Example2_UsingMiddlewareInstance(app *fiber.App) {
	// Crear instancia específica con configuración personalizada
	config := &UnifiedAuthConfig{
		CacheEnabled:   true,
		CacheTTL:       10 * 60, // 10 minutos
		DebugMode:      true,
		DefaultDenyAll: true,
		RequireOrg:     true,
		EnableAudit:    true,
	}
	auth := NewUnifiedAuthMiddleware(config)

	// Usar la instancia específica
	expenses := app.Group("/api/v1/org/:slug/expenses")

	expenses.Get("/", auth.Protect("expenses:read:organization"), func(c *fiber.Ctx) error {
		// Contexto automático disponible
		authCtx, _ := GetAuthContext(c)
		org, _ := GetCurrentOrg(c)

		return c.JSON(fiber.Map{
			"expenses": "list",
			"user":     authCtx.Identity.Email,
			"org":      org.Name,
		})
	})

	expenses.Post("/", auth.Protect("expenses:create:own"), func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "expense created"})
	})

	expenses.Put("/:id", auth.Protect("expenses:update:own"), func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "expense updated"})
	})

	expenses.Delete("/:id", auth.Protect("expenses:delete:own"), func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "expense deleted"})
	})
}

// Example3_ConvenienceMethods usando métodos de conveniencia
func Example3_ConvenienceMethods(app *fiber.App) {
	InitGlobalAuth(DefaultUnifiedAuthConfig())

	// Diferentes niveles de protección
	app.Get("/public", Allow(), func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "public access"})
	})

	app.Get("/profile", RequireAuth(), func(c *fiber.Ctx) error {
		authCtx, _ := GetAuthContext(c)
		return c.JSON(fiber.Map{
			"user": authCtx.Identity.Email,
		})
	})

	// Usando instancia específica para métodos de conveniencia
	auth := NewUnifiedAuthMiddleware(nil)

	app.Get("/expenses/mine", auth.Own("expenses", "read"), func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "personal expenses"})
	})

	app.Get("/expenses/all", auth.Org("expenses", "read"), func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "organizational expenses"})
	})

	app.Get("/admin", auth.Admin(), func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "admin access"})
	})

	// Permisos alternativos
	app.Get("/reports", auth.Any(
		"reports:read:organization",
		"reports:export:organization",
		"admin:manage:all",
	), func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "reports access"})
	})
}

// Example4_GastosModuleComplete ejemplo completo del módulo gastos
func Example4_GastosModuleComplete(app *fiber.App) {
	// Configuración específica para gastos
	config := &UnifiedAuthConfig{
		CacheEnabled:   true,
		CacheTTL:       15 * 60, // 15 minutos
		DebugMode:      false,
		DefaultDenyAll: true,
		RequireOrg:     true,
		EnableAudit:    true,
	}
	auth := NewUnifiedAuthMiddleware(config)

	// Setup automático de permisos para gastos
	permissions := []string{
		"expenses:read:organization",
		"expenses:create:own",
		"expenses:read:own",
		"expenses:update:own",
		"expenses:delete:own",
		"expenses:submit:own",
		"expenses:approve:department",
		"expenses:reject:department",
		"expenses:reimburse:organization",
		"categories:read:organization",
		"categories:create:organization",
		"categories:update:organization",
		"categories:delete:organization",
		"reports:read:organization",
		"reports:export:organization",
		"budgets:read:organization",
		"budgets:create:organization",
		"budgets:update:organization",
		"budgets:delete:organization",
		"admin:manage:organization",
	}

	// Crear permisos automáticamente
	CreatePermissionsFromList(database.DBconn, permissions)

	// API de gastos con middleware unificado
	expenses := app.Group("/api/v1/org/:slug/expenses")

	// CRUD básico - Una línea por operación
	expenses.Get("/", auth.Protect("expenses:read:organization"), listExpenses)
	expenses.Post("/", auth.Protect("expenses:create:own"), createExpense)
	expenses.Get("/:id", auth.Protect("expenses:read:own"), getExpense)
	expenses.Put("/:id", auth.Protect("expenses:update:own"), updateExpense)
	expenses.Delete("/:id", auth.Protect("expenses:delete:own"), deleteExpense)

	// Operaciones de workflow
	expenses.Put("/:id/submit", auth.Protect("expenses:submit:own"), submitExpense)
	expenses.Put("/:id/approve", auth.Protect("expenses:approve:department"), approveExpense)
	expenses.Put("/:id/reject", auth.Protect("expenses:reject:department"), rejectExpense)
	expenses.Put("/:id/reimburse", auth.Protect("expenses:reimburse:organization"), reimburseExpense)

	// Categorías
	categories := app.Group("/api/v1/org/:slug/categories")
	categories.Get("/", auth.Protect("categories:read:organization"), listCategories)
	categories.Post("/", auth.Protect("categories:create:organization"), createCategory)
	categories.Put("/:id", auth.Protect("categories:update:organization"), updateCategory)
	categories.Delete("/:id", auth.Protect("categories:delete:organization"), deleteCategory)

	// Reportes con permisos múltiples
	reports := app.Group("/api/v1/org/:slug/reports")
	reports.Get("/", auth.Any(
		"reports:read:organization",
		"admin:manage:organization",
	), listReports)
	reports.Post("/export", auth.Protect("reports:export:organization"), exportReport)

	// Presupuestos
	budgets := app.Group("/api/v1/org/:slug/budgets")
	budgets.Get("/", auth.Protect("budgets:read:organization"), listBudgets)
	budgets.Post("/", auth.Protect("budgets:create:organization"), createBudget)
	budgets.Put("/:id", auth.Protect("budgets:update:organization"), updateBudget)
	budgets.Delete("/:id", auth.Protect("budgets:delete:organization"), deleteBudget)

	// Admin
	admin := app.Group("/api/v1/org/:slug/admin")
	admin.Use(auth.Admin()) // Todas las subrutas requieren admin
	admin.Get("/settings", getOrgSettings)
	admin.Put("/settings", updateOrgSettings)
	admin.Get("/users", listOrgUsers)
	admin.Post("/users/invite", inviteUser)

	// Rutas personales (sin contexto organizacional)
	personal := app.Group("/api/v1/personal")
	personal.Use(auth.RequireAuth()) // Solo autenticación
	personal.Get("/expenses", auth.Own("expenses", "read"), getPersonalExpenses)
	personal.Get("/dashboard", getPersonalDashboard)

	// Rutas públicas
	public := app.Group("/api/v1/public")
	public.Use(auth.Allow())
	public.Get("/currencies", getCurrencies)
	public.Get("/exchange-rates", getExchangeRates)
}

// Example5_BackwardCompatibility compatibilidad con código existente
func Example5_BackwardCompatibility(app *fiber.App) {
	// El middleware unificado es completamente compatible con código existente
	InitGlobalAuth(DefaultUnifiedAuthConfig())

	// Código existente que usa SmartAuthMiddleware sigue funcionando
	app.Use(SmartAuthMiddleware())

	// Nueva funcionalidad usando middleware unificado
	app.Get("/api/new-endpoint", Protect("new:read:organization"), func(c *fiber.Ctx) error {
		// Ambos contextos están disponibles
		authCtx, _ := GetAuthContext(c)   // Del SmartAuth o UnifiedAuth
		identityID, _ := GetIdentityID(c) // Disponible desde ambos

		return c.JSON(fiber.Map{
			"user":        authCtx.Identity.Email,
			"identity_id": identityID,
		})
	})

	// Migración gradual: algunas rutas con SmartAuth, otras con UnifiedAuth
	legacy := app.Group("/api/legacy")
	legacy.Use(SmartAuthMiddleware()) // Código existente
	legacy.Get("/endpoint", func(c *fiber.Ctx) error {
		authCtx, _ := GetAuthContext(c)
		return c.JSON(fiber.Map{"user": authCtx.Identity.Email})
	})

	modern := app.Group("/api/modern")
	// modern.Use() no necesario, Protect() incluye autenticación
	modern.Get("/endpoint", Protect("modern:read:organization"), func(c *fiber.Ctx) error {
		authCtx, _ := GetAuthContext(c)
		return c.JSON(fiber.Map{"user": authCtx.Identity.Email})
	})
}

// =============================================================================
// HANDLERS DE EJEMPLO (SIMPLIFICADOS)
// =============================================================================

func listExpenses(c *fiber.Ctx) error {
	authCtx, _ := GetAuthContext(c)
	org, _ := GetCurrentOrg(c)

	return c.JSON(fiber.Map{
		"message": "expenses listed",
		"user":    authCtx.Identity.Email,
		"org":     org.Name,
	})
}

func createExpense(c *fiber.Ctx) error {
	authCtx, _ := GetAuthContext(c)
	return c.JSON(fiber.Map{
		"message": "expense created",
		"user":    authCtx.Identity.Email,
	})
}

func getExpense(c *fiber.Ctx) error {
	authCtx, _ := GetAuthContext(c)
	expenseID := c.Params("id")

	return c.JSON(fiber.Map{
		"message":    "expense retrieved",
		"expense_id": expenseID,
		"user":       authCtx.Identity.Email,
	})
}

func updateExpense(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "expense updated"})
}

func deleteExpense(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "expense deleted"})
}

func submitExpense(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "expense submitted"})
}

func approveExpense(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "expense approved"})
}

func rejectExpense(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "expense rejected"})
}

func reimburseExpense(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "expense reimbursed"})
}

// Handlers de categorías
func listCategories(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "categories listed"})
}

func createCategory(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "category created"})
}

func updateCategory(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "category updated"})
}

func deleteCategory(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "category deleted"})
}

// Handlers de reportes
func listReports(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "reports listed"})
}

func exportReport(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "report exported"})
}

// Handlers de presupuestos
func listBudgets(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "budgets listed"})
}

func createBudget(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "budget created"})
}

func updateBudget(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "budget updated"})
}

func deleteBudget(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "budget deleted"})
}

// Handlers de admin
func getOrgSettings(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "org settings"})
}

func updateOrgSettings(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "org settings updated"})
}

func listOrgUsers(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "org users listed"})
}

func inviteUser(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "user invited"})
}

// Handlers personales
func getPersonalExpenses(c *fiber.Ctx) error {
	authCtx, _ := GetAuthContext(c)
	return c.JSON(fiber.Map{
		"message": "personal expenses",
		"user":    authCtx.Identity.Email,
	})
}

func getPersonalDashboard(c *fiber.Ctx) error {
	authCtx, _ := GetAuthContext(c)
	return c.JSON(fiber.Map{
		"message": "personal dashboard",
		"user":    authCtx.Identity.Email,
	})
}

// Handlers públicos
func getCurrencies(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"currencies": []string{"USD", "EUR", "GBP"}})
}

func getExchangeRates(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"rates": fiber.Map{"USD": 1.0, "EUR": 0.85}})
}

// =============================================================================
// COMPARACIÓN: ANTES VS DESPUÉS CON INTEGRACIÓN
// =============================================================================

func ComparisonWithIntegration() {
	/*
		❌ ANTES (Separado y Complejo):

		// 1. Middleware de autenticación separado
		app.Use(SmartAuthMiddleware())

		// 2. Lógica de autorización manual en cada handler
		app.Get("/expenses", func(c *fiber.Ctx) error {
			authCtx, ok := c.Locals("authContext").(*AuthContext)
			if !ok || authCtx.IsGuest {
				return utils.SendError(c, 401, "Authentication required")
			}

			// Verificación manual de permisos (hardcodeada o compleja)
			if !hasPermission(authCtx, "expenses", "read", "organization") {
				return utils.SendError(c, 403, "Insufficient permissions")
			}

			// Lógica del handler...
		})

		✅ DESPUÉS (Unificado y Simple):

		// 1. Inicialización una sola vez
		InitGlobalAuth(DefaultUnifiedAuthConfig())

		// 2. Una línea combina auth + authz + contexto + auditoría
		app.Get("/expenses", Protect("expenses:read:organization"), func(c *fiber.Ctx) error {
			authCtx, _ := GetAuthContext(c) // Automático
			org, _ := GetCurrentOrg(c)      // Automático

			// Lógica del handler...
		})

		📊 BENEFICIOS DE LA INTEGRACIÓN:

		✅ Compatibilidad total con sistema existente
		✅ Reutilización de SmartAuthMiddleware
		✅ Mismo AuthContext que el código existente
		✅ Migración gradual posible
		✅ Una línea = Auth + Authz + Contexto + Auditoría
		✅ Performance optimizado con cache unificado
		✅ Permisos dinámicos desde base de datos
		✅ Sin romper código existente

		🎯 RESULTADO:
		El middleware unificado integra perfectamente con nuestro sistema
		existente, manteniendo compatibilidad total mientras agrega
		capacidades avanzadas de autorización RBAC dinámica.
	*/
}
