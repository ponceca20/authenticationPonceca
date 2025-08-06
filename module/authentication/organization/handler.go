package organization

import (
	"practicev2/module/authentication/middleware"
	"practicev2/module/authentication/utils"

	"github.com/gofiber/fiber/v2"
)

// OrganizationHandler handles the HTTP requests for organizations.
type OrganizationHandler struct {
	service OrganizationService
}

// NewOrganizationHandler creates a new instance of OrganizationHandler.
func NewOrganizationHandler(service OrganizationService) *OrganizationHandler {
	return &OrganizationHandler{service: service}
}

// CreateOrganization is the handler for creating a new organization.
func (h *OrganizationHandler) CreateOrganization(c *fiber.Ctx) error {
	var dto OrganizationRegistrationDTO
	if err := c.BodyParser(&dto); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	// The validation is handled inside the service, but for a quicker failure,
	// we could also call utils.ValidateStruct here.
	// Let's rely on the service's validation for now to keep the handler clean.

	org, err := h.service.CreateOrganization(&dto)
	if err != nil {
		// A more sophisticated error handling could check the error type
		// and return different status codes (e.g., 409 for conflict, 400 for bad request).
		return utils.SendError(c, fiber.StatusConflict, err.Error())
	}

	return utils.SendSuccess(c, fiber.StatusCreated, ToOrganizationResponseDTO(org), "Organization created successfully")
}

// GetOrganization is the handler for retrieving a single organization by its slug.
func (h *OrganizationHandler) GetOrganization(c *fiber.Ctx) error {
	slug := c.Params("slug")

	// Check authentication context first
	authCtx := c.Locals("authContext")
	if authCtx == nil {
		return utils.SendError(c, fiber.StatusUnauthorized, "Authentication required")
	}

	ctx, ok := authCtx.(*middleware.AuthContext)
	if !ok || ctx.IsGuest {
		return utils.SendError(c, fiber.StatusUnauthorized, "Authentication required")
	}

	// Try to get organization
	org, err := h.service.GetOrganizationBySlug(slug)
	if err != nil {
		return utils.SendError(c, fiber.StatusNotFound, "Organization not found", err)
	}

	// Check if user has access to this organization
	hasAccess := false
	for _, membership := range ctx.Memberships {
		if membership.Organization.Slug == slug {
			hasAccess = true
			break
		}
	}

	if !hasAccess {
		return utils.SendError(c, fiber.StatusForbidden, "Access denied to this organization")
	}

	return utils.SendSuccess(c, fiber.StatusOK, ToOrganizationResponseDTO(org))
}

// ListOrganizations is the handler for listing all organizations.
func (h *OrganizationHandler) ListOrganizations(c *fiber.Ctx) error {
	orgs, err := h.service.GetAllOrganizations()
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to retrieve organizations", err)
	}

	// Convert models to DTOs for public response
	var orgDTOs []OrganizationResponseDTO
	for _, org := range orgs {
		orgDTOs = append(orgDTOs, ToOrganizationResponseDTO(&org))
	}

	return utils.SendSuccess(c, fiber.StatusOK, orgDTOs)
}

// UpdateOrganization is the handler for updating an organization's details.
func (h *OrganizationHandler) UpdateOrganization(c *fiber.Ctx) error {
	// This route should be protected and have the org context.
	// For now, we get the slug from the URL params.
	slug := c.Params("slug")
	if slug == "" {
		return utils.SendError(c, fiber.StatusBadRequest, "Organization slug is required")
	}

	var dto UpdateOrganizationDTO
	if err := c.BodyParser(&dto); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	org, err := h.service.UpdateOrganization(slug, &dto)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to update organization", err)
	}

	return utils.SendSuccess(c, fiber.StatusOK, ToOrganizationResponseDTO(org), "Organization updated successfully")
}

// DeleteOrganization is the handler for deleting an organization.
func (h *OrganizationHandler) DeleteOrganization(c *fiber.Ctx) error {
	slug := c.Params("slug")
	if slug == "" {
		return utils.SendError(c, fiber.StatusBadRequest, "Organization slug is required")
	}

	err := h.service.DeleteOrganization(slug)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to delete organization", err)
	}

	return utils.SendSuccess(c, fiber.StatusNoContent, nil)
}

// RemoveMember is the handler for removing a user from an organization.
func (h *OrganizationHandler) RemoveMember(c *fiber.Ctx) error {
	orgSlug := c.Params("slug")
	if orgSlug == "" {
		return utils.SendError(c, fiber.StatusBadRequest, "Organization slug is required")
	}
	userID := c.Params("userId")
	if userID == "" {
		return utils.SendError(c, fiber.StatusBadRequest, "User ID is required")
	}

	// In a real app, you'd get the org ID from the slug
	// For now, we'll assume the service can handle the slug.
	// A better approach would be a middleware to load the org from the slug.
	if err := h.service.RemoveMember(orgSlug, userID); err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to remove member", err)
	}

	return utils.SendSuccess(c, fiber.StatusNoContent, nil)
}

// AddMember is the handler for adding a user to an organization.
func (h *OrganizationHandler) AddMember(c *fiber.Ctx) error {
	orgSlug := c.Params("slug")
	if orgSlug == "" {
		return utils.SendError(c, fiber.StatusBadRequest, "Organization slug is required")
	}

	var dto struct {
		UserID string `json:"user_id" validate:"required,uuid"`
		RoleID string `json:"role_id" validate:"required,uuid"`
	}

	if err := c.BodyParser(&dto); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	if errs := utils.ValidateStruct(&dto); errs != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": errs})
	}

	if err := h.service.AddMember(orgSlug, dto.UserID, dto.RoleID); err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to add member", err)
	}

	return utils.SendSuccess(c, fiber.StatusCreated, nil, "Member added successfully")
}

// ListMembers is the handler for listing organization members.
func (h *OrganizationHandler) ListMembers(c *fiber.Ctx) error {
	orgSlug := c.Params("slug")
	if orgSlug == "" {
		return utils.SendError(c, fiber.StatusBadRequest, "Organization slug is required")
	}

	members, err := h.service.ListMembers(orgSlug)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to retrieve members", err)
	}

	// Convert to DTOs
	var memberDTOs []MemberResponseDTO
	for _, member := range members {
		memberDTOs = append(memberDTOs, ToMemberResponseDTO(&member))
	}

	return utils.SendSuccess(c, fiber.StatusOK, memberDTOs)
}

// ChangeUserRole is the handler for changing a user's role in an organization.
func (h *OrganizationHandler) ChangeUserRole(c *fiber.Ctx) error {
	orgSlug := c.Params("slug")
	userID := c.Params("userId")
	if orgSlug == "" || userID == "" {
		return utils.SendError(c, fiber.StatusBadRequest, "Organization slug and user ID are required")
	}

	var dto struct {
		NewRoleID string `json:"new_role_id" validate:"required,uuid"`
	}

	if err := c.BodyParser(&dto); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	if errs := utils.ValidateStruct(&dto); errs != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": errs})
	}

	if err := h.service.ChangeUserRole(orgSlug, userID, dto.NewRoleID); err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to change user role", err)
	}

	return utils.SendSuccess(c, fiber.StatusOK, nil, "User role changed successfully")
}
