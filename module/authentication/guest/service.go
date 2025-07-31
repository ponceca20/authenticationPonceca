package guest

import (
	"fmt"
	"practicev2/module/authentication/models"
	"practicev2/module/authentication/utils"
	"time"

	"github.com/google/uuid"
)

// GuestService defines the interface for guest session business logic.
type GuestService interface {
	CreateGuestSession() (*models.GuestSession, error)
	GetSession(token string) (*models.GuestSession, error)
	UpdateCart(token, cartData string) (*models.GuestSession, error)
}

type guestService struct {
	repo GuestRepository
}

// NewGuestService creates a new instance of GuestService.
func NewGuestService(repo GuestRepository) GuestService {
	return &guestService{repo: repo}
}

// CreateGuestSession creates a new, secure guest session.
func (s *guestService) CreateGuestSession() (*models.GuestSession, error) {
	token, err := utils.GenerateSecureToken(32)
	if err != nil {
		return nil, fmt.Errorf("failed to generate session token: %w", err)
	}

	session := &models.GuestSession{
		ID:           uuid.New().String(),
		SessionToken: token,
		ExpiresAt:    time.Now().Add(24 * time.Hour), // Session expires in 24 hours
		LastActivity: time.Now(),
	}

	if err := s.repo.CreateSession(session); err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return session, nil
}

// GetSession retrieves a guest session by its token.
func (s *guestService) GetSession(token string) (*models.GuestSession, error) {
	return s.repo.FindSessionByToken(token)
}

// UpdateCart updates the cart data for a given session.
func (s *guestService) UpdateCart(token, cartData string) (*models.GuestSession, error) {
	session, err := s.repo.FindSessionByToken(token)
	if err != nil {
		return nil, fmt.Errorf("invalid or expired session token")
	}

	session.CartData = cartData
	session.LastActivity = time.Now()

	if err := s.repo.UpdateSession(session); err != nil {
		return nil, fmt.Errorf("failed to update session: %w", err)
	}

	return session, nil
}
