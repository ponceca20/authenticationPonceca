package user

import (
	"fmt"
	"practicev2/module/authentication/models"
	"practicev2/module/authentication/utils"
	"time"

	"github.com/google/uuid"
)

// UserService defines the interface for user-related business logic.
type UserService interface {
	ListOrgUsers(orgID string) ([]UserDTO, error)
	GetOrgUser(orgID, identityID string) (*UserDTO, error)
	UpdateUser(orgID, identityID string, dto *UpdateUserDTO) (*UserDTO, error)
	UpdateUserStatus(orgID, identityID string, dto *UpdateUserStatusDTO) error
	AssignRoleToUser(orgID, identityID, roleID string) error
	BulkCreateUsers(orgID string, dto *BulkCreateUserDTO) error
	BulkUpdateUsers(orgID string, dto *BulkUpdateUserDTO) error
	BulkDeleteUsers(orgID string, dto *BulkDeleteUserDTO) error
}

type userService struct {
	repo UserRepository
}

// NewUserService creates a new instance of UserService.
func NewUserService(repo UserRepository) UserService {
	return &userService{repo: repo}
}

// ListOrgUsers retrieves all users for an organization and returns them as DTOs.
func (s *userService) ListOrgUsers(orgID string) ([]UserDTO, error) {
	memberships, err := s.repo.FindMembershipsByOrganizationID(orgID)
	if err != nil {
		return nil, err
	}

	var userDTOs []UserDTO
	for _, m := range memberships {
		// Ensure Identity and Role are not nil before accessing them
		if m.Identity.ID != "" && m.Role.ID != "" {
			userDTOs = append(userDTOs, ToUserDTO(&m))
		}
	}

	return userDTOs, nil
}

// GetOrgUser retrieves a specific user in an organization and returns it as a DTO.
func (s *userService) GetOrgUser(orgID, identityID string) (*UserDTO, error) {
	membership, err := s.repo.FindMembershipByUserInOrganization(orgID, identityID)
	if err != nil {
		return nil, err
	}

	dto := ToUserDTO(membership)
	return &dto, nil
}

// UpdateUser updates a user's role or department within an organization.
func (s *userService) UpdateUser(orgID, identityID string, dto *UpdateUserDTO) (*UserDTO, error) {
	membership, err := s.repo.FindMembershipByUserInOrganization(orgID, identityID)
	if err != nil {
		return nil, err
	}

	if dto.RoleID != "" {
		membership.RoleID = dto.RoleID
	}
	if dto.Department != "" {
		membership.Department = dto.Department
	}

	if err := s.repo.UpdateMembership(membership); err != nil {
		return nil, err
	}

	// We need to refetch the membership to get the updated Role details if it changed
	updatedMembership, err := s.repo.FindMembershipByUserInOrganization(orgID, identityID)
	if err != nil {
		return nil, err
	}

	responseDTO := ToUserDTO(updatedMembership)
	return &responseDTO, nil
}

// UpdateUserStatus updates a user's active status within an organization.
func (s *userService) UpdateUserStatus(orgID, identityID string, dto *UpdateUserStatusDTO) error {
	membership, err := s.repo.FindMembershipByUserInOrganization(orgID, identityID)
	if err != nil {
		return err
	}

	membership.IsActive = dto.IsActive
	// Optionally set ActiveUntil if deactivating
	return s.repo.UpdateMembership(membership)
}

// AssignRoleToUser assigns a role to a user by updating their membership.
func (s *userService) AssignRoleToUser(orgID, identityID, roleID string) error {
	_, err := s.UpdateUser(orgID, identityID, &UpdateUserDTO{RoleID: roleID})
	return err
}

// BulkCreateUsers creates multiple users in an organization.
func (s *userService) BulkCreateUsers(orgID string, dto *BulkCreateUserDTO) error {
	var identities []*models.Identity
	var memberships []*models.OrganizationalMembership

	for _, userDTO := range dto.Users {
		hashedPassword, err := utils.HashPassword(userDTO.Password)
		if err != nil {
			// In a real app, you might want to collect errors and continue,
			// or fail the whole batch.
			return fmt.Errorf("failed to hash password for user %s: %w", userDTO.Email, err)
		}

		identityID := uuid.New().String()
		identities = append(identities, &models.Identity{
			ID:           identityID,
			FirstName:    userDTO.FirstName,
			LastName:     userDTO.LastName,
			Email:        userDTO.Email,
			PasswordHash: hashedPassword,
		})

		memberships = append(memberships, &models.OrganizationalMembership{
			ID:             uuid.New().String(),
			IdentityID:     identityID,
			OrganizationID: orgID,
			RoleID:         userDTO.RoleID,
			Department:     userDTO.Department,
			IsActive:       true,
			ActiveFrom:     time.Now(),
		})
	}

	return s.repo.BulkCreateUsers(identities, memberships)
}

// BulkUpdateUsers updates multiple users in an organization.
func (s *userService) BulkUpdateUsers(orgID string, dto *BulkUpdateUserDTO) error {
	// In a real app, you'd add more validation here, e.g.,
	// ensuring all user IDs in the DTO belong to the specified orgID.
	return s.repo.BulkUpdateMemberships(orgID, dto.Users)
}

// BulkDeleteUsers deletes multiple users from an organization.
func (s *userService) BulkDeleteUsers(orgID string, dto *BulkDeleteUserDTO) error {
	// Similarly, you'd validate that the user making the request has permission
	// to delete all the users in the list.
	return s.repo.BulkDeleteMemberships(orgID, dto.UserIDs)
}
