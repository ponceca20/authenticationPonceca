package organization_config

import (
	"encoding/json"
	"errors"
	"fmt"

	"practicev2/module/authentication/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

//=============================================================================
// SERVICIO DE CONFIGURACIÓN DINÁMICA DE MÓDULOS ORGANIZACIONALES
//=============================================================================

// OrganizationConfigService gestiona la configuración de módulos por organización
type OrganizationConfigService struct {
	db   *gorm.DB
	repo OrganizationConfigRepository
}

// ModuleConfiguration representa la configuración de un módulo
type ModuleConfiguration struct {
	ModuleName   string                 `json:"module_name"`
	IsEnabled    bool                   `json:"is_enabled"`
	Resources    map[string][]string    `json:"resources"` // {"expenses": ["create", "read"]}
	DefaultRoles map[string]RoleConfig  `json:"default_roles"`
	Settings     map[string]interface{} `json:"settings"`
	Version      string                 `json:"version"`
}

// RoleConfig representa la configuración de un rol
type RoleConfig struct {
	DisplayName    string   `json:"display_name"`
	HierarchyLevel int      `json:"hierarchy_level"`
	Permissions    []string `json:"permissions"` // ["expenses:create:own", "expenses:read:department"]
	Description    string   `json:"description"`
}

// NewOrganizationConfigService crea una nueva instancia del servicio
func NewOrganizationConfigService(db *gorm.DB) *OrganizationConfigService {
	return &OrganizationConfigService{
		db:   db,
		repo: NewOrganizationConfigRepository(db),
	}
}

//=============================================================================
// MÉTODOS PRINCIPALES DEL SERVICIO
//=============================================================================

// ConfigureModulesForOrganization configura múltiples módulos para una organización
func (s *OrganizationConfigService) ConfigureModulesForOrganization(
	orgID string,
	moduleConfigs []ModuleConfiguration,
	installerID string,
) error {
	// Verificar permisos del usuario para configurar módulos
	if !HasModulePermission(installerID, orgID, "org_modules", "configure") {
		return errors.New("insufficient permissions to configure modules")
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		for _, config := range moduleConfigs {
			if err := s.configureModule(tx, orgID, config, installerID); err != nil {
				return fmt.Errorf("error configurando módulo %s: %w", config.ModuleName, err)
			}
		}
		return nil
	})
}

// ConfigureModule configura un módulo específico para una organización
func (s *OrganizationConfigService) ConfigureModule(
	orgID string,
	config ModuleConfiguration,
	installerID string,
) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		return s.configureModule(tx, orgID, config, installerID)
	})
}

// GetOrganizationModules obtiene todos los módulos configurados de una organización
func (s *OrganizationConfigService) GetOrganizationModules(orgID string) ([]models.OrganizationModuleConfig, error) {
	return s.repo.GetByOrganizationID(orgID)
}

// GetEnabledModules obtiene solo los módulos habilitados de una organización
func (s *OrganizationConfigService) GetEnabledModules(orgID string) ([]models.OrganizationModuleConfig, error) {
	return s.repo.GetEnabledByOrganizationID(orgID)
}

// EnableModule habilita un módulo para una organización
func (s *OrganizationConfigService) EnableModule(orgID, moduleName string) error {
	return s.repo.UpdateModuleStatus(orgID, moduleName, true)
}

// DisableModule deshabilita un módulo para una organización
func (s *OrganizationConfigService) DisableModule(orgID, moduleName string) error {
	return s.repo.UpdateModuleStatus(orgID, moduleName, false)
}

//=============================================================================
// MÉTODOS ESPECÍFICOS POR MÓDULO
//=============================================================================

// SetupExpensesModule configura el módulo de gastos para una organización
func (s *OrganizationConfigService) SetupExpensesModule(orgID, installerID string) error {
	config := ModuleConfiguration{
		ModuleName: "expenses",
		IsEnabled:  true,
		Resources: map[string][]string{
			"expenses":   {"create", "read", "update", "delete", "approve", "reject"},
			"budgets":    {"create", "read", "update", "approve", "monitor", "report"},
			"categories": {"create", "read", "update", "delete", "assign"},
			"approvals":  {"view", "approve", "reject", "delegate"},
		},
		DefaultRoles: map[string]RoleConfig{
			"expense_manager": {
				DisplayName:    "Gerente de Gastos",
				HierarchyLevel: 80,
				Permissions:    []string{"expenses:*:organization", "budgets:*:organization", "categories:*:organization"},
				Description:    "Control total sobre gastos y presupuestos organizacionales",
			},
			"expense_approver": {
				DisplayName:    "Aprobador de Gastos",
				HierarchyLevel: 70,
				Permissions:    []string{"expenses:approve:department", "expenses:read:organization", "budgets:read:organization"},
				Description:    "Puede aprobar gastos departamentales",
			},
			"expense_creator": {
				DisplayName:    "Creador de Gastos",
				HierarchyLevel: 40,
				Permissions:    []string{"expenses:create:own", "expenses:read:own", "budgets:read:department"},
				Description:    "Puede crear y ver sus propios gastos",
			},
		},
		Settings: map[string]interface{}{
			"approval_workflow":      "manager_approval",
			"max_amount_no_approval": 1000,
			"currency":               "USD",
			"require_receipts":       true,
		},
		Version: "1.0.0",
	}

	return s.ConfigureModule(orgID, config, installerID)
}

// SetupInventoryModule configura el módulo de inventario para una organización
func (s *OrganizationConfigService) SetupInventoryModule(orgID, installerID string) error {
	config := ModuleConfiguration{
		ModuleName: "inventory",
		IsEnabled:  true,
		Resources: map[string][]string{
			"inventory":  {"create", "read", "update", "delete", "transfer", "adjust", "audit"},
			"warehouses": {"create", "read", "update", "delete", "manage", "assign"},
			"stock":      {"read", "update", "reserve", "release", "count", "audit"},
			"transfers":  {"create", "read", "approve", "cancel", "receive", "dispatch"},
			"purchases":  {"create", "read", "update", "approve", "receive", "cancel"},
			"suppliers":  {"create", "read", "update", "delete", "evaluate", "block"},
		},
		DefaultRoles: map[string]RoleConfig{
			"warehouse_manager": {
				DisplayName:    "Gerente de Almacén",
				HierarchyLevel: 80,
				Permissions:    []string{"inventory:*:organization", "warehouses:*:organization", "transfers:*:organization"},
				Description:    "Control total sobre inventario y almacenes",
			},
			"inventory_supervisor": {
				DisplayName:    "Supervisor de Inventario",
				HierarchyLevel: 60,
				Permissions:    []string{"inventory:read:organization", "inventory:update:department", "stock:*:department"},
				Description:    "Supervisa inventario departamental",
			},
			"stock_operator": {
				DisplayName:    "Operador de Stock",
				HierarchyLevel: 40,
				Permissions:    []string{"inventory:read:department", "inventory:update:own", "stock:count:department"},
				Description:    "Operaciones básicas de stock",
			},
		},
		Settings: map[string]interface{}{
			"multi_warehouse": true,
			"lot_tracking":    false,
			"serial_tracking": false,
			"auto_reorder":    true,
			"reorder_point":   10,
		},
		Version: "1.0.0",
	}

	return s.ConfigureModule(orgID, config, installerID)
}

//=============================================================================
// MÉTODOS PRIVADOS
//=============================================================================

// configureModule configura un módulo específico (método interno)
func (s *OrganizationConfigService) configureModule(
	tx *gorm.DB,
	orgID string,
	config ModuleConfiguration,
	installerID string,
) error {
	// 1. Serializar configuración JSON
	resourcesJSON, err := json.Marshal(config.Resources)
	if err != nil {
		return fmt.Errorf("error serializando recursos: %w", err)
	}

	defaultRolesJSON, err := json.Marshal(config.DefaultRoles)
	if err != nil {
		return fmt.Errorf("error serializando roles: %w", err)
	}

	settingsJSON, err := json.Marshal(config.Settings)
	if err != nil {
		return fmt.Errorf("error serializando configuración: %w", err)
	}

	// 2. Crear configuración del módulo
	moduleConfig := models.OrganizationModuleConfig{
		ID:             uuid.New().String(),
		OrganizationID: orgID,
		ModuleName:     config.ModuleName,
		IsEnabled:      config.IsEnabled,
		Resources:      string(resourcesJSON),
		DefaultRoles:   string(defaultRolesJSON),
		Settings:       string(settingsJSON),
		Version:        config.Version,
		InstalledBy:    installerID,
	}

	// 3. Upsert configuración del módulo
	if err := tx.Where("organization_id = ? AND module_name = ?", orgID, config.ModuleName).
		FirstOrCreate(&moduleConfig).Error; err != nil {
		return fmt.Errorf("error creando configuración de módulo: %w", err)
	}

	// 4. Crear permisos necesarios en la base de datos
	if err := s.createModulePermissions(tx, config.Resources); err != nil {
		return fmt.Errorf("error creando permisos: %w", err)
	}

	// 5. Crear roles predefinidos para la organización
	if err := s.createModuleRoles(tx, orgID, config.ModuleName, config.DefaultRoles, config.Resources); err != nil {
		return fmt.Errorf("error creando roles: %w", err)
	}

	return nil
}

// createModulePermissions crea todos los permisos necesarios para un módulo
func (s *OrganizationConfigService) createModulePermissions(
	tx *gorm.DB,
	resources map[string][]string,
) error {
	scopes := []string{"own", "department", "organization", "all"}

	for resource, actions := range resources {
		for _, action := range actions {
			for _, scope := range scopes {
				permission := models.Permission{
					Resource: resource,
					Action:   action,
					Scope:    scope,
				}

				// Crear permiso si no existe
				if err := tx.Where("resource = ? AND action = ? AND scope = ?",
					resource, action, scope).FirstOrCreate(&permission).Error; err != nil {
					return fmt.Errorf("error creando permiso %s:%s:%s: %w",
						resource, action, scope, err)
				}
			}
		}
	}

	return nil
}

// createModuleRoles crea roles predefinidos para un módulo
func (s *OrganizationConfigService) createModuleRoles(
	tx *gorm.DB,
	orgID, moduleName string,
	roleConfigs map[string]RoleConfig,
	resources map[string][]string,
) error {
	for roleName, roleConfig := range roleConfigs {
		// 1. Crear el rol
		role := models.Role{
			ID:             uuid.New().String(),
			OrganizationID: orgID,
			Name:           fmt.Sprintf("%s_%s", moduleName, roleName),
			DisplayName:    roleConfig.DisplayName,
			Description:    roleConfig.Description,
			HierarchyLevel: roleConfig.HierarchyLevel,
			IsSystemRole:   false,
		}

		if err := tx.Where("organization_id = ? AND name = ?", orgID, role.Name).
			FirstOrCreate(&role).Error; err != nil {
			return fmt.Errorf("error creando rol %s: %w", role.Name, err)
		}

		// 2. Asignar permisos al rol
		if err := s.assignPermissionsToRole(tx, role.ID, roleConfig.Permissions); err != nil {
			return fmt.Errorf("error asignando permisos al rol %s: %w", role.Name, err)
		}
	}

	return nil
}

// assignPermissionsToRole asigna permisos específicos a un rol
func (s *OrganizationConfigService) assignPermissionsToRole(
	tx *gorm.DB,
	roleID string,
	permissions []string,
) error {
	// Limpiar permisos existentes
	if err := tx.Exec("DELETE FROM role_permission WHERE role_id = ?", roleID).Error; err != nil {
		return fmt.Errorf("error limpiando permisos existentes: %w", err)
	}

	// Asignar nuevos permisos
	for _, permStr := range permissions {
		parts := parsePermissionString(permStr)
		if len(parts) != 3 {
			continue // Skip malformed permissions
		}

		// Buscar el permiso
		var permission models.Permission
		if err := tx.Where("resource = ? AND action = ? AND scope = ?",
			parts[0], parts[1], parts[2]).First(&permission).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				continue // Skip non-existent permissions
			}
			return fmt.Errorf("error buscando permiso: %w", err)
		}

		// Crear asociación role-permission
		if err := tx.Exec("INSERT INTO role_permission (role_id, permission_id) VALUES (?, ?) ON CONFLICT DO NOTHING",
			roleID, permission.ID).Error; err != nil {
			return fmt.Errorf("error asignando permiso: %w", err)
		}
	}

	return nil
}

//=============================================================================
// UTILIDADES
//=============================================================================

// parsePermissionString parsea un string de permiso como "expenses:create:own"
func parsePermissionString(permStr string) []string {
	// Expandir wildcards
	if permStr == "*:*:*" {
		return []string{"all", "all", "all"}
	}

	// Parsear formato resource:action:scope
	parts := make([]string, 3)
	tokens := splitPermission(permStr, ":")

	if len(tokens) >= 1 {
		parts[0] = tokens[0]
	}
	if len(tokens) >= 2 {
		parts[1] = tokens[1]
	}
	if len(tokens) >= 3 {
		parts[2] = tokens[2]
	}

	return parts
}

// splitPermission divide un string por delimitador
func splitPermission(str, delim string) []string {
	var result []string
	current := ""

	for _, char := range str {
		if string(char) == delim {
			result = append(result, current)
			current = ""
		} else {
			current += string(char)
		}
	}

	if current != "" {
		result = append(result, current)
	}

	return result
}

// GetAvailableModules retorna lista de módulos disponibles
func GetAvailableModules() map[string]ModuleConfiguration {
	return map[string]ModuleConfiguration{
		"expenses": {
			ModuleName: "expenses",
			Resources:  models.ModuleResources["expenses"],
			Version:    "1.0.0",
		},
		"inventory": {
			ModuleName: "inventory",
			Resources:  models.ModuleResources["inventory"],
			Version:    "1.0.0",
		},
		"ecommerce": {
			ModuleName: "ecommerce",
			Resources:  models.ModuleResources["ecommerce"],
			Version:    "1.0.0",
		},
		"education": {
			ModuleName: "education",
			Resources:  models.ModuleResources["education"],
			Version:    "1.0.0",
		},
		"reports": {
			ModuleName: "reports",
			Resources:  models.ModuleResources["reports"],
			Version:    "1.0.0",
		},
	}
}
