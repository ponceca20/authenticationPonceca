package audit

import (
	"math"
	"practicev2/module/authentication/models"
	"time"

	"github.com/google/uuid"
)

// PaginatedAuditResponse is a struct for paginated audit log responses.
type PaginatedAuditResponse struct {
	Data       []AuditLogResponseDTO `json:"data"`
	Total      int64                 `json:"total"`
	Page       int                   `json:"page"`
	PageSize   int                   `json:"page_size"`
	TotalPages int                   `json:"total_pages"`
}

// AuditService defines the interface for audit-related business logic.
type AuditService interface {
	ListAuditLogs(orgID string, query *AuditQueryDTO) (*PaginatedAuditResponse, error)
	LogAction(identityID string, action string, resource string, resourceID string, status string, organizationID *string, details string, ipAddress string, userAgent string) error
	GetAuditLogsByIdentityID(identityID string) ([]models.AuditLog, error)
}

type auditService struct {
	repo AuditRepository
}

// NewAuditService creates a new instance of AuditService.
func NewAuditService(repo AuditRepository) AuditService {
	return &auditService{repo: repo}
}

// ListAuditLogs retrieves a paginated list of audit logs.
func (s *auditService) ListAuditLogs(orgID string, query *AuditQueryDTO) (*PaginatedAuditResponse, error) {
	logs, total, err := s.repo.ListAuditLogs(orgID, query)
	if err != nil {
		return nil, err
	}

	var dtos []AuditLogResponseDTO
	for _, log := range logs {
		dtos = append(dtos, ToAuditLogResponseDTO(&log))
	}

	return &PaginatedAuditResponse{
		Data:       dtos,
		Total:      total,
		Page:       query.Page,
		PageSize:   query.PageSize,
		TotalPages: int(math.Ceil(float64(total) / float64(query.PageSize))),
	}, nil
}

// LogAction logs an action performed by a user to the audit log.
func (s *auditService) LogAction(identityID string, action string, resource string, resourceID string, status string, organizationID *string, details string, ipAddress string, userAgent string) error {
	auditLog := &models.AuditLog{
		ID:             uuid.New().String(),
		IdentityID:     identityID,
		Action:         action,
		Resource:       resource,
		ResourceID:     resourceID,
		Status:         status,
		OrganizationID: organizationID,
		Details:        details,
		IPAddress:      ipAddress,
		UserAgent:      userAgent,
		Timestamp:      time.Now(),
	}

	return s.repo.CreateAuditLog(auditLog)
}

// GetAuditLogsByIdentityID retrieves all audit logs for a specific identity.
func (s *auditService) GetAuditLogsByIdentityID(identityID string) ([]models.AuditLog, error) {
	return s.repo.GetAuditLogsByIdentityID(identityID)
}
