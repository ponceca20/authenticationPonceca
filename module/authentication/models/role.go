package models

import "time"

// Role defines a set of permissions within an organization.
type Role struct {
	ID             string       `json:"id" gorm:"primaryKey;type:varchar(36)"`
	OrganizationID string       `json:"organization_id" gorm:"type:varchar(36);not null;uniqueIndex:idx_org_role_name,priority:1"`
	Organization   Organization `json:"organization" gorm:"foreignKey:OrganizationID"`
	Name           string       `json:"name" gorm:"not null;size:100;uniqueIndex:idx_org_role_name,priority:2"`
	DisplayName    string       `json:"display_name" gorm:"size:255"`
	Description    string       `json:"description" gorm:"size:500"`
	HierarchyLevel int          `json:"hierarchy_level" gorm:"default:0;type:smallint"`
	IsSystemRole   bool         `json:"is_system_role" gorm:"default:false"`
	Permissions    []Permission `json:"permissions" gorm:"many2many:role_permission;"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
	DeletedAt      *time.Time   `json:"deleted_at,omitempty" gorm:"index"`
}

// Permission defines a specific action that can be performed on a resource.
type Permission struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Resource  string    `json:"resource" gorm:"not null;size:100;uniqueIndex:idx_resource_action_scope"` // e.g., "users", "products"
	Action    string    `json:"action" gorm:"not null;size:100;uniqueIndex:idx_resource_action_scope"`   // e.g., "create", "read", "update", "delete"
	Scope     string    `json:"scope" gorm:"not null;size:50;uniqueIndex:idx_resource_action_scope"`     // e.g., "own", "department", "organization", "all"
	CreatedAt time.Time `json:"created_at"`
}

// SystemResources defines the available resources and their actions in the system.
var SystemResources = map[string][]string{
	"users":        {"create", "read", "update", "delete", "invite", "suspend"},
	"students":     {"read", "update", "grade", "report", "communicate"},
	"teachers":     {"read", "update", "assign", "evaluate"},
	"courses":      {"create", "read", "update", "delete", "enroll"},
	"grades":       {"create", "read", "update", "approve", "publish"},
	"finances":     {"read", "create", "update", "approve", "report"},
	"departments":  {"create", "read", "update", "delete", "manage"},
	"products":     {"create", "read", "update", "delete", "price", "inventory"},
	"orders":       {"create", "read", "update", "process", "refund"},
	"customers":    {"read", "update", "communicate", "discount"},
	"audit":        {"read", "export"},
	"settings":     {"read", "update"},
	"integrations": {"read", "configure"},
}

// AuthorizationScopes defines the different levels of access.
var AuthorizationScopes = []string{
	"own",
	"department",
	"organization",
	"all",
}
