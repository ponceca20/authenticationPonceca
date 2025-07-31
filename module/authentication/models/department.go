package models

import "time"

// Department represents a subdivision within an organization.
type Department struct {
	ID             string       `json:"id" gorm:"primaryKey;type:varchar(36)"`
	OrganizationID string       `json:"organization_id" gorm:"index;type:varchar(36)"`
	Organization   Organization `json:"organization" gorm:"foreignKey:OrganizationID"`
	ParentID       *string      `json:"parent_id,omitempty" gorm:"index;type:varchar(36)"` // For hierarchical structures
	Name           string       `json:"name" gorm:"not null;size:255"`
	Description    string       `json:"description" gorm:"size:500"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
	DeletedAt      *time.Time   `json:"deleted_at,omitempty" gorm:"index"`
}
