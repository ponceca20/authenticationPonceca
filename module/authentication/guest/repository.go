package guest

import (
	"practicev2/module/authentication/models"

	"gorm.io/gorm"
)

// GuestRepository defines the interface for database operations related to guest sessions.
type GuestRepository interface {
	CreateSession(session *models.GuestSession) error
	FindSessionByToken(token string) (*models.GuestSession, error)
	UpdateSession(session *models.GuestSession) error
	DeleteSession(token string) error
}

type guestRepository struct {
	db *gorm.DB
}

// NewGuestRepository creates a new instance of GuestRepository.
func NewGuestRepository(db *gorm.DB) GuestRepository {
	return &guestRepository{db: db}
}

// CreateSession saves a new guest session to the database.
func (r *guestRepository) CreateSession(session *models.GuestSession) error {
	return r.db.Create(session).Error
}

// FindSessionByToken finds a guest session by its unique session token.
func (r *guestRepository) FindSessionByToken(token string) (*models.GuestSession, error) {
	var session models.GuestSession
	if err := r.db.Where("session_token = ? AND expires_at > NOW()", token).First(&session).Error; err != nil {
		return nil, err
	}
	return &session, nil
}

// UpdateSession updates an existing guest session.
func (r *guestRepository) UpdateSession(session *models.GuestSession) error {
	return r.db.Save(session).Error
}

// DeleteSession removes a guest session from the database.
func (r *guestRepository) DeleteSession(token string) error {
	return r.db.Where("session_token = ?", token).Delete(&models.GuestSession{}).Error
}
