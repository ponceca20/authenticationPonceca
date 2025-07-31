package invitation

import (
	"practicev2/module/authentication/middleware"
	"practicev2/module/authentication/utils"

	"github.com/gofiber/fiber/v2"
)


// InvitationHandler handles the HTTP requests for invitations.
type InvitationHandler struct {
	service InvitationService
}

// NewInvitationHandler creates a new instance of InvitationHandler.
func NewInvitationHandler(service InvitationService) *InvitationHandler {
	return &InvitationHandler{service: service}
}

// CreateInvitation is the handler for sending a new invitation.
func (h *InvitationHandler) CreateInvitation(c *fiber.Ctx) error {
	authCtx, ok := c.Locals("authContext").(*middleware.AuthContext)
	if !ok || authCtx.IsGuest || authCtx.CurrentOrg == nil {
		return utils.SendError(c, fiber.StatusUnauthorized, "Authentication and organization context required")
	}

	// TODO: Add permission check: if !authCtx.HasPermission("invitations.create") ...

	var dto InvitationDTO
	if err := c.BodyParser(&dto); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	if errs := utils.ValidateStruct(&dto); errs != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": errs})
	}

	invitation, err := h.service.CreateInvitation(authCtx.CurrentOrg.ID, authCtx.Identity.ID, &dto)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to create invitation", err)
	}

	response := InvitationResponseDTO{
		ID:        invitation.ID,
		Email:     invitation.Email,
		RoleID:    invitation.RoleID,
		Status:    invitation.Status,
		ExpiresAt: invitation.ExpiresAt.String(),
	}

	return utils.SendSuccess(c, fiber.StatusCreated, response, "Invitation sent successfully")
}

// AcceptInvitation is the handler for accepting an invitation.
func (h *InvitationHandler) AcceptInvitation(c *fiber.Ctx) error {
	var dto AcceptInvitationDTO
	if err := c.BodyParser(&dto); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	// The token from the URL is the primary identifier, but we can also get it from the DTO.
	// Let's assume the token is in the DTO for a cleaner API.
	// dto.Token = c.Params("token")

	if errs := utils.ValidateStruct(&dto); errs != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": errs})
	}

	if err := h.service.AcceptInvitation(&dto); err != nil {
		return utils.SendError(c, fiber.StatusConflict, "Failed to accept invitation", err)
	}

	return utils.SendSuccess(c, fiber.StatusOK, nil, "Invitation accepted successfully. You can now log in.")
}

// ListInvitations lists all pending invitations for an organization.
func (h *InvitationHandler) ListInvitations(c *fiber.Ctx) error {
	authCtx, ok := c.Locals("authContext").(*middleware.AuthContext)
	if !ok || authCtx.CurrentOrg == nil {
		return utils.SendError(c, fiber.StatusForbidden, "Organization context not found")
	}

	invitations, err := h.service.ListInvitations(authCtx.CurrentOrg.ID)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to retrieve invitations", err)
	}

	var dtos []InvitationResponseDTO
	for _, inv := range invitations {
		dtos = append(dtos, InvitationResponseDTO{
			ID:        inv.ID,
			Email:     inv.Email,
			RoleID:    inv.RoleID,
			Status:    inv.Status,
			ExpiresAt: inv.ExpiresAt.String(),
		})
	}

	return utils.SendSuccess(c, fiber.StatusOK, dtos)
}

// CancelInvitation cancels a pending invitation.
func (h *InvitationHandler) CancelInvitation(c *fiber.Ctx) error {
	authCtx, ok := c.Locals("authContext").(*middleware.AuthContext)
	if !ok || authCtx.CurrentOrg == nil {
		return utils.SendError(c, fiber.StatusForbidden, "Organization context not found")
	}

	invID := c.Params("id")
	if invID == "" {
		return utils.SendError(c, fiber.StatusBadRequest, "Invitation ID is required")
	}

	if err := h.service.CancelInvitation(authCtx.CurrentOrg.ID, invID); err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to cancel invitation", err)
	}

	return utils.SendSuccess(c, fiber.StatusNoContent, nil)
}

// GetInvitation retrieves a single invitation.
func (h *InvitationHandler) GetInvitation(c *fiber.Ctx) error {
	authCtx, ok := c.Locals("authContext").(*middleware.AuthContext)
	if !ok || authCtx.CurrentOrg == nil {
		return utils.SendError(c, fiber.StatusForbidden, "Organization context not found")
	}

	invID := c.Params("id")
	if invID == "" {
		return utils.SendError(c, fiber.StatusBadRequest, "Invitation ID is required")
	}

	invitation, err := h.service.GetInvitation(authCtx.CurrentOrg.ID, invID)
	if err != nil {
		return utils.SendError(c, fiber.StatusNotFound, "Invitation not found", err)
	}

	response := InvitationResponseDTO{
		ID:        invitation.ID,
		Email:     invitation.Email,
		RoleID:    invitation.RoleID,
		Status:    invitation.Status,
		ExpiresAt: invitation.ExpiresAt.String(),
	}

	return utils.SendSuccess(c, fiber.StatusOK, response)
}
