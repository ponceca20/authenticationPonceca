package admin

// AdminService defines the interface for admin business logic.
type AdminService interface {
	GetSystemStats() (*SystemStatsDTO, error)
}

type adminService struct {
	repo AdminRepository
}

// NewAdminService creates a new instance of AdminService.
func NewAdminService(repo AdminRepository) AdminService {
	return &adminService{repo: repo}
}

// GetSystemStats compiles system statistics.
func (s *adminService) GetSystemStats() (*SystemStatsDTO, error) {
	totalUsers, _ := s.repo.CountUsers()
	totalOrgs, _ := s.repo.CountOrganizations()

	dto := &SystemStatsDTO{
		TotalUsers:         int(totalUsers),
		TotalOrganizations: int(totalOrgs),
		ActiveSessions:     0, // This would require cache/session store logic
	}
	return dto, nil
}
