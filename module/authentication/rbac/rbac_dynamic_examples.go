package rbac

import (
	"fmt"
	"practicev2/module/authentication/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// =============================================================================
// EJEMPLOS DE USO DEL SISTEMA RBAC DINÁMICO
// =============================================================================

// DemoRBACDynamicUsage demuestra el uso completo del sistema RBAC dinámico
func DemoRBACDynamicUsage(db *gorm.DB) (*DynamicRBACEngine, error) {
	fmt.Println("🎯 Demo: Sistema RBAC Dinámico Completo")
	fmt.Println("=====================================")

	// 1. CONFIGURACIÓN INICIAL CON BUILDER
	fmt.Println("1️⃣ Configurando RBAC dinámico...")

	engine := NewRBACBuilder(db, nil).
		WithDebugMode(true).
		WithCaching(true, 5*60). // 5 minutos de cache
		WithAuditLog(true).
		WithAutoSeed(true). // Inicializar permisos automáticamente
		AddQuickModule("gastos", map[string][]string{
			"expenses":   {"create", "read", "update", "delete", "approve"},
			"categories": {"create", "read", "update", "delete"},
			"budgets":    {"create", "read", "update", "approve"},
			"reports":    {"create", "read", "export"},
		}).
		Build()

	fmt.Println("✅ Motor RBAC dinámico configurado")

	// 2. DEMOSTRAR PERMISOS EN BASE DE DATOS
	if err := demoPermissionsInDatabase(db); err != nil {
		return nil, fmt.Errorf("error demonstrating database permissions: %w", err)
	}

	// 3. DEMOSTRAR ROLES Y ASIGNACIONES
	if err := demoRolesAndAssignments(db, engine); err != nil {
		return nil, fmt.Errorf("error demonstrating roles: %w", err)
	}

	// 4. DEMOSTRAR MIDDLEWARE EN ACCIÓN
	demoMiddlewareProtection(engine)

	fmt.Println("🎉 Demo completado exitosamente!")
	return engine, nil
}

// =============================================================================
// DEMOSTRACIÓN DE PERMISOS EN BASE DE DATOS
// =============================================================================

func demoPermissionsInDatabase(db *gorm.DB) error {
	fmt.Println("\n2️⃣ Verificando permisos en base de datos...")

	// Contar permisos totales
	var totalPermissions int64
	if err := db.Model(&models.Permission{}).Count(&totalPermissions).Error; err != nil {
		return err
	}

	fmt.Printf("📊 Total de permisos en BD: %d\n", totalPermissions)

	// Mostrar algunos permisos de ejemplo
	var samplePermissions []models.Permission
	if err := db.Limit(5).Find(&samplePermissions).Error; err != nil {
		return err
	}

	fmt.Println("🔐 Ejemplos de permisos:")
	for _, perm := range samplePermissions {
		fmt.Printf("   - %s:%s:%s (ID: %d)\n",
			perm.Resource, perm.Action, perm.Scope, perm.ID)
	}

	// Verificar roles de sistema
	var systemRoles []models.Role
	if err := db.Where("is_system_role = ?", true).Find(&systemRoles).Error; err != nil {
		return err
	}

	fmt.Printf("👥 Roles de sistema disponibles: %d\n", len(systemRoles))
	for _, role := range systemRoles {
		fmt.Printf("   - %s (%s) - Nivel: %d\n",
			role.Name, role.DisplayName, role.HierarchyLevel)
	}

	return nil
}

// =============================================================================
// DEMOSTRACIÓN DE ROLES Y ASIGNACIONES
// =============================================================================

func demoRolesAndAssignments(db *gorm.DB, engine *DynamicRBACEngine) error {
	fmt.Println("\n3️⃣ Demostrando gestión de roles...")

	// Crear una organización de ejemplo
	org := &models.Organization{
		ID:          "demo-org-001",
		Name:        "Demo Company",
		Slug:        "demo-company",
		Type:        "company",
		Description: "Organización de demostración",
		IsActive:    true,
	}

	// Verificar si la organización ya existe
	var existingOrg models.Organization
	if err := db.Where("id = ?", org.ID).First(&existingOrg).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Crear la organización con roles por defecto
			seeder := NewRBACSeeder(db, DefaultRBACConfig())
			if err := seeder.CreateOrganizationWithDefaultRoles(org); err != nil {
				return fmt.Errorf("failed to create demo organization: %w", err)
			}
			fmt.Printf("🏢 Organización '%s' creada con roles por defecto\n", org.Name)
		} else {
			return err
		}
	} else {
		fmt.Printf("🏢 Usando organización existente: '%s'\n", existingOrg.Name)
		org = &existingOrg
	}

	// Mostrar roles de la organización
	var orgRoles []models.Role
	if err := db.Where("organization_id = ?", org.ID).Find(&orgRoles).Error; err != nil {
		return err
	}

	fmt.Printf("📋 Roles en la organización: %d\n", len(orgRoles))
	for _, role := range orgRoles {
		// Contar permisos del rol usando la tabla intermedia
		var permissionCount int64
		if err := db.Table("role_permission").Where("role_id = ?", role.ID).Count(&permissionCount).Error; err != nil {
			return err
		}
		fmt.Printf("   - %s: %d permisos\n", role.DisplayName, permissionCount)
	}

	return nil
}

// =============================================================================
// DEMOSTRACIÓN DE MIDDLEWARE DE PROTECCIÓN
// =============================================================================

func demoMiddlewareProtection(engine *DynamicRBACEngine) {
	fmt.Println("\n4️⃣ Configurando protección de rutas...")

	// Crear una aplicación Fiber de demostración
	app := fiber.New()

	// Rutas protegidas con diferentes niveles de acceso
	api := app.Group("/api/v1")

	// Rutas de gastos con diferentes permisos
	gastos := api.Group("/gastos")
	gastos.Get("/", engine.SimpleProtect("expenses:read:organization"), demoHandler("Listar gastos"))
	gastos.Post("/", engine.SimpleProtect("expenses:create:organization"), demoHandler("Crear gasto"))
	gastos.Put("/:id", engine.SimpleProtect("expenses:update:own"), demoHandler("Actualizar gasto"))
	gastos.Delete("/:id", engine.SimpleProtect("expenses:delete:own"), demoHandler("Eliminar gasto"))
	gastos.Post("/:id/approve", engine.SimpleProtect("expenses:approve:department"), demoHandler("Aprobar gasto"))

	// Rutas de administración
	admin := api.Group("/admin")
	admin.Get("/users", engine.SimpleProtect("users:read:organization"), demoHandler("Listar usuarios"))
	admin.Post("/users", engine.SimpleProtect("users:create:organization"), demoHandler("Crear usuario"))
	admin.Get("/roles", engine.SimpleProtect("roles:read:organization"), demoHandler("Listar roles"))

	// Rutas con múltiples permisos (ProtectAny)
	reports := api.Group("/reports")
	reports.Get("/expenses",
		engine.ProtectAny("reports:read:department", "expenses:read:organization"),
		demoHandler("Reporte de gastos"))

	fmt.Println("🛡️ Rutas protegidas configuradas:")
	fmt.Println("   GET    /api/v1/gastos           -> expenses:read:organization")
	fmt.Println("   POST   /api/v1/gastos           -> expenses:create:organization")
	fmt.Println("   PUT    /api/v1/gastos/:id       -> expenses:update:own")
	fmt.Println("   DELETE /api/v1/gastos/:id       -> expenses:delete:own")
	fmt.Println("   POST   /api/v1/gastos/:id/approve -> expenses:approve:department")
	fmt.Println("   GET    /api/v1/admin/users      -> users:read:organization")
	fmt.Println("   POST   /api/v1/admin/users      -> users:create:organization")
	fmt.Println("   GET    /api/v1/admin/roles      -> roles:read:organization")
	fmt.Println("   GET    /api/v1/reports/expenses -> reports:read:department OR expenses:read:organization")
}

// demoHandler crea un handler de demostración
func demoHandler(action string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message":  fmt.Sprintf("✅ Acceso autorizado para: %s", action),
			"user":     c.Locals("rbac_user"),
			"resource": c.Locals("rbac_resource"),
			"action":   c.Locals("rbac_action"),
			"scope":    c.Locals("rbac_scope"),
		})
	}
}

// =============================================================================
// EJEMPLO DE VERIFICACIÓN MANUAL DE PERMISOS
// =============================================================================

// DemoManualPermissionCheck demuestra verificación manual de permisos
func DemoManualPermissionCheck(engine *DynamicRBACEngine, userID, orgID string) {
	fmt.Println("\n🔍 Demo: Verificación manual de permisos")
	fmt.Println("========================================")

	// Simular diferentes verificaciones de permisos
	permissions := []string{
		"expenses:create:own",
		"expenses:read:department",
		"expenses:approve:organization",
		"users:create:organization",
		"reports:export:all",
	}

	fmt.Printf("Usuario: %s, Organización: %s\n", userID, orgID)
	fmt.Println("Verificando permisos:")

	for _, perm := range permissions {
		// Esta es una demostración conceptual
		fmt.Printf("   %s -> [Requiere implementación con contexto real]\n", perm)
	}
}

// =============================================================================
// EJEMPLO DE GESTIÓN DINÁMICA DE PERMISOS
// =============================================================================

// DemoDynamicPermissionManagement demuestra gestión dinámica de permisos
func DemoDynamicPermissionManagement(db *gorm.DB, engine *DynamicRBACEngine) error {
	fmt.Println("\n⚡ Demo: Gestión dinámica de permisos")
	fmt.Println("====================================")

	// 1. Agregar nuevos permisos en tiempo de ejecución
	fmt.Println("1. Agregando nuevos permisos...")

	newPermissions := []models.Permission{
		{Resource: "invoices", Action: "create", Scope: "own"},
		{Resource: "invoices", Action: "read", Scope: "department"},
		{Resource: "invoices", Action: "approve", Scope: "organization"},
		{Resource: "contracts", Action: "create", Scope: "department"},
		{Resource: "contracts", Action: "sign", Scope: "organization"},
	}

	for _, perm := range newPermissions {
		var existing models.Permission
		result := db.Where("resource = ? AND action = ? AND scope = ?",
			perm.Resource, perm.Action, perm.Scope).FirstOrCreate(&existing, perm)

		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected > 0 {
			fmt.Printf("   ✅ Creado: %s:%s:%s\n", perm.Resource, perm.Action, perm.Scope)
		} else {
			fmt.Printf("   ⏭️  Existe: %s:%s:%s\n", perm.Resource, perm.Action, perm.Scope)
		}
	}

	// 2. Refrescar cache del engine
	fmt.Println("2. Refrescando cache del sistema...")
	// Cache se limpia automáticamente en la próxima verificación

	// 3. Verificar que los nuevos permisos están disponibles
	fmt.Println("3. Verificando nuevos permisos...")
	var count int64
	if err := db.Model(&models.Permission{}).Where("resource IN ?",
		[]string{"invoices", "contracts"}).Count(&count).Error; err != nil {
		return err
	}

	fmt.Printf("   📊 Permisos de facturas/contratos: %d\n", count)

	return nil
}

// =============================================================================
// UTILIDADES DE TESTING
// =============================================================================

// CreateTestOrganizationWithUser crea una organización y usuario de prueba
func CreateTestOrganizationWithUser(db *gorm.DB, seeder *RBACSeeder) (*models.Organization, *models.Identity, error) {
	// Crear organización de prueba
	org := &models.Organization{
		ID:          "test-org-001",
		Name:        "Test Organization",
		Slug:        "test-organization",
		Type:        "company",
		Description: "Organización de prueba para RBAC",
		IsActive:    true,
	}

	if err := seeder.CreateOrganizationWithDefaultRoles(org); err != nil {
		return nil, nil, err
	}

	// Crear usuario de prueba
	user := &models.Identity{
		ID:            "test-user-001",
		Email:         "test@example.com",
		PasswordHash:  "hashed_password",
		FirstName:     "Test",
		LastName:      "User",
		EmailVerified: true,
	}

	if err := db.Create(user).Error; err != nil {
		return nil, nil, err
	}

	// Asignar rol de empleado al usuario
	var employeeRole models.Role
	if err := db.Where("name = ? AND organization_id = ?",
		"employee", org.ID).First(&employeeRole).Error; err != nil {
		return nil, nil, err
	}

	membership := &models.OrganizationalMembership{
		ID:             "test-membership-001",
		IdentityID:     user.ID,
		OrganizationID: org.ID,
		RoleID:         employeeRole.ID,
		IsActive:       true,
	}

	if err := db.Create(membership).Error; err != nil {
		return nil, nil, err
	}

	return org, user, nil
}
