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
