package audit

import "math"

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
