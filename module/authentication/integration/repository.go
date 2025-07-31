package integration

import "gorm.io/gorm"

// IntegrationRepository defines the interface for integration-related DB operations.
type IntegrationRepository interface {
	// Empty for now, as validation is mostly a JWT concern.
}

type integrationRepository struct {
	db *gorm.DB
}

// NewIntegrationRepository creates a new instance of IntegrationRepository.
func NewIntegrationRepository(db *gorm.DB) IntegrationRepository {
	return &integrationRepository{db: db}
}
