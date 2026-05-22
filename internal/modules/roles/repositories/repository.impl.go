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
	"github.com/open-suite/authorization/internal/shared"
)

const tableName = "roles"

var ErrNotFound = errors.New("role not found")

type RoleRepositoryImpl struct {
	db  *database.Database
	log *logger.LayerLogger
	sb  sq.StatementBuilderType
}

func NewRoleRepository(db *database.Database, appLogger *logger.Logger) RoleRepository {
	return &RoleRepositoryImpl{
		db:  db,
		log: appLogger.Layer("repository.roles"),
		sb:  sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *RoleRepositoryImpl) Find(ctx context.Context, params shared.ListParams) ([]entities.Role, error) {
	builder := r.sb.Select(columns()...).
		From(tableName).
		OrderBy("id DESC").
		Limit(params.Limit).
		Offset(params.Offset)

	if params.Search != "" {
		pattern := "%" + params.Search + "%"
		builder = builder.Where(sq.Or{
			sq.Expr("LOWER(code) LIKE ?", pattern),
			sq.Expr("LOWER(name) LIKE ?", pattern),
			sq.Expr("LOWER(description) LIKE ?", pattern),
			sq.Expr("LOWER(scope) LIKE ?", pattern),
		})
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[entities.Role])
}

func (r *RoleRepositoryImpl) FindByID(ctx context.Context, id int64) (*entities.Role, error) {
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

	item, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entities.Role])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (r *RoleRepositoryImpl) Create(ctx context.Context, entity entities.Role) (*entities.Role, error) {
	values := map[string]any{
		"organization_id": entity.OrganizationId,
		"app_id":          entity.AppId,
		"code":            entity.Code,
		"name":            entity.Name,
		"description":     entity.Description,
		"scope":           entity.Scope,
		"is_system":       entity.IsSystem,
		"status":          entity.Status,
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

	created, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entities.Role])
	if err != nil {
		return nil, err
	}

	return &created, nil
}

func (r *RoleRepositoryImpl) Update(ctx context.Context, id int64, data map[string]any) (*entities.Role, error) {
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

	updated, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entities.Role])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return &updated, nil
}

func (r *RoleRepositoryImpl) Delete(ctx context.Context, id int64) error {
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
		"code",
		"name",
		"description",
		"scope",
		"is_system",
		"status",
		"created_at",
		"updated_at",
	}
}

func columnList() string {
	return "id, organization_id, app_id, code, name, description, scope, is_system, status, created_at, updated_at"
}

func canUpdateTimestamp() bool {
	for _, column := range columns() {
		if column == "updated_at" {
			return true
		}
	}
	return false
}
