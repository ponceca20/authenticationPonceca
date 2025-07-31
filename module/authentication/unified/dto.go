package unified

// UnifiedDashboardDTO represents the data for a user's unified dashboard.
type UnifiedDashboardDTO struct {
	WelcomeMessage   string        `json:"welcome_message"`
	RecentActivities []ActivityDTO `json:"recent_activities"`
	PendingTasks     []TaskDTO     `json:"pending_tasks"`
}

// ActivityDTO represents a single recent activity item.
type ActivityDTO struct {
	Description string `json:"description"`
	Timestamp   string `json:"timestamp"`
}

// TaskDTO represents a single pending task for the user.
type TaskDTO struct {
	Title string `json:"title"`
	Link  string `json:"link"`
}
