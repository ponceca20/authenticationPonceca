package role

import (
	"practicev2/module/authentication/middleware"
	"practicev2/module/authentication/user"
	"practicev2/module/authentication/utils"

	"github.com/gofiber/fiber/v2"
)

// RoleHandler handles the HTTP requests for roles and permissions.
type RoleHandler struct {
	roleService RoleService
	userService user.UserService
}

// NewRoleHandler creates a new instance of RoleHandler.
func NewRoleHandler(roleService RoleService, userService user.UserService) *RoleHandler {
	return &RoleHandler{roleService: roleService, userService: userService}
}

// CreateRole is the handler for creating a new role in an organization.
func (h *RoleHandler) CreateRole(c *fiber.Ctx) error {
	authCtx, ok := c.Locals("authContext").(*middleware.AuthContext)

	if !ok || authCtx.CurrentOrg == nil {
		return utils.SendError(c, fiber.StatusForbidden, "Organization context not found")
	}

	// TODO: Add permission check: if !authCtx.HasPermission("roles.create") ...

	var dto RoleDTO
	if err := c.BodyParser(&dto); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	role, err := h.roleService.CreateRole(authCtx.CurrentOrg.ID, &dto)
	if err != nil {
		return utils.SendError(c, fiber.StatusConflict, "Failed to create role", err)
	}

	return utils.SendSuccess(c, fiber.StatusCreated, ToRoleResponseDTO(role))
}

// ListRoles is the handler for listing all roles in an organization.
func (h *RoleHandler) ListRoles(c *fiber.Ctx) error {
	authCtx, ok := c.Locals("authContext").(*middleware.AuthContext)
	if !ok || authCtx.CurrentOrg == nil {
		return utils.SendError(c, fiber.StatusForbidden, "Organization context not found")
	}

	// TODO: Add permission check: if !authCtx.HasPermission("roles.read") ...

	roles, err := h.roleService.ListRoles(authCtx.CurrentOrg.ID)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to retrieve roles", err)
	}

	return utils.SendSuccess(c, fiber.StatusOK, roles)
}

// GetRole is the handler for retrieving a specific role by ID.
func (h *RoleHandler) GetRole(c *fiber.Ctx) error {
	authCtx, ok := c.Locals("authContext").(*middleware.AuthContext)
	if !ok || authCtx.CurrentOrg == nil {
		return utils.SendError(c, fiber.StatusForbidden, "Organization context not found")
	}

	roleID := c.Params("id")
	if roleID == "" {
		return utils.SendError(c, fiber.StatusBadRequest, "Role ID is required")
	}

	// TODO: Add permission check: if !authCtx.HasPermission("roles.read") ...

	role, err := h.roleService.GetRole(authCtx.CurrentOrg.ID, roleID)
	if err != nil {
		return utils.SendError(c, fiber.StatusNotFound, "Role not found", err)
	}

	return utils.SendSuccess(c, fiber.StatusOK, ToRoleResponseDTO(role))
}

// UpdateRole is the handler for updating a role.
func (h *RoleHandler) UpdateRole(c *fiber.Ctx) error {
	authCtx, ok := c.Locals("authContext").(*middleware.AuthContext)
	if !ok || authCtx.CurrentOrg == nil {
		return utils.SendError(c, fiber.StatusForbidden, "Organization context not found")
	}

	roleID := c.Params("id")
	if roleID == "" {
		return utils.SendError(c, fiber.StatusBadRequest, "Role ID is required")
	}

	// TODO: Add permission check: if !authCtx.HasPermission("roles.update") ...

	// First check if role exists before validating the body
	_, err := h.roleService.GetRole(authCtx.CurrentOrg.ID, roleID)
	if err != nil {
		return utils.SendError(c, fiber.StatusNotFound, "Role not found")
	}

	var dto RoleDTO
	if err := c.BodyParser(&dto); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	role, err := h.roleService.UpdateRole(authCtx.CurrentOrg.ID, roleID, &dto)
	if err != nil {
		if err.Error() == "role not found" {
			return utils.SendError(c, fiber.StatusNotFound, "Role not found")
		}
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to update role", err)
	}

	return utils.SendSuccess(c, fiber.StatusOK, ToRoleResponseDTO(role))
}

// DeleteRole is the handler for deleting a role.
func (h *RoleHandler) DeleteRole(c *fiber.Ctx) error {
	authCtx, ok := c.Locals("authContext").(*middleware.AuthContext)
	if !ok || authCtx.CurrentOrg == nil {
		return utils.SendError(c, fiber.StatusForbidden, "Organization context not found")
	}

	roleID := c.Params("id")
	if roleID == "" {
		return utils.SendError(c, fiber.StatusBadRequest, "Role ID is required")
	}

	// TODO: Add permission check: if !authCtx.HasPermission("roles.delete") ...

	if err := h.roleService.DeleteRole(authCtx.CurrentOrg.ID, roleID); err != nil {
		if err.Error() == "role not found" {
			return utils.SendError(c, fiber.StatusNotFound, "Role not found")
		}
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to delete role", err)
	}

	return utils.SendSuccess(c, fiber.StatusNoContent, nil)
}

// ListUsersInRole is the handler for listing users assigned to a specific role.
func (h *RoleHandler) ListUsersInRole(c *fiber.Ctx) error {
	authCtx, ok := c.Locals("authContext").(*middleware.AuthContext)
	if !ok || authCtx.CurrentOrg == nil {
		return utils.SendError(c, fiber.StatusForbidden, "Organization context not found")
	}

	roleID := c.Params("id")
	if roleID == "" {
		return utils.SendError(c, fiber.StatusBadRequest, "Role ID is required")
	}

	// TODO: Add permission check: if !authCtx.HasPermission("roles.read") ...

	users, err := h.roleService.ListUsersInRole(authCtx.CurrentOrg.ID, roleID)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to retrieve users for role", err)
	}

	return utils.SendSuccess(c, fiber.StatusOK, users)
}

// AssignRoleToUser is the handler for assigning a role to a user.
func (h *RoleHandler) AssignRoleToUser(c *fiber.Ctx) error {
	authCtx, ok := c.Locals("authContext").(*middleware.AuthContext)
	if !ok || authCtx.CurrentOrg == nil {
		return utils.SendError(c, fiber.StatusForbidden, "Organization context not found")
	}

	roleID := c.Params("id")
	userID := c.Params("userId")

	// TODO: Add permission check: if !authCtx.HasPermission("roles.assign") ...

	if err := h.userService.AssignRoleToUser(authCtx.CurrentOrg.ID, userID, roleID); err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to assign role", err)
	}

	return utils.SendSuccess(c, fiber.StatusOK, nil, "Role assigned successfully.")
}
