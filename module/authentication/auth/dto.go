package auth

import (
	"practicev2/module/authentication/models"
	"time"
)

// RegisterDTO defines the data structure for a user registration request.
type RegisterDTO struct {
	FirstName string `json:"first_name" validate:"required,min=2,max=100"`
	LastName  string `json:"last_name" validate:"required,min=2,max=100"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8,max=100"`
}

// LoginDTO defines the data structure for a user login request.
type LoginDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	OrgSlug  string `json:"org_slug,omitempty"` // Optional for organization-specific context
}

// IdentityProfileDTO is a subset of the Identity model for public responses.
type IdentityProfileDTO struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Avatar    string `json:"avatar,omitempty"`
}

// ContextSummaryDTO provides a summary of an available user context.
type ContextSummaryDTO struct {
	Type string `json:"type"` // "organization" or "customer"
	ID   string `json:"id"`
	Name string `json:"name"`
	Role string `json:"role,omitempty"`
}

// TokenResponseDTO defines the structure of the response when tokens are generated.
type TokenResponseDTO struct {
	AccessToken  string              `json:"access_token"`
	RefreshToken string              `json:"refresh_token"`
	ExpiresAt    time.Time           `json:"expires_at"`
	Identity     IdentityProfileDTO  `json:"identity"`
	Contexts     []ContextSummaryDTO `json:"contexts"`
}

// ToIdentityProfileDTO converts an Identity model to a public DTO.
func ToIdentityProfileDTO(identity *models.Identity) IdentityProfileDTO {
	return IdentityProfileDTO{
		ID:        identity.ID,
		Email:     identity.Email,
		FirstName: identity.FirstName,
		LastName:  identity.LastName,
		Avatar:    identity.Avatar,
	}
}

// RefreshTokenDTO defines the data structure for a token refresh request.
type RefreshTokenDTO struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// ForgotPasswordDTO defines the data structure for a forgot password request.
type ForgotPasswordDTO struct {
	Email string `json:"email" validate:"required,email"`
}

// ResetPasswordDTO defines the data structure for a password reset request.
type ResetPasswordDTO struct {
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"new_password" validate:"required"` // Validation handled by service
}

// ChangePasswordDTO defines the data structure for changing password (authenticated user).
type ChangePasswordDTO struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8,max=100"`
}

// VerifyEmailDTO defines the data structure for email verification.
type VerifyEmailDTO struct {
	Token string `json:"token" validate:"required"`
}

// ResendVerificationDTO defines the data structure for resending verification email.
type ResendVerificationDTO struct {
	Email string `json:"email" validate:"required,email"`
}

// ConvertGuestDTO defines the data for converting a guest session to a full user.
type ConvertGuestDTO struct {
	GuestSessionToken string `json:"guest_session_token" validate:"required"`
	FirstName         string `json:"first_name" validate:"required"`
	LastName          string `json:"last_name" validate:"required"`
	Email             string `json:"email" validate:"required,email"`
	Phone             string `json:"phone,omitempty"`
	Password          string `json:"password" validate:"required"`
}
