package organization

import (
	"practicev2/module/authentication/models"

	"gorm.io/gorm"
)

// OrganizationRepository defines the interface for database operations related to organizations.
type OrganizationRepository interface {
	CreateOrganizationInTransaction(org *models.Organization, founder *models.Identity, role *models.Role, membership *models.OrganizationalMembership) error
	FindOrganizationBySlug(slug string) (*models.Organization, error)
	FindOrganizationByName(name string) (*models.Organization, error)
	ListOrganizations() ([]models.Organization, error)
	UpdateOrganization(org *models.Organization) error
	DeleteOrganization(org *models.Organization) error
	RemoveMembership(orgID, userID string) error
}

type organizationRepository struct {
	db *gorm.DB
}

// NewOrganizationRepository creates a new instance of OrganizationRepository.
func NewOrganizationRepository(db *gorm.DB) OrganizationRepository {
	return &organizationRepository{db: db}
}

// CreateOrganizationInTransaction creates an organization, its founder, the admin role, and the membership link in a single DB transaction.
func (r *organizationRepository) CreateOrganizationInTransaction(org *models.Organization, founder *models.Identity, role *models.Role, membership *models.OrganizationalMembership) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Create the organization
		if err := tx.Create(org).Error; err != nil {
			return err
		}

		// 2. Create the founder's identity
		if err := tx.Create(founder).Error; err != nil {
			return err
		}

		// 3. Create the role for the new organization
		role.OrganizationID = org.ID
		if err := tx.Create(role).Error; err != nil {
			return err
		}

		// 4. Create the membership linking founder, org, and role
		membership.IdentityID = founder.ID
		membership.OrganizationID = org.ID
		membership.RoleID = role.ID
		if err := tx.Create(membership).Error; err != nil {
			return err
		}

		return nil
	})
}

// FindOrganizationBySlug finds an organization by its unique slug.
func (r *organizationRepository) FindOrganizationBySlug(slug string) (*models.Organization, error) {
	var org models.Organization
	if err := r.db.Where("slug = ?", slug).First(&org).Error; err != nil {
		return nil, err
	}
	return &org, nil
}

// FindOrganizationByName finds an organization by its name.
func (r *organizationRepository) FindOrganizationByName(name string) (*models.Organization, error) {
	var org models.Organization
	if err := r.db.Where("name = ?", name).First(&org).Error; err != nil {
		return nil, err
	}
	return &org, nil
}

// ListOrganizations retrieves all organizations from the database.
func (r *organizationRepository) ListOrganizations() ([]models.Organization, error) {
	var orgs []models.Organization
	if err := r.db.Find(&orgs).Error; err != nil {
		return nil, err
	}
	return orgs, nil
}

// UpdateOrganization saves the changes to an organization model.
func (r *organizationRepository) UpdateOrganization(org *models.Organization) error {
	return r.db.Save(org).Error
}

// DeleteOrganization performs a soft delete on an organization.
func (r *organizationRepository) DeleteOrganization(org *models.Organization) error {
	return r.db.Delete(org).Error
}

// RemoveMembership performs a soft delete on an organizational membership.
func (r *organizationRepository) RemoveMembership(orgID, userID string) error {
	return r.db.Where("organization_id = ? AND identity_id = ?", orgID, userID).Delete(&models.OrganizationalMembership{}).Error
}
