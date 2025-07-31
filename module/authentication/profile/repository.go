package profile

import (
	"practicev2/module/authentication/models"

	"gorm.io/gorm"
)

// ProfileRepository defines the interface for database operations related to user profiles.
type ProfileRepository interface {
	FindIdentityByID(identityID string) (*models.Identity, error)
	UpdateIdentity(identity *models.Identity) error
	FindOrCreateUserProfile(profile *models.UserProfile) error
	UpdateUserProfile(profile *models.UserProfile) error
}

type profileRepository struct {
	db *gorm.DB
}

// NewProfileRepository creates a new instance of ProfileRepository.
func NewProfileRepository(db *gorm.DB) ProfileRepository {
	return &profileRepository{db: db}
}

// FindIdentityByID retrieves an identity by its ID.
func (r *profileRepository) FindIdentityByID(identityID string) (*models.Identity, error) {
	var identity models.Identity
	if err := r.db.Where("id = ?", identityID).First(&identity).Error; err != nil {
		return nil, err
	}
	return &identity, nil
}

// UpdateIdentity saves changes to an Identity model.
func (r *profileRepository) UpdateIdentity(identity *models.Identity) error {
	return r.db.Save(identity).Error
}

// FindOrCreateUserProfile finds a user profile by identity ID or creates it if it doesn't exist.
func (r *profileRepository) FindOrCreateUserProfile(profile *models.UserProfile) error {
	// FirstOrCreate finds the first record that matches given conditions,
	// or create a new one with the given conditions if none found.
	return r.db.Where(models.UserProfile{IdentityID: profile.IdentityID}).FirstOrCreate(profile).Error
}

// UpdateUserProfile saves changes to a UserProfile model.
func (r *profileRepository) UpdateUserProfile(profile *models.UserProfile) error {
	return r.db.Save(profile).Error
}
