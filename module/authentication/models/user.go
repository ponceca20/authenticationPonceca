package models

// The User model is now represented by the central Identity model.
// This file is kept for structural alignment with the guide.
// In this architecture, any 'User' is an 'Identity'.

// For specific contexts, we use other models that link to the Identity:
// - OrganizationalMembership for employees/members of an organization.
// - CustomerProfile for e-commerce customers.

// If a distinct User model that extends Identity is needed later,
// it can be defined here. For now, we refer directly to Identity.
type User = Identity
