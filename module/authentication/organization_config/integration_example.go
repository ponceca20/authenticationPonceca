package organization_config

import (
	"fmt"
	"practicev2/module/authentication/middleware"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// =============================================================================
// EJEMPLO DE INTEGRACIÓN COMPLETA CON RBAC
// =============================================================================

// ExampleCompleteIntegration muestra cómo usar el módulo completo con RBAC
func ExampleCompleteIntegration() {
	// Crear app Fiber
	app := fiber.New()

	// Simular conexión a base de datos (en tu caso sería database.DBconn)
	var db *gorm.DB // Tu instancia real de GORM

	// 1. INICIALIZAR EL SISTEMA RBAC GLOBAL
	middleware.InitGlobalAuth(middleware.DefaultUnifiedAuthConfig())

	// 2. INICIALIZAR EL MÓDULO DE CONFIGURACIÓN ORGANIZACIONAL
	err := InitializeOrganizationConfigModule(db)
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize organization config module: %v", err))
	}

	// 3. REGISTRAR LAS RUTAS DEL MÓDULO
	RegisterOrganizationConfigModule(app, db)

	// 4. EJEMPLO DE RUTAS QUE SE CREAN AUTOMÁTICAMENTE:
	fmt.Println("Las siguientes rutas están disponibles:")
	fmt.Println("GET    /api/v1/org/:slug/config/modules                    - Obtener módulos de la organización")
	fmt.Println("POST   /api/v1/org/:slug/config/modules/install           - Instalar un módulo")
	fmt.Println("PUT    /api/v1/org/:slug/config/modules/:module/enable    - Habilitar módulo")
	fmt.Println("PUT    /api/v1/org/:slug/config/modules/:module/disable   - Deshabilitar módulo")
	fmt.Println("GET    /api/v1/org/:slug/config/modules/:module           - Obtener detalles del módulo")
	fmt.Println("GET    /api/v1/modules/available                          - Obtener módulos disponibles")

	// 5. INICIAR EL SERVIDOR
	// app.Listen(":8080")
}

// =============================================================================
// EJEMPLO DE USO MANUAL DE FUNCIONES RBAC
// =============================================================================

// ExampleManualRBACUsage muestra cómo usar las funciones RBAC manualmente
func ExampleManualRBACUsage(userID, orgID string) {
	// 1. Verificar permisos específicos
	canInstall := HasModulePermission(userID, orgID, "org_modules", "install")
	fmt.Printf("User %s can install modules in org %s: %v\n", userID, orgID, canInstall)

	canRead := HasModulePermission(userID, orgID, "org_modules", "read")
	fmt.Printf("User %s can read modules in org %s: %v\n", userID, orgID, canRead)

	// 2. Obtener todos los permisos del usuario
	permissions, err := GetUserModulePermissions(userID, orgID)
	if err != nil {
		fmt.Printf("Error getting permissions: %v\n", err)
		return
	}

	fmt.Printf("User %s has the following module permissions:\n", userID)
	for _, perm := range permissions {
		fmt.Printf("  - %s\n", perm)
	}
}

// =============================================================================
// EJEMPLO DE CONFIGURACIÓN EN MAIN.GO
// =============================================================================

/*
func main() {
	// Configuración básica
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Conexión a base de datos
	db := database.DBconn // Tu conexión real

	// Configurar CORS
	app.Use(cors.New())

	// 1. INICIALIZAR MIDDLEWARE DE AUTENTICACIÓN GLOBAL
	middleware.InitGlobalAuth(middleware.DefaultUnifiedAuthConfig())

	// 2. INICIALIZAR MÓDULO DE CONFIGURACIÓN ORGANIZACIONAL
	err := organization_config.InitializeOrganizationConfigModule(db)
	if err != nil {
		log.Fatalf("Failed to initialize organization config module: %v", err)
	}

	// 3. CONFIGURAR RUTAS BASE DE AUTENTICACIÓN
	auth.SetupAuthRoutes(app.Group("/api/v1"), db)
	organization.SetupOrganizationRoutes(app.Group("/api/v1"), db)

	// 4. REGISTRAR EL MÓDULO DE CONFIGURACIÓN ORGANIZACIONAL
	organization_config.RegisterOrganizationConfigModule(app, db)

	// 5. OTROS MÓDULOS...
	// gastos.SetupExpenseRoutes(app.Group("/api/v1"), db)

	// Iniciar servidor
	log.Fatal(app.Listen(":8080"))
}
*/

// =============================================================================
// EJEMPLO DE LLAMADAS API
// =============================================================================

/*
EJEMPLOS DE LLAMADAS API (usando curl o Postman):

1. INSTALAR MÓDULO DE GASTOS PARA UNA ORGANIZACIÓN:
POST /api/v1/org/mi-empresa/config/modules/install
Authorization: Bearer <jwt-token>
Content-Type: application/json

{
    "module_name": "expenses"
}

2. OBTENER MÓDULOS CONFIGURADOS DE LA ORGANIZACIÓN:
GET /api/v1/org/mi-empresa/config/modules
Authorization: Bearer <jwt-token>

3. HABILITAR MÓDULO:
PUT /api/v1/org/mi-empresa/config/modules/expenses/enable
Authorization: Bearer <jwt-token>

4. OBTENER MÓDULOS DISPONIBLES:
GET /api/v1/modules/available
Authorization: Bearer <jwt-token>

5. OBTENER DETALLES DE UN MÓDULO:
GET /api/v1/org/mi-empresa/config/modules/expenses
Authorization: Bearer <jwt-token>

RESPUESTA ESPERADA:
{
    "success": true,
    "message": "Organization modules retrieved successfully",
    "data": {
        "modules": [
            {
                "id": 1,
                "module_name": "expenses",
                "is_enabled": true,
                "permissions": {
                    "expense": ["read", "create", "update", "delete"],
                    "category": ["read", "create"]
                },
                "configuration": {
                    "max_expense_amount": 10000,
                    "require_approval": true
                },
                "created_at": "2025-08-06T10:30:00Z",
                "updated_at": "2025-08-06T10:30:00Z"
            }
        ],
        "permissions": [
            "org_modules:read:organization",
            "org_modules:install:organization",
            "org_modules:enable:organization"
        ],
        "count": 1
    }
}
*/
