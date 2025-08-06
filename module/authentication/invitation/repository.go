package invitation

import (
	"practicev2/module/authentication/models"
	"time"

	"gorm.io/gorm"
)

// InvitationRepository defines the interface for database operations related to invitations.
type InvitationRepository interface {
	CreateInvitation(invitation *models.Invitation) error
	FindInvitationByToken(token string) (*models.Invitation, error)
	AcceptInvitation(invitation *models.Invitation, identity *models.Identity, membership *models.OrganizationalMembership) error
	ListInvitationsByOrganization(orgID string) ([]models.Invitation, error)
	FindInvitationByID(orgID, invID string) (*models.Invitation, error)
	DeleteInvitation(invitation *models.Invitation) error
	ValidateRoleInOrganization(orgID, roleID string) (*models.Role, error)
	UpdateInvitation(invitation *models.Invitation) error
}

type invitationRepository struct {
	db *gorm.DB
}

// NewInvitationRepository creates a new instance of InvitationRepository.
func NewInvitationRepository(db *gorm.DB) InvitationRepository {
	return &invitationRepository{db: db}
}

// CreateInvitation saves a new invitation to the database.
func (r *invitationRepository) CreateInvitation(invitation *models.Invitation) error {
	return r.db.Create(invitation).Error
}

// FindInvitationByToken finds an invitation by its unique token.
// It also preloads the organization and role details.
func (r *invitationRepository) FindInvitationByToken(token string) (*models.Invitation, error) {
	var invitation models.Invitation
	err := r.db.
		Preload("Organization").
		Preload("Role").
		Where("token = ? AND status = 'pending' AND expires_at > ?", token, time.Now()).
		First(&invitation).Error
	return &invitation, err
}

// AcceptInvitation creates the user's membership and marks the invitation as accepted
// in a single database transaction.
func (r *invitationRepository) AcceptInvitation(invitation *models.Invitation, identity *models.Identity, membership *models.OrganizationalMembership) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Use FirstOrCreate to handle cases where the user might already exist
		// from a previous (but different) invitation or direct registration.
		if err := tx.Where(models.Identity{Email: identity.Email}).FirstOrCreate(identity).Error; err != nil {
			return err
		}

		// Create the membership link
		membership.IdentityID = identity.ID
		membership.OrganizationID = invitation.OrganizationID
		membership.RoleID = invitation.RoleID
		if err := tx.Create(membership).Error; err != nil {
			return err
		}

		// Mark invitation as accepted instead of deleting, for auditing purposes.
		invitation.Status = "accepted"
		if err := tx.Save(invitation).Error; err != nil {
			return err
		}

		return nil
	})
}

// ListInvitationsByOrganization lists all pending invitations for an organization.
func (r *invitationRepository) ListInvitationsByOrganization(orgID string) ([]models.Invitation, error) {
	var invitations []models.Invitation
	err := r.db.
		Preload("Role").
		Where("organization_id = ? AND status = 'pending'", orgID).
		Find(&invitations).Error
	return invitations, err
}

// FindInvitationByID finds a pending invitation by its ID.
func (r *invitationRepository) FindInvitationByID(orgID, invID string) (*models.Invitation, error) {
	var invitation models.Invitation
	err := r.db.
		Where("organization_id = ? AND id = ? AND status = 'pending'", orgID, invID).
		First(&invitation).Error
	return &invitation, err
}

// DeleteInvitation permanently deletes an invitation.
func (r *invitationRepository) DeleteInvitation(invitation *models.Invitation) error {
	// Use Unscoped() to perform a hard delete, as this is for a pending invitation.
	return r.db.Unscoped().Delete(invitation).Error
}

// ValidateRoleInOrganization checks if a role belongs to the specified organization.
func (r *invitationRepository) ValidateRoleInOrganization(orgID, roleID string) (*models.Role, error) {
	var role models.Role
	err := r.db.Where("id = ? AND organization_id = ?", roleID, orgID).First(&role).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &role, nil
}

// UpdateInvitation saves changes to an invitation.
func (r *invitationRepository) UpdateInvitation(invitation *models.Invitation) error {
	return r.db.Save(invitation).Error
}
