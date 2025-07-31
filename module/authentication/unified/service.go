package unified

import (
	"fmt"
	"practicev2/module/authentication/models"
)

// UnifiedService defines the interface for unified business logic.
type UnifiedService interface {
	GetDashboardData(identity *models.Identity) (*UnifiedDashboardDTO, error)
}

type unifiedService struct {
	repo UnifiedRepository
}

// NewUnifiedService creates a new instance of UnifiedService.
func NewUnifiedService(repo UnifiedRepository) UnifiedService {
	return &unifiedService{repo: repo}
}

// GetDashboardData compiles data for the unified dashboard.
func (s *unifiedService) GetDashboardData(identity *models.Identity) (*UnifiedDashboardDTO, error) {
	// This is a placeholder implementation.
	// A real implementation would call repository methods to get real data.
	dto := &UnifiedDashboardDTO{
		WelcomeMessage: fmt.Sprintf("Welcome back, %s!", identity.FirstName),
		RecentActivities: []ActivityDTO{
			{Description: "You logged in.", Timestamp: "Just now"},
		},
		PendingTasks: []TaskDTO{
			{Title: "Complete your profile", Link: "/api/v1/profiles/me"},
		},
	}
	return dto, nil
}
