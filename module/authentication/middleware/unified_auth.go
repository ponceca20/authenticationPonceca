package middleware

import (
	"fmt"
	"strings"
	"time"

	"practicev2/database"
	"practicev2/module/authentication/auth"
	"practicev2/module/authentication/models"
	"practicev2/module/authentication/rbac"
	"practicev2/module/authentication/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// =============================================================================
// MIDDLEWARE UNIFICADO DINÁMICO - INTEGRADO CON SISTEMA EXISTENTE
// =============================================================================
// Este middleware combina nuestro SmartAuthMiddleware existente con autorización
// RBAC dinámica en una sola línea, manteniendo compatibilidad total.

// UnifiedAuthConfig configuración del middleware unificado
type UnifiedAuthConfig struct {
	CacheEnabled   bool
	CacheTTL       time.Duration
	DebugMode      bool
	DefaultDenyAll bool
	RequireOrg     bool
	AllowAnonymous bool
	EnableAudit    bool
}

// UnifiedAuthMiddleware middleware principal que combina Auth + Authz
type UnifiedAuthMiddleware struct {
	db              *gorm.DB
	jwtService      *utils.JWTService
	authRepo        auth.AuthRepository
	rbacEngine      *rbac.DynamicRBACEngine
	config          *UnifiedAuthConfig
	permissionCache map[string]*PermissionCacheEntry
}

// PermissionCacheEntry entrada del cache de permisos
type PermissionCacheEntry struct {
	Allowed   bool
	ExpiresAt time.Time
}

// NewUnifiedAuthMiddleware crea el middleware unificado usando nuestros componentes existentes
func NewUnifiedAuthMiddleware(config *UnifiedAuthConfig) *UnifiedAuthMiddleware {
	if config == nil {
		config = DefaultUnifiedAuthConfig()
	}

	db := database.DBconn

	// Inicializar motor RBAC dinámico
	rbacEngine := rbac.NewDynamicRBACEngine(db, rbac.DefaultRBACConfig())

	return &UnifiedAuthMiddleware{
		db:              db,
		jwtService:      utils.NewJWTService(),
		authRepo:        auth.NewAuthRepository(db),
		rbacEngine:      rbacEngine,
		config:          config,
		permissionCache: make(map[string]*PermissionCacheEntry),
	}
}

// NewUnifiedAuthMiddlewareWithDB crea el middleware unificado con una base de datos específica
func NewUnifiedAuthMiddlewareWithDB(config *UnifiedAuthConfig, db *gorm.DB) *UnifiedAuthMiddleware {
	if config == nil {
		config = DefaultUnifiedAuthConfig()
	}

	// Inicializar motor RBAC dinámico con la DB específica
	rbacEngine := rbac.NewDynamicRBACEngine(db, rbac.DefaultRBACConfig())

	return &UnifiedAuthMiddleware{
		db:              db,
		jwtService:      utils.NewJWTService(),
		authRepo:        auth.NewAuthRepository(db),
		rbacEngine:      rbacEngine,
		config:          config,
		permissionCache: make(map[string]*PermissionCacheEntry),
	}
}

// DefaultUnifiedAuthConfig configuración por defecto
func DefaultUnifiedAuthConfig() *UnifiedAuthConfig {
	return &UnifiedAuthConfig{
		CacheEnabled:   true,
		CacheTTL:       15 * time.Minute,
		DebugMode:      false,
		DefaultDenyAll: true,
		RequireOrg:     true,
		AllowAnonymous: false,
		EnableAudit:    true,
	}
}

// =============================================================================
// MIDDLEWARE PRINCIPAL - UNA LÍNEA PARA AUTH + AUTHZ
// =============================================================================

// Protect método principal que combina autenticación + autorización
// Uso: app.Get("/api/expenses", auth.Protect("expenses:read:organization"), handler)
func (m *UnifiedAuthMiddleware) Protect(permission string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 1. AUTENTICACIÓN usando nuestro sistema existente
		authCtx, err := m.authenticateRequest(c)
		if err != nil || authCtx.IsGuest {
			return m.handleAuthError(c, "authentication_required", err)
		}

		// 2. AUTORIZACIÓN dinámica desde base de datos
		allowed, err := m.checkPermissionDynamic(authCtx, permission, c)
		if err != nil {
			return m.handleAuthError(c, "authorization_error", err)
		}

		if !allowed {
			return m.handleAuthError(c, "insufficient_permissions",
				fmt.Errorf("permission %s denied", permission))
		}

		// 3. CONTEXTO - Inyectar contexto enriquecido (ya disponible por SmartAuth)
		// El contexto ya está disponible en c.Locals("authContext")

		// 4. AUDITORÍA automática
		if m.config.EnableAudit {
			go m.logAccess(authCtx, permission, c.Path(), true)
		}

		return c.Next()
	}
}

// =============================================================================
// MÉTODOS DE CONVENIENCIA SIMPLIFICADOS
// =============================================================================

// Allow para rutas públicas
func (m *UnifiedAuthMiddleware) Allow() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Crear contexto de invitado usando nuestro sistema
		c.Locals("authContext", &AuthContext{IsGuest: true})
		if m.config.EnableAudit {
			go m.logAccess(nil, "public", c.Path(), true)
		}
		return c.Next()
	}
}

// RequireAuth solo autenticación usando nuestro SmartAuthMiddleware
func (m *UnifiedAuthMiddleware) RequireAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authCtx, err := m.authenticateRequest(c)
		if err != nil || authCtx.IsGuest {
			return m.handleAuthError(c, "authentication_required", err)
		}
		return c.Next()
	}
}

// Admin requiere permisos de administrador
func (m *UnifiedAuthMiddleware) Admin() fiber.Handler {
	return m.Protect("admin:manage:all")
}

// Own para recursos propios del usuario
func (m *UnifiedAuthMiddleware) Own(resource, action string) fiber.Handler {
	return m.Protect(fmt.Sprintf("%s:%s:own", resource, action))
}

// Org para recursos a nivel organizacional
func (m *UnifiedAuthMiddleware) Org(resource, action string) fiber.Handler {
	return m.Protect(fmt.Sprintf("%s:%s:organization", resource, action))
}

// Any para permisos alternativos (OR lógico)
func (m *UnifiedAuthMiddleware) Any(permissions ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authCtx, err := m.authenticateRequest(c)
		if err != nil || authCtx.IsGuest {
			return m.handleAuthError(c, "authentication_required", err)
		}

		// Verificar si tiene al menos uno de los permisos
		for _, permission := range permissions {
			if allowed, _ := m.checkPermissionDynamic(authCtx, permission, c); allowed {
				if m.config.EnableAudit {
					go m.logAccess(authCtx, permission, c.Path(), true)
				}
				return c.Next()
			}
		}

		return m.handleAuthError(c, "insufficient_permissions",
			fmt.Errorf("none of the required permissions found"))
	}
}

// =============================================================================
// INTEGRACIÓN CON SISTEMA DE AUTENTICACIÓN EXISTENTE
// =============================================================================

// authenticateRequest usa nuestro SmartAuthMiddleware existente
func (m *UnifiedAuthMiddleware) authenticateRequest(c *fiber.Ctx) (*AuthContext, error) {
	// Reutilizar la lógica del SmartAuthMiddleware existente
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return &AuthContext{IsGuest: true}, nil
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return nil, fmt.Errorf("malformed JWT")
	}
	tokenString := parts[1]

	claims, err := m.jwtService.ValidateToken(tokenString)
	if err != nil {
		return nil, fmt.Errorf("invalid or expired JWT: %w", err)
	}

	// Extraer orgSlug de la URL (misma lógica que SmartAuth)
	orgSlug := m.extractOrgSlug(c)

	// Construir contexto usando nuestro sistema existente
	authCtx, err := m.buildAuthContext(claims, orgSlug)
	if err != nil {
		return nil, fmt.Errorf("failed to build auth context: %w", err)
	}

	// Inyectar en el contexto de Fiber (compatibilidad con código existente)
	c.Locals("authContext", authCtx)
	c.Locals("identity_id", claims.IdentityID)

	return authCtx, nil
}

// extractOrgSlug extrae el slug de organización de la URL
func (m *UnifiedAuthMiddleware) extractOrgSlug(c *fiber.Ctx) string {
	// Primero intentar desde parámetros de ruta
	orgSlug := c.Params("slug")
	if orgSlug != "" {
		return orgSlug
	}

	// Luego desde la URL directamente (misma lógica que SmartAuth)
	path := c.OriginalURL()
	if strings.Contains(path, "/org/") {
		parts := strings.Split(path, "/org/")
		if len(parts) > 1 {
			slugPart := strings.Split(parts[1], "/")[0]
			if slugPart != "" {
				return slugPart
			}
		}
	}

	return ""
}

// buildAuthContext construye el contexto usando nuestro sistema existente
func (m *UnifiedAuthMiddleware) buildAuthContext(claims *utils.UnifiedClaims, orgSlug string) (*AuthContext, error) {
	if claims.IdentityID == "" {
		return nil, fmt.Errorf("invalid token: missing identity_id")
	}

	// Usar nuestro repositorio existente
	identity, memberships, customer, err := m.authRepo.GetFullIdentityContext(claims.IdentityID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve full identity context: %w", err)
	}

	ctx := &AuthContext{
		Identity:        identity,
		Memberships:     memberships,
		CustomerProfile: customer,
		IsGuest:         false,
	}

	// Encontrar membresía organizacional si se especifica
	if orgSlug != "" {
		var currentMembership *models.OrganizationalMembership

		for _, m := range memberships {
			if m.Organization.Slug == orgSlug {
				currentMembership = &m
				break
			}
		}

		if currentMembership == nil {
			return nil, fmt.Errorf("user is not a member of the specified organization")
		}

		ctx.CurrentOrg = &currentMembership.Organization
		ctx.CurrentRole = &currentMembership.Role
	}

	return ctx, nil
}

// =============================================================================
// AUTORIZACIÓN DINÁMICA RBAC
// =============================================================================

// checkPermissionDynamic verifica permisos usando el motor RBAC dinámico
func (m *UnifiedAuthMiddleware) checkPermissionDynamic(authCtx *AuthContext, permission string, c *fiber.Ctx) (bool, error) {
	// Parsear permiso: "resource:action:scope"
	parts := strings.Split(permission, ":")
	if len(parts) != 3 {
		return false, fmt.Errorf("invalid permission format: %s", permission)
	}

	resource, action, scope := parts[0], parts[1], parts[2]

	// Obtener ID de organización del contexto
	orgID := ""
	if authCtx.CurrentOrg != nil {
		orgID = authCtx.CurrentOrg.ID
	}

	// Cache check primero
	if m.config.CacheEnabled {
		cacheKey := fmt.Sprintf("%s:%s:%s:%s:%s", authCtx.Identity.ID, orgID, resource, action, scope)
		if entry, exists := m.permissionCache[cacheKey]; exists {
			if time.Now().Before(entry.ExpiresAt) {
				return entry.Allowed, nil
			}
			// Cache expirado, eliminar entrada
			delete(m.permissionCache, cacheKey)
		}
	}

	// Usar el motor RBAC dinámico para verificar permisos
	allowed := m.rbacEngine.CheckUserPermission(authCtx.Identity.ID, orgID, resource, action, scope)

	// Cache el resultado
	if m.config.CacheEnabled {
		cacheKey := fmt.Sprintf("%s:%s:%s:%s:%s", authCtx.Identity.ID, orgID, resource, action, scope)
		m.permissionCache[cacheKey] = &PermissionCacheEntry{
			Allowed:   allowed,
			ExpiresAt: time.Now().Add(m.config.CacheTTL),
		}
	}

	return allowed, nil
}

// =============================================================================
// UTILIDADES Y HELPERS
// =============================================================================

// handleAuthError maneja errores de autenticación/autorización
func (m *UnifiedAuthMiddleware) handleAuthError(c *fiber.Ctx, errorType string, err error) error {
	// Log del error para auditoría
	if m.config.EnableAudit {
		go m.logAccess(nil, "error", c.Path(), false)
	}

	status := fiber.StatusUnauthorized
	message := "Authentication required"

	switch errorType {
	case "insufficient_permissions":
		status = fiber.StatusForbidden
		message = "Insufficient permissions"
	case "authorization_error":
		status = fiber.StatusForbidden
		message = "Authorization error"
	}

	// Usar nuestro sistema de errores existente
	return utils.SendError(c, status, message, err)
}

// logAccess registra el acceso para auditoría
func (m *UnifiedAuthMiddleware) logAccess(authCtx *AuthContext, permission, path string, success bool) {
	auditLog := &models.AuditLog{
		ID:        uuid.New().String(), // Generar ID único
		Action:    "api_access",
		Resource:  permission,
		Details:   fmt.Sprintf("Path: %s, Success: %v", path, success),
		Timestamp: time.Now(),
	}

	if authCtx != nil && !authCtx.IsGuest {
		auditLog.IdentityID = authCtx.Identity.ID
		if authCtx.CurrentOrg != nil {
			orgID := authCtx.CurrentOrg.ID
			auditLog.OrganizationID = &orgID
		}
	}

	m.db.Create(auditLog)
}

// =============================================================================
// INSTANCIA GLOBAL PARA CONVENIENCIA
// =============================================================================

// GlobalAuth instancia global del middleware
var GlobalAuth *UnifiedAuthMiddleware

// InitGlobalAuth inicializa el middleware global
func InitGlobalAuth(config *UnifiedAuthConfig) {
	GlobalAuth = NewUnifiedAuthMiddleware(config)
}

// Protect función global de conveniencia
func Protect(permission string) fiber.Handler {
	if GlobalAuth == nil {
		panic("GlobalAuth not initialized. Call InitGlobalAuth first.")
	}
	return GlobalAuth.Protect(permission)
}

// RequireAuth función global de conveniencia
func RequireAuth() fiber.Handler {
	if GlobalAuth == nil {
		panic("GlobalAuth not initialized. Call InitGlobalAuth first.")
	}
	return GlobalAuth.RequireAuth()
}

// Allow función global de conveniencia
func Allow() fiber.Handler {
	if GlobalAuth == nil {
		panic("GlobalAuth not initialized. Call InitGlobalAuth first.")
	}
	return GlobalAuth.Allow()
}

// =============================================================================
// UTILIDADES PARA SETUP AUTOMÁTICO DE PERMISOS
// =============================================================================

// CreatePermissionsFromList crea permisos automáticamente usando el sistema RBAC
func CreatePermissionsFromList(db *gorm.DB, permissions []string) error {
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}

	// Crear permisos directamente usando GORM (método simple)
	for _, permStr := range permissions {
		parts := strings.Split(permStr, ":")
		if len(parts) != 3 {
			continue // Skip malformed permissions
		}

		permission := models.Permission{
			Resource: parts[0],
			Action:   parts[1],
			Scope:    parts[2],
		}

		// Upsert del permiso usando la tabla singular
		if err := db.Where("resource = ? AND action = ? AND scope = ?",
			parts[0], parts[1], parts[2]).FirstOrCreate(&permission).Error; err != nil {
			return fmt.Errorf("failed to create permission %s: %w", permStr, err)
		}
	}
	return nil
}

// =============================================================================
// COMPATIBILIDAD CON MIDDLEWARE EXISTENTE
// =============================================================================

// GetAuthContext helper para extraer el contexto de autenticación
func GetAuthContext(c *fiber.Ctx) (*AuthContext, bool) {
	ctx, ok := c.Locals("authContext").(*AuthContext)
	return ctx, ok
}

// GetIdentityID helper para extraer el ID de identidad
func GetIdentityID(c *fiber.Ctx) (string, bool) {
	id, ok := c.Locals("identity_id").(string)
	return id, ok
}

// IsGuest verifica si el usuario es invitado
func IsGuest(c *fiber.Ctx) bool {
	ctx, ok := GetAuthContext(c)
	return !ok || ctx.IsGuest
}

// GetCurrentOrg obtiene la organización actual del contexto
func GetCurrentOrg(c *fiber.Ctx) (*models.Organization, bool) {
	ctx, ok := GetAuthContext(c)
	if !ok || ctx.IsGuest || ctx.CurrentOrg == nil {
		return nil, false
	}
	return ctx.CurrentOrg, true
}

// GetCurrentRole obtiene el rol actual del contexto
func GetCurrentRole(c *fiber.Ctx) (*models.Role, bool) {
	ctx, ok := GetAuthContext(c)
	if !ok || ctx.IsGuest || ctx.CurrentRole == nil {
		return nil, false
	}
	return ctx.CurrentRole, true
}
