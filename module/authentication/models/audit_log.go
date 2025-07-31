package models

import "time"

// AuditLog records an action performed by a user in the system.
type AuditLog struct {
	ID             string    `json:"id" gorm:"primaryKey;type:varchar(36)"`
	OrganizationID *string   `json:"organization_id,omitempty" gorm:"index;type:varchar(36)"` // Optional, for system-level actions
	IdentityID     string    `json:"identity_id" gorm:"index;type:varchar(36);not null"`
	Identity       Identity  `json:"identity" gorm:"foreignKey:IdentityID"`
	Action         string    `json:"action" gorm:"not null;size:255;index"` // e.g., "user.login", "organization.create"
	Resource       string    `json:"resource" gorm:"size:100;index"`        // e.g., "organization"
	ResourceID     string    `json:"resource_id" gorm:"index;size:100"`     // e.g., the ID of the created organization
	Status         string    `json:"status" gorm:"size:50"`                 // e.g., "success", "failure"
	IPAddress      string    `json:"ip_address,omitempty" gorm:"size:45"`
	UserAgent      string    `json:"user_agent,omitempty" gorm:"size:1000"`
	Details        string    `json:"details,omitempty" gorm:"type:json"`
	Timestamp      time.Time `json:"timestamp" gorm:"not null;index"`
}
