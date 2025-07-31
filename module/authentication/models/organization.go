package models

import "time"

// Organization represents a multi-tenant entity like a company or a school.
type Organization struct {
	ID          string     `json:"id" gorm:"primaryKey;type:varchar(36)"`
	Name        string     `json:"name" gorm:"not null"`
	Slug        string     `json:"slug" gorm:"uniqueIndex;size:191"`
	Type        string     `json:"type" gorm:"not null"` // e.g., "company", "educational_institution"
	Description string     `json:"description"`
	Avatar      string     `json:"avatar,omitempty"`
	Website     string     `json:"website,omitempty"`
	Phone       string     `json:"phone,omitempty"`
	Address     string     `json:"address,omitempty"`
	IsActive    bool       `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}
