package organization

import (
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
	org, err := h.service.GetOrganizationBySlug(slug)
	if err != nil {
		return utils.SendError(c, fiber.StatusNotFound, "Organization not found", err)
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
