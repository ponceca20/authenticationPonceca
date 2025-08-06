package profile

import (
	"fmt"
	"practicev2/module/authentication/middleware"
	"practicev2/module/authentication/utils"
	"strings"

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

	// Validate the DTO
	if errs := utils.ValidateStruct(&dto); errs != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"errors": errs,
		})
	}

	updatedProfile, err := h.service.UpdateProfile(authCtx.Identity.ID, &dto)
	if err != nil {
		// Log the actual error for debugging
		fmt.Printf("Profile update error: %v\n", err)

		// Handle specific error types
		errMsg := err.Error()

		// User not found
		if strings.Contains(errMsg, "user not found") {
			return utils.SendError(c, fiber.StatusNotFound, "User not found", err)
		}

		// Validation errors
		if strings.Contains(errMsg, "validation failed") {
			return utils.SendError(c, fiber.StatusBadRequest, "Validation failed", err)
		}

		// Database/transaction errors
		if strings.Contains(errMsg, "failed to update") || strings.Contains(errMsg, "failed to find") {
			return utils.SendError(c, fiber.StatusInternalServerError, "Failed to update profile", err)
		}

		// Default to bad request for other errors
		return utils.SendError(c, fiber.StatusBadRequest, "Failed to update profile", err)
	}

	return utils.SendSuccess(c, fiber.StatusOK, updatedProfile, "Profile updated successfully")
}
