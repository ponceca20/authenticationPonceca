package guest

import "time"

// GuestSessionDTO is the public representation of a guest session.
type GuestSessionDTO struct {
	SessionToken string    `json:"session_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	CartData     string    `json:"cart_data,omitempty"`
}

// UpdateCartDTO defines the structure for updating a guest's cart.
type UpdateCartDTO struct {
	CartData string `json:"cart_data" validate:"required,json"`
}
