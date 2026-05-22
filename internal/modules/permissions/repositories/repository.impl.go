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

const tableName = "permissions"

var ErrNotFound = errors.New("permission not found")

type PermissionRepositoryImpl struct {
	db  *database.Database
	log *logger.LayerLogger
	sb  sq.StatementBuilderType
}

func NewPermissionRepository(db *database.Database, appLogger *logger.Logger) PermissionRepository {
	return &PermissionRepositoryImpl{
		db:  db,
		log: appLogger.Layer("repository.permissions"),
		sb:  sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *PermissionRepositoryImpl) Find(ctx context.Context, params shared.ListParams) ([]entities.Permission, error) {
	builder := r.sb.Select(columns()...).
		From(tableName).
		OrderBy("id DESC").
		Limit(params.Limit).
		Offset(params.Offset)

	builder = applySearch(builder, params)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[entities.Permission])
}

func (r *PermissionRepositoryImpl) FindByID(ctx context.Context, id int64) (*entities.Permission, error) {
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

	item, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entities.Permission])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (r *PermissionRepositoryImpl) FindByAppID(ctx context.Context, appID int64, params shared.ListParams) ([]entities.Permission, error) {
	builder := r.sb.Select(prefixedColumns("p")...).
		From(tableName + " p").
		Join("apps a ON a.id = p.app_id").
		Where(sq.Eq{"a.id": appID}).
		OrderBy("p.id DESC").
		Limit(params.Limit).
		Offset(params.Offset)

	builder = applyPrefixedSearch(builder, params, "p")

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[entities.Permission])
}

func (r *PermissionRepositoryImpl) FindByAppCode(ctx context.Context, appCode string, params shared.ListParams) ([]entities.Permission, error) {
	builder := r.sb.Select(prefixedColumns("p")...).
		From(tableName + " p").
		Join("apps a ON a.id = p.app_id").
		Where(sq.Expr("LOWER(a.code) = LOWER(?)", appCode)).
		OrderBy("p.id DESC").
		Limit(params.Limit).
		Offset(params.Offset)

	builder = applyPrefixedSearch(builder, params, "p")

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[entities.Permission])
}

func (r *PermissionRepositoryImpl) FindByCode(ctx context.Context, code string) (*entities.Permission, error) {
	query, args, err := r.sb.Select(columns()...).
		From(tableName).
		Where(sq.Expr("LOWER(code) = LOWER(?)", code)).
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

	item, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entities.Permission])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (r *PermissionRepositoryImpl) Create(ctx context.Context, entity entities.Permission) (*entities.Permission, error) {
	values := map[string]any{
		"app_id":      entity.AppId,
		"module_id":   entity.ModuleId,
		"action_id":   entity.ActionId,
		"code":        entity.Code,
		"name":        entity.Name,
		"description": entity.Description,
		"risk_level":  entity.RiskLevel,
		"is_system":   entity.IsSystem,
		"status":      entity.Status,
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

	created, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entities.Permission])
	if err != nil {
		return nil, err
	}

	return &created, nil
}

func (r *PermissionRepositoryImpl) CreateBulk(ctx context.Context, items []entities.Permission) ([]entities.Permission, error) {
	if len(items) == 0 {
		return []entities.Permission{}, nil
	}

	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	builder := r.sb.Insert(tableName).
		Columns(
			"app_id",
			"module_id",
			"action_id",
			"code",
			"name",
			"description",
			"risk_level",
			"is_system",
			"status",
		).
		Suffix("RETURNING " + columnList())

	for _, entity := range items {
		builder = builder.Values(
			entity.AppId,
			entity.ModuleId,
			entity.ActionId,
			entity.Code,
			entity.Name,
			entity.Description,
			entity.RiskLevel,
			entity.IsSystem,
			entity.Status,
		)
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := tx.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	created, err := pgx.CollectRows(rows, pgx.RowToStructByName[entities.Permission])
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return created, nil
}

func (r *PermissionRepositoryImpl) Update(ctx context.Context, id int64, data map[string]any) (*entities.Permission, error) {
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

	updated, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entities.Permission])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return &updated, nil
}

func (r *PermissionRepositoryImpl) Delete(ctx context.Context, id int64) error {
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
		"app_id",
		"module_id",
		"action_id",
		"code",
		"name",
		"description",
		"risk_level",
		"is_system",
		"status",
		"created_at",
		"updated_at",
	}
}

func prefixedColumns(prefix string) []string {
	values := columns()
	for index, column := range values {
		values[index] = prefix + "." + column
	}
	return values
}

func columnList() string {
	return "id, app_id, module_id, action_id, code, name, description, risk_level, is_system, status, created_at, updated_at"
}

func applySearch(builder sq.SelectBuilder, params shared.ListParams) sq.SelectBuilder {
	return applyPrefixedSearch(builder, params, "")
}

func applyPrefixedSearch(builder sq.SelectBuilder, params shared.ListParams, prefix string) sq.SelectBuilder {
	if params.Search == "" {
		return builder
	}

	columnPrefix := ""
	if prefix != "" {
		columnPrefix = prefix + "."
	}
	pattern := "%" + params.Search + "%"
	return builder.Where(sq.Or{
		sq.Expr("LOWER("+columnPrefix+"code) LIKE ?", pattern),
		sq.Expr("LOWER("+columnPrefix+"name) LIKE ?", pattern),
	})
}

func canUpdateTimestamp() bool {
	for _, column := range columns() {
		if column == "updated_at" {
			return true
		}
	}
	return false
}
