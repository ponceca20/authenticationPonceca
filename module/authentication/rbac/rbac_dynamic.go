package rbac

import (
	"fmt"
	"strings"
	"time"

	"practicev2/module/authentication/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// =============================================================================
// SISTEMA RBAC DINÁMICO - EVOLUCIÓN INTELIGENTE
// =============================================================================

// DynamicRBACEngine extiende el RBACEngine con capacidades dinámicas
type DynamicRBACEngine struct {
	*RBACEngine
	resourceManager *ResourceManager
	permissionSync  *PermissionSynchronizer
}

// ResourceManager gestiona recursos y permisos desde base de datos
type ResourceManager struct {
	db                *gorm.DB
	registeredModules map[string]*ModuleDefinition
	autoSync          bool
	syncInterval      time.Duration
}

// ModuleDefinition define un módulo y sus recursos
type ModuleDefinition struct {
	Name         string                       `json:"name"`
	Version      string                       `json:"version"`
	Resources    map[string]*ResourceMetadata `json:"resources"`
	Routes       []*RouteDefinition           `json:"routes"`
	RegisteredAt time.Time                    `json:"registered_at"`
}

// ResourceMetadata define metadata de un recurso
type ResourceMetadata struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Actions     []string          `json:"actions"`
	Scopes      []string          `json:"scopes"`
	Module      string            `json:"module"`
	Category    string            `json:"category,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// RouteDefinition define una ruta y sus permisos requeridos
type RouteDefinition struct {
	Method      string            `json:"method"`
	Path        string            `json:"path"`
	Resource    string            `json:"resource"`
	Action      string            `json:"action"`
	Scope       string            `json:"scope"`
	Description string            `json:"description"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// PermissionSynchronizer mantiene sincronizados permisos entre código y BD
type PermissionSynchronizer struct {
	db           *gorm.DB
	lastSync     time.Time
	syncInterval time.Duration
	autoSync     bool
}

// =============================================================================
// CONSTRUCTOR DEL SISTEMA DINÁMICO
// =============================================================================

// NewDynamicRBACEngine crea un nuevo motor RBAC dinámico
func NewDynamicRBACEngine(db *gorm.DB, config *RBACConfig) *DynamicRBACEngine {
	// Crear el engine base
	baseEngine := NewRBACEngine(db, config)

	// Crear el motor dinámico
	dynamicEngine := &DynamicRBACEngine{
		RBACEngine: baseEngine,
		resourceManager: &ResourceManager{
			db:                db,
			registeredModules: make(map[string]*ModuleDefinition),
			autoSync:          true,
			syncInterval:      5 * time.Minute,
		},
		permissionSync: &PermissionSynchronizer{
			db:           db,
			syncInterval: 2 * time.Minute,
			autoSync:     true,
		},
	}

	// Inicializar sincronización automática
	if dynamicEngine.permissionSync.autoSync {
		go dynamicEngine.startPermissionSync()
	}

	return dynamicEngine
}

// =============================================================================
// REGISTRO DINÁMICO DE MÓDULOS
// =============================================================================

// RegisterModule registra un módulo completo con sus recursos y rutas
func (e *DynamicRBACEngine) RegisterModule(moduleName string, definition *ModuleDefinition) error {
	definition.Name = moduleName
	definition.RegisteredAt = time.Now()

	// 1. Registrar en memoria
	e.resourceManager.registeredModules[moduleName] = definition

	// 2. Sincronizar con base de datos
	if err := e.syncModuleToDatabase(definition); err != nil {
		return fmt.Errorf("failed to sync module to database: %w", err)
	}

	// 3. Registrar rutas en el engine base
	for _, resource := range definition.Resources {
		baseResource := &ResourceDefinition{
			Name:        resource.Name,
			Module:      moduleName,
			Description: resource.Description,
			Actions:     resource.Actions,
			Scopes:      resource.Scopes,
			Routes:      e.convertRoutes(definition.Routes, resource.Name),
		}
		e.RBACEngine.RegisterResource(baseResource)
	}

	if e.config.DebugMode {
		fmt.Printf("[RBAC-Dynamic] Módulo '%s' registrado con %d recursos\n",
			moduleName, len(definition.Resources))
	}

	return nil
}

// QuickRegisterModule método simplificado para registro rápido
func (e *DynamicRBACEngine) QuickRegisterModule(moduleName string, resources map[string][]string) error {
	definition := &ModuleDefinition{
		Name:      moduleName,
		Version:   "1.0.0",
		Resources: make(map[string]*ResourceMetadata),
		Routes:    []*RouteDefinition{},
	}

	// Convertir recursos simples a ResourceMetadata
	for resourceName, actions := range resources {
		definition.Resources[resourceName] = &ResourceMetadata{
			Name:        resourceName,
			Description: fmt.Sprintf("%s management", resourceName),
			Actions:     actions,
			Scopes:      []string{"own", "department", "organization", "all"},
			Module:      moduleName,
		}
	}

	return e.RegisterModule(moduleName, definition)
}

// =============================================================================
// MIDDLEWARE DINÁMICO UNIFICADO
// =============================================================================

// Protect es el middleware unificado que combina autenticación + autorización dinámica
func (e *DynamicRBACEngine) Protect(resource, action, scope string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		startTime := time.Now()

		// 1. Autenticación JWT
		userCtx, err := e.getUserContext(c)
		if err != nil {
			return e.handleAuthError(c, "authentication_failed", err)
		}

		// 2. Sistema admin bypass
		if userCtx.IsSystemAdmin {
			c.Locals("rbac_user", userCtx)
			return c.Next()
		}

		// 3. Determinar organización activa
		orgID := e.determineActiveOrganization(c, userCtx)

		// 4. Verificación de permisos DINÁMICA
		allowed, err := e.checkDynamicPermission(userCtx, orgID, resource, action, scope)
		if err != nil {
			return e.handleAuthError(c, "authorization_failed", err)
		}

		if !allowed {
			if e.config.EnableAuditLog {
				go e.logAccessDenied(userCtx, c, resource, action, scope)
			}
			return e.handleAuthError(c, "insufficient_permissions",
				fmt.Errorf("access denied to %s:%s:%s", resource, action, scope))
		}

		// 5. Inyectar contexto
		c.Locals("rbac_user", userCtx)
		c.Locals("rbac_org_id", orgID)
		c.Locals("rbac_resource", resource)
		c.Locals("rbac_action", action)
		c.Locals("rbac_scope", scope)

		// 6. Auditoría de éxito
		if e.config.EnableAuditLog {
			go e.logAccessGranted(userCtx, c, resource, action, scope, time.Since(startTime))
		}

		return c.Next()
	}
}

// ProtectAny permite múltiples combinaciones de permisos
func (e *DynamicRBACEngine) ProtectAny(permissions ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userCtx, err := e.getUserContext(c)
		if err != nil {
			return e.handleAuthError(c, "authentication_failed", err)
		}

		if userCtx.IsSystemAdmin {
			c.Locals("rbac_user", userCtx)
			return c.Next()
		}

		orgID := e.determineActiveOrganization(c, userCtx)

		// Verificar si tiene alguno de los permisos especificados
		for _, permission := range permissions {
			parts := strings.Split(permission, ":")
			if len(parts) == 3 {
				if allowed, _ := e.checkDynamicPermission(userCtx, orgID, parts[0], parts[1], parts[2]); allowed {
					c.Locals("rbac_user", userCtx)
					c.Locals("rbac_org_id", orgID)
					return c.Next()
				}
			}
		}

		return e.handleAuthError(c, "insufficient_permissions",
			fmt.Errorf("none of the required permissions found"))
	}
}

// =============================================================================
// VERIFICACIÓN DINÁMICA DE PERMISOS
// =============================================================================

// checkDynamicPermission verifica permisos usando datos de base de datos
func (e *DynamicRBACEngine) checkDynamicPermission(userCtx *UserContext, orgID, resource, action, scope string) (bool, error) {
	// 1. Cache check primero
	cacheKey := fmt.Sprintf("dyn:%s:%s:%s:%s:%s", userCtx.Identity.ID, orgID, resource, action, scope)
	if e.config.CacheEnabled {
		if allowed, found := e.getFromCache(cacheKey); found {
			return allowed, nil
		}
	}

	// 2. Verificar permisos desde base de datos
	allowed := false

	// Si no hay contexto organizacional, verificar permisos globales
	if orgID == "" {
		allowed = e.checkGlobalPermissions(userCtx.Identity.ID, resource, action, scope)
	} else {
		allowed = e.checkOrganizationalPermissions(userCtx.Identity.ID, orgID, resource, action, scope)
	}

	// 3. Cache el resultado
	if e.config.CacheEnabled {
		e.setCache(cacheKey, allowed)
	}

	return allowed, nil
}

// checkOrganizationalPermissions verifica permisos en contexto organizacional
func (e *DynamicRBACEngine) checkOrganizationalPermissions(identityID, orgID, resource, action, scope string) bool {
	var count int64

	// Query optimizada que une todas las tablas necesarias con nombres correctos (singular)
	query := `
		SELECT COUNT(*) 
		FROM organizational_membership om
		JOIN role_permission rp ON om.role_id = rp.role_id
		JOIN permission p ON rp.permission_id = p.id
		WHERE om.identity_id = ? 
		  AND om.organization_id = ? 
		  AND om.is_active = true
		  AND p.resource = ? 
		  AND p.action = ?
		  AND (p.scope = ? OR p.scope = 'all')
		  AND om.deleted_at IS NULL
	`

	err := e.db.Raw(query, identityID, orgID, resource, action, scope).Count(&count).Error
	if err != nil {
		if e.config.DebugMode {
			fmt.Printf("[RBAC-Dynamic] Error checking organizational permissions: %v\n", err)
		}
		return false
	}

	return count > 0
}

// checkGlobalPermissions verifica permisos globales (sin contexto organizacional)
func (e *DynamicRBACEngine) checkGlobalPermissions(identityID, resource, action, scope string) bool {
	var count int64

	// Verificar si el usuario tiene permisos globales (roles de sistema) - tablas singulares
	query := `
		SELECT COUNT(*) 
		FROM organizational_membership om
		JOIN role_permission rp ON om.role_id = rp.role_id
		JOIN permission p ON rp.permission_id = p.id
		JOIN role r ON om.role_id = r.id
		WHERE om.identity_id = ? 
		  AND om.is_active = true
		  AND r.is_system_role = true
		  AND p.resource = ? 
		  AND p.action = ?
		  AND (p.scope = ? OR p.scope = 'all')
		  AND om.deleted_at IS NULL
	`

	err := e.db.Raw(query, identityID, resource, action, scope).Count(&count).Error
	if err != nil {
		if e.config.DebugMode {
			fmt.Printf("[RBAC-Dynamic] Error checking global permissions: %v\n", err)
		}
		return false
	}

	return count > 0
}

// =============================================================================
// SINCRONIZACIÓN CON BASE DE DATOS
// =============================================================================

// syncModuleToDatabase sincroniza un módulo con la base de datos
func (e *DynamicRBACEngine) syncModuleToDatabase(definition *ModuleDefinition) error {
	tx := e.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. Crear/actualizar permisos para cada recurso del módulo
	for _, resource := range definition.Resources {
		for _, action := range resource.Actions {
			for _, scope := range resource.Scopes {
				permission := models.Permission{
					Resource: resource.Name,
					Action:   action,
					Scope:    scope,
				}

				// Upsert del permiso
				if err := tx.Where("resource = ? AND action = ? AND scope = ?",
					resource.Name, action, scope).FirstOrCreate(&permission).Error; err != nil {
					tx.Rollback()
					return fmt.Errorf("failed to create permission %s:%s:%s: %w",
						resource.Name, action, scope, err)
				}
			}
		}
	}

	return tx.Commit().Error
}

// startPermissionSync inicia la sincronización automática de permisos
func (e *DynamicRBACEngine) startPermissionSync() {
	ticker := time.NewTicker(e.permissionSync.syncInterval)
	defer ticker.Stop()

	for range ticker.C {
		if err := e.syncPermissionsFromDatabase(); err != nil && e.config.DebugMode {
			fmt.Printf("[RBAC-Dynamic] Error syncing permissions: %v\n", err)
		}
	}
}

// syncPermissionsFromDatabase sincroniza permisos desde la base de datos
func (e *DynamicRBACEngine) syncPermissionsFromDatabase() error {
	// Limpiar cache cuando hay cambios en permisos
	e.permissionCache.mu.Lock()
	e.permissionCache.cache = make(map[string]*CacheEntry)
	e.permissionCache.mu.Unlock()

	e.permissionSync.lastSync = time.Now()

	if e.config.DebugMode {
		fmt.Printf("[RBAC-Dynamic] Permissions synced at %v\n", e.permissionSync.lastSync)
	}

	return nil
}

// =============================================================================
// MÉTODOS HELPER PARA USO SIMPLE
// =============================================================================

// SimpleProtect método simplificado para protección básica
func (e *DynamicRBACEngine) SimpleProtect(permission string) fiber.Handler {
	parts := strings.Split(permission, ":")
	if len(parts) != 3 {
		panic(fmt.Sprintf("Invalid permission format: %s. Expected: resource:action:scope", permission))
	}
	return e.Protect(parts[0], parts[1], parts[2])
}

// convertRoutes convierte RouteDefinition a RoutePermission para compatibilidad
func (e *DynamicRBACEngine) convertRoutes(routes []*RouteDefinition, resourceName string) []*RoutePermission {
	var converted []*RoutePermission
	for _, route := range routes {
		if route.Resource == resourceName {
			converted = append(converted, &RoutePermission{
				Path:        route.Path,
				Method:      route.Method,
				Resource:    route.Resource,
				Action:      route.Action,
				Scope:       route.Scope,
				Description: route.Description,
			})
		}
	}
	return converted
}

// =============================================================================
// API HELPERS PARA GESTIÓN DINÁMICA
// =============================================================================

// GetUserPermissions obtiene todos los permisos de un usuario (para APIs de gestión)
func (e *DynamicRBACEngine) GetUserPermissions(identityID string, orgID string) ([]models.Permission, error) {
	var permissions []models.Permission

	query := e.db.Table("permissions p").
		Select("DISTINCT p.*").
		Joins("JOIN role_permission rp ON p.id = rp.permission_id").
		Joins("JOIN organizational_membership om ON rp.role_id = om.role_id").
		Where("om.identity_id = ? AND om.is_active = true", identityID)

	if orgID != "" {
		query = query.Where("om.organization_id = ?", orgID)
	}

	err := query.Find(&permissions).Error
	return permissions, err
}

// RefreshUserCache refresca el cache de un usuario específico
func (e *DynamicRBACEngine) RefreshUserCache(identityID string) {
	e.InvalidateUserCache(identityID)
	if e.config.DebugMode {
		fmt.Printf("[RBAC-Dynamic] Cache refreshed for user: %s\n", identityID)
	}
}

// CheckUserPermission verifica si un usuario tiene un permiso específico
func (e *DynamicRBACEngine) CheckUserPermission(identityID, orgID, resource, action, scope string) bool {
	// Usar directamente el método de verificación organizacional
	if orgID != "" {
		return e.checkOrganizationalPermissions(identityID, orgID, resource, action, scope)
	}

	// Si no hay orgID, verificar permisos globales
	return e.checkGlobalPermissions(identityID, resource, action, scope)
}
