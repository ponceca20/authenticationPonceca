package unified

import "gorm.io/gorm"

// UnifiedRepository defines the interface for cross-submodule database queries.
type UnifiedRepository interface {
	// In a real implementation, you would have methods here that perform
	// complex joins and queries across multiple tables (audit logs, invitations, etc.)
	// For now, we'll leave it empty as the logic is not specified.
}

type unifiedRepository struct {
	db *gorm.DB
}

// NewUnifiedRepository creates a new instance of UnifiedRepository.
func NewUnifiedRepository(db *gorm.DB) UnifiedRepository {
	return &unifiedRepository{db: db}
}
