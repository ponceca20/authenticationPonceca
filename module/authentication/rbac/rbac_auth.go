package rbac

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"practicev2/module/authentication/models"
	"practicev2/module/authentication/utils"

	"github.com/gofiber/fiber/v2"
)

// =============================================================================
// CONTEXTO DE USUARIO Y AUTENTICACIÓN
// =============================================================================

// getUserContext extrae y construye el contexto completo del usuario
func (e *RBACEngine) getUserContext(c *fiber.Ctx) (*UserContext, error) {
	// 1. Extraer token del header Authorization
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return nil, fmt.Errorf("authorization header missing")
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		return nil, fmt.Errorf("invalid authorization header format")
	}

	// 2. Validar y parsear JWT
	jwtService := utils.NewJWTService()
	claims, err := jwtService.ValidateToken(tokenString)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	// 3. Extraer datos del token
	identityID := claims.IdentityID
	if identityID == "" {
		return nil, fmt.Errorf("invalid identity_id in token")
	}

	email := claims.Email
	if email == "" {
		return nil, fmt.Errorf("invalid email in token")
	}

	// 4. Cargar datos completos del usuario
	userCtx, err := e.loadUserContext(identityID, email)
	if err != nil {
		return nil, fmt.Errorf("failed to load user context: %w", err)
	}

	// 5. Procesar membresías desde el token
	if len(claims.Memberships) > 0 {
		userCtx.Memberships = e.processMembershipsFromClaims(claims.Memberships, identityID)
	}

	// 6. Construir mapa de permisos
	userCtx.Permissions = e.buildPermissionsMap(userCtx.Memberships)

	// 7. Verificar si es admin del sistema
	userCtx.IsSystemAdmin = e.isSystemAdmin(userCtx)

	return userCtx, nil
}

// loadUserContext carga el contexto base del usuario desde la base de datos
func (e *RBACEngine) loadUserContext(identityID, email string) (*UserContext, error) {
	// Cargar Identity
	var identity models.Identity
	if err := e.db.Where("id = ?", identityID).First(&identity).Error; err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Cargar membresías organizacionales con roles y permisos
	var memberships []models.OrganizationalMembership
	err := e.db.Preload("Organization").
		Preload("Role").
		Preload("Role.Permissions").
		Where("identity_id = ? AND is_active = ?", identityID, true).
		Find(&memberships).Error
	if err != nil {
		return nil, fmt.Errorf("failed to load memberships: %w", err)
	}

	// Cargar perfil de cliente si existe
	var customerProfile models.CustomerProfile
	customerProfileExists := e.db.Where("identity_id = ?", identityID).First(&customerProfile).Error == nil

	userCtx := &UserContext{
		Identity:    &identity,
		Memberships: memberships,
	}

	if customerProfileExists {
		userCtx.CustomerProfile = &customerProfile
	}

	return userCtx, nil
}

// processMembershipsFromClaims procesa las membresías desde los claims del JWT
func (e *RBACEngine) processMembershipsFromClaims(membershipClaims []utils.MembershipClaim, identityID string) []models.OrganizationalMembership {
	var memberships []models.OrganizationalMembership

	for _, membershipClaim := range membershipClaims {
		if membershipClaim.OrganizationID != "" && membershipClaim.Role != "" {
			// Cargar datos completos desde la base de datos
			var membership models.OrganizationalMembership
			err := e.db.Preload("Organization").
				Preload("Role").
				Preload("Role.Permissions").
				Where("identity_id = ? AND organization_id = ? AND is_active = ?",
					identityID, membershipClaim.OrganizationID, true).
				First(&membership).Error

			if err == nil {
				memberships = append(memberships, membership)
			}
		}
	}

	return memberships
}

// buildPermissionsMap construye el mapa de permisos del usuario
func (e *RBACEngine) buildPermissionsMap(memberships []models.OrganizationalMembership) map[string]map[string]map[string]bool {
	permissions := make(map[string]map[string]map[string]bool)

	for _, membership := range memberships {
		for _, permission := range membership.Role.Permissions {
			// Inicializar estructura si no existe
			if permissions[permission.Resource] == nil {
				permissions[permission.Resource] = make(map[string]map[string]bool)
			}
			if permissions[permission.Resource][permission.Action] == nil {
				permissions[permission.Resource][permission.Action] = make(map[string]bool)
			}

			// Marcar permiso como permitido
			permissions[permission.Resource][permission.Action][permission.Scope] = true
		}
	}

	return permissions
}

// isSystemAdmin verifica si el usuario es administrador del sistema
func (e *RBACEngine) isSystemAdmin(userCtx *UserContext) bool {
	// Verificar si tiene rol de system admin en alguna organización
	for _, membership := range userCtx.Memberships {
		if membership.Role.IsSystemRole && membership.Role.Name == "system_admin" {
			return true
		}
	}
	return false
}

// determineActiveOrganization determina la organización activa para la request
func (e *RBACEngine) determineActiveOrganization(c *fiber.Ctx, userCtx *UserContext) string {
	// 1. Verificar header X-Organization-ID
	if orgID := c.Get("X-Organization-ID"); orgID != "" {
		// Verificar que el usuario tiene membresía en esta organización
		for _, membership := range userCtx.Memberships {
			if membership.OrganizationID == orgID && membership.IsActive {
				return orgID
			}
		}
	}

	// 2. Verificar parámetro en la URL (ej: /api/v1/org/:orgSlug/...)
	if orgSlug := c.Params("orgSlug"); orgSlug != "" {
		// Buscar organización por slug
		for _, membership := range userCtx.Memberships {
			if membership.Organization.Slug == orgSlug && membership.IsActive {
				return membership.OrganizationID
			}
		}
	}

	// 3. Verificar query parameter
	if orgID := c.Query("org_id"); orgID != "" {
		for _, membership := range userCtx.Memberships {
			if membership.OrganizationID == orgID && membership.IsActive {
				return orgID
			}
		}
	}

	// 4. Usar la primera organización activa como fallback
	if len(userCtx.Memberships) > 0 {
		for _, membership := range userCtx.Memberships {
			if membership.IsActive {
				return membership.OrganizationID
			}
		}
	}

	return ""
}

// =============================================================================
// MANEJO DE ERRORES DE AUTENTICACIÓN
// =============================================================================

// handleAuthError maneja errores de autenticación de forma consistente
func (e *RBACEngine) handleAuthError(c *fiber.Ctx, errorType string, err error) error {
	response := fiber.Map{
		"status":  "error",
		"message": "Authentication failed",
		"error":   errorType,
	}

	statusCode := fiber.StatusUnauthorized

	switch errorType {
	case "invalid_token":
		response["message"] = "Invalid or expired token"
		response["details"] = "Please login again"
		statusCode = fiber.StatusUnauthorized

	case "insufficient_permissions":
		response["message"] = "Insufficient permissions"
		response["details"] = "You don't have permission to access this resource"
		statusCode = fiber.StatusForbidden

	case "insufficient_role":
		response["message"] = "Required role not found"
		response["details"] = "You don't have the required role to access this resource"
		statusCode = fiber.StatusForbidden

	case "no_organization":
		response["message"] = "Organization context required"
		response["details"] = "This endpoint requires organization context"
		statusCode = fiber.StatusBadRequest

	case "permission_check_failed":
		response["message"] = "Permission check failed"
		response["details"] = "Unable to verify permissions"
		statusCode = fiber.StatusInternalServerError
	}

	if e.config.DebugMode && err != nil {
		response["debug_error"] = err.Error()
	}

	// Log del error para auditoría
	if e.config.EnableAuditLog {
		go e.logAuthError(c, errorType, err)
	}

	return c.Status(statusCode).JSON(response)
}

// =============================================================================
// LOGGING Y AUDITORÍA
// =============================================================================

// logAccessGranted registra acceso concedido
func (e *RBACEngine) logAccessGranted(userCtx *UserContext, c *fiber.Ctx, resource, action, scope string, duration time.Duration) {
	auditLog := &models.AuditLog{
		ID:         generateUUID(),
		IdentityID: userCtx.Identity.ID,
		Action:     fmt.Sprintf("access.granted.%s.%s", resource, action),
		Resource:   resource,
		ResourceID: scope,
		Status:     "success",
		IPAddress:  c.IP(),
		UserAgent:  c.Get("User-Agent"),
		Timestamp:  time.Now(),
	}

	details := map[string]interface{}{
		"method":     c.Method(),
		"path":       c.Path(),
		"resource":   resource,
		"action":     action,
		"scope":      scope,
		"duration":   duration.String(),
		"user_email": userCtx.Identity.Email,
	}

	if userCtx.ActiveOrg != nil {
		auditLog.OrganizationID = &userCtx.ActiveOrg.ID
		details["organization"] = userCtx.ActiveOrg.Name
	}

	if detailsJSON, err := json.Marshal(details); err == nil {
		auditLog.Details = string(detailsJSON)
	}

	e.db.Create(auditLog)
}

// logAccessDenied registra acceso denegado
func (e *RBACEngine) logAccessDenied(userCtx *UserContext, c *fiber.Ctx, resource, action, scope string) {
	auditLog := &models.AuditLog{
		ID:         generateUUID(),
		IdentityID: userCtx.Identity.ID,
		Action:     fmt.Sprintf("access.denied.%s.%s", resource, action),
		Resource:   resource,
		ResourceID: scope,
		Status:     "failure",
		IPAddress:  c.IP(),
		UserAgent:  c.Get("User-Agent"),
		Timestamp:  time.Now(),
	}

	details := map[string]interface{}{
		"method":     c.Method(),
		"path":       c.Path(),
		"resource":   resource,
		"action":     action,
		"scope":      scope,
		"user_email": userCtx.Identity.Email,
		"reason":     "insufficient_permissions",
	}

	if userCtx.ActiveOrg != nil {
		auditLog.OrganizationID = &userCtx.ActiveOrg.ID
		details["organization"] = userCtx.ActiveOrg.Name
	}

	if detailsJSON, err := json.Marshal(details); err == nil {
		auditLog.Details = string(detailsJSON)
	}

	e.db.Create(auditLog)
}

// logAuthError registra errores de autenticación
func (e *RBACEngine) logAuthError(c *fiber.Ctx, errorType string, err error) {
	// Obtener IP de forma segura para tests
	var ipAddress string = "unknown"
	var userAgent string = "unknown"

	if c != nil {
		// Usar defer/recover para manejar panics en tests
		func() {
			defer func() {
				if r := recover(); r != nil {
					// Ignore panic in tests
				}
			}()
			if ip := c.IP(); ip != "" {
				ipAddress = ip
			}
			if ua := c.Get("User-Agent"); ua != "" {
				userAgent = ua
			}
		}()
	}

	auditLog := &models.AuditLog{
		ID:        generateUUID(),
		Action:    fmt.Sprintf("auth.error.%s", errorType),
		Resource:  "authentication",
		Status:    "failure",
		IPAddress: ipAddress,
		UserAgent: userAgent,
		Timestamp: time.Now(),
	}

	details := map[string]interface{}{
		"error_type": errorType,
	}

	// Obtener información de forma segura para tests
	if c != nil {
		if method := c.Method(); method != "" {
			details["method"] = method
		}
		if path := c.Path(); path != "" {
			details["path"] = path
		}
		// OriginalURL puede fallar en tests, usar de forma segura
		func() {
			defer func() {
				if r := recover(); r != nil {
					// Ignore panic in tests
				}
			}()
			if url := c.OriginalURL(); url != "" {
				details["endpoint"] = url
			}
		}()
	}

	if err != nil {
		details["error_message"] = err.Error()
	}

	if detailsJSON, err := json.Marshal(details); err == nil {
		auditLog.Details = string(detailsJSON)
	}

	e.db.Create(auditLog)
}

// =============================================================================
// UTILIDADES
// =============================================================================

// generateUUID genera un UUID simple (puedes usar una librería más robusta)
func generateUUID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// =============================================================================
// MIDDLEWARE HELPERS ADICIONALES
// =============================================================================

// AllowSelf permite acceso solo a recursos propios del usuario
func (e *RBACEngine) AllowSelf() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userCtx, err := GetUserFromContext(c)
		if err != nil {
			return e.handleAuthError(c, "invalid_token", err)
		}

		// Verificar que el recurso pertenece al usuario
		resourceUserID := c.Params("userID")
		if resourceUserID == "" {
			resourceUserID = c.Params("id")
		}

		if resourceUserID != userCtx.Identity.ID {
			return e.handleAuthError(c, "insufficient_permissions",
				fmt.Errorf("can only access own resources"))
		}

		return c.Next()
	}
}

// AllowOrganizationMembers permite acceso solo a miembros de la organización
func (e *RBACEngine) AllowOrganizationMembers() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userCtx, err := GetUserFromContext(c)
		if err != nil {
			return e.handleAuthError(c, "invalid_token", err)
		}

		orgID := GetOrgIDFromContext(c)
		if orgID == "" {
			return e.handleAuthError(c, "no_organization",
				fmt.Errorf("organization context required"))
		}

		// Verificar membresía activa
		isMember := false
		for _, membership := range userCtx.Memberships {
			if membership.OrganizationID == orgID && membership.IsActive {
				isMember = true
				break
			}
		}

		if !isMember {
			return e.handleAuthError(c, "insufficient_permissions",
				fmt.Errorf("not a member of this organization"))
		}

		return c.Next()
	}
}
