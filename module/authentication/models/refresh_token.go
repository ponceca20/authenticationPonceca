package models

import "time"

// RefreshToken stores a refresh token for an identity.
type RefreshToken struct {
	ID         string    `json:"id" gorm:"primaryKey;type:varchar(36)"`
	IdentityID string    `json:"identity_id" gorm:"index;type:varchar(36)"`
	Identity   Identity  `json:"identity" gorm:"foreignKey:IdentityID"`
	Token      string    `json:"-" gorm:"uniqueIndex;size:512"`
	ExpiresAt  time.Time `json:"expires_at"`
	IsRevoked  bool      `json:"is_revoked" gorm:"default:false"`
	CreatedAt  time.Time `json:"created_at"`
}
