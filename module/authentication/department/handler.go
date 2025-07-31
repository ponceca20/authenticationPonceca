package department

import (
	"practicev2/module/authentication/middleware"
	"practicev2/module/authentication/utils"

	"github.com/gofiber/fiber/v2"
)

// DepartmentHandler handles the HTTP requests for departments.
type DepartmentHandler struct {
	service DepartmentService
}

// NewDepartmentHandler creates a new instance of DepartmentHandler.
func NewDepartmentHandler(service DepartmentService) *DepartmentHandler {
	return &DepartmentHandler{service: service}
}

// CreateDepartment is the handler for creating a new department.
func (h *DepartmentHandler) CreateDepartment(c *fiber.Ctx) error {
	authCtx, ok := c.Locals("authContext").(*middleware.AuthContext)
	if !ok || authCtx.CurrentOrg == nil {
		return utils.SendError(c, fiber.StatusForbidden, "Organization context not found")
	}

	var dto DepartmentDTO
	if err := c.BodyParser(&dto); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	dept, err := h.service.CreateDepartment(authCtx.CurrentOrg.ID, &dto)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to create department", err)
	}

	return utils.SendSuccess(c, fiber.StatusCreated, ToDepartmentResponseDTO(dept))
}

// ListDepartments lists all departments in an organization.
func (h *DepartmentHandler) ListDepartments(c *fiber.Ctx) error {
	authCtx, ok := c.Locals("authContext").(*middleware.AuthContext)
	if !ok || authCtx.CurrentOrg == nil {
		return utils.SendError(c, fiber.StatusForbidden, "Organization context not found")
	}

	depts, err := h.service.ListDepartments(authCtx.CurrentOrg.ID)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to retrieve departments", err)
	}

	var dtos []DepartmentResponseDTO
	for _, dept := range depts {
		dtos = append(dtos, ToDepartmentResponseDTO(&dept))
	}
	return utils.SendSuccess(c, fiber.StatusOK, dtos)
}

// GetDepartment retrieves a single department.
func (h *DepartmentHandler) GetDepartment(c *fiber.Ctx) error {
	authCtx, ok := c.Locals("authContext").(*middleware.AuthContext)
	if !ok || authCtx.CurrentOrg == nil {
		return utils.SendError(c, fiber.StatusForbidden, "Organization context not found")
	}
	deptID := c.Params("id")

	dept, err := h.service.GetDepartment(authCtx.CurrentOrg.ID, deptID)
	if err != nil {
		return utils.SendError(c, fiber.StatusNotFound, "Department not found", err)
	}
	return utils.SendSuccess(c, fiber.StatusOK, ToDepartmentResponseDTO(dept))
}

// UpdateDepartment updates a department.
func (h *DepartmentHandler) UpdateDepartment(c *fiber.Ctx) error {
	authCtx, ok := c.Locals("authContext").(*middleware.AuthContext)
	if !ok || authCtx.CurrentOrg == nil {
		return utils.SendError(c, fiber.StatusForbidden, "Organization context not found")
	}
	deptID := c.Params("id")

	var dto DepartmentDTO
	if err := c.BodyParser(&dto); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	dept, err := h.service.UpdateDepartment(authCtx.CurrentOrg.ID, deptID, &dto)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to update department", err)
	}
	return utils.SendSuccess(c, fiber.StatusOK, ToDepartmentResponseDTO(dept))
}

// DeleteDepartment deletes a department.
func (h *DepartmentHandler) DeleteDepartment(c *fiber.Ctx) error {
	authCtx, ok := c.Locals("authContext").(*middleware.AuthContext)
	if !ok || authCtx.CurrentOrg == nil {
		return utils.SendError(c, fiber.StatusForbidden, "Organization context not found")
	}
	deptID := c.Params("id")

	if err := h.service.DeleteDepartment(authCtx.CurrentOrg.ID, deptID); err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to delete department", err)
	}
	return utils.SendSuccess(c, fiber.StatusNoContent, nil)
}

// ListUsersInDepartment lists all users in a department.
func (h *DepartmentHandler) ListUsersInDepartment(c *fiber.Ctx) error {
	authCtx, ok := c.Locals("authContext").(*middleware.AuthContext)
	if !ok || authCtx.CurrentOrg == nil {
		return utils.SendError(c, fiber.StatusForbidden, "Organization context not found")
	}
	deptID := c.Params("id")

	users, err := h.service.ListUsersInDepartment(authCtx.CurrentOrg.ID, deptID)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to retrieve users for department", err)
	}

	return utils.SendSuccess(c, fiber.StatusOK, users)
}
