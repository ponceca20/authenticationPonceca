package models

// User represents the primary user entity in the system.
// This is an alias to Identity to maintain backward compatibility and provide
// a more intuitive interface for general user operations.
//
// Architecture Notes:
// - User is the primary interface for authentication and basic user operations
// - Identity contains the core user data and credentials
// - For specialized contexts, use:
//   * OrganizationalMembership for employees/members within organizations
//   * CustomerProfile for e-commerce specific data
//   * UserProfile for extended non-essential user information
//
// This design provides:
// - Single source of truth for user identity
// - Flexible extension for different user contexts
// - Clean separation of concerns
// - Easy migration path if User needs to become a distinct model
type User = Identity

// TableName ensures User uses the same table as Identity
// This maintains database consistency while providing the User interface
func (User) TableName() string {
	return "identity"
}
