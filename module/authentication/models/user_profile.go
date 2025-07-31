package models

import "time"

// UserProfile contains extended, non-essential information for an identity.
type UserProfile struct {
	ID         string   `json:"id" gorm:"primaryKey;type:varchar(36)"`
	IdentityID string   `json:"identity_id" gorm:"uniqueIndex;type:varchar(36)"`
	Identity   Identity `json:"identity" gorm:"foreignKey:IdentityID"`
	Bio        string   `json:"bio,omitempty" gorm:"type:text"`
	Location   string   `json:"location,omitempty"`
	Website    string   `json:"website,omitempty"`
	Socials    string   `json:"socials,omitempty" gorm:"type:text"` // JSON blob for social links
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
