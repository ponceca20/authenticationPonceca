package models

import "time"

// GuestSession represents a temporary session for a user who has not registered.
type GuestSession struct {
	ID           string `json:"id" gorm:"primaryKey;type:varchar(36)"`
	SessionToken string `json:"session_token" gorm:"uniqueIndex;size:512;not null"`

	// Temporary guest data
	Email     string `json:"email,omitempty" gorm:"size:191"`
	FirstName string `json:"first_name,omitempty" gorm:"size:100"`
	LastName  string `json:"last_name,omitempty" gorm:"size:100"`
	Phone     string `json:"phone,omitempty" gorm:"size:20"`

	// Activity tracking
	CartData     string    `json:"cart_data,omitempty" gorm:"type:json"` // JSON blob for cart
	LastActivity time.Time `json:"last_activity" gorm:"not null"`
	IPAddress    string    `json:"ip_address,omitempty" gorm:"size:45"` // IPv6 support
	UserAgent    string    `json:"user_agent,omitempty" gorm:"size:1000"`

	// Auto-expiration
	ExpiresAt time.Time `json:"expires_at" gorm:"not null;index"`
	CreatedAt time.Time `json:"created_at"`
}
