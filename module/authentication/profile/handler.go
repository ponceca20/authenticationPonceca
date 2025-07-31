package profile

import (
	"practicev2/module/authentication/middleware"
	"practicev2/module/authentication/utils"

	"github.com/gofiber/fiber/v2"
)

// ProfileHandler handles the HTTP requests for user profiles.
type ProfileHandler struct {
	service ProfileService
}

// NewProfileHandler creates a new instance of ProfileHandler.
func NewProfileHandler(service ProfileService) *ProfileHandler {
	return &ProfileHandler{service: service}
}

// GetMyProfile is the handler for retrieving the authenticated user's own profile.
func (h *ProfileHandler) GetMyProfile(c *fiber.Ctx) error {
	authCtx, ok := c.Locals("authContext").(*middleware.AuthContext)
	if !ok || authCtx.IsGuest || authCtx.Identity == nil {
		return utils.SendError(c, fiber.StatusUnauthorized, "Authentication required")
	}

	profile, err := h.service.GetProfile(authCtx.Identity.ID)
	if err != nil {
		return utils.SendError(c, fiber.StatusNotFound, "Profile not found", err)
	}

	return utils.SendSuccess(c, fiber.StatusOK, profile)
}

// UpdateMyProfile is the handler for updating the authenticated user's own profile.
func (h *ProfileHandler) UpdateMyProfile(c *fiber.Ctx) error {
	authCtx, ok := c.Locals("authContext").(*middleware.AuthContext)
	if !ok || authCtx.IsGuest || authCtx.Identity == nil {
		return utils.SendError(c, fiber.StatusUnauthorized, "Authentication required")
	}

	var dto ProfileDTO
	if err := c.BodyParser(&dto); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	updatedProfile, err := h.service.UpdateProfile(authCtx.Identity.ID, &dto)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to update profile", err)
	}

	return utils.SendSuccess(c, fiber.StatusOK, updatedProfile, "Profile updated successfully")
}
