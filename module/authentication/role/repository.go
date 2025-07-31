package role

import (
	"practicev2/module/authentication/models"

	"gorm.io/gorm"
)

// RoleRepository defines the interface for database operations related to roles.
type RoleRepository interface {
	CreateRole(role *models.Role, permissions []*models.Permission) error
	FindRoleByID(orgID, roleID string) (*models.Role, error)
	FindRoleByName(orgID string, name string) (*models.Role, error)
	ListRolesByOrganizationID(orgID string) ([]models.Role, error)
	UpdateRole(role *models.Role, permissions []*models.Permission) error
	DeleteRole(role *models.Role) error
	FindOrCreatePermissions(permissionsToFind []*models.Permission) ([]*models.Permission, error)
	FindMembershipsByRole(orgID, roleID string) ([]models.OrganizationalMembership, error)
}

type roleRepository struct {
	db *gorm.DB
}

// NewRoleRepository creates a new instance of RoleRepository.
func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{db: db}
}

// CreateRole creates a new role and associates it with the given permissions in a transaction.
func (r *roleRepository) CreateRole(role *models.Role, permissions []*models.Permission) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(role).Error; err != nil {
			return err
		}

		if len(permissions) > 0 {
			if err := tx.Model(role).Association("Permissions").Append(permissions); err != nil {
				return err
			}
		}

		return nil
	})
}

// FindMembershipsByRole finds all memberships associated with a specific role.
func (r *roleRepository) FindMembershipsByRole(orgID, roleID string) ([]models.OrganizationalMembership, error) {
	var memberships []models.OrganizationalMembership
	err := r.db.
		Preload("Identity").
		Where("organization_id = ? AND role_id = ?", orgID, roleID).
		Find(&memberships).Error
	return memberships, err
}

// DeleteRole performs a soft delete on a role.
// It also clears the permissions association.
func (r *roleRepository) DeleteRole(role *models.Role) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(role).Association("Permissions").Clear(); err != nil {
			return err
		}
		return tx.Delete(role).Error
	})
}

// FindRoleByName finds a role by name within a specific organization.
func (r *roleRepository) FindRoleByName(orgID, name string) (*models.Role, error) {
	var role models.Role
	if err := r.db.Where("organization_id = ? AND name = ?", orgID, name).First(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

// ListRolesByOrganizationID retrieves all roles for a given organization, preloading their permissions.
func (r *roleRepository) ListRolesByOrganizationID(orgID string) ([]models.Role, error) {
	var roles []models.Role
	err := r.db.
		Preload("Permissions").
		Where("organization_id = ?", orgID).
		Find(&roles).Error
	return roles, err
}

// FindOrCreatePermissions finds permissions by their unique composite key (resource, action, scope)
// or creates them if they don't exist. This prevents duplicate permission entries.
func (r *roleRepository) FindOrCreatePermissions(permissionsToFind []*models.Permission) ([]*models.Permission, error) {
	var foundPermissions []*models.Permission
	for _, p := range permissionsToFind {
		var found models.Permission
		// First try to find existing permission
		err := r.db.Where("resource = ? AND action = ? AND scope = ?", p.Resource, p.Action, p.Scope).First(&found).Error
		if err == nil {
			// Found existing permission
			foundPermissions = append(foundPermissions, &found)
		} else {
			// Create new permission (ID will be auto-generated)
			newPerm := &models.Permission{
				Resource: p.Resource,
				Action:   p.Action,
				Scope:    p.Scope,
			}
			if err := r.db.Create(newPerm).Error; err != nil {
				return nil, err
			}
			foundPermissions = append(foundPermissions, newPerm)
		}
	}
	return foundPermissions, nil
}

// FindRoleByID finds a role by its ID within a specific organization.
func (r *roleRepository) FindRoleByID(orgID, roleID string) (*models.Role, error) {
	var role models.Role
	if err := r.db.Where("organization_id = ? AND id = ?", orgID, roleID).First(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

// UpdateRole updates a role and replaces its permissions in a transaction.
func (r *roleRepository) UpdateRole(role *models.Role, permissions []*models.Permission) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(role).Error; err != nil {
			return err
		}

		if err := tx.Model(role).Association("Permissions").Replace(permissions); err != nil {
			return err
		}

		return nil
	})
}
