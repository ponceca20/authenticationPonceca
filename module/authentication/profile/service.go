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
	// Validate input DTO
	if errs := utils.ValidateStruct(dto); errs != nil {
		return nil, fmt.Errorf("validation failed: %v", errs)
	}

	// Validate that identityID is not empty
	if identityID == "" {
		return nil, fmt.Errorf("identity ID cannot be empty")
	}

	var updatedIdentity *models.Identity
	var updatedProfile *models.UserProfile

	err := s.db.Transaction(func(tx *gorm.DB) error {
		// Use a transaction-specific repository instance
		txRepo := NewProfileRepository(tx)

		// Find the identity first
		identity, err := txRepo.FindIdentityByID(identityID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return fmt.Errorf("user not found with ID: %s", identityID)
			}
			return fmt.Errorf("error finding user: %w", err)
		}

		// Update Identity fields only if they are not empty
		updated := false
		if dto.FirstName != "" && dto.FirstName != identity.FirstName {
			identity.FirstName = dto.FirstName
			updated = true
		}
		if dto.LastName != "" && dto.LastName != identity.LastName {
			identity.LastName = dto.LastName
			updated = true
		}
		if dto.Avatar != "" && dto.Avatar != identity.Avatar {
			identity.Avatar = dto.Avatar
			updated = true
		}
		if dto.Phone != "" && dto.Phone != identity.Phone {
			identity.Phone = dto.Phone
			updated = true
		}

		// Only update if there were changes
		if updated {
			if err := txRepo.UpdateIdentity(identity); err != nil {
				return fmt.Errorf("failed to update identity: %w", err)
			}
		}
		updatedIdentity = identity

		// Find or create and then update UserProfile
		profile := &models.UserProfile{IdentityID: identityID}
		if err := txRepo.FindOrCreateUserProfile(profile); err != nil {
			return fmt.Errorf("failed to find or create user profile: %w", err)
		}

		// Update profile fields
		profileUpdated := false
		if dto.Bio != profile.Bio {
			profile.Bio = dto.Bio
			profileUpdated = true
		}
		if dto.Location != profile.Location {
			profile.Location = dto.Location
			profileUpdated = true
		}
		if dto.Website != profile.Website {
			profile.Website = dto.Website
			profileUpdated = true
		}

		// Only update if there were changes
		if profileUpdated {
			if err := txRepo.UpdateUserProfile(profile); err != nil {
				return fmt.Errorf("failed to update user profile: %w", err)
			}
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
