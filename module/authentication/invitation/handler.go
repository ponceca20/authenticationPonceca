package invitation

import (
	"practicev2/module/authentication/middleware"
	"practicev2/module/authentication/utils"
	"strings"

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
		// Handle specific error types
		errMsg := err.Error()

		// Role not found
		if strings.Contains(errMsg, "role not found") || strings.Contains(errMsg, "invalid role") {
			return utils.SendError(c, fiber.StatusBadRequest, "Invalid role specified", err)
		}

		// User already exists
		if strings.Contains(errMsg, "user already exists") || strings.Contains(errMsg, "email already exists") {
			return utils.SendError(c, fiber.StatusConflict, "User with this email already exists", err)
		}

		// Validation errors
		if strings.Contains(errMsg, "validation failed") {
			return utils.SendError(c, fiber.StatusBadRequest, "Validation failed", err)
		}

		// Default to internal server error
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
		errMsg := err.Error()
		if strings.Contains(errMsg, "invitation not found") || strings.Contains(errMsg, "record not found") || strings.Contains(errMsg, "not found") {
			return utils.SendError(c, fiber.StatusNotFound, "Invitation not found")
		}
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

// ResendInvitation resends an existing invitation with a new token.
func (h *InvitationHandler) ResendInvitation(c *fiber.Ctx) error {
	authCtx := c.Locals("authContext")
	if authCtx == nil {
		return utils.SendError(c, fiber.StatusUnauthorized, "Authentication required")
	}

	ctx, ok := authCtx.(*middleware.AuthContext)
	if !ok || ctx.IsGuest {
		return utils.SendError(c, fiber.StatusUnauthorized, "Authentication required")
	}

	if ctx.CurrentOrg == nil {
		return utils.SendError(c, fiber.StatusForbidden, "Organization context not found")
	}

	invID := c.Params("id")
	if invID == "" {
		return utils.SendError(c, fiber.StatusBadRequest, "Invitation ID is required")
	}

	if err := h.service.ResendInvitation(ctx.CurrentOrg.ID, invID); err != nil {
		errMsg := err.Error()
		if strings.Contains(errMsg, "invitation not found") || strings.Contains(errMsg, "not found") {
			return utils.SendError(c, fiber.StatusNotFound, "Invitation not found")
		}
		if strings.Contains(errMsg, "no longer pending") {
			return utils.SendError(c, fiber.StatusBadRequest, "Invitation is no longer pending")
		}
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to resend invitation", err)
	}

	return utils.SendSuccess(c, fiber.StatusOK, nil, "Invitation resent successfully")
}

// VerifyInvitationToken verifies an invitation token and returns invitation info.
func (h *InvitationHandler) VerifyInvitationToken(c *fiber.Ctx) error {
	token := c.Params("token")
	if token == "" {
		return utils.SendError(c, fiber.StatusBadRequest, "Invitation token is required")
	}

	invitation, err := h.service.VerifyInvitationToken(token)
	if err != nil {
		errMsg := err.Error()
		if strings.Contains(errMsg, "invalid") || strings.Contains(errMsg, "expired") {
			return utils.SendError(c, fiber.StatusBadRequest, err.Error())
		}
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to verify invitation", err)
	}

	// Return safe information about the invitation
	response := map[string]interface{}{
		"organization_name": invitation.Organization.Name,
		"role_name":         invitation.Role.DisplayName,
		"email":             invitation.Email,
		"expires_at":        invitation.ExpiresAt,
	}

	return utils.SendSuccess(c, fiber.StatusOK, response)
}

// AcceptInvitationByToken accepts an invitation using token from URL parameters.
func (h *InvitationHandler) AcceptInvitationByToken(c *fiber.Ctx) error {
	token := c.Params("token")
	if token == "" {
		return utils.SendError(c, fiber.StatusBadRequest, "Invitation token is required")
	}

	var dto AcceptInvitationByTokenDTO
	if err := c.BodyParser(&dto); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	if errs := utils.ValidateStruct(&dto); errs != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": errs})
	}

	if err := h.service.AcceptInvitationByToken(token, &dto); err != nil {
		errMsg := err.Error()
		if strings.Contains(errMsg, "invalid") || strings.Contains(errMsg, "expired") {
			return utils.SendError(c, fiber.StatusBadRequest, err.Error())
		}
		return utils.SendError(c, fiber.StatusConflict, "Failed to accept invitation", err)
	}

	return utils.SendSuccess(c, fiber.StatusOK, nil, "Invitation accepted successfully. You can now log in.")
}
