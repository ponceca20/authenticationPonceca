package utils

import (
	"fmt"
	"practicev2/config"
	"practicev2/module/authentication/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// MembershipClaim holds the claims related to a single organizational membership.
type MembershipClaim struct {
	OrganizationID string `json:"org_id"`
	Role           string `json:"role"`
}

// CustomerClaim holds the claims related to a customer profile.
type CustomerClaim struct {
	ProfileID string `json:"profile_id"`
}

// GuestClaim holds the claims related to a guest session.
type GuestClaim struct {
	SessionID string `json:"session_id"`
}

// UnifiedClaims represents the custom claims for the JWT.
// It contains the core identity and all possible contexts (memberships, customer, guest).
type UnifiedClaims struct {
	IdentityID   string            `json:"identity_id"`
	Email        string            `json:"email"`
	Memberships  []MembershipClaim `json:"memberships,omitempty"`
	Customer     *CustomerClaim    `json:"customer,omitempty"`
	Guest        *GuestClaim       `json:"guest,omitempty"`
	jwt.RegisteredClaims
}

// JWTService provides methods for handling JWTs.
type JWTService struct {
	accessSecret  string
	refreshSecret string
	issuer        string
	accessTTL     time.Duration
	refreshTTL    time.Duration
}

// NewJWTService creates a new JWTService with configuration from the config package.
func NewJWTService() *JWTService {
	return &JWTService{
		accessSecret:  config.GetJWTAccessSecret(),
		refreshSecret: config.GetJWTRefreshSecret(),
		issuer:        config.GetJWTIssuer(),
		accessTTL:     config.GetJWTAccessTTL(),
		refreshTTL:    config.GetJWTRefreshTTL(),
	}
}

// GetAccessTTL returns the configured time-to-live for an access token.
func (s *JWTService) GetAccessTTL() time.Duration {
	return s.accessTTL
}

// GetRefreshTTL returns the configured time-to-live for a refresh token.
func (s *JWTService) GetRefreshTTL() time.Duration {
	return s.refreshTTL
}

// GenerateTokenPair generates both an access and a refresh token for a given identity and contexts.
func (s *JWTService) GenerateTokenPair(identity *models.Identity, memberships []models.OrganizationalMembership, customer *models.CustomerProfile) (string, string, error) {
	// Generate claims
	claims := &UnifiedClaims{
		IdentityID: identity.ID,
		Email:      identity.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.accessTTL)),
			Issuer:    s.issuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	for _, m := range memberships {
		claims.Memberships = append(claims.Memberships, MembershipClaim{
			OrganizationID: m.OrganizationID,
			Role:           m.Role.Name,
		})
	}

	if customer != nil {
		claims.Customer = &CustomerClaim{ProfileID: customer.ID}
	}

	// Generate access token
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessString, err := accessToken.SignedString([]byte(s.accessSecret))
	if err != nil {
		return "", "", fmt.Errorf("failed to sign access token: %w", err)
	}

	// Generate refresh token
	refreshClaims := &UnifiedClaims{
		IdentityID: identity.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.refreshTTL)),
			Issuer:    s.issuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshString, err := refreshToken.SignedString([]byte(s.refreshSecret))
	if err != nil {
		return "", "", fmt.Errorf("failed to sign refresh token: %w", err)
	}

	return accessString, refreshString, nil
}

// ValidateToken validates a JWT string and returns the claims if the token is valid.
func (s *JWTService) ValidateToken(tokenString string) (*UnifiedClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UnifiedClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// Differentiate secret based on claims (e.g. presence of memberships) if needed,
		// for now we assume access token is always passed here for validation.
		return []byte(s.accessSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(*UnifiedClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// ValidateRefreshToken validates a refresh token string.
func (s *JWTService) ValidateRefreshToken(tokenString string) (*UnifiedClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UnifiedClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.refreshSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse refresh token: %w", err)
	}

	if claims, ok := token.Claims.(*UnifiedClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid refresh token")
}
