package models

import "time"

// Identity represents the central, unified identity of a person in the system.
// It holds the core credentials and personal information.
type Identity struct {
	ID                  string     `json:"id" gorm:"primaryKey;type:varchar(36)"`
	Email               string     `json:"email" gorm:"uniqueIndex;size:191;not null"`
	PasswordHash        string     `json:"-" gorm:"column:password_hash;size:255;not null"`
	EmailVerified       bool       `json:"email_verified" gorm:"default:false"`
	EmailVerifiedAt     *time.Time `json:"email_verified_at,omitempty"`
	FirstName           string     `json:"first_name" gorm:"size:100;not null"`
	LastName            string     `json:"last_name" gorm:"size:100;not null"`
	Avatar              string     `json:"avatar,omitempty" gorm:"size:500"`
	Phone               string     `json:"phone,omitempty" gorm:"size:20"`
	DateOfBirth         *time.Time `json:"date_of_birth,omitempty"`
	LastLoginAt         *time.Time `json:"last_login_at,omitempty"`
	FailedLoginAttempts int        `json:"-" gorm:"default:0;type:smallint"`
	LockedUntil         *time.Time `json:"locked_until,omitempty"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
	DeletedAt           *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}
