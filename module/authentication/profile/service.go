package profile

import (
	"fmt"
	"practicev2/module/authentication/models"
	"practicev2/module/authentication/utils"

	"gorm.io/gorm"
)

// ProfileService defines the interface for profile-related business logic.
type ProfileService interface {
	GetProfile(identityID string) (*ProfileResponseDTO, error)
	UpdateProfile(identityID string, dto *ProfileDTO) (*ProfileResponseDTO, error)
}

type profileService struct {
	repo ProfileRepository
	db   *gorm.DB // For transactions
}

// NewProfileService creates a new instance of ProfileService.
func NewProfileService(repo ProfileRepository, db *gorm.DB) ProfileService {
	return &profileService{repo: repo, db: db}
}

// GetProfile retrieves a user's full profile.
func (s *profileService) GetProfile(identityID string) (*ProfileResponseDTO, error) {
	identity, err := s.repo.FindIdentityByID(identityID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	profile := &models.UserProfile{IdentityID: identityID}
	// This will find the profile or initialize a new one if not found.
	// We don't care about the error here, as a missing profile is not an error for a GET request.
	_ = s.repo.FindOrCreateUserProfile(profile)

	dto := ToProfileResponseDTO(identity, profile)
	return &dto, nil
}

// UpdateProfile updates a user's core and extended profile information.
func (s *profileService) UpdateProfile(identityID string, dto *ProfileDTO) (*ProfileResponseDTO, error) {
	if errs := utils.ValidateStruct(dto); errs != nil {
		return nil, fmt.Errorf("validation failed: %v", errs)
	}

	var updatedIdentity *models.Identity
	var updatedProfile *models.UserProfile

	err := s.db.Transaction(func(tx *gorm.DB) error {
		// Use a transaction-specific repository instance
		txRepo := NewProfileRepository(tx)

		identity, err := txRepo.FindIdentityByID(identityID)
		if err != nil {
			return fmt.Errorf("user not found: %w", err)
		}

		// Update Identity fields
		if dto.FirstName != "" {
			identity.FirstName = dto.FirstName
		}
		if dto.LastName != "" {
			identity.LastName = dto.LastName
		}
		if dto.Avatar != "" {
			identity.Avatar = dto.Avatar
		}
		if dto.Phone != "" {
			identity.Phone = dto.Phone
		}
		if err := txRepo.UpdateIdentity(identity); err != nil {
			return err
		}
		updatedIdentity = identity

		// Find or create and then update UserProfile
		profile := &models.UserProfile{IdentityID: identityID}
		if err := txRepo.FindOrCreateUserProfile(profile); err != nil {
			return err
		}

		profile.Bio = dto.Bio
		profile.Location = dto.Location
		profile.Website = dto.Website
		if err := txRepo.UpdateUserProfile(profile); err != nil {
			return err
		}
		updatedProfile = profile

		return nil
	})

	if err != nil {
		return nil, err
	}

	responseDTO := ToProfileResponseDTO(updatedIdentity, updatedProfile)
	return &responseDTO, nil
}
