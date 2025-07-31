package utils

import (
	"practicev2/config"
	"practicev2/module/authentication/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestJWTService_GenerateAndValidateToken(t *testing.T) {
	// Setup
	config.Init() // Ensure config is loaded with defaults
	jwtService := NewJWTService()
	identity := &models.Identity{
		ID:    "user-123",
		Email: "test@example.com",
	}
	memberships := []models.OrganizationalMembership{
		{
			OrganizationID: "org-456",
			Role:           models.Role{Name: "admin"},
		},
	}
	customerProfile := &models.CustomerProfile{
		ID: "cust-789",
	}

	// Act: Generate tokens
	accessToken, refreshToken, err := jwtService.GenerateTokenPair(identity, memberships, customerProfile)

	// Assert: Generation
	assert.NoError(t, err)
	assert.NotEmpty(t, accessToken)
	assert.NotEmpty(t, refreshToken)

	// Act: Validate Access Token
	claims, err := jwtService.ValidateToken(accessToken)

	// Assert: Validation
	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, identity.ID, claims.IdentityID)
	assert.Equal(t, identity.Email, claims.Email)
	assert.Equal(t, jwtService.issuer, claims.Issuer)
	assert.Len(t, claims.Memberships, 1)
	assert.Equal(t, "org-456", claims.Memberships[0].OrganizationID)
	assert.Equal(t, "admin", claims.Memberships[0].Role)
	assert.NotNil(t, claims.Customer)
	assert.Equal(t, "cust-789", claims.Customer.ProfileID)
}

func TestJWTService_ExpiredToken(t *testing.T) {
	// Setup: Create a service with a very short TTL
	jwtService := &JWTService{
		accessSecret: "test-secret",
		issuer:       "test-issuer",
		accessTTL:    -1 * time.Minute, // Expired in the past
	}
	identity := &models.Identity{ID: "user-123", Email: "test@example.com"}

	// Act: Generate token
	accessToken, _, err := jwtService.GenerateTokenPair(identity, nil, nil)
	assert.NoError(t, err)

	// Act: Validate the expired token
	// We need a new service instance with the correct secret to validate
	validationService := &JWTService{accessSecret: "test-secret"}
	claims, err := validationService.ValidateToken(accessToken)

	// Assert: Validation fails
	assert.Error(t, err)
	assert.Nil(t, claims)
	assert.Contains(t, err.Error(), "token is expired")
}

func TestJWTService_InvalidSignature(t *testing.T) {
	// Setup
	jwtService1 := NewJWTService()
	jwtService2 := &JWTService{ // Service with a different secret
		accessSecret: "a-different-secret",
	}
	identity := &models.Identity{ID: "user-123", Email: "test@example.com"}

	// Act: Generate a token with the first service
	accessToken, _, err := jwtService1.GenerateTokenPair(identity, nil, nil)
	assert.NoError(t, err)

	// Act: Validate the token with the second service (wrong secret)
	claims, err := jwtService2.ValidateToken(accessToken)

	// Assert: Validation fails
	assert.Error(t, err)
	assert.Nil(t, claims)
	assert.Contains(t, err.Error(), "signature is invalid")
}
