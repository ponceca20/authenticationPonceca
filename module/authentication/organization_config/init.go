package organization_config

import (
	"practicev2/module/authentication/middleware"

	"gorm.io/gorm"
)

// =============================================================================
// INICIALIZACIÓN DEL MÓDULO DE CONFIGURACIÓN ORGANIZACIONAL
// =============================================================================

// InitializeOrganizationConfigModule inicializa el módulo con RBAC integrado
func InitializeOrganizationConfigModule(db *gorm.DB) error {
	// 1. Verificar que GlobalAuth esté inicializado
	if middleware.GlobalAuth == nil {
		middleware.InitGlobalAuth(middleware.DefaultUnifiedAuthConfig())
	}

	// 2. Registrar permisos específicos del módulo en RBAC
	err := registerModulePermissions()
	if err != nil {
		return err
	}

	// 3. Ejecutar migraciones si es necesario
	err = runModuleMigrations(db)
	if err != nil {
		return err
	}

	return nil
}

// registerModulePermissions registra los permisos específicos del módulo
func registerModulePermissions() error {
	if middleware.GlobalAuth == nil {
		return nil // No hay GlobalAuth inicializado
	}

	rbacEngine := middleware.GlobalAuth.GetRBACEngine()
	if rbacEngine == nil {
		return nil // No hay motor RBAC
	}

	// Recursos específicos para configuración de módulos organizacionales
	moduleResources := map[string][]string{
		"org_modules":    {"read", "configure", "install", "uninstall", "enable", "disable", "manage"},
		"module_config":  {"read", "create", "update", "delete"},
		"module_catalog": {"read", "browse"}, // Para catálogo de módulos disponibles
	}

	// Registrar usando QuickRegisterModule
	return rbacEngine.QuickRegisterModule("organization_config", moduleResources)
}

// runModuleMigrations ejecuta las migraciones del módulo si es necesario
func runModuleMigrations(db *gorm.DB) error {
	// Las migraciones ya se ejecutaron cuando se creó OrganizationModuleConfig
	// en all_models.go, pero aquí podríamos agregar migraciones específicas
	// del módulo si las necesitamos en el futuro

	return nil
}

// =============================================================================
// FUNCIONES DE UTILIDAD PARA RBAC
// =============================================================================

// HasModulePermission verifica si el usuario actual tiene permiso para una acción específica del módulo
func HasModulePermission(userID, orgID, resource, action string) bool {
	if middleware.GlobalAuth == nil {
		return false
	}

	rbacEngine := middleware.GlobalAuth.GetRBACEngine()
	if rbacEngine == nil {
		return false
	}

	// Usar defer/recover para manejar panics de la DB
	var result bool
	func() {
		defer func() {
			if r := recover(); r != nil {
				// Si hay un panic (como DB nil), devolver false
				result = false
			}
		}()

		// Usar el método CheckUserPermission del motor RBAC
		result = rbacEngine.CheckUserPermission(userID, orgID, resource, action, "organization")
	}()

	return result
}

// GetUserModulePermissions obtiene todos los permisos del usuario para este módulo
func GetUserModulePermissions(userID, orgID string) ([]string, error) {
	// Inicializar slice vacío por defecto
	userPermissions := make([]string, 0)

	if middleware.GlobalAuth == nil {
		return userPermissions, nil
	}

	rbacEngine := middleware.GlobalAuth.GetRBACEngine()
	if rbacEngine == nil {
		return userPermissions, nil
	}

	// Lista de permisos a verificar
	permissionsToCheck := []string{
		"org_modules:read:organization",
		"org_modules:configure:organization",
		"org_modules:install:organization",
		"org_modules:uninstall:organization",
		"org_modules:enable:organization",
		"org_modules:disable:organization",
		"org_modules:manage:organization",
		"module_config:read:organization",
		"module_config:create:organization",
		"module_config:update:organization",
		"module_config:delete:organization",
		"module_catalog:read:all",
	}

	for _, perm := range permissionsToCheck {
		// Parsear permiso (formato: resource:action:scope)
		parts := parsePermission(perm)
		if len(parts) == 3 {
			resource, action, scope := parts[0], parts[1], parts[2]

			// Usar defer/recover para manejar panics de la DB
			func() {
				defer func() {
					if r := recover(); r != nil {
						// Si hay un panic (como DB nil), simplemente continuar
						return
					}
				}()

				if rbacEngine.CheckUserPermission(userID, orgID, resource, action, scope) {
					userPermissions = append(userPermissions, perm)
				}
			}()
		}
	}

	return userPermissions, nil
} // parsePermission parsea un string de permiso en sus componentes
func parsePermission(permission string) []string {
	// Implementación simple para parsear "resource:action:scope"
	parts := make([]string, 0, 3)
	current := ""

	for _, char := range permission {
		if char == ':' {
			if current != "" {
				parts = append(parts, current)
				current = ""
			}
		} else {
			current += string(char)
		}
	}

	if current != "" {
		parts = append(parts, current)
	}

	return parts
}

// =============================================================================
// CONFIGURACIÓN DE PERMISOS POR ROL
// =============================================================================

// SetupDefaultRolePermissions configura permisos por defecto para roles comunes
func SetupDefaultRolePermissions(db *gorm.DB) error {
	// Esta función podría ser llamada durante la inicialización para configurar
	// permisos por defecto del módulo para roles como "admin", "manager", etc.

	// Ejemplo de implementación futura:
	// - Admin de organización: todos los permisos del módulo
	// - Manager: solo lectura y configuración básica
	// - Usuario regular: solo lectura

	return nil
}
