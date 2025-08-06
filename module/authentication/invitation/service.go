package invitation

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"practicev2/module/authentication/models"
	"practicev2/module/authentication/utils"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// InvitationService defines the interface for invitation business logic.
type InvitationService interface {
	CreateInvitation(orgID, inviterID string, dto *InvitationDTO) (*models.Invitation, error)
	AcceptInvitation(dto *AcceptInvitationDTO) error
	ListInvitations(orgID string) ([]models.Invitation, error)
	CancelInvitation(orgID, invID string) error
	GetInvitation(orgID, invID string) (*models.Invitation, error)
	ResendInvitation(orgID, invID string) error
	VerifyInvitationToken(token string) (*models.Invitation, error)
	AcceptInvitationByToken(token string, dto *AcceptInvitationByTokenDTO) error
}

type invitationService struct {
	repo InvitationRepository
}

// NewInvitationService creates a new instance of InvitationService.
func NewInvitationService(repo InvitationRepository) InvitationService {
	return &invitationService{repo: repo}
}

// CreateInvitation handles the logic for creating a new user invitation.
func (s *invitationService) CreateInvitation(orgID, inviterID string, dto *InvitationDTO) (*models.Invitation, error) {
	// Validate input
	if orgID == "" || inviterID == "" {
		return nil, fmt.Errorf("organization ID and inviter ID are required")
	}

	if dto.RoleID == "" {
		return nil, fmt.Errorf("role ID is required")
	}

	if dto.Email == "" {
		return nil, fmt.Errorf("email is required")
	}

	// Validate that the role belongs to the organization
	_, err := s.repo.ValidateRoleInOrganization(orgID, dto.RoleID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("role not found in organization")
		}
		return nil, fmt.Errorf("failed to validate role: %w", err)
	}

	// Check if user already exists in the organization
	// This should be implemented in the repository if needed

	token, err := generateSecureToken(32)
	if err != nil {
		return nil, fmt.Errorf("failed to generate invitation token: %w", err)
	}

	invitation := &models.Invitation{
		ID:             uuid.New().String(),
		OrganizationID: orgID,
		InviterID:      inviterID,
		Email:          dto.Email,
		RoleID:         dto.RoleID,
		Token:          token,
		Status:         "pending",
		ExpiresAt:      time.Now().Add(7 * 24 * time.Hour), // Invitation expires in 7 days
	}

	if err := s.repo.CreateInvitation(invitation); err != nil {
		return nil, fmt.Errorf("failed to save invitation: %w", err)
	}

	// In a real application, you would send an email to the user here
	// with a link like: https://yourapp.com/accept-invitation?token=...

	return invitation, nil
}

// GetInvitation retrieves a single invitation by its ID.
func (s *invitationService) GetInvitation(orgID, invID string) (*models.Invitation, error) {
	return s.repo.FindInvitationByID(orgID, invID)
}

// ListInvitations lists all pending invitations for an organization.
func (s *invitationService) ListInvitations(orgID string) ([]models.Invitation, error) {
	return s.repo.ListInvitationsByOrganization(orgID)
}

// CancelInvitation cancels a pending invitation.
func (s *invitationService) CancelInvitation(orgID, invID string) error {
	invitation, err := s.repo.FindInvitationByID(orgID, invID)
	if err != nil {
		return fmt.Errorf("invitation not found: %w", err)
	}

	// In a real app, you might want to change the status to 'canceled'
	// instead of hard deleting, for auditing.
	return s.repo.DeleteInvitation(invitation)
}

// AcceptInvitation handles the logic for a user accepting an invitation.
func (s *invitationService) AcceptInvitation(dto *AcceptInvitationDTO) error {
	// Find the invitation
	invitation, err := s.repo.FindInvitationByToken(dto.Token)
	if err != nil {
		return fmt.Errorf("invalid or expired invitation token")
	}

	// Hash the user's password
	hashedPassword, err := utils.HashPassword(dto.Password)
	if err != nil {
		return fmt.Errorf("failed to hash password")
	}

	// Prepare the new identity and membership
	identity := &models.Identity{
		ID:              uuid.New().String(),
		Email:           invitation.Email,
		FirstName:       dto.FirstName,
		LastName:        dto.LastName,
		PasswordHash:    hashedPassword,
		EmailVerified:   true, // Email is verified by the act of receiving the invitation
		EmailVerifiedAt: &time.Time{},
	}
	*identity.EmailVerifiedAt = time.Now()

	membership := &models.OrganizationalMembership{
		ID:         uuid.New().String(),
		IsActive:   true,
		ActiveFrom: time.Now(),
	}

	// Finalize the acceptance in a transaction
	return s.repo.AcceptInvitation(invitation, identity, membership)
}

// ResendInvitation resends an existing invitation by generating a new token.
func (s *invitationService) ResendInvitation(orgID, invID string) error {
	invitation, err := s.repo.FindInvitationByID(orgID, invID)
	if err != nil {
		return fmt.Errorf("invitation not found: %w", err)
	}

	// Check if invitation is still valid
	if invitation.Status != "pending" {
		return fmt.Errorf("invitation is no longer pending")
	}

	// Generate new token and extend expiry
	newToken, err := generateSecureToken(32)
	if err != nil {
		return fmt.Errorf("failed to generate new token: %w", err)
	}

	invitation.Token = newToken
	invitation.ExpiresAt = time.Now().Add(7 * 24 * time.Hour) // Extend for another 7 days

	if err := s.repo.UpdateInvitation(invitation); err != nil {
		return fmt.Errorf("failed to update invitation: %w", err)
	}

	// In a real application, you would send a new email here
	return nil
}

// VerifyInvitationToken validates an invitation token and returns invitation info.
func (s *invitationService) VerifyInvitationToken(token string) (*models.Invitation, error) {
	invitation, err := s.repo.FindInvitationByToken(token)
	if err != nil {
		return nil, fmt.Errorf("invalid or expired invitation token")
	}

	// Check if invitation is still pending and not expired
	if invitation.Status != "pending" {
		return nil, fmt.Errorf("invitation is no longer pending")
	}

	if time.Now().After(invitation.ExpiresAt) {
		return nil, fmt.Errorf("invitation has expired")
	}

	return invitation, nil
}

// AcceptInvitationByToken accepts an invitation using token from URL parameters.
func (s *invitationService) AcceptInvitationByToken(token string, dto *AcceptInvitationByTokenDTO) error {
	// Verify the token first
	invitation, err := s.VerifyInvitationToken(token)
	if err != nil {
		return err
	}

	// Hash the user's password
	hashedPassword, err := utils.HashPassword(dto.Password)
	if err != nil {
		return fmt.Errorf("failed to hash password")
	}

	// Prepare the new identity and membership
	identity := &models.Identity{
		ID:              uuid.New().String(),
		Email:           invitation.Email,
		FirstName:       dto.FirstName,
		LastName:        dto.LastName,
		PasswordHash:    hashedPassword,
		EmailVerified:   true, // Email is verified by the act of receiving the invitation
		EmailVerifiedAt: &time.Time{},
	}
	*identity.EmailVerifiedAt = time.Now()

	membership := &models.OrganizationalMembership{
		ID:         uuid.New().String(),
		IsActive:   true,
		ActiveFrom: time.Now(),
	}

	// Finalize the acceptance in a transaction
	return s.repo.AcceptInvitation(invitation, identity, membership)
}

// generateSecureToken creates a random, URL-safe string.
func generateSecureToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
