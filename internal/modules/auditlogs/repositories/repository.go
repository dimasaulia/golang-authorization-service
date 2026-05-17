package repositories

import (
	"context"

	"github.com/open-suite/authorization/internal/entities"
)

type AuditLogRepository interface {
	Find(ctx context.Context, limit uint64, offset uint64) ([]entities.AuditLog, error)
	FindByID(ctx context.Context, id int64) (*entities.AuditLog, error)
	Create(ctx context.Context, entity entities.AuditLog) (*entities.AuditLog, error)
	Update(ctx context.Context, id int64, data map[string]any) (*entities.AuditLog, error)
	Delete(ctx context.Context, id int64) error
}
