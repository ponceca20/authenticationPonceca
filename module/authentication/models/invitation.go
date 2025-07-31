package models

import "time"

// Invitation represents an invitation for a user to join an organization.
type Invitation struct {
	ID             string       `json:"id"gorm:"primaryKey;type:varchar(36)"`
	OrganizationID string       `json:"organization_id" gorm:"index;type:varchar(36)"`
	Organization   Organization `json:"organization" gorm:"foreignKey:OrganizationID"`
	InviterID      string       `json:"inviter_id" gorm:"index;type:varchar(36)"` // The user who sent the invitation
	Inviter        Identity     `json:"inviter" gorm:"foreignKey:InviterID"`
	Email          string       `json:"email" gorm:"not null"`
	RoleID         string       `json:"role_id" gorm:"index;type:varchar(36)"`
	Role           Role         `json:"role" gorm:"foreignKey:RoleID"`
	Token          string       `json:"-" gorm:"uniqueIndex;size:191"`
	Status         string       `json:"status" gorm:"default:'pending'"` // e.g., pending, accepted, expired
	ExpiresAt      time.Time    `json:"expires_at"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
}
