package admin

// SystemStatsDTO represents high-level system statistics.
type SystemStatsDTO struct {
	TotalUsers         int `json:"total_users"`
	TotalOrganizations int `json:"total_organizations"`
	ActiveSessions     int `json:"active_sessions"`
}
