package rbac

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"practicev2/module/authentication/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// =============================================================================
// SISTEMA RBAC INTELIGENTE - MOTOR PRINCIPAL
// =============================================================================

// RBACEngine es el motor principal del sistema RBAC
type RBACEngine struct {
	db               *gorm.DB
	permissionCache  *PermissionCache
	organizationCtx  *OrganizationContext
	resourceRegistry *ResourceRegistry
	config           *RBACConfig
}

// RBACConfig contiene la configuración del sistema RBAC
type RBACConfig struct {
	CacheEnabled     bool          `json:"cache_enabled"`
	CacheTTL         time.Duration `json:"cache_ttl"`
	DefaultDenyAll   bool          `json:"default_deny_all"`
	EnableAuditLog   bool          `json:"enable_audit_log"`
	DebugMode        bool          `json:"debug_mode"`
	OrganizationMode bool          `json:"organization_mode"` // true para multi-tenant
}

// PermissionCache maneja el cache de permisos para performance
type PermissionCache struct {
	mu    sync.RWMutex
	cache map[string]*CacheEntry
	ttl   time.Duration
}

type CacheEntry struct {
	permissions map[string]bool
	expiresAt   time.Time
	identityID  string
	orgID       string
}

// OrganizationContext maneja el contexto organizacional
type OrganizationContext struct {
	mu              sync.RWMutex
	activeOrgs      map[string]*models.Organization
	userMemberships map[string][]models.OrganizationalMembership
	lastSync        time.Time
	syncInterval    time.Duration
}

// ResourceRegistry mantiene registro de todos los recursos del sistema
type ResourceRegistry struct {
	mu        sync.RWMutex
	resources map[string]*ResourceDefinition
	routes    map[string]*RoutePermission
}

// ResourceDefinition define un recurso del sistema
type ResourceDefinition struct {
	Name        string             `json:"name"`
	Module      string             `json:"module"`
	Description string             `json:"description"`
	Actions     []string           `json:"actions"`
	Scopes      []string           `json:"scopes"`
	Routes      []*RoutePermission `json:"routes"`
	CreatedAt   time.Time          `json:"created_at"`
}

// RoutePermission define los permisos requeridos para una ruta
type RoutePermission struct {
	Path         string   `json:"path"`
	Method       string   `json:"method"`
	Resource     string   `json:"resource"`
	Action       string   `json:"action"`
	Scope        string   `json:"scope"`
	RequiredRole []string `json:"required_role,omitempty"` // roles específicos si se necesita
	Module       string   `json:"module"`
	Description  string   `json:"description"`
}

// UserContext contiene el contexto completo del usuario autenticado
type UserContext struct {
	Identity        *models.Identity                      `json:"identity"`
	ActiveOrg       *models.Organization                  `json:"active_org,omitempty"`
	Memberships     []models.OrganizationalMembership     `json:"memberships"`
	Permissions     map[string]map[string]map[string]bool `json:"permissions"` // resource -> action -> scope -> allowed
	CustomerProfile *models.CustomerProfile               `json:"customer_profile,omitempty"`
	IsSystemAdmin   bool                                  `json:"is_system_admin"`
	AccessLevel     string                                `json:"access_level"`
}

// =============================================================================
// CONSTRUCTOR Y CONFIGURACIÓN
// =============================================================================

// NewRBACEngine crea una nueva instancia del motor RBAC
func NewRBACEngine(db *gorm.DB, config *RBACConfig) *RBACEngine {
	if config == nil {
		config = DefaultRBACConfig()
	}

	engine := &RBACEngine{
		db:     db,
		config: config,
		permissionCache: &PermissionCache{
			cache: make(map[string]*CacheEntry),
			ttl:   config.CacheTTL,
		},
		organizationCtx: &OrganizationContext{
			activeOrgs:      make(map[string]*models.Organization),
			userMemberships: make(map[string][]models.OrganizationalMembership),
			syncInterval:    15 * time.Minute,
		},
		resourceRegistry: &ResourceRegistry{
			resources: make(map[string]*ResourceDefinition),
			routes:    make(map[string]*RoutePermission),
		},
	}

	// Cargar recursos del sistema
	engine.loadSystemResources()

	// Iniciar sincronización en background
	go engine.startBackgroundSync()

	return engine
}

// DefaultRBACConfig retorna la configuración por defecto
func DefaultRBACConfig() *RBACConfig {
	return &RBACConfig{
		CacheEnabled:     true,
		CacheTTL:         10 * time.Minute,
		DefaultDenyAll:   true,
		EnableAuditLog:   true,
		DebugMode:        false,
		OrganizationMode: true,
	}
}

// =============================================================================
// REGISTRO DE RECURSOS
// =============================================================================

// RegisterResource registra un nuevo recurso en el sistema
func (e *RBACEngine) RegisterResource(resource *ResourceDefinition) error {
	e.resourceRegistry.mu.Lock()
	defer e.resourceRegistry.mu.Unlock()

	resource.CreatedAt = time.Now()
	e.resourceRegistry.resources[resource.Name] = resource

	// Registrar las rutas asociadas
	for _, route := range resource.Routes {
		routeKey := fmt.Sprintf("%s:%s", route.Method, route.Path)
		route.Resource = resource.Name
		route.Module = resource.Module
		e.resourceRegistry.routes[routeKey] = route
	}

	if e.config.DebugMode {
		fmt.Printf("[RBAC] Recurso registrado: %s (módulo: %s) con %d rutas\n",
			resource.Name, resource.Module, len(resource.Routes))
	}

	return nil
}

// RegisterModuleResources registra múltiples recursos de un módulo
func (e *RBACEngine) RegisterModuleResources(module string, resources []*ResourceDefinition) error {
	for _, resource := range resources {
		resource.Module = module
		if err := e.RegisterResource(resource); err != nil {
			return fmt.Errorf("error registrando recurso %s del módulo %s: %w", resource.Name, module, err)
		}
	}

	if e.config.DebugMode {
		fmt.Printf("[RBAC] Módulo '%s' registrado con %d recursos\n", module, len(resources))
	}

	return nil
}

// =============================================================================
// MIDDLEWARE PRINCIPAL - LA MAGIA SUCEDE AQUÍ
// =============================================================================

// ProtectRoute es el middleware principal que protege las rutas
func (e *RBACEngine) ProtectRoute(resource, action, scope string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		startTime := time.Now()

		// 1. Obtener contexto del usuario
		userCtx, err := e.getUserContext(c)
		if err != nil {
			return e.handleAuthError(c, "invalid_token", err)
		}

		// 2. Verificar si es acceso de sistema (bypass para admins globales)
		if userCtx.IsSystemAdmin {
			if e.config.DebugMode {
				fmt.Printf("[RBAC] System admin bypass: %s\n", userCtx.Identity.Email)
			}
			c.Locals("rbac_user", userCtx)
			return c.Next()
		}

		// 3. Determinar organización activa
		orgID := e.determineActiveOrganization(c, userCtx)

		// 4. Verificar permisos (con cache)
		allowed, err := e.checkPermission(userCtx, orgID, resource, action, scope)
		if err != nil {
			return e.handleAuthError(c, "permission_check_failed", err)
		}

		if !allowed {
			// Log del intento de acceso denegado
			if e.config.EnableAuditLog {
				go e.logAccessDenied(userCtx, c, resource, action, scope)
			}
			return e.handleAuthError(c, "insufficient_permissions",
				fmt.Errorf("access denied to %s:%s:%s", resource, action, scope))
		}

		// 5. Agregar contexto a la request
		c.Locals("rbac_user", userCtx)
		c.Locals("rbac_org_id", orgID)
		c.Locals("rbac_permissions", userCtx.Permissions)

		// 6. Log de acceso exitoso (opcional)
		if e.config.EnableAuditLog {
			go e.logAccessGranted(userCtx, c, resource, action, scope, time.Since(startTime))
		}

		return c.Next()
	}
}

// RequirePermission es un middleware más simple para casos específicos
func (e *RBACEngine) RequirePermission(permission string) fiber.Handler {
	parts := strings.Split(permission, ":")
	if len(parts) != 3 {
		panic(fmt.Sprintf("Invalid permission format: %s. Expected format: resource:action:scope", permission))
	}

	return e.ProtectRoute(parts[0], parts[1], parts[2])
}

// RequireRole middleware que requiere un rol específico
func (e *RBACEngine) RequireRole(roleName string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userCtx, err := e.getUserContext(c)
		if err != nil {
			return e.handleAuthError(c, "invalid_token", err)
		}

		// Verificar si el usuario tiene el rol en alguna organización
		hasRole := false
		for _, membership := range userCtx.Memberships {
			if membership.Role.Name == roleName && membership.IsActive {
				hasRole = true
				break
			}
		}

		if !hasRole {
			return e.handleAuthError(c, "insufficient_role",
				fmt.Errorf("role '%s' required", roleName))
		}

		c.Locals("rbac_user", userCtx)
		return c.Next()
	}
}

// RequireOrganization middleware que requiere contexto organizacional
func (e *RBACEngine) RequireOrganization() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userCtx, err := e.getUserContext(c)
		if err != nil {
			return e.handleAuthError(c, "invalid_token", err)
		}

		orgID := e.determineActiveOrganization(c, userCtx)
		if orgID == "" {
			return e.handleAuthError(c, "no_organization",
				fmt.Errorf("organization context required"))
		}

		c.Locals("rbac_user", userCtx)
		c.Locals("rbac_org_id", orgID)
		return c.Next()
	}
}

// =============================================================================
// VERIFICACIÓN DE PERMISOS CON CACHE
// =============================================================================

// checkPermission verifica si un usuario tiene un permiso específico
func (e *RBACEngine) checkPermission(userCtx *UserContext, orgID, resource, action, scope string) (bool, error) {
	// Generar clave de cache
	cacheKey := fmt.Sprintf("%s:%s:%s:%s:%s", userCtx.Identity.ID, orgID, resource, action, scope)

	// Verificar cache primero
	if e.config.CacheEnabled {
		if allowed, found := e.getFromCache(cacheKey); found {
			return allowed, nil
		}
	}

	// Verificar permisos desde la base de datos
	allowed := false

	// 1. Verificar permisos directos en el contexto del usuario
	if resourcePerms, exists := userCtx.Permissions[resource]; exists {
		if actionPerms, exists := resourcePerms[action]; exists {
			if scopeAllowed, exists := actionPerms[scope]; exists && scopeAllowed {
				allowed = true
			}
			// También verificar si tiene permiso 'all' scope
			if allScopeAllowed, exists := actionPerms["all"]; exists && allScopeAllowed {
				allowed = true
			}
		}
	}

	// 2. Si no tiene permisos directos, verificar por rol y organización
	if !allowed && orgID != "" {
		allowed = e.checkRoleBasedPermission(userCtx, orgID, resource, action, scope)
	}

	// 3. Cachear el resultado
	if e.config.CacheEnabled {
		e.setCache(cacheKey, allowed)
	}

	return allowed, nil
}

// checkRoleBasedPermission verifica permisos basados en roles
func (e *RBACEngine) checkRoleBasedPermission(userCtx *UserContext, orgID, resource, action, scope string) bool {
	// Buscar membresía activa en la organización
	var activeMembership *models.OrganizationalMembership
	for _, membership := range userCtx.Memberships {
		if membership.OrganizationID == orgID && membership.IsActive {
			activeMembership = &membership
			break
		}
	}

	if activeMembership == nil {
		return false
	}

	// Verificar permisos del rol
	for _, permission := range activeMembership.Role.Permissions {
		if permission.Resource == resource &&
			permission.Action == action &&
			(permission.Scope == scope || permission.Scope == "all") {
			return true
		}
	}

	return false
}

// =============================================================================
// GESTIÓN DE CACHE
// =============================================================================

func (e *RBACEngine) getFromCache(key string) (bool, bool) {
	e.permissionCache.mu.RLock()
	defer e.permissionCache.mu.RUnlock()

	entry, exists := e.permissionCache.cache[key]
	if !exists || time.Now().After(entry.expiresAt) {
		return false, false
	}

	// La clave en cache contiene resource:action:scope
	parts := strings.Split(key, ":")
	if len(parts) >= 3 {
		permKey := strings.Join(parts[len(parts)-3:], ":")
		if allowed, exists := entry.permissions[permKey]; exists {
			return allowed, true
		}
	}

	return false, false
}

func (e *RBACEngine) setCache(key string, allowed bool) {
	e.permissionCache.mu.Lock()
	defer e.permissionCache.mu.Unlock()

	parts := strings.Split(key, ":")
	if len(parts) < 5 {
		return
	}

	identityID := parts[0]
	orgID := parts[1]
	permKey := strings.Join(parts[2:], ":")

	cacheKey := fmt.Sprintf("%s:%s", identityID, orgID)

	entry, exists := e.permissionCache.cache[cacheKey]
	if !exists {
		entry = &CacheEntry{
			permissions: make(map[string]bool),
			expiresAt:   time.Now().Add(e.permissionCache.ttl),
			identityID:  identityID,
			orgID:       orgID,
		}
		e.permissionCache.cache[cacheKey] = entry
	}

	entry.permissions[permKey] = allowed
}

// InvalidateUserCache invalida el cache de un usuario específico
func (e *RBACEngine) InvalidateUserCache(identityID string) {
	e.permissionCache.mu.Lock()
	defer e.permissionCache.mu.Unlock()

	for key := range e.permissionCache.cache {
		if strings.HasPrefix(key, identityID+":") {
			delete(e.permissionCache.cache, key)
		}
	}
}

// =============================================================================
// UTILIDADES Y HELPERS
// =============================================================================

// GetUserFromContext obtiene el contexto del usuario desde Fiber context
func GetUserFromContext(c *fiber.Ctx) (*UserContext, error) {
	user := c.Locals("rbac_user")
	if user == nil {
		return nil, fmt.Errorf("user context not found")
	}

	userCtx, ok := user.(*UserContext)
	if !ok {
		return nil, fmt.Errorf("invalid user context type")
	}

	return userCtx, nil
}

// GetOrgIDFromContext obtiene el ID de organización del contexto
func GetOrgIDFromContext(c *fiber.Ctx) string {
	orgID := c.Locals("rbac_org_id")
	if orgID == nil {
		return ""
	}

	if id, ok := orgID.(string); ok {
		return id
	}

	return ""
}

// HasPermission verifica si el usuario actual tiene un permiso específico
func (e *RBACEngine) HasPermission(c *fiber.Ctx, resource, action, scope string) bool {
	userCtx, err := GetUserFromContext(c)
	if err != nil {
		return false
	}

	orgID := GetOrgIDFromContext(c)
	allowed, err := e.checkPermission(userCtx, orgID, resource, action, scope)
	if err != nil {
		return false
	}

	return allowed
}

// loadSystemResources carga los recursos del sistema desde la base de datos
func (e *RBACEngine) loadSystemResources() {
	// Cargar recursos base del sistema de autenticación
	authResources := []*ResourceDefinition{
		{
			Name:        "users",
			Module:      "authentication",
			Description: "User management",
			Actions:     []string{"create", "read", "update", "delete", "invite", "suspend"},
			Scopes:      []string{"own", "department", "organization", "all"},
		},
		{
			Name:        "organizations",
			Module:      "authentication",
			Description: "Organization management",
			Actions:     []string{"create", "read", "update", "delete", "manage"},
			Scopes:      []string{"own", "all"},
		},
		{
			Name:        "roles",
			Module:      "authentication",
			Description: "Role and permission management",
			Actions:     []string{"create", "read", "update", "delete", "assign"},
			Scopes:      []string{"organization", "all"},
		},
	}

	for _, resource := range authResources {
		e.RegisterResource(resource)
	}
}

// startBackgroundSync inicia la sincronización en background
func (e *RBACEngine) startBackgroundSync() {
	ticker := time.NewTicker(e.organizationCtx.syncInterval)
	defer ticker.Stop()

	for range ticker.C {
		e.syncOrganizationData()
	}
}

// syncOrganizationData sincroniza datos organizacionales
func (e *RBACEngine) syncOrganizationData() {
	if e.config.DebugMode {
		fmt.Println("[RBAC] Sincronizando datos organizacionales...")
	}

	// Limpiar cache expirado
	e.cleanExpiredCache()

	// Actualizar timestamp
	e.organizationCtx.mu.Lock()
	e.organizationCtx.lastSync = time.Now()
	e.organizationCtx.mu.Unlock()
}

// cleanExpiredCache limpia entradas expiradas del cache
func (e *RBACEngine) cleanExpiredCache() {
	e.permissionCache.mu.Lock()
	defer e.permissionCache.mu.Unlock()

	now := time.Now()
	for key, entry := range e.permissionCache.cache {
		if now.After(entry.expiresAt) {
			delete(e.permissionCache.cache, key)
		}
	}
}
