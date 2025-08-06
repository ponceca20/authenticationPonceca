package invitation

// InvitationDTO is used to create a new invitation.
type InvitationDTO struct {
	Email  string `json:"email" validate:"required,email"`
	RoleID string `json:"role_id" validate:"required"`
}

// AcceptInvitationDTO is used to accept an invitation.
type AcceptInvitationDTO struct {
	Token     string `json:"token" validate:"required"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Password  string `json:"password" validate:"required,min=8"`
}

// AcceptInvitationByTokenDTO is used to accept an invitation using only the token from URL.
type AcceptInvitationByTokenDTO struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Password  string `json:"password" validate:"required,min=8"`
}

// InvitationResponseDTO is a public representation of a sent invitation.
type InvitationResponseDTO struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	RoleID    string `json:"role_id"`
	Status    string `json:"status"`
	ExpiresAt string `json:"expires_at"`
}
