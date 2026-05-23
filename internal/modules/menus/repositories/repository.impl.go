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

const tableName = "menus"

var ErrNotFound = errors.New("menu not found")

type MenuRepositoryImpl struct {
	db  *database.Database
	log *logger.LayerLogger
	sb  sq.StatementBuilderType
}

func NewMenuRepository(db *database.Database, appLogger *logger.Logger) MenuRepository {
	return &MenuRepositoryImpl{
		db:  db,
		log: appLogger.Layer("repository.menus"),
		sb:  sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *MenuRepositoryImpl) Find(ctx context.Context, params shared.ListParams) ([]entities.Menu, error) {
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

	return pgx.CollectRows(rows, pgx.RowToStructByName[entities.Menu])
}

func (r *MenuRepositoryImpl) FindByID(ctx context.Context, id int64) (*entities.Menu, error) {
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

	item, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entities.Menu])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (r *MenuRepositoryImpl) FindByCode(ctx context.Context, code string) (*entities.Menu, error) {
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

	item, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entities.Menu])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (r *MenuRepositoryImpl) FindByAppID(ctx context.Context, appID int64, params shared.ListParams) ([]entities.Menu, error) {
	builder := r.sb.Select(prefixedColumns("m")...).
		From(tableName + " m").
		Join("apps a ON a.id = m.app_id").
		Where(sq.Eq{"a.id": appID}).
		OrderBy("m.id DESC").
		Limit(params.Limit).
		Offset(params.Offset)

	builder = applyPrefixedSearch(builder, params, "m")

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[entities.Menu])
}

func (r *MenuRepositoryImpl) FindByAppCode(ctx context.Context, appCode string, params shared.ListParams) ([]entities.Menu, error) {
	builder := r.sb.Select(prefixedColumns("m")...).
		From(tableName + " m").
		Join("apps a ON a.id = m.app_id").
		Where(sq.Expr("LOWER(a.code) = LOWER(?)", appCode)).
		OrderBy("m.id DESC").
		Limit(params.Limit).
		Offset(params.Offset)

	builder = applyPrefixedSearch(builder, params, "m")

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[entities.Menu])
}

func (r *MenuRepositoryImpl) Create(ctx context.Context, entity entities.Menu) (*entities.Menu, error) {
	values := map[string]any{
		"app_id":                 entity.AppId,
		"module_id":              entity.ModuleId,
		"parent_id":              entity.ParentId,
		"code":                   entity.Code,
		"name":                   entity.Name,
		"route_path":             entity.RoutePath,
		"sort_order":             entity.SortOrder,
		"required_permission_id": entity.RequiredPermissionId,
		"status":                 entity.Status,
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

	created, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entities.Menu])
	if err != nil {
		return nil, err
	}

	return &created, nil
}

func (r *MenuRepositoryImpl) CreateBulk(ctx context.Context, items []entities.Menu) ([]entities.Menu, error) {
	if len(items) == 0 {
		return []entities.Menu{}, nil
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
			"parent_id",
			"code",
			"name",
			"route_path",
			"sort_order",
			"required_permission_id",
			"status",
		).
		Suffix("RETURNING " + columnList())

	for _, entity := range items {
		builder = builder.Values(
			entity.AppId,
			entity.ModuleId,
			entity.ParentId,
			entity.Code,
			entity.Name,
			entity.RoutePath,
			entity.SortOrder,
			entity.RequiredPermissionId,
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

	created, err := pgx.CollectRows(rows, pgx.RowToStructByName[entities.Menu])
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return created, nil
}

func (r *MenuRepositoryImpl) Update(ctx context.Context, id int64, data map[string]any) (*entities.Menu, error) {
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

	updated, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entities.Menu])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return &updated, nil
}

func (r *MenuRepositoryImpl) Delete(ctx context.Context, id int64) error {
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
		"parent_id",
		"code",
		"name",
		"route_path",
		"sort_order",
		"required_permission_id",
		"status",
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
	return "id, app_id, module_id, parent_id, code, name, route_path, sort_order, required_permission_id, status"
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
		sq.Expr("LOWER("+columnPrefix+"route_path) LIKE ?", pattern),
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
