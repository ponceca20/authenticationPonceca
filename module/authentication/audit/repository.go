package audit

import (
	"practicev2/module/authentication/models"

	"gorm.io/gorm"
)

// AuditRepository defines the interface for database operations related to audit logs.
type AuditRepository interface {
	ListAuditLogs(orgID string, query *AuditQueryDTO) ([]models.AuditLog, int64, error)
}

type auditRepository struct {
	db *gorm.DB
}

// NewAuditRepository creates a new instance of AuditRepository.
func NewAuditRepository(db *gorm.DB) AuditRepository {
	return &auditRepository{db: db}
}

// ListAuditLogs retrieves a paginated and filtered list of audit logs.
func (r *auditRepository) ListAuditLogs(orgID string, query *AuditQueryDTO) ([]models.AuditLog, int64, error) {
	var logs []models.AuditLog
	var total int64

	tx := r.db.Model(&models.AuditLog{}).Where("organization_id = ?", orgID)

	// Apply filters from the query DTO
	if query.UserID != "" {
		tx = tx.Where("identity_id = ?", query.UserID)
	}
	if query.Action != "" {
		tx = tx.Where("action LIKE ?", "%"+query.Action+"%")
	}
	if query.Resource != "" {
		tx = tx.Where("resource = ?", query.Resource)
	}
	if !query.StartDate.IsZero() {
		tx = tx.Where("timestamp >= ?", query.StartDate)
	}
	if !query.EndDate.IsZero() {
		tx = tx.Where("timestamp <= ?", query.EndDate)
	}

	// Count total records for pagination
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.PageSize <= 0 {
		query.PageSize = 10
	}
	offset := (query.Page - 1) * query.PageSize
	tx = tx.Offset(offset).Limit(query.PageSize)

	// Execute the query
	err := tx.Preload("Identity").Order("timestamp desc").Find(&logs).Error
	if err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}
