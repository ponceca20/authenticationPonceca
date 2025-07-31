package role

import (
	"errors"
	"fmt"
	"practicev2/module/authentication/models"
	"practicev2/module/authentication/user"
	"practicev2/module/authentication/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RoleService defines the interface for role-related business logic.
type RoleService interface {
	CreateRole(orgID string, dto *RoleDTO) (*models.Role, error)
	ListRoles(orgID string) ([]RoleResponseDTO, error)
	UpdateRole(orgID, roleID string, dto *RoleDTO) (*models.Role, error)
	DeleteRole(orgID, roleID string) error
	ListUsersInRole(orgID, roleID string) ([]user.UserDTO, error)
}

type roleService struct {
	repo RoleRepository
}

// NewRoleService creates a new instance of RoleService.
func NewRoleService(repo RoleRepository) RoleService {
	return &roleService{repo: repo}
}

// CreateRole handles the logic for creating a new role and its permissions.
func (s *roleService) CreateRole(orgID string, dto *RoleDTO) (*models.Role, error) {
	if errs := utils.ValidateStruct(dto); errs != nil {
		return nil, fmt.Errorf("validation failed: %v", errs)
	}

	if _, err := s.repo.FindRoleByName(orgID, dto.Name); !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("a role with this name already exists in the organization")
	}

	// Flatten permissions from DTO to a list of model objects
	var permissionsToCreate []*models.Permission
	for _, pDTO := range dto.Permissions {
		for _, action := range pDTO.Actions {
			permissionsToCreate = append(permissionsToCreate, &models.Permission{
				Resource: pDTO.Resource,
				Action:   action,
				Scope:    pDTO.Scope,
			})
		}
	}

	// Find existing permissions or create new ones to avoid duplicates
	dbPermissions, err := s.repo.FindOrCreatePermissions(permissionsToCreate)
	if err != nil {
		return nil, fmt.Errorf("failed to find or create permissions: %w", err)
	}

	// Prepare the role model
	role := &models.Role{
		ID:             uuid.New().String(),
		OrganizationID: orgID,
		Name:           dto.Name,
		DisplayName:    dto.DisplayName,
		Description:    dto.Description,
		HierarchyLevel: dto.HierarchyLevel,
		IsSystemRole:   false, // Only system can create system roles
	}

	// Create the role and its associations in a transaction
	if err := s.repo.CreateRole(role, dbPermissions); err != nil {
		return nil, fmt.Errorf("failed to create role: %w", err)
	}

	return role, nil
}

// ListUsersInRole retrieves all users assigned to a specific role.
func (s *roleService) ListUsersInRole(orgID, roleID string) ([]user.UserDTO, error) {
	memberships, err := s.repo.FindMembershipsByRole(orgID, roleID)
	if err != nil {
		return nil, err
	}

	var userDTOs []user.UserDTO
	for _, m := range memberships {
		if m.Identity.ID != "" {
			userDTOs = append(userDTOs, user.ToUserDTO(&m))
		}
	}

	return userDTOs, nil
}

// DeleteRole handles the logic for deleting a role.
func (s *roleService) DeleteRole(orgID, roleID string) error {
	role, err := s.repo.FindRoleByID(orgID, roleID)
	if err != nil {
		return errors.New("role not found")
	}

	// Prevent deletion of system roles
	if role.IsSystemRole {
		return errors.New("cannot delete a system role")
	}

	// TODO: Add logic to check if any user is currently assigned this role.
	// If so, prevent deletion or re-assign them to a default role.

	return s.repo.DeleteRole(role)
}

// ListRoles retrieves all roles for an organization.
func (s *roleService) ListRoles(orgID string) ([]RoleResponseDTO, error) {
	roles, err := s.repo.ListRolesByOrganizationID(orgID)
	if err != nil {
		return nil, err
	}

	var dtos []RoleResponseDTO
	for _, role := range roles {
		dtos = append(dtos, ToRoleResponseDTO(&role))
	}

	return dtos, nil
}

// UpdateRole handles the logic for updating a role and its permissions.
func (s *roleService) UpdateRole(orgID, roleID string, dto *RoleDTO) (*models.Role, error) {
	if errs := utils.ValidateStruct(dto); errs != nil {
		return nil, fmt.Errorf("validation failed: %v", errs)
	}

	role, err := s.repo.FindRoleByID(orgID, roleID)
	if err != nil {
		return nil, errors.New("role not found")
	}

	// Prevent modification of system roles
	if role.IsSystemRole {
		return nil, errors.New("cannot modify a system role")
	}

	// Flatten permissions from DTO
	var permissionsToSet []*models.Permission
	for _, pDTO := range dto.Permissions {
		for _, action := range pDTO.Actions {
			permissionsToSet = append(permissionsToSet, &models.Permission{
				Resource: pDTO.Resource,
				Action:   action,
				Scope:    pDTO.Scope,
			})
		}
	}

	dbPermissions, err := s.repo.FindOrCreatePermissions(permissionsToSet)
	if err != nil {
		return nil, fmt.Errorf("failed to find or create permissions: %w", err)
	}

	// Update role fields
	role.Name = dto.Name
	role.DisplayName = dto.DisplayName
	role.Description = dto.Description
	role.HierarchyLevel = dto.HierarchyLevel

	if err := s.repo.UpdateRole(role, dbPermissions); err != nil {
		return nil, fmt.Errorf("failed to update role: %w", err)
	}

	return role, nil
}
