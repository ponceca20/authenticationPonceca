package user

import (
	"practicev2/module/authentication/models"

	"gorm.io/gorm"
)

// UserRepository defines the interface for database operations related to organizational users.
type UserRepository interface {
	FindMembershipsByOrganizationID(orgID string) ([]models.OrganizationalMembership, error)
	FindMembershipByUserInOrganization(orgID, identityID string) (*models.OrganizationalMembership, error)
	UpdateMembership(membership *models.OrganizationalMembership) error
	BulkCreateUsers(identities []*models.Identity, memberships []*models.OrganizationalMembership) error
	BulkUpdateMemberships(orgID string, users []BulkUpdateUserItemDTO) error
	BulkDeleteMemberships(orgID string, userIDs []string) error
}

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new instance of UserRepository.
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// FindMembershipsByOrganizationID retrieves all memberships for a given organization ID.
// It preloads the associated Identity and Role for each membership.
func (r *userRepository) FindMembershipsByOrganizationID(orgID string) ([]models.OrganizationalMembership, error) {
	var memberships []models.OrganizationalMembership
	err := r.db.
		Preload("Identity").
		Preload("Role").
		Where("organization_id = ?", orgID).
		Find(&memberships).Error

	if err != nil {
		return nil, err
	}
	return memberships, nil
}

// FindMembershipByUserInOrganization retrieves a specific user's membership within an organization.
func (r *userRepository) FindMembershipByUserInOrganization(orgID, identityID string) (*models.OrganizationalMembership, error) {
	var membership models.OrganizationalMembership
	err := r.db.
		Preload("Identity").
		Preload("Role").
		Where("organization_id = ? AND identity_id = ?", orgID, identityID).
		First(&membership).Error

	if err != nil {
		return nil, err
	}
	return &membership, nil
}

// UpdateMembership saves changes to a membership model.
func (r *userRepository) UpdateMembership(membership *models.OrganizationalMembership) error {
	return r.db.Save(membership).Error
}

// BulkCreateUsers creates multiple users and their memberships in a transaction.
func (r *userRepository) BulkCreateUsers(identities []*models.Identity, memberships []*models.OrganizationalMembership) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.CreateInBatches(identities, 100).Error; err != nil {
			return err
		}
		if err := tx.CreateInBatches(memberships, 100).Error; err != nil {
			return err
		}
		return nil
	})
}

// BulkUpdateMemberships updates multiple memberships in a transaction.
func (r *userRepository) BulkUpdateMemberships(orgID string, users []BulkUpdateUserItemDTO) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, user := range users {
			updates := make(map[string]interface{})
			if user.RoleID != "" {
				updates["role_id"] = user.RoleID
			}
			if user.Department != "" {
				updates["department"] = user.Department
			}
			if user.IsActive != nil {
				updates["is_active"] = *user.IsActive
			}

			if len(updates) > 0 {
				if err := tx.Model(&models.OrganizationalMembership{}).Where("organization_id = ? AND identity_id = ?", orgID, user.ID).Updates(updates).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}

// BulkDeleteMemberships soft deletes multiple memberships.
func (r *userRepository) BulkDeleteMemberships(orgID string, userIDs []string) error {
	return r.db.Where("organization_id = ? AND identity_id IN ?", orgID, userIDs).Delete(&models.OrganizationalMembership{}).Error
}
