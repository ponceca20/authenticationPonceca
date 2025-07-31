package models

import "time"

// OrganizationalMembership links an Identity to an Organization with a specific Role and context.
type OrganizationalMembership struct {
	ID             string       `json:"id" gorm:"primaryKey;type:varchar(36)"`
	IdentityID     string       `json:"identity_id" gorm:"index;type:varchar(36);not null"`
	Identity       Identity     `json:"identity" gorm:"foreignKey:IdentityID"`
	OrganizationID string       `json:"organization_id" gorm:"index;type:varchar(36);not null"`
	Organization   Organization `json:"organization" gorm:"foreignKey:OrganizationID"`
	RoleID         string       `json:"role_id" gorm:"index;type:varchar(36);not null"`
	Role           Role         `json:"role" gorm:"foreignKey:RoleID"`

	// Context-specific data
	Department string `json:"department,omitempty"`               // For employees
	Grade      string `json:"grade,omitempty"`                    // For students
	Subject    string `json:"subject,omitempty"`                  // For teachers
	StudentID  string `json:"student_id,omitempty" gorm:"index"`  // Unique student ID
	EmployeeID string `json:"employee_id,omitempty" gorm:"index"` // Unique employee ID

	// Temporal control
	ActiveFrom  time.Time  `json:"active_from" gorm:"not null"`
	ActiveUntil *time.Time `json:"active_until,omitempty"`
	IsActive    bool       `json:"is_active" gorm:"default:true"`

	// Standard model fields
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}
