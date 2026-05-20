package repository

import (
	"github.com/Lapnes/pos-kopitiam-v3nf/internal/models"
	"gorm.io/gorm"
)

type AuditRepository interface {
	Create(log *models.AuditLog) error
	FindByTableAndRecord(tableName, recordID string) ([]models.AuditLog, error)
}

type auditRepository struct {
	db *gorm.DB
}

func NewAuditRepository(db *gorm.DB) AuditRepository {
	return &auditRepository{db: db}
}

func (r *auditRepository) Create(log *models.AuditLog) error {
	return r.db.Create(log).Error
}

func (r *auditRepository) FindByTableAndRecord(tableName, recordID string) ([]models.AuditLog, error) {
	var logs []models.AuditLog
	err := r.db.Where("table_name = ? AND record_id = ?", tableName, recordID).
		Order("created_at desc").Find(&logs).Error
	return logs, err
}