package organization_config

import (
	"fmt"

	"practicev2/module/authentication/middleware"
	"practicev2/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

//=============================================================================
// HANDLER LIMPIO PARA CONFIGURACIÓN DE MÓDULOS ORGANIZACIONALES
//=============================================================================

// OrganizationConfigHandler maneja las operaciones HTTP para configuración de módulos
type OrganizationConfigHandler struct {
	service *OrganizationConfigService
}

// NewOrganizationConfigHandler crea una nueva instancia del handler
func NewOrganizationConfigHandler(db *gorm.DB) *OrganizationConfigHandler {
	return &OrganizationConfigHandler{
		service: NewOrganizationConfigService(db),
	}
}

//=============================================================================
// ENDPOINTS PRINCIPALES
//=============================================================================

// GetOrganizationModules obtiene todos los módulos configurados de una organización
func (h *OrganizationConfigHandler) GetOrganizationModules(c *fiber.Ctx) error {
	// Obtener contexto de autenticación
	authCtx, ok := middleware.GetAuthContext(c)
	if !ok {
		return utils.SendErrorResponse(c, 401, "Authentication required", nil)
	}

	// Obtener organización del contexto (establecido por middleware)
	org, ok := middleware.GetCurrentOrg(c)
	if !ok {
		return utils.SendErrorResponse(c, 404, "Organization not found", nil)
	}

	// Verificar permisos adicionales si es necesario
	userPermissions, err := GetUserModulePermissions(authCtx.Identity.ID, org.ID)
	if err != nil {
		return utils.SendErrorResponse(c, 500, "Failed to check permissions", err)
	}

	// Obtener módulos configurados
	configs, err := h.service.GetOrganizationModules(org.ID)
	if err != nil {
		return utils.SendErrorResponse(c, 500, "Failed to get organization modules", err)
	}

	// Convertir a DTOs
	responses := ToModuleConfigResponseList(configs)

	return utils.SendSuccessResponse(c, "Organization modules retrieved successfully", fiber.Map{
		"modules":     responses,
		"permissions": userPermissions, // Incluir permisos del usuario para el frontend
		"count":       len(responses),
	})
}

// InstallModule instala un módulo predefinido para una organización
func (h *OrganizationConfigHandler) InstallModule(c *fiber.Ctx) error {
	// Obtener contexto de autenticación
	authCtx, ok := middleware.GetAuthContext(c)
	if !ok {
		return utils.SendErrorResponse(c, 401, "Authentication required", nil)
	}

	// Parsear request
	var request ModuleInstallationRequest
	if err := c.BodyParser(&request); err != nil {
		return utils.SendErrorResponse(c, 400, "Invalid request body", err)
	}

	// Validación básica
	if request.ModuleName == "" {
		return utils.SendErrorResponse(c, 400, "Module name is required", nil)
	}

	// Obtener contexto de organización
	org, ok := middleware.GetCurrentOrg(c)
	if !ok {
		return utils.SendErrorResponse(c, 404, "Organization not found", nil)
	}

	// Verificar permisos específicos para instalación
	if !HasModulePermission(authCtx.Identity.ID, org.ID, "org_modules", "install") {
		return utils.SendErrorResponse(c, 403, "Insufficient permissions to install modules", nil)
	}

	// Obtener ID de usuario del contexto
	identityID, ok := middleware.GetIdentityID(c)
	if !ok {
		return utils.SendErrorResponse(c, 401, "User identity required", nil)
	}

	// Instalar módulo específico según el tipo
	var err error
	switch request.ModuleName {
	case "expenses":
		err = h.service.SetupExpensesModule(org.ID, identityID)
	case "inventory":
		err = h.service.SetupInventoryModule(org.ID, identityID)
	default:
		return utils.SendErrorResponse(c, 400, fmt.Sprintf("Module '%s' not available", request.ModuleName), nil)
	}

	if err != nil {
		return utils.SendErrorResponse(c, 500, "Failed to install module", err)
	}

	return utils.SendSuccessResponse(c, fmt.Sprintf("Module '%s' installed successfully", request.ModuleName), fiber.Map{
		"module":          request.ModuleName,
		"organization_id": org.ID,
	})
}

// GetAvailableModules obtiene lista de módulos disponibles para instalar
func (h *OrganizationConfigHandler) GetAvailableModules(c *fiber.Ctx) error {
	modules := GetAvailableModulesResponse()

	return utils.SendSuccessResponse(c, "Available modules retrieved successfully", fiber.Map{
		"modules": modules,
		"count":   len(modules),
	})
}

// EnableModule habilita un módulo para una organización
func (h *OrganizationConfigHandler) EnableModule(c *fiber.Ctx) error {
	moduleName := c.Params("module")
	if moduleName == "" {
		return utils.SendErrorResponse(c, 400, "Module name required", nil)
	}

	org, ok := middleware.GetCurrentOrg(c)
	if !ok {
		return utils.SendErrorResponse(c, 404, "Organization not found", nil)
	}

	if err := h.service.EnableModule(org.ID, moduleName); err != nil {
		return utils.SendErrorResponse(c, 500, "Failed to enable module", err)
	}

	return utils.SendSuccessResponse(c, fmt.Sprintf("Module '%s' enabled successfully", moduleName), fiber.Map{
		"module":          moduleName,
		"organization_id": org.ID,
		"is_enabled":      true,
	})
}

// DisableModule deshabilita un módulo para una organización
func (h *OrganizationConfigHandler) DisableModule(c *fiber.Ctx) error {
	moduleName := c.Params("module")
	if moduleName == "" {
		return utils.SendErrorResponse(c, 400, "Module name required", nil)
	}

	org, ok := middleware.GetCurrentOrg(c)
	if !ok {
		return utils.SendErrorResponse(c, 404, "Organization not found", nil)
	}

	if err := h.service.DisableModule(org.ID, moduleName); err != nil {
		return utils.SendErrorResponse(c, 500, "Failed to disable module", err)
	}

	return utils.SendSuccessResponse(c, fmt.Sprintf("Module '%s' disabled successfully", moduleName), fiber.Map{
		"module":          moduleName,
		"organization_id": org.ID,
		"is_enabled":      false,
	})
}

// GetModuleDetails obtiene detalles de un módulo específico
func (h *OrganizationConfigHandler) GetModuleDetails(c *fiber.Ctx) error {
	moduleName := c.Params("module")
	if moduleName == "" {
		return utils.SendErrorResponse(c, 400, "Module name required", nil)
	}

	org, ok := middleware.GetCurrentOrg(c)
	if !ok {
		return utils.SendErrorResponse(c, 404, "Organization not found", nil)
	}

	// Obtener configuración del módulo
	config, err := h.service.repo.GetByModuleName(org.ID, moduleName)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.SendErrorResponse(c, 404, "Module not found", nil)
		}
		return utils.SendErrorResponse(c, 500, "Failed to get module details", err)
	}

	response := ToModuleConfigResponse(config)

	return utils.SendSuccessResponse(c, "Module details retrieved successfully", response)
}
