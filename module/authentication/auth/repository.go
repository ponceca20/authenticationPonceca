package auth

import (
	"practicev2/module/authentication/models"
	"time"

	"gorm.io/gorm"
)

// AuthRepository defines the interface for database operations related to authentication.
type AuthRepository interface {
	CreateIdentity(identity *models.Identity) error
	FindIdentityByEmail(email string) (*models.Identity, error)
	FindIdentityByID(id string) (*models.Identity, error)
	UpdateIdentity(identity *models.Identity) error
	GetFullIdentityContext(identityID string) (*models.Identity, []models.OrganizationalMembership, *models.CustomerProfile, error)
	SaveRefreshToken(token *models.RefreshToken) error
	FindRefreshToken(tokenValue string) (*models.RefreshToken, error)
	RevokeRefreshToken(token *models.RefreshToken) error
	CreatePasswordResetToken(token *models.PasswordResetToken) error
	FindPasswordResetToken(tokenValue string) (*models.PasswordResetToken, error)
	DeletePasswordResetToken(token *models.PasswordResetToken) error
	CreateEmailVerificationToken(token *models.EmailVerificationToken) error
	FindEmailVerificationToken(tokenValue string) (*models.EmailVerificationToken, error)
	DeleteEmailVerificationToken(token *models.EmailVerificationToken) error
}

// authRepository is the implementation of AuthRepository.
type authRepository struct {
	db *gorm.DB
}

// NewAuthRepository creates a new instance of authRepository.
func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{db: db}
}

// CreateIdentity saves a new Identity to the database.
func (r *authRepository) CreateIdentity(identity *models.Identity) error {
	return r.db.Create(identity).Error
}

// FindIdentityByEmail finds an identity by their email address.
func (r *authRepository) FindIdentityByEmail(email string) (*models.Identity, error) {
	var identity models.Identity
	if err := r.db.Where("email = ?", email).First(&identity).Error; err != nil {
		return nil, err // gorm.ErrRecordNotFound if not found
	}
	return &identity, nil
}

// FindIdentityByID finds an identity by their ID.
func (r *authRepository) FindIdentityByID(id string) (*models.Identity, error) {
	var identity models.Identity
	if err := r.db.Where("id = ?", id).First(&identity).Error; err != nil {
		return nil, err
	}
	return &identity, nil
}

// GetFullIdentityContext retrieves an identity and all their related contexts.
func (r *authRepository) GetFullIdentityContext(identityID string) (*models.Identity, []models.OrganizationalMembership, *models.CustomerProfile, error) {
	var identity models.Identity
	if err := r.db.Where("id = ?", identityID).First(&identity).Error; err != nil {
		return nil, nil, nil, err
	}

	var memberships []models.OrganizationalMembership
	// Preload the Role and Organization to get complete membership context
	r.db.Preload("Role").Preload("Organization").Where("identity_id = ? AND is_active = ?", identityID, true).Find(&memberships)

	var customerProfile models.CustomerProfile
	err := r.db.Where("identity_id = ?", identityID).First(&customerProfile).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		// Log or handle error, but don't fail the whole context retrieval
		// if the customer profile is just missing.
		return &identity, memberships, nil, nil
	}

	var customer *models.CustomerProfile
	if err == nil {
		customer = &customerProfile
	}

	return &identity, memberships, customer, nil
}

// SaveRefreshToken saves a new refresh token to the database.
func (r *authRepository) SaveRefreshToken(token *models.RefreshToken) error {
	return r.db.Create(token).Error
}

// FindRefreshToken finds a refresh token by its value.
func (r *authRepository) FindRefreshToken(tokenValue string) (*models.RefreshToken, error) {
	var token models.RefreshToken
	// Also preload the Identity to avoid another DB call
	err := r.db.Preload("Identity").Where("token = ? AND is_revoked = false AND expires_at > ?", tokenValue, time.Now()).First(&token).Error
	return &token, err
}

// RevokeRefreshToken marks a refresh token as revoked.
func (r *authRepository) RevokeRefreshToken(token *models.RefreshToken) error {
	token.IsRevoked = true
	return r.db.Save(token).Error
}

// UpdateIdentity saves changes to an Identity model.
func (r *authRepository) UpdateIdentity(identity *models.Identity) error {
	return r.db.Save(identity).Error
}

// --- Password Reset ---

// CreatePasswordResetToken saves a new password reset token.
func (r *authRepository) CreatePasswordResetToken(token *models.PasswordResetToken) error {
	return r.db.Create(token).Error
}

// FindPasswordResetToken finds a password reset token by its value.
func (r *authRepository) FindPasswordResetToken(tokenValue string) (*models.PasswordResetToken, error) {
	var token models.PasswordResetToken
	err := r.db.Where("token = ? AND expires_at > ?", tokenValue, time.Now()).First(&token).Error
	return &token, err
}

// DeletePasswordResetToken deletes a password reset token after it has been used.
func (r *authRepository) DeletePasswordResetToken(token *models.PasswordResetToken) error {
	return r.db.Delete(token).Error
}

// --- Email Verification ---

// CreateEmailVerificationToken saves a new email verification token.
func (r *authRepository) CreateEmailVerificationToken(token *models.EmailVerificationToken) error {
	return r.db.Create(token).Error
}

// FindEmailVerificationToken finds an email verification token by its value.
func (r *authRepository) FindEmailVerificationToken(tokenValue string) (*models.EmailVerificationToken, error) {
	var token models.EmailVerificationToken
	err := r.db.Where("token = ? AND expires_at > ?", tokenValue, time.Now()).First(&token).Error
	return &token, err
}

// DeleteEmailVerificationToken deletes an email verification token after it has been used.
func (r *authRepository) DeleteEmailVerificationToken(token *models.EmailVerificationToken) error {
	return r.db.Delete(token).Error
}
