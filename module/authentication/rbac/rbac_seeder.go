package rbac

import (
	"fmt"
	"practicev2/module/authentication/models"
	"strings"

	"gorm.io/gorm"
)

// =============================================================================
// SISTEMA DE INICIALIZACIÓN DE PERMISOS DINÁMICOS
// =============================================================================

// RBACSeeder gestiona la inicialización de permisos y roles en base de datos
type RBACSeeder struct {
	db     *gorm.DB
	config *RBACConfig
}

// NewRBACSeeder crea un nuevo seeder de RBAC
func NewRBACSeeder(db *gorm.DB, config *RBACConfig) *RBACSeeder {
	return &RBACSeeder{
		db:     db,
		config: config,
	}
}

// =============================================================================
// INICIALIZACIÓN COMPLETA DEL SISTEMA RBAC
// =============================================================================

// InitializeRBACSystem inicializa todo el sistema RBAC con permisos y roles básicos
func (s *RBACSeeder) InitializeRBACSystem() error {
	fmt.Println("🔐 Inicializando sistema RBAC dinámico...")

	// 1. Crear permisos básicos del sistema
	if err := s.SeedSystemPermissions(); err != nil {
		return fmt.Errorf("failed to seed system permissions: %w", err)
	}

	// 2. Crear roles de sistema
	if err := s.SeedSystemRoles(); err != nil {
		return fmt.Errorf("failed to seed system roles: %w", err)
	}

	fmt.Println("✅ Sistema RBAC inicializado correctamente")
	return nil
}

// =============================================================================
// CREACIÓN DE PERMISOS DINÁMICOS
// =============================================================================

// SeedSystemPermissions crea todos los permisos básicos del sistema
func (s *RBACSeeder) SeedSystemPermissions() error {
	// Definir todos los recursos y acciones del sistema
	systemResources := map[string][]string{
		// Gestión de usuarios
		"users":    {"create", "read", "update", "delete", "invite", "suspend", "activate"},
		"profiles": {"read", "update"},

		// Organizaciones
		"organizations": {"create", "read", "update", "delete", "manage"},
		"departments":   {"create", "read", "update", "delete", "manage"},
		"roles":         {"create", "read", "update", "delete", "assign", "revoke"},
		"permissions":   {"read", "assign", "revoke"},

		// Membresías
		"memberships": {"create", "read", "update", "delete", "transfer"},
		"invitations": {"create", "read", "update", "delete", "send", "accept", "decline"},

		// Auditoría y seguridad
		"audit":    {"read", "export", "delete"},
		"security": {"read", "configure", "monitor"},

		// Configuración del sistema
		"settings":     {"read", "update", "system_config"},
		"integrations": {"read", "configure", "enable", "disable"},

		// E-commerce básico
		"customers": {"create", "read", "update", "delete", "communicate"},
		"products":  {"create", "read", "update", "delete", "price", "inventory"},
		"orders":    {"create", "read", "update", "process", "refund", "cancel"},

		// Financiero/Gastos
		"expenses":   {"create", "read", "update", "delete", "approve", "report"},
		"budgets":    {"create", "read", "update", "delete", "monitor"},
		"categories": {"create", "read", "update", "delete", "manage"},
		"reports":    {"create", "read", "export", "schedule"},
	}

	// Scopes disponibles
	scopes := []string{"own", "department", "organization", "all"}

	permissionCount := 0
	createdCount := 0

	for resource, actions := range systemResources {
		for _, action := range actions {
			for _, scope := range scopes {
				permissionCount++

				permission := models.Permission{
					Resource: resource,
					Action:   action,
					Scope:    scope,
				}

				// Usar FirstOrCreate para evitar duplicados
				var existingPermission models.Permission
				result := s.db.Where("resource = ? AND action = ? AND scope = ?",
					resource, action, scope).FirstOrCreate(&existingPermission, permission)

				if result.Error != nil {
					return fmt.Errorf("failed to create permission %s:%s:%s: %w",
						resource, action, scope, result.Error)
				}

				if result.RowsAffected > 0 {
					createdCount++
				}
			}
		}
	}

	if s.config.DebugMode {
		fmt.Printf("📊 Permisos: %d total, %d creados, %d existentes\n",
			permissionCount, createdCount, permissionCount-createdCount)
	}

	return nil
}

// =============================================================================
// CREACIÓN DE ROLES DE SISTEMA
// =============================================================================

// SeedSystemRoles crea los roles básicos del sistema
func (s *RBACSeeder) SeedSystemRoles() error {
	// Definir roles de sistema con sus permisos
	systemRoles := map[string]*SystemRoleDefinition{
		"super_admin": {
			Name:           "super_admin",
			DisplayName:    "Super Administrador",
			Description:    "Acceso completo al sistema",
			HierarchyLevel: 100,
			IsSystemRole:   true,
			Permissions:    []string{"*:*:all"}, // Todos los permisos
		},
		"system_admin": {
			Name:           "system_admin",
			DisplayName:    "Administrador del Sistema",
			Description:    "Administración general del sistema",
			HierarchyLevel: 90,
			IsSystemRole:   true,
			Permissions: []string{
				"users:*:organization", "organizations:*:all", "departments:*:all",
				"roles:*:organization", "permissions:*:organization", "settings:*:organization",
				"audit:read:organization", "security:*:organization",
			},
		},
		"org_admin": {
			Name:           "org_admin",
			DisplayName:    "Administrador de Organización",
			Description:    "Administrador completo de la organización",
			HierarchyLevel: 80,
			IsSystemRole:   true,
			Permissions: []string{
				"users:*:organization", "departments:*:organization", "roles:*:organization",
				"memberships:*:organization", "invitations:*:organization",
				"customers:*:organization", "products:*:organization", "orders:*:organization",
				"expenses:*:organization", "budgets:*:organization", "reports:*:organization",
			},
		},
		"manager": {
			Name:           "manager",
			DisplayName:    "Gerente",
			Description:    "Gestión departamental y de equipos",
			HierarchyLevel: 60,
			IsSystemRole:   true, // Cambiado a true para ser rol plantilla
			Permissions: []string{
				"users:read:department", "users:update:department",
				"departments:read:department", "memberships:read:department",
				"customers:*:department", "products:read:department", "orders:*:department",
				"expenses:*:department", "budgets:read:department", "reports:create:department",
			},
		},
		"employee": {
			Name:           "employee",
			DisplayName:    "Empleado",
			Description:    "Acceso básico para empleados",
			HierarchyLevel: 40,
			IsSystemRole:   true, // Cambiado a true para ser rol plantilla
			Permissions: []string{
				"users:read:own", "profiles:*:own",
				"customers:read:organization", "products:read:organization", "orders:read:organization",
				"expenses:create:own", "expenses:read:own", "expenses:update:own",
			},
		},
		"customer": {
			Name:           "customer",
			DisplayName:    "Cliente",
			Description:    "Acceso para clientes externos",
			HierarchyLevel: 20,
			IsSystemRole:   true, // Cambiado a true para ser rol plantilla
			Permissions: []string{
				"profiles:read:own", "profiles:update:own",
				"orders:create:own", "orders:read:own", "orders:update:own",
			},
		},
		"guest": {
			Name:           "guest",
			DisplayName:    "Invitado",
			Description:    "Acceso mínimo para invitados",
			HierarchyLevel: 10,
			IsSystemRole:   true, // Cambiado a true para ser rol plantilla
			Permissions: []string{
				"products:read:all",
			},
		},
	}

	createdRoles := 0

	for roleName, roleDef := range systemRoles {
		// Crear rol de sistema (sin organización específica)
		role := &models.Role{
			ID:             generateID(),
			Name:           roleDef.Name,
			DisplayName:    roleDef.DisplayName,
			Description:    roleDef.Description,
			HierarchyLevel: roleDef.HierarchyLevel,
			IsSystemRole:   roleDef.IsSystemRole,
			// OrganizationID se deja vacío para roles de sistema
		}

		// Buscar o crear el rol
		var existingRole models.Role
		result := s.db.Where("name = ? AND organization_id IS NULL AND is_system_role = ?",
			roleDef.Name, roleDef.IsSystemRole).FirstOrCreate(&existingRole, role)

		if result.Error != nil {
			return fmt.Errorf("failed to create system role %s: %w", roleName, result.Error)
		}

		if result.RowsAffected > 0 {
			createdRoles++
		}

		// Asignar permisos al rol
		if err := s.assignPermissionsToRole(&existingRole, roleDef.Permissions); err != nil {
			return fmt.Errorf("failed to assign permissions to role %s: %w", roleName, err)
		}
	}

	if s.config.DebugMode {
		fmt.Printf("👥 Roles de sistema: %d total, %d creados\n", len(systemRoles), createdRoles)
	}

	return nil
}

// SystemRoleDefinition define un rol de sistema
type SystemRoleDefinition struct {
	Name           string
	DisplayName    string
	Description    string
	HierarchyLevel int
	IsSystemRole   bool
	Permissions    []string
}

// =============================================================================
// ASIGNACIÓN DE PERMISOS A ROLES
// =============================================================================

// assignPermissionsToRole asigna permisos a un rol
func (s *RBACSeeder) assignPermissionsToRole(role *models.Role, permissionStrings []string) error {
	for _, permStr := range permissionStrings {
		// Manejar permiso wildcard
		if permStr == "*:*:all" {
			if err := s.assignAllPermissionsToRole(role); err != nil {
				return err
			}
			continue
		}

		// Expandir wildcards específicos
		permissions, err := s.expandPermissionString(permStr)
		if err != nil {
			return fmt.Errorf("failed to expand permission %s: %w", permStr, err)
		}

		// Asignar cada permiso expandido
		for _, permission := range permissions {
			if err := s.assignSinglePermissionToRole(role, permission); err != nil {
				return err
			}
		}
	}

	return nil
}

// expandPermissionString expande strings con wildcards a permisos específicos
func (s *RBACSeeder) expandPermissionString(permStr string) ([]models.Permission, error) {
	parts := strings.Split(permStr, ":")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid permission format: %s", permStr)
	}

	resource, action, scope := parts[0], parts[1], parts[2]
	var permissions []models.Permission

	// Obtener permisos que coincidan con el patrón
	query := s.db.Where("1=1")

	if resource != "*" {
		query = query.Where("resource = ?", resource)
	}
	if action != "*" {
		query = query.Where("action = ?", action)
	}
	if scope != "*" {
		query = query.Where("scope = ?", scope)
	}

	if err := query.Find(&permissions).Error; err != nil {
		return nil, err
	}

	return permissions, nil
}

// assignAllPermissionsToRole asigna todos los permisos del sistema a un rol
func (s *RBACSeeder) assignAllPermissionsToRole(role *models.Role) error {
	var allPermissions []models.Permission
	if err := s.db.Find(&allPermissions).Error; err != nil {
		return err
	}

	for _, permission := range allPermissions {
		if err := s.assignSinglePermissionToRole(role, permission); err != nil {
			return err
		}
	}

	return nil
}

// assignSinglePermissionToRole asigna un permiso específico a un rol
func (s *RBACSeeder) assignSinglePermissionToRole(role *models.Role, permission models.Permission) error {
	// Verificar si ya existe la asociación
	var count int64
	err := s.db.Table("role_permission").
		Where("role_id = ? AND permission_id = ?", role.ID, permission.ID).
		Count(&count).Error

	if err != nil {
		return err
	}

	// Si no existe, crear la asociación
	if count == 0 {
		err = s.db.Exec("INSERT INTO role_permission (role_id, permission_id) VALUES (?, ?)",
			role.ID, permission.ID).Error
		if err != nil {
			return err
		}
	}

	return nil
}

// =============================================================================
// UTILIDADES DE GESTIÓN
// =============================================================================

// GetSystemRoleByName obtiene un rol de sistema por nombre
func (s *RBACSeeder) GetSystemRoleByName(name string) (*models.Role, error) {
	return s.GetSystemRoleByNameWithDB(s.db, name)
}

// GetSystemRoleByNameWithDB obtiene un rol de sistema por nombre usando una instancia específica de DB
func (s *RBACSeeder) GetSystemRoleByNameWithDB(db *gorm.DB, name string) (*models.Role, error) {
	var role models.Role
	err := db.Where("name = ? AND is_system_role = true", name).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// CreateOrganizationWithDefaultRoles crea una organización con roles por defecto
func (s *RBACSeeder) CreateOrganizationWithDefaultRoles(org *models.Organization) error {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. Crear la organización
	if err := tx.Create(org).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 2. Crear roles específicos para la organización
	orgRoles := []string{"org_admin", "manager", "employee", "customer"}

	for _, roleName := range orgRoles {
		// Obtener definición del rol de sistema (usando s.db directamente, no la transacción)
		systemRole, err := s.GetSystemRoleByNameWithDB(s.db, roleName)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("system role %s not found: %w", roleName, err)
		}

		// Crear rol específico para la organización
		orgRole := &models.Role{
			ID:             generateID(),
			OrganizationID: org.ID,
			Name:           systemRole.Name,
			DisplayName:    systemRole.DisplayName,
			Description:    systemRole.Description,
			HierarchyLevel: systemRole.HierarchyLevel,
			IsSystemRole:   false,
		}

		if err := tx.Create(orgRole).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to create org role %s: %w", roleName, err)
		}

		// Copiar permisos del rol de sistema
		var systemPermissions []models.Permission
		if err := s.db.Model(systemRole).Association("Permissions").Find(&systemPermissions); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to get system role permissions: %w", err)
		}

		if err := tx.Model(orgRole).Association("Permissions").Append(systemPermissions); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to assign permissions to org role: %w", err)
		}
	}

	return tx.Commit().Error
}

// RefreshSystemPermissions actualiza permisos del sistema
func (s *RBACSeeder) RefreshSystemPermissions() error {
	fmt.Println("🔄 Actualizando permisos del sistema...")

	if err := s.SeedSystemPermissions(); err != nil {
		return err
	}

	fmt.Println("✅ Permisos del sistema actualizados")
	return nil
}
