package models

import "time"

// AuditLog records an action performed by a user in the system.
type AuditLog struct {
	ID             string        `json:"id" gorm:"primaryKey;type:varchar(36)"`
	OrganizationID *string       `json:"organization_id,omitempty" gorm:"index;type:varchar(36)"` // Optional, for system-level actions
	IdentityID     string        `json:"identity_id" gorm:"index;type:varchar(36)"`
	Identity       Identity      `json:"identity" gorm:"foreignKey:IdentityID"`
	Action         string        `json:"action" gorm:"not null"`        // e.g., "user.login", "organization.create"
	Resource       string        `json:"resource"`                      // e.g., "organization"
	ResourceID     string        `json:"resource_id"`                   // e.g., the ID of the created organization
	Status         string        `json:"status"`                        // e.g., "success", "failure"
	IPAddress      string        `json:"ip_address,omitempty"`
	UserAgent      string        `json:"user_agent,omitempty"`
	Details        string        `json:"details,omitempty" gorm:"type:text"`
	Timestamp      time.Time     `json:"timestamp" gorm:"not null"`
}
