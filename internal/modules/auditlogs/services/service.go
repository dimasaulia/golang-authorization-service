package services

import (
	"context"

	"github.com/open-suite/authorization/internal/entities"
	"github.com/open-suite/authorization/internal/modules/auditlogs/dto"
	"github.com/open-suite/authorization/internal/shared"
)

type AuditLogService interface {
	Find(ctx context.Context, params shared.ListParams) ([]entities.AuditLog, error)
	FindByID(ctx context.Context, id int64) (*entities.AuditLog, error)
	Create(ctx context.Context, request dto.CreateAuditLogRequest) (*entities.AuditLog, error)
	Update(ctx context.Context, id int64, request dto.UpdateAuditLogRequest) (*entities.AuditLog, error)
	Delete(ctx context.Context, id int64) error
}
