package admin

import (
	"practicev2/module/authentication/utils"

	"github.com/gofiber/fiber/v2"
)

// AdminHandler handles requests for system-level administration.
type AdminHandler struct {
	service AdminService
}

// NewAdminHandler creates a new instance of AdminHandler.
func NewAdminHandler(service AdminService) *AdminHandler {
	return &AdminHandler{service: service}
}

// GetSystemStats is the handler for the system stats endpoint.
func (h *AdminHandler) GetSystemStats(c *fiber.Ctx) error {
	// TODO: Add super-admin permission check here.

	stats, err := h.service.GetSystemStats()
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to get system stats", err)
	}
	return utils.SendSuccess(c, fiber.StatusOK, stats)
}
