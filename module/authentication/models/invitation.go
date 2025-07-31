package models

import "time"

// Invitation represents an invitation for a user to join an organization.
type Invitation struct {
	ID             string       `json:"id" gorm:"primaryKey;type:varchar(36)"`
	OrganizationID string       `json:"organization_id" gorm:"index;type:varchar(36);not null"`
	Organization   Organization `json:"organization" gorm:"foreignKey:OrganizationID"`
	InviterID      string       `json:"inviter_id" gorm:"index;type:varchar(36);not null"` // The user who sent the invitation
	Inviter        Identity     `json:"inviter" gorm:"foreignKey:InviterID"`
	Email          string       `json:"email" gorm:"not null;size:191"`
	RoleID         string       `json:"role_id" gorm:"index;type:varchar(36);not null"`
	Role           Role         `json:"role" gorm:"foreignKey:RoleID"`
	Token          string       `json:"-" gorm:"uniqueIndex;size:191;not null"`
	Status         string       `json:"status" gorm:"default:'pending';size:50"` // e.g., pending, accepted, expired
	ExpiresAt      time.Time    `json:"expires_at" gorm:"not null"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
}
