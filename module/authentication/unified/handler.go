package unified

import (
	"practicev2/module/authentication/middleware"
	"practicev2/module/authentication/utils"

	"github.com/gofiber/fiber/v2"
)

// UnifiedHandler handles requests for unified, cross-context resources.
type UnifiedHandler struct {
	service UnifiedService
}

// NewUnifiedHandler creates a new instance of UnifiedHandler.
func NewUnifiedHandler(service UnifiedService) *UnifiedHandler {
	return &UnifiedHandler{service: service}
}

// GetDashboard is the handler for the unified dashboard endpoint.
func (h *UnifiedHandler) GetDashboard(c *fiber.Ctx) error {
	authCtx, ok := c.Locals("authContext").(*middleware.AuthContext)
	if !ok || authCtx.IsGuest {
		return utils.SendError(c, fiber.StatusUnauthorized, "User not authenticated")
	}

	dashboardData, err := h.service.GetDashboardData(authCtx.Identity)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to get dashboard data", err)
	}

	return utils.SendSuccess(c, fiber.StatusOK, dashboardData)
}
