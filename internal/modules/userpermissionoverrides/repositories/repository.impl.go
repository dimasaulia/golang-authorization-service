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

const tableName = "user_permission_overrides"

var ErrNotFound = errors.New("user permission override not found")

type UserPermissionOverrideRepositoryImpl struct {
	db  *database.Database
	log *logger.LayerLogger
	sb  sq.StatementBuilderType
}

func NewUserPermissionOverrideRepository(db *database.Database, appLogger *logger.Logger) UserPermissionOverrideRepository {
	return &UserPermissionOverrideRepositoryImpl{
		db:  db,
		log: appLogger.Layer("repository.user_permission_overrides"),
		sb:  sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *UserPermissionOverrideRepositoryImpl) Find(ctx context.Context, limit uint64, offset uint64) ([]entities.UserPermissionOverride, error) {
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

	return pgx.CollectRows(rows, pgx.RowToStructByName[entities.UserPermissionOverride])
}

func (r *UserPermissionOverrideRepositoryImpl) FindByID(ctx context.Context, id int64) (*entities.UserPermissionOverride, error) {
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

	item, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entities.UserPermissionOverride])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (r *UserPermissionOverrideRepositoryImpl) Create(ctx context.Context, entity entities.UserPermissionOverride) (*entities.UserPermissionOverride, error) {
	values := map[string]any{
		"user_id":       entity.UserId,
		"permission_id": entity.PermissionId,
		"effect":        entity.Effect,
		"reason":        entity.Reason,
		"expires_at":    entity.ExpiresAt,
		"created_by":    entity.CreatedBy,
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

	created, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entities.UserPermissionOverride])
	if err != nil {
		return nil, err
	}

	return &created, nil
}

func (r *UserPermissionOverrideRepositoryImpl) Update(ctx context.Context, id int64, data map[string]any) (*entities.UserPermissionOverride, error) {
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

	updated, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entities.UserPermissionOverride])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return &updated, nil
}

func (r *UserPermissionOverrideRepositoryImpl) Delete(ctx context.Context, id int64) error {
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
		"user_id",
		"permission_id",
		"effect",
		"reason",
		"expires_at",
		"created_by",
		"created_at",
	}
}

func columnList() string {
	return "id, user_id, permission_id, effect, reason, expires_at, created_by, created_at"
}

func canUpdateTimestamp() bool {
	for _, column := range columns() {
		if column == "updated_at" {
			return true
		}
	}
	return false
}
