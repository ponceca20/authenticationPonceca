package integration

import (
	"practicev2/module/authentication/utils"

	"github.com/gofiber/fiber/v2"
)

// IntegrationHandler handles requests for integration with other services.
type IntegrationHandler struct {
	service IntegrationService
}

// NewIntegrationHandler creates a new instance of IntegrationHandler.
func NewIntegrationHandler(service IntegrationService) *IntegrationHandler {
	return &IntegrationHandler{service: service}
}

// ValidateToken is the handler for the token validation endpoint.
func (h *IntegrationHandler) ValidateToken(c *fiber.Ctx) error {
	var dto ValidateTokenRequestDTO
	if err := c.BodyParser(&dto); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	response, _ := h.service.ValidateToken(dto.Token)
	return utils.SendSuccess(c, fiber.StatusOK, response)
}
