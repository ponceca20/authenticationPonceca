package user

import (
	"practicev2/module/authentication/middleware"
	"practicev2/module/authentication/utils"

	"github.com/gofiber/fiber/v2"
)

// UserHandler handles the HTTP requests for organizational users.
type UserHandler struct {
	service UserService
}

// NewUserHandler creates a new instance of UserHandler.
func NewUserHandler(service UserService) *UserHandler {
	return &UserHandler{service: service}
}

// ListUsers is the handler for listing all users in an organization.
// This route should be protected by middleware that ensures the user is part of the org.
func (h *UserHandler) ListUsers(c *fiber.Ctx) error {
	authCtx, ok := c.Locals("authContext").(*middleware.AuthContext)
	if !ok || authCtx.CurrentOrg == nil {
		return utils.SendError(c, fiber.StatusForbidden, "Organization context not found")
	}

	// TODO: Add permission check here, e.g., check if authCtx.CurrentRole has 'users.read' permission.

	users, err := h.service.ListOrgUsers(authCtx.CurrentOrg.ID)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to retrieve users", err)
	}

	return utils.SendSuccess(c, fiber.StatusOK, users)
}

// GetUser is the handler for retrieving a single user from an organization.
func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	authCtx, ok := c.Locals("authContext").(*middleware.AuthContext)
	if !ok || authCtx.CurrentOrg == nil {
		return utils.SendError(c, fiber.StatusForbidden, "Organization context not found")
	}

	userID := c.Params("id")
	if userID == "" {
		return utils.SendError(c, fiber.StatusBadRequest, "User ID is required")
	}

	// TODO: Add permission check here. A user might only be able to see their own profile,
	// unless they have special permissions.

	user, err := h.service.GetOrgUser(authCtx.CurrentOrg.ID, userID)
	if err != nil {
		return utils.SendError(c, fiber.StatusNotFound, "User not found in this organization", err)
	}

	return utils.SendSuccess(c, fiber.StatusOK, user)
}

// UpdateUser is the handler for updating a user's details in an organization.
func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	authCtx, ok := c.Locals("authContext").(*middleware.AuthContext)
	if !ok || authCtx.CurrentOrg == nil {
		return utils.SendError(c, fiber.StatusForbidden, "Organization context not found")
	}

	userID := c.Params("id")
	if userID == "" {
		return utils.SendError(c, fiber.StatusBadRequest, "User ID is required")
	}

	// TODO: Add permission check: if !authCtx.HasPermission("users.update") ...

	var dto UpdateUserDTO
	if err := c.BodyParser(&dto); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	updatedUser, err := h.service.UpdateUser(authCtx.CurrentOrg.ID, userID, &dto)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to update user", err)
	}

	return utils.SendSuccess(c, fiber.StatusOK, updatedUser)
}

// UpdateUserStatus is the handler for changing a user's active status.
func (h *UserHandler) UpdateUserStatus(c *fiber.Ctx) error {
	authCtx, ok := c.Locals("authContext").(*middleware.AuthContext)
	if !ok || authCtx.CurrentOrg == nil {
		return utils.SendError(c, fiber.StatusForbidden, "Organization context not found")
	}

	userID := c.Params("id")
	if userID == "" {
		return utils.SendError(c, fiber.StatusBadRequest, "User ID is required")
	}

	// TODO: Add permission check: if !authCtx.HasPermission("users.update.status") ...

	var dto UpdateUserStatusDTO
	if err := c.BodyParser(&dto); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	if err := h.service.UpdateUserStatus(authCtx.CurrentOrg.ID, userID, &dto); err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to update user status", err)
	}

	return utils.SendSuccess(c, fiber.StatusOK, nil, "User status updated successfully.")
}

// BulkCreateUsers is the handler for creating multiple users at once.
func (h *UserHandler) BulkCreateUsers(c *fiber.Ctx) error {
	authCtx, ok := c.Locals("authContext").(*middleware.AuthContext)
	if !ok || authCtx.CurrentOrg == nil {
		return utils.SendError(c, fiber.StatusForbidden, "Organization context not found")
	}

	// TODO: Add permission check: if !authCtx.HasPermission("users.bulk.create") ...

	var dto BulkCreateUserDTO
	if err := c.BodyParser(&dto); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	if errs := utils.ValidateStruct(&dto); errs != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": errs})
	}

	if err := h.service.BulkCreateUsers(authCtx.CurrentOrg.ID, &dto); err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to bulk create users", err)
	}

	return utils.SendSuccess(c, fiber.StatusCreated, nil, "Users created successfully.")
}

// BulkUpdateUsers is the handler for updating multiple users at once.
func (h *UserHandler) BulkUpdateUsers(c *fiber.Ctx) error {
	authCtx, ok := c.Locals("authContext").(*middleware.AuthContext)
	if !ok || authCtx.CurrentOrg == nil {
		return utils.SendError(c, fiber.StatusForbidden, "Organization context not found")
	}

	var dto BulkUpdateUserDTO
	if err := c.BodyParser(&dto); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	if errs := utils.ValidateStruct(&dto); errs != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": errs})
	}

	if err := h.service.BulkUpdateUsers(authCtx.CurrentOrg.ID, &dto); err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to bulk update users", err)
	}

	return utils.SendSuccess(c, fiber.StatusOK, nil, "Users updated successfully.")
}

// BulkDeleteUsers is the handler for deleting multiple users at once.
func (h *UserHandler) BulkDeleteUsers(c *fiber.Ctx) error {
	authCtx, ok := c.Locals("authContext").(*middleware.AuthContext)
	if !ok || authCtx.CurrentOrg == nil {
		return utils.SendError(c, fiber.StatusForbidden, "Organization context not found")
	}

	var dto BulkDeleteUserDTO
	if err := c.BodyParser(&dto); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	if errs := utils.ValidateStruct(&dto); errs != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": errs})
	}

	if err := h.service.BulkDeleteUsers(authCtx.CurrentOrg.ID, &dto); err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to bulk delete users", err)
	}

	return utils.SendSuccess(c, fiber.StatusNoContent, nil)
}
