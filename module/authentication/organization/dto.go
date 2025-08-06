package organization

import (
	"errors"
	"practicev2/module/authentication/auth"
	"practicev2/module/authentication/models"
	"strings"

	"github.com/go-playground/validator/v10"
)

// OrganizationRegistrationDTO defines the data structure for registering a new organization
// along with its founding user.
type OrganizationRegistrationDTO struct {
	// Founder's identity data
	Identity auth.RegisterDTO `json:"identity" validate:"required"`

	// Organization data
	Name        string `json:"name" validate:"required,min=2,max=100"`
	Type        string `json:"type" validate:"required,oneof=company educational_institution"`
	Description string `json:"description,omitempty"`
	Website     string `json:"website,omitempty" validate:"omitempty,url"`
	Address     string `json:"address,omitempty"`
	Phone       string `json:"phone,omitempty"`
}

// SanitizeAndValidate cleans and validates the DTO.
func (dto *OrganizationRegistrationDTO) SanitizeAndValidate() error {
	// Sanitization
	dto.Name = strings.TrimSpace(dto.Name)
	dto.Identity.FirstName = strings.TrimSpace(dto.Identity.FirstName)
	dto.Identity.LastName = strings.TrimSpace(dto.Identity.LastName)
	dto.Identity.Email = strings.ToLower(strings.TrimSpace(dto.Identity.Email))

	// Validation
	validate := validator.New()
	if err := validate.Struct(dto); err != nil {
		// This returns the first validation error found.
		// For a more detailed response, the handler can use utils.ValidateStruct
		return err
	}

	// Custom business logic validation
	if dto.Type != "company" && dto.Type != "educational_institution" {
		return errors.New("invalid organization type")
	}

	return nil
}

// OrganizationResponseDTO is a public representation of an organization.
type OrganizationResponseDTO struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Type        string `json:"type"`
	Description string `json:"description,omitempty"`
	Avatar      string `json:"avatar,omitempty"`
	Website     string `json:"website,omitempty"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// ToOrganizationResponseDTO converts an Organization model to a public DTO.
func ToOrganizationResponseDTO(org *models.Organization) OrganizationResponseDTO {
	return OrganizationResponseDTO{
		ID:          org.ID,
		Name:        org.Name,
		Slug:        org.Slug,
		Type:        org.Type,
		Description: org.Description,
		Avatar:      org.Avatar,
		Website:     org.Website,
		CreatedAt:   org.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   org.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// UpdateOrganizationDTO defines the data structure for updating an organization.
type UpdateOrganizationDTO struct {
	Name        string `json:"name,omitempty" validate:"omitempty,min=2,max=100"`
	Description string `json:"description,omitempty"`
	Website     string `json:"website,omitempty" validate:"omitempty,url"`
	Address     string `json:"address,omitempty"`
	Phone       string `json:"phone,omitempty"`
	Avatar      string `json:"avatar,omitempty" validate:"omitempty,url"`
}

// MemberResponseDTO defines the simplified response structure for organization members.
type MemberResponseDTO struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	RoleID    string `json:"role_id"`
	IsActive  bool   `json:"is_active"`
	JoinedAt  string `json:"joined_at"`
}

// ToMemberResponseDTO converts an OrganizationalMembership to a simplified DTO.
func ToMemberResponseDTO(membership *models.OrganizationalMembership) MemberResponseDTO {
	joinedAt := ""
	if !membership.ActiveFrom.IsZero() {
		joinedAt = membership.ActiveFrom.Format("2006-01-02T15:04:05Z07:00")
	}

	return MemberResponseDTO{
		ID:        membership.ID,
		FirstName: membership.Identity.FirstName,
		LastName:  membership.Identity.LastName,
		Email:     membership.Identity.Email,
		Role:      membership.Role.Name,
		RoleID:    membership.RoleID,
		IsActive:  membership.IsActive,
		JoinedAt:  joinedAt,
	}
}
