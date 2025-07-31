package user

import (
	"practicev2/module/authentication/models"
	"time"
)

// UserDTO represents a user within the context of an organization for API responses.
type UserDTO struct {
	ID             string    `json:"id"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	Email          string    `json:"email"`
	Avatar         string    `json:"avatar,omitempty"`
	Role           string    `json:"role"`
	Department     string    `json:"department,omitempty"`
	IsActive       bool      `json:"is_active"`
	JoinedAt       time.Time `json:"joined_at"`
}

// ToUserDTO converts an Identity and its corresponding Membership into a public UserDTO.
func ToUserDTO(membership *models.OrganizationalMembership) UserDTO {
	return UserDTO{
		ID:         membership.Identity.ID,
		FirstName:  membership.Identity.FirstName,
		LastName:   membership.Identity.LastName,
		Email:      membership.Identity.Email,
		Avatar:     membership.Identity.Avatar,
		Role:       membership.Role.Name,
		Department: membership.Department,
		IsActive:   membership.IsActive,
		JoinedAt:   membership.CreatedAt,
	}
}

// UpdateUserDTO defines the data for updating a user's organizational context.
type UpdateUserDTO struct {
	RoleID     string `json:"role_id,omitempty"`
	Department string `json:"department,omitempty"`
}

// UpdateUserStatusDTO defines the data for changing a user's status.
type UpdateUserStatusDTO struct {
	IsActive bool `json:"is_active"`
}

// BulkCreateUserItemDTO defines a single user in a bulk creation request.
type BulkCreateUserItemDTO struct {
	FirstName  string `json:"first_name" validate:"required"`
	LastName   string `json:"last_name" validate:"required"`
	Email      string `json:"email" validate:"required,email"`
	Password   string `json:"password" validate:"required,min=8"`
	RoleID     string `json:"role_id" validate:"required"`
	Department string `json:"department,omitempty"`
}

// BulkCreateUserDTO is the top-level DTO for bulk user creation.
type BulkCreateUserDTO struct {
	Users []BulkCreateUserItemDTO `json:"users" validate:"required,min=1,dive"`
}

// BulkUpdateUserItemDTO defines a single user in a bulk update request.
type BulkUpdateUserItemDTO struct {
	ID         string `json:"id" validate:"required"`
	RoleID     string `json:"role_id,omitempty"`
	Department string `json:"department,omitempty"`
	IsActive   *bool  `json:"is_active,omitempty"`
}

// BulkUpdateUserDTO is the top-level DTO for bulk user updates.
type BulkUpdateUserDTO struct {
	Users []BulkUpdateUserItemDTO `json:"users" validate:"required,min=1,dive"`
}

// BulkDeleteUserDTO is the top-level DTO for bulk user deletion.
type BulkDeleteUserDTO struct {
	UserIDs []string `json:"user_ids" validate:"required,min=1"`
}
