package audit

import (
	"practicev2/module/authentication/middleware"
	"practicev2/module/authentication/utils"

	"github.com/gofiber/fiber/v2"
)

// AuditHandler handles the HTTP requests for audit logs.
type AuditHandler struct {
	service AuditService
}

// NewAuditHandler creates a new instance of AuditHandler.
func NewAuditHandler(service AuditService) *AuditHandler {
	return &AuditHandler{service: service}
}

// ListAuditLogs is the handler for querying audit logs.
func (h *AuditHandler) ListAuditLogs(c *fiber.Ctx) error {
	authCtx, ok := c.Locals("authContext").(*middleware.AuthContext)
	if !ok || authCtx.CurrentOrg == nil {
		return utils.SendError(c, fiber.StatusForbidden, "Organization context not found")
	}

	// TODO: Add permission check: if !authCtx.HasPermission("audit.read") ...

	// Parse query parameters into the DTO
	var query AuditQueryDTO
	if err := c.QueryParser(&query); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Invalid query parameters", err)
	}

	paginatedResponse, err := h.service.ListAuditLogs(authCtx.CurrentOrg.ID, &query)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to retrieve audit logs", err)
	}

	return utils.SendSuccess(c, fiber.StatusOK, paginatedResponse)
}
