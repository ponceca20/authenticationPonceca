package models

import "time"

// PasswordResetToken stores a token for the password reset process.
type PasswordResetToken struct {
	ID         string    `json:"id" gorm:"primaryKey;type:varchar(36)"`
	IdentityID string    `json:"identity_id" gorm:"index;type:varchar(36)"`
	Identity   Identity  `json:"identity" gorm:"foreignKey:IdentityID"`
	Token      string    `json:"-" gorm:"uniqueIndex;size:512"`
	ExpiresAt  time.Time `json:"expires_at"`
	CreatedAt  time.Time `json:"created_at"`
}
