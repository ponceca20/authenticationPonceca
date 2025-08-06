package auth

import (
	"errors"
	"fmt"
	"practicev2/module/authentication/models"
	"practicev2/module/authentication/utils"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AuthService defines the interface for authentication business logic.
type AuthService interface {
	Register(dto *RegisterDTO) (*models.Identity, error)
	Login(dto *LoginDTO) (*TokenResponseDTO, error)
	RefreshToken(dto *RefreshTokenDTO) (*TokenResponseDTO, error)
	Logout(dto *RefreshTokenDTO) error
	ForgotPassword(dto *ForgotPasswordDTO) error
	ResetPassword(dto *ResetPasswordDTO) error
	ChangePassword(identityID string, dto *ChangePasswordDTO) error
	VerifyEmail(dto *VerifyEmailDTO) error
	ResendVerification(dto *ResendVerificationDTO) error
}

type authService struct {
	repo       AuthRepository
	jwtService *utils.JWTService
}

// NewAuthService creates a new instance of AuthService.
func NewAuthService(repo AuthRepository, jwtService *utils.JWTService) AuthService {
	return &authService{
		repo:       repo,
		jwtService: jwtService,
	}
}

// Register handles the business logic for creating a new identity.
func (s *authService) Register(dto *RegisterDTO) (*models.Identity, error) {
	// Validate password complexity
	if err := utils.ValidatePassword(dto.Password); err != nil {
		return nil, err
	}

	// Check if identity already exists
	existing, err := s.repo.FindIdentityByEmail(dto.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err // Database error
	}
	if existing != nil {
		return nil, errors.New("email already in use")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(dto.Password)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Create new identity
	newIdentity := &models.Identity{
		ID:           uuid.New().String(),
		FirstName:    dto.FirstName,
		LastName:     dto.LastName,
		Email:        dto.Email,
		PasswordHash: hashedPassword,
	}

	if err := s.repo.CreateIdentity(newIdentity); err != nil {
		return nil, errors.New("failed to create identity")
	}

	return newIdentity, nil
}

// Login handles the business logic for authenticating an identity and generating tokens.
func (s *authService) Login(dto *LoginDTO) (*TokenResponseDTO, error) {
	const maxLoginAttempts = 5
	const lockoutDuration = 15 * time.Minute

	identity, err := s.repo.FindIdentityByEmail(dto.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid credentials")
		}
		return nil, err
	}

	// Check if account is currently locked
	if identity.LockedUntil != nil && identity.LockedUntil.After(time.Now()) {
		return nil, fmt.Errorf("account is locked until %v", identity.LockedUntil)
	}

	// Check password
	if !utils.CheckPasswordHash(dto.Password, identity.PasswordHash) {
		identity.FailedLoginAttempts++
		if identity.FailedLoginAttempts >= maxLoginAttempts {
			lockedUntil := time.Now().Add(lockoutDuration)
			identity.LockedUntil = &lockedUntil
		}
		if err := s.repo.UpdateIdentity(identity); err != nil {
			// Log the update error but still return invalid credentials
			fmt.Printf("Error updating identity after failed login: %v\n", err)
		}
		return nil, errors.New("invalid credentials")
	}

	// Password is correct, reset failure counter if needed
	if identity.FailedLoginAttempts > 0 || identity.LockedUntil != nil {
		identity.FailedLoginAttempts = 0
		identity.LockedUntil = nil
		if err := s.repo.UpdateIdentity(identity); err != nil {
			// Log the update error but proceed with login
			fmt.Printf("Error resetting failed login attempts: %v\n", err)
		}
	}

	// Update LastLoginAt timestamp
	now := time.Now()
	identity.LastLoginAt = &now
	if err := s.repo.UpdateIdentity(identity); err != nil {
		fmt.Printf("Error updating last login time: %v\n", err)
	}

	// Get user contexts for token generation
	fullIdentity, memberships, customer, err := s.repo.GetFullIdentityContext(identity.ID)
	if err != nil {
		return nil, errors.New("failed to retrieve user context")
	}

	return s._generateTokenResponse(fullIdentity, memberships, customer)
}

// RefreshToken handles the logic for refreshing an access token.
func (s *authService) RefreshToken(dto *RefreshTokenDTO) (*TokenResponseDTO, error) {
	// Validate the refresh token itself (is it a valid JWT?)
	claims, err := s.jwtService.ValidateRefreshToken(dto.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	// Find the token in the database to ensure it's not revoked
	dbToken, err := s.repo.FindRefreshToken(dto.RefreshToken)
	if err != nil {
		return nil, errors.New("refresh token not found or has been revoked")
	}

	// Revoke the used refresh token
	if err := s.repo.RevokeRefreshToken(dbToken); err != nil {
		return nil, errors.New("failed to revoke old refresh token")
	}

	// Get the full context for the user
	identity, memberships, customer, err := s.repo.GetFullIdentityContext(claims.IdentityID)
	if err != nil {
		return nil, errors.New("failed to retrieve user context for token refresh")
	}

	// Generate a new token response
	return s._generateTokenResponse(identity, memberships, customer)
}

// _generateTokenResponse is a private helper to create a token response from a user context.
func (s *authService) _generateTokenResponse(identity *models.Identity, memberships []models.OrganizationalMembership, customer *models.CustomerProfile) (*TokenResponseDTO, error) {
	accessToken, refreshToken, err := s.jwtService.GenerateTokenPair(identity, memberships, customer)
	if err != nil {
		return nil, errors.New("failed to generate tokens")
	}

	tokenResponse := &TokenResponseDTO{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(s.jwtService.GetAccessTTL()),
		Identity:     ToIdentityProfileDTO(identity),
	}

	refreshTokenModel := &models.RefreshToken{
		ID:         uuid.New().String(),
		IdentityID: identity.ID,
		Token:      refreshToken,
		ExpiresAt:  time.Now().Add(s.jwtService.GetRefreshTTL()),
	}
	if err := s.repo.SaveRefreshToken(refreshTokenModel); err != nil {
		return nil, errors.New("failed to save refresh token")
	}

	for _, m := range memberships {
		tokenResponse.Contexts = append(tokenResponse.Contexts, ContextSummaryDTO{
			Type: "organization",
			ID:   m.OrganizationID,
			Name: m.Organization.Name,
			Role: m.Role.Name,
		})
	}

	if customer != nil {
		tokenResponse.Contexts = append(tokenResponse.Contexts, ContextSummaryDTO{
			Type: "customer",
			ID:   customer.ID,
			Name: "E-commerce Profile",
		})
	}

	return tokenResponse, nil
}

// Logout handles revoking a refresh token.
func (s *authService) Logout(dto *RefreshTokenDTO) error {
	dbToken, err := s.repo.FindRefreshToken(dto.RefreshToken)
	if err != nil {
		// If token is not found, it's already effectively logged out.
		return nil
	}
	return s.repo.RevokeRefreshToken(dbToken)
}

// ForgotPassword handles the logic for initiating a password reset.
func (s *authService) ForgotPassword(dto *ForgotPasswordDTO) error {
	identity, err := s.repo.FindIdentityByEmail(dto.Email)
	if err != nil {
		// Don't reveal if the user exists or not.
		return nil
	}

	// Generate a secure, non-JWT token for password reset.
	token, err := utils.GenerateSecureToken(32)
	if err != nil {
		return errors.New("could not generate reset token")
	}

	resetToken := &models.PasswordResetToken{
		ID:         uuid.New().String(),
		IdentityID: identity.ID,
		Token:      token,
		ExpiresAt:  time.Now().Add(15 * time.Minute), // Short expiry
	}

	if err := s.repo.CreatePasswordResetToken(resetToken); err != nil {
		return errors.New("could not save reset token")
	}

	// In a real app, email the token to the user.
	// log.Printf("Password reset token for %s: %s", identity.Email, token)

	return nil
}

// ResetPassword handles the logic for resetting a password with a valid token.
func (s *authService) ResetPassword(dto *ResetPasswordDTO) error {
	resetToken, err := s.repo.FindPasswordResetToken(dto.Token)
	if err != nil {
		return errors.New("invalid or expired password reset token")
	}

	identity, err := s.repo.FindIdentityByID(resetToken.IdentityID)
	if err != nil {
		return errors.New("user associated with token not found")
	}

	newHashedPassword, err := utils.HashPassword(dto.NewPassword)
	if err != nil {
		return errors.New("failed to hash new password")
	}

	identity.PasswordHash = newHashedPassword
	if err := s.repo.UpdateIdentity(identity); err != nil {
		return errors.New("failed to update password")
	}

	// Delete the token so it can't be used again
	_ = s.repo.DeletePasswordResetToken(resetToken)

	return nil
}

// ChangePassword handles changing the password for an authenticated user.
func (s *authService) ChangePassword(identityID string, dto *ChangePasswordDTO) error {
	// Validate new password complexity
	if err := utils.ValidatePassword(dto.NewPassword); err != nil {
		return err
	}

	// Find the identity
	identity, err := s.repo.FindIdentityByID(identityID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return errors.New("failed to retrieve user")
	}

	// Verify current password
	if !utils.CheckPasswordHash(dto.CurrentPassword, identity.PasswordHash) {
		return errors.New("current password is incorrect")
	}

	// Ensure new password is different from current password
	if utils.CheckPasswordHash(dto.NewPassword, identity.PasswordHash) {
		return errors.New("new password must be different from current password")
	}

	// Hash new password
	newHashedPassword, err := utils.HashPassword(dto.NewPassword)
	if err != nil {
		return errors.New("failed to hash new password")
	}

	// Update password
	identity.PasswordHash = newHashedPassword
	if err := s.repo.UpdateIdentity(identity); err != nil {
		return errors.New("failed to update password")
	}

	return nil
}

// VerifyEmail handles the verification of email address using a token.
func (s *authService) VerifyEmail(dto *VerifyEmailDTO) error {
	// Find the verification token
	verificationToken, err := s.repo.FindEmailVerificationToken(dto.Token)
	if err != nil {
		return errors.New("invalid or expired verification token")
	}

	// Find the identity
	identity, err := s.repo.FindIdentityByID(verificationToken.IdentityID)
	if err != nil {
		return errors.New("user associated with token not found")
	}

	// Mark email as verified
	now := time.Now()
	identity.EmailVerifiedAt = &now
	if err := s.repo.UpdateIdentity(identity); err != nil {
		return errors.New("failed to verify email")
	}

	// Delete the verification token
	_ = s.repo.DeleteEmailVerificationToken(verificationToken)

	return nil
}

// ResendVerification handles resending email verification token.
func (s *authService) ResendVerification(dto *ResendVerificationDTO) error {
	identity, err := s.repo.FindIdentityByEmail(dto.Email)
	if err != nil {
		// Don't reveal if the user exists or not
		return nil
	}

	// Check if email is already verified
	if identity.EmailVerifiedAt != nil {
		return errors.New("email is already verified")
	}

	// Generate a secure token for email verification
	token, err := utils.GenerateSecureToken(32)
	if err != nil {
		return errors.New("could not generate verification token")
	}

	verificationToken := &models.EmailVerificationToken{
		ID:         uuid.New().String(),
		IdentityID: identity.ID,
		Token:      token,
		ExpiresAt:  time.Now().Add(24 * time.Hour), // 24 hour expiry
	}

	if err := s.repo.CreateEmailVerificationToken(verificationToken); err != nil {
		return errors.New("could not save verification token")
	}

	// In a real app, email the token to the user
	// log.Printf("Email verification token for %s: %s", identity.Email, token)

	return nil
}
