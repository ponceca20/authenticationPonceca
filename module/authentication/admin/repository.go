package admin

import (
	"practicev2/module/authentication/models"

	"gorm.io/gorm"
)

// AdminRepository defines the interface for admin-related DB operations.
type AdminRepository interface {
	CountUsers() (int64, error)
	CountOrganizations() (int64, error)
}

type adminRepository struct {
	db *gorm.DB
}

// NewAdminRepository creates a new instance of AdminRepository.
func NewAdminRepository(db *gorm.DB) AdminRepository {
	return &adminRepository{db: db}
}

// CountUsers counts all non-deleted users.
func (r *adminRepository) CountUsers() (int64, error) {
	var count int64
	err := r.db.Model(&models.Identity{}).Count(&count).Error
	return count, err
}

// CountOrganizations counts all non-deleted organizations.
func (r *adminRepository) CountOrganizations() (int64, error) {
	var count int64
	err := r.db.Model(&models.Organization{}).Count(&count).Error
	return count, err
}
