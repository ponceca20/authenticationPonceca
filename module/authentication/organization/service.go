package organization

import (
	"errors"
	"practicev2/module/authentication/auth"
	"practicev2/module/authentication/models"
	"practicev2/module/authentication/utils"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// OrganizationService defines the interface for organization business logic.
type OrganizationService interface {
	CreateOrganization(dto *OrganizationRegistrationDTO) (*models.Organization, error)
	GetOrganizationBySlug(slug string) (*models.Organization, error)
	GetAllOrganizations() ([]models.Organization, error)
	UpdateOrganization(orgID string, dto *UpdateOrganizationDTO) (*models.Organization, error)
	DeleteOrganization(orgID string) error
	RemoveMember(orgID, userID string) error
	AddMember(orgSlug, userID, roleID string) error
	ListMembers(orgSlug string) ([]models.OrganizationalMembership, error)
	ChangeUserRole(orgSlug, userID, newRoleID string) error
}

type organizationService struct {
	orgRepo  OrganizationRepository
	authRepo auth.AuthRepository // To check for existing users
}

// NewOrganizationService creates a new instance of OrganizationService.
func NewOrganizationService(orgRepo OrganizationRepository, authRepo auth.AuthRepository) OrganizationService {
	return &organizationService{orgRepo: orgRepo, authRepo: authRepo}
}

// CreateOrganization handles the business logic for creating a new organization.
func (s *organizationService) CreateOrganization(dto *OrganizationRegistrationDTO) (*models.Organization, error) {
	if err := dto.SanitizeAndValidate(); err != nil {
		return nil, err
	}

	// Check for duplicate organization name
	if _, err := s.orgRepo.FindOrganizationByName(dto.Name); !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("organization name already exists")
	}

	// Check for duplicate user email
	if _, err := s.authRepo.FindIdentityByEmail(dto.Identity.Email); !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("email is already registered")
	}

	// Hash founder's password
	hashedPassword, err := utils.HashPassword(dto.Identity.Password)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Prepare models
	org := &models.Organization{
		ID:          uuid.New().String(),
		Name:        dto.Name,
		Slug:        slugify(dto.Name),
		Type:        dto.Type,
		Description: dto.Description,
		Website:     dto.Website,
		Address:     dto.Address,
		Phone:       dto.Phone,
	}

	founder := &models.Identity{
		ID:           uuid.New().String(),
		FirstName:    dto.Identity.FirstName,
		LastName:     dto.Identity.LastName,
		Email:        dto.Identity.Email,
		PasswordHash: hashedPassword,
	}

	// Define the founder's role
	adminRole := &models.Role{
		ID:             uuid.New().String(),
		OrganizationID: org.ID,
		Name:           "owner",
		DisplayName:    "Owner",
		Description:    "Full access to the organization.",
		HierarchyLevel: 100,
		IsSystemRole:   true,
	}

	membership := &models.OrganizationalMembership{
		ID:         uuid.New().String(),
		IsActive:   true,
		ActiveFrom: time.Now(),
	}

	// Create everything in a transaction
	err = s.orgRepo.CreateOrganizationInTransaction(org, founder, adminRole, membership)
	if err != nil {
		return nil, errors.New("failed to create organization")
	}

	return org, nil
}

func (s *organizationService) GetOrganizationBySlug(slug string) (*models.Organization, error) {
	return s.orgRepo.FindOrganizationBySlug(slug)
}

func (s *organizationService) GetAllOrganizations() ([]models.Organization, error) {
	return s.orgRepo.ListOrganizations()
}

// UpdateOrganization handles the logic for updating an organization's details.
func (s *organizationService) UpdateOrganization(orgID string, dto *UpdateOrganizationDTO) (*models.Organization, error) {
	if errs := utils.ValidateStruct(dto); errs != nil {
		return nil, errors.New("invalid update data")
	}

	org, err := s.orgRepo.FindOrganizationBySlug(orgID) // Assuming orgID is the slug for now
	if err != nil {
		return nil, errors.New("organization not found")
	}

	// Update fields if they are provided in the DTO
	if dto.Name != "" {
		org.Name = dto.Name
		org.Slug = slugify(dto.Name) // Re-slugify if name changes
	}
	if dto.Description != "" {
		org.Description = dto.Description
	}
	if dto.Website != "" {
		org.Website = dto.Website
	}
	if dto.Address != "" {
		org.Address = dto.Address
	}
	if dto.Phone != "" {
		org.Phone = dto.Phone
	}
	if dto.Avatar != "" {
		org.Avatar = dto.Avatar
	}

	if err := s.orgRepo.UpdateOrganization(org); err != nil {
		return nil, errors.New("failed to update organization")
	}

	return org, nil
}

// slugify creates a URL-friendly slug from a string.
func slugify(s string) string {
	s = strings.ToLower(s)
	re := regexp.MustCompile(`[^a-z0-9]+`)
	s = re.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	// You might want to add a random suffix to ensure uniqueness in a real app
	return s
}

// DeleteOrganization handles the logic for soft-deleting an organization.
func (s *organizationService) DeleteOrganization(orgID string) error {
	org, err := s.orgRepo.FindOrganizationBySlug(orgID) // Assuming orgID is the slug
	if err != nil {
		return errors.New("organization not found")
	}

	return s.orgRepo.DeleteOrganization(org)
}

// RemoveMember handles the logic for removing a user from an organization.
func (s *organizationService) RemoveMember(orgID, userID string) error {
	// The repository method handles the soft delete.
	// We could add more logic here, e.g., checking if the user being removed
	// is the last owner of the organization.
	return s.orgRepo.RemoveMembership(orgID, userID)
}

// AddMember adds a user to an organization with a specific role.
func (s *organizationService) AddMember(orgSlug, userID, roleID string) error {
	// Find the organization by slug
	org, err := s.orgRepo.FindOrganizationBySlug(orgSlug)
	if err != nil {
		return errors.New("organization not found")
	}

	// Verify the user exists
	_, err = s.authRepo.FindIdentityByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	// Create the membership
	membership := &models.OrganizationalMembership{
		ID:             uuid.New().String(),
		IdentityID:     userID,
		OrganizationID: org.ID,
		RoleID:         roleID,
		IsActive:       true,
		ActiveFrom:     time.Now(),
	}

	return s.orgRepo.CreateMembership(membership)
}

// ListMembers returns all members of an organization.
func (s *organizationService) ListMembers(orgSlug string) ([]models.OrganizationalMembership, error) {
	// Find the organization by slug
	org, err := s.orgRepo.FindOrganizationBySlug(orgSlug)
	if err != nil {
		return nil, errors.New("organization not found")
	}

	return s.orgRepo.ListMemberships(org.ID)
}

// ChangeUserRole changes a user's role within an organization.
func (s *organizationService) ChangeUserRole(orgSlug, userID, newRoleID string) error {
	// Find the organization by slug
	org, err := s.orgRepo.FindOrganizationBySlug(orgSlug)
	if err != nil {
		return errors.New("organization not found")
	}

	// Find the existing membership
	membership, err := s.orgRepo.FindMembership(org.ID, userID)
	if err != nil {
		return errors.New("user is not a member of this organization")
	}

	// Update the role
	membership.RoleID = newRoleID
	return s.orgRepo.UpdateMembership(membership)
}
