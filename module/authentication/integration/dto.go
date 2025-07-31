package integration

// ValidateTokenRequestDTO defines the request for the token validation endpoint.
type ValidateTokenRequestDTO struct {
	Token string `json:"token" validate:"required"`
}

// ValidateTokenResponseDTO defines the response for a valid token.
type ValidateTokenResponseDTO struct {
	IsValid bool   `json:"is_valid"`
	UserID  string `json:"user_id,omitempty"`
	Email   string `json:"email,omitempty"`
	// Add other claims as needed
}
