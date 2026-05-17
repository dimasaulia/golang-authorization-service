package services

import (
	"context"

	"github.com/open-suite/authorization/internal/entities"
	"github.com/open-suite/authorization/internal/modules/auditlogs/dto"
)

type AuditLogService interface {
	Find(ctx context.Context, limit uint64, offset uint64) ([]entities.AuditLog, error)
	FindByID(ctx context.Context, id int64) (*entities.AuditLog, error)
	Create(ctx context.Context, request dto.CreateAuditLogRequest) (*entities.AuditLog, error)
	Update(ctx context.Context, id int64, request dto.UpdateAuditLogRequest) (*entities.AuditLog, error)
	Delete(ctx context.Context, id int64) error
}
