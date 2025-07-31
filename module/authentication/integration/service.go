package integration

import "practicev2/module/authentication/utils"

// IntegrationService defines the interface for integration business logic.
type IntegrationService interface {
	ValidateToken(token string) (*ValidateTokenResponseDTO, error)
}

type integrationService struct {
	jwtService *utils.JWTService
}

// NewIntegrationService creates a new instance of IntegrationService.
func NewIntegrationService(jwtService *utils.JWTService) IntegrationService {
	return &integrationService{jwtService: jwtService}
}

// ValidateToken validates a token for external services.
func (s *integrationService) ValidateToken(token string) (*ValidateTokenResponseDTO, error) {
	claims, err := s.jwtService.ValidateToken(token)
	if err != nil {
		return &ValidateTokenResponseDTO{IsValid: false}, nil
	}

	return &ValidateTokenResponseDTO{
		IsValid: true,
		UserID:  claims.IdentityID,
		Email:   claims.Email,
	}, nil
}
