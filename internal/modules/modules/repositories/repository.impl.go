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

const tableName = "modules"

var ErrNotFound = errors.New("module not found")

type ModuleRepositoryImpl struct {
	db  *database.Database
	log *logger.LayerLogger
	sb  sq.StatementBuilderType
}

func NewModuleRepository(db *database.Database, appLogger *logger.Logger) ModuleRepository {
	return &ModuleRepositoryImpl{
		db:  db,
		log: appLogger.Layer("repository.modules"),
		sb:  sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *ModuleRepositoryImpl) Find(ctx context.Context, params shared.ListParams) ([]entities.Module, error) {
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

	return pgx.CollectRows(rows, pgx.RowToStructByName[entities.Module])
}

func (r *ModuleRepositoryImpl) FindByID(ctx context.Context, id int64) (*entities.Module, error) {
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

	item, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entities.Module])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (r *ModuleRepositoryImpl) FindByAppID(ctx context.Context, appID int64, params shared.ListParams) ([]entities.Module, error) {
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

	return pgx.CollectRows(rows, pgx.RowToStructByName[entities.Module])
}

func (r *ModuleRepositoryImpl) FindByAppCode(ctx context.Context, appCode string, params shared.ListParams) ([]entities.Module, error) {
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

	return pgx.CollectRows(rows, pgx.RowToStructByName[entities.Module])
}

func (r *ModuleRepositoryImpl) FindByCode(ctx context.Context, code string) (*entities.Module, error) {
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

	item, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entities.Module])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (r *ModuleRepositoryImpl) Create(ctx context.Context, entity entities.Module) (*entities.Module, error) {
	values := map[string]any{
		"app_id": entity.AppId,
		"code":   entity.Code,
		"name":   entity.Name,
		"status": entity.Status,
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

	created, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entities.Module])
	if err != nil {
		return nil, err
	}

	return &created, nil
}

func (r *ModuleRepositoryImpl) Update(ctx context.Context, id int64, data map[string]any) (*entities.Module, error) {
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

	updated, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entities.Module])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return &updated, nil
}

func (r *ModuleRepositoryImpl) Delete(ctx context.Context, id int64) error {
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
		"code",
		"name",
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
	return "id, app_id, code, name, status"
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
