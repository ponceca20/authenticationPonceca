package guest

import (
	"practicev2/module/authentication/utils"

	"github.com/gofiber/fiber/v2"
)

// GuestHandler handles the HTTP requests for guest sessions.
type GuestHandler struct {
	service GuestService
}

// NewGuestHandler creates a new instance of GuestHandler.
func NewGuestHandler(service GuestService) *GuestHandler {
	return &GuestHandler{service: service}
}

// CreateSession is the handler for creating a new guest session.
func (h *GuestHandler) CreateSession(c *fiber.Ctx) error {
	session, err := h.service.CreateGuestSession()
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to create guest session", err)
	}

	// The session token should be sent back to the client,
	// often in a secure, HttpOnly cookie.
	c.Cookie(&fiber.Cookie{
		Name:     "guest_session_token",
		Value:    session.SessionToken,
		Expires:  session.ExpiresAt,
		HTTPOnly: true,
		Secure:   true, // Set to false if not using HTTPS in dev
		SameSite: "Lax",
	})

	dto := GuestSessionDTO{
		SessionToken: session.SessionToken,
		ExpiresAt:    session.ExpiresAt,
	}

	return utils.SendSuccess(c, fiber.StatusCreated, dto, "Guest session created")
}

// GetSession retrieves the current guest session using a token from a cookie.
func (h *GuestHandler) GetSession(c *fiber.Ctx) error {
	token := c.Cookies("guest_session_token")
	if token == "" {
		return utils.SendError(c, fiber.StatusBadRequest, "Guest session token not found")
	}

	session, err := h.service.GetSession(token)
	if err != nil {
		return utils.SendError(c, fiber.StatusNotFound, "Guest session not found or expired", err)
	}

	dto := GuestSessionDTO{
		SessionToken: session.SessionToken,
		ExpiresAt:    session.ExpiresAt,
		CartData:     session.CartData,
	}

	return utils.SendSuccess(c, fiber.StatusOK, dto)
}

// UpdateCart updates the cart data for the current guest session.
func (h *GuestHandler) UpdateCart(c *fiber.Ctx) error {
	token := c.Cookies("guest_session_token")
	if token == "" {
		return utils.SendError(c, fiber.StatusBadRequest, "Guest session token not found")
	}

	var dto UpdateCartDTO
	if err := c.BodyParser(&dto); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	if errors := utils.ValidateStruct(&dto); errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": errors})
	}

	_, err := h.service.UpdateCart(token, dto.CartData)
	if err != nil {
		return utils.SendError(c, fiber.StatusNotFound, err.Error())
	}

	return utils.SendSuccess(c, fiber.StatusOK, nil, "Cart updated successfully")
}
