package profile

import "practicev2/module/authentication/models"

// ProfileDTO is used to update a user's profile information.
// It includes fields from both the Identity and UserProfile models.
type ProfileDTO struct {
	FirstName string `json:"first_name" validate:"omitempty,min=2,max=100"`
	LastName  string `json:"last_name" validate:"omitempty,min=2,max=100"`
	Avatar    string `json:"avatar,omitempty" validate:"omitempty,url"`
	Phone     string `json:"phone,omitempty"`
	Bio       string `json:"bio,omitempty"`
	Location  string `json:"location,omitempty"`
	Website   string `json:"website,omitempty" validate:"omitempty,url"`
}

// ProfileResponseDTO is a public representation of a user's full profile.
type ProfileResponseDTO struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Avatar    string `json:"avatar,omitempty"`
	Phone     string `json:"phone,omitempty"`
	Bio       string `json:"bio,omitempty"`
	Location  string `json:"location,omitempty"`
	Website   string `json:"website,omitempty"`
}

// ToProfileResponseDTO converts an Identity and UserProfile model into a public DTO.
func ToProfileResponseDTO(identity *models.Identity, profile *models.UserProfile) ProfileResponseDTO {
	dto := ProfileResponseDTO{
		ID:        identity.ID,
		Email:     identity.Email,
		FirstName: identity.FirstName,
		LastName:  identity.LastName,
		Avatar:    identity.Avatar,
		Phone:     identity.Phone,
	}
	if profile != nil {
		dto.Bio = profile.Bio
		dto.Location = profile.Location
		dto.Website = profile.Website
	}
	return dto
}
