package models

import "time"

// GuestSession represents a temporary session for a user who has not registered.
type GuestSession struct {
	ID           string `json:"id" gorm:"primaryKey;type:varchar(36)"`
	SessionToken string `json:"session_token" gorm:"uniqueIndex;size:191"`

	// Temporary guest data
	Email     string `json:"email,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Phone     string `json:"phone,omitempty"`

	// Activity tracking
	CartData     string    `json:"cart_data,omitempty" gorm:"type:text"` // JSON blob for cart
	LastActivity time.Time `json:"last_activity"`
	IPAddress    string    `json:"ip_address,omitempty"`
	UserAgent    string    `json:"user_agent,omitempty"`

	// Auto-expiration
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}
