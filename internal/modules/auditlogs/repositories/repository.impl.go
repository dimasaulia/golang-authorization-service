package repositories

import (
	"context"
	"errors"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"

	"github.com/open-suite/authorization/internal/entities"
	"github.com/open-suite/authorization/internal/platform/database"
	"github.com/open-suite/authorization/internal/platform/logger"
)

const tableName = "audit_logs"

var ErrNotFound = errors.New("audit log not found")

type AuditLogRepositoryImpl struct {
	db  *database.Database
	log *logger.LayerLogger
	sb  sq.StatementBuilderType
}

func NewAuditLogRepository(db *database.Database, appLogger *logger.Logger) AuditLogRepository {
	return &AuditLogRepositoryImpl{
		db:  db,
		log: appLogger.Layer("repository.audit_logs"),
		sb:  sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *AuditLogRepositoryImpl) Find(ctx context.Context, limit uint64, offset uint64) ([]entities.AuditLog, error) {
	query, args, err := r.sb.Select(columns()...).
		From(tableName).
		OrderBy("id DESC").
		Limit(limit).
		Offset(offset).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[entities.AuditLog])
}

func (r *AuditLogRepositoryImpl) FindByID(ctx context.Context, id int64) (*entities.AuditLog, error) {
	query, args, err := r.sb.Select(columns()...).
		From(tableName).
		Where(sq.Eq{"id": id}).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	item, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entities.AuditLog])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (r *AuditLogRepositoryImpl) Create(ctx context.Context, entity entities.AuditLog) (*entities.AuditLog, error) {
	values := map[string]any{
		"organization_id": entity.OrganizationId,
		"app_id":          entity.AppId,
		"actor_user_id":   entity.ActorUserId,
		"target_user_id":  entity.TargetUserId,
		"permission_id":   entity.PermissionId,
		"action":          entity.Action,
		"resource_type":   entity.ResourceType,
		"resource_id":     entity.ResourceId,
		"result":          entity.Result,
		"metadata_json":   entity.MetadataJson,
		"ip_address":      entity.IpAddress,
		"user_agent":      entity.UserAgent,
	}

	query, args, err := r.sb.Insert(tableName).
		SetMap(values).
		Suffix("RETURNING " + columnList()).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	created, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entities.AuditLog])
	if err != nil {
		return nil, err
	}

	return &created, nil
}

func (r *AuditLogRepositoryImpl) Update(ctx context.Context, id int64, data map[string]any) (*entities.AuditLog, error) {
	if len(data) == 0 {
		return r.FindByID(ctx, id)
	}
	if canUpdateTimestamp() {
		data["updated_at"] = time.Now()
	}

	query, args, err := r.sb.Update(tableName).
		SetMap(data).
		Where(sq.Eq{"id": id}).
		Suffix("RETURNING " + columnList()).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	updated, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entities.AuditLog])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return &updated, nil
}

func (r *AuditLogRepositoryImpl) Delete(ctx context.Context, id int64) error {
	query, args, err := r.sb.Delete(tableName).
		Where(sq.Eq{"id": id}).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return err
	}

	var deletedID int64
	if err := r.db.Pool.QueryRow(ctx, query, args...).Scan(&deletedID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}
		return err
	}

	return nil
}

func columns() []string {
	return []string{
		"id",
		"organization_id",
		"app_id",
		"actor_user_id",
		"target_user_id",
		"permission_id",
		"action",
		"resource_type",
		"resource_id",
		"result",
		"metadata_json",
		"ip_address",
		"user_agent",
		"created_at",
	}
}

func columnList() string {
	return "id, organization_id, app_id, actor_user_id, target_user_id, permission_id, action, resource_type, resource_id, result, metadata_json, ip_address, user_agent, created_at"
}

func canUpdateTimestamp() bool {
	for _, column := range columns() {
		if column == "updated_at" {
			return true
		}
	}
	return false
}
