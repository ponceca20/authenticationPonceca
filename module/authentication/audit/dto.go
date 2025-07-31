package audit

import (
	"practicev2/module/authentication/models"
	"time"
)

// AuditLogResponseDTO is a public representation of an audit log entry.
type AuditLogResponseDTO struct {
	ID             string    `json:"id"`
	OrganizationID *string   `json:"organization_id,omitempty"`
	IdentityID     string    `json:"identity_id"`
	UserEmail      string    `json:"user_email"`
	Action         string    `json:"action"`
	Resource       string    `json:"resource,omitempty"`
	ResourceID     string    `json:"resource_id,omitempty"`
	Status         string    `json:"status"`
	IPAddress      string    `json:"ip_address,omitempty"`
	Details        string    `json:"details,omitempty"`
	Timestamp      time.Time `json:"timestamp"`
}

// AuditQueryDTO represents the query parameters for filtering audit logs.
// These would be parsed from the URL query string in the handler.
type AuditQueryDTO struct {
	UserID    string    `query:"user_id"`
	Action    string    `query:"action"`
	Resource  string    `query:"resource"`
	StartDate time.Time `query:"start_date"`
	EndDate   time.Time `query:"end_date"`
	Page      int       `query:"page"`
	PageSize  int       `query:"page_size"`
}

// ToAuditLogResponseDTO converts an AuditLog model to a public DTO.
func ToAuditLogResponseDTO(log *models.AuditLog) AuditLogResponseDTO {
	return AuditLogResponseDTO{
		ID:             log.ID,
		OrganizationID: log.OrganizationID,
		IdentityID:     log.IdentityID,
		UserEmail:      log.Identity.Email, // Assumes Identity is preloaded
		Action:         log.Action,
		Resource:       log.Resource,
		ResourceID:     log.ResourceID,
		Status:         log.Status,
		IPAddress:      log.IPAddress,
		Details:        log.Details,
		Timestamp:      log.Timestamp,
	}
}
