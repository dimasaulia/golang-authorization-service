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

const tableName = "role_permissions"

var ErrNotFound = errors.New("role permission not found")

type RolePermissionRepositoryImpl struct {
	db  *database.Database
	log *logger.LayerLogger
	sb  sq.StatementBuilderType
}

func NewRolePermissionRepository(db *database.Database, appLogger *logger.Logger) RolePermissionRepository {
	return &RolePermissionRepositoryImpl{
		db:  db,
		log: appLogger.Layer("repository.role_permissions"),
		sb:  sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *RolePermissionRepositoryImpl) Find(ctx context.Context, params shared.ListParams) ([]entities.RolePermission, error) {
	query, args, err := r.sb.Select(columns()...).
		From(tableName).
		OrderBy("id DESC").
		Limit(params.Limit).
		Offset(params.Offset).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[entities.RolePermission])
}

func (r *RolePermissionRepositoryImpl) FindByAppID(ctx context.Context, appID int64, params shared.ListParams) ([]entities.RolePermissionDetail, error) {
	builder := r.detailBuilder(params).
		Where(sq.Eq{"a.id": appID})

	return r.findDetails(ctx, builder)
}

func (r *RolePermissionRepositoryImpl) FindByAppCode(ctx context.Context, appCode string, params shared.ListParams) ([]entities.RolePermissionDetail, error) {
	builder := r.detailBuilder(params).
		Where(sq.Expr("LOWER(a.code) = LOWER(?)", appCode))

	return r.findDetails(ctx, builder)
}

func (r *RolePermissionRepositoryImpl) FindByRole(ctx context.Context, roleID int64, params shared.ListParams) ([]entities.RolePermissionDetail, error) {
	builder := r.detailBuilder(params).
		Where(sq.Eq{"rp.role_id": roleID})

	return r.findDetails(ctx, builder)
}

func (r *RolePermissionRepositoryImpl) FindRoleIDByCode(ctx context.Context, roleCode string) (int64, error) {
	query, args, err := r.sb.Select("id").
		From("roles").
		Where(sq.Expr("LOWER(code) = LOWER(?)", roleCode)).
		Limit(1).
		ToSql()
	if err != nil {
		return 0, err
	}

	var roleID int64
	if err := r.db.Pool.QueryRow(ctx, query, args...).Scan(&roleID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, ErrNotFound
		}
		return 0, err
	}

	return roleID, nil
}

func (r *RolePermissionRepositoryImpl) FindRoleSummaries(ctx context.Context, params shared.ListParams) ([]entities.RolePermissionSummary, error) {
	return r.findRoleSummaries(ctx, r.roleSummaryBuilder(params))
}

func (r *RolePermissionRepositoryImpl) FindRoleSummariesByAppID(ctx context.Context, appID int64, params shared.ListParams) ([]entities.RolePermissionSummary, error) {
	builder := r.roleSummaryBuilder(params).
		Where(sq.Eq{"a.id": appID})

	return r.findRoleSummaries(ctx, builder)
}

func (r *RolePermissionRepositoryImpl) FindRoleSummariesByAppCode(ctx context.Context, appCode string, params shared.ListParams) ([]entities.RolePermissionSummary, error) {
	builder := r.roleSummaryBuilder(params).
		Where(sq.Expr("LOWER(a.code) = LOWER(?)", appCode))

	return r.findRoleSummaries(ctx, builder)
}

func (r *RolePermissionRepositoryImpl) FindAvailablePermissionsByAppID(ctx context.Context, appID int64, params shared.ListParams) ([]entities.AvailablePermissionRow, error) {
	builder := r.availablePermissionBuilder(params).
		Where(sq.Eq{"a.id": appID})

	return r.findAvailablePermissions(ctx, builder)
}

func (r *RolePermissionRepositoryImpl) FindAvailablePermissionsByAppCode(ctx context.Context, appCode string, params shared.ListParams) ([]entities.AvailablePermissionRow, error) {
	builder := r.availablePermissionBuilder(params).
		Where(sq.Expr("LOWER(a.code) = LOWER(?)", appCode))

	return r.findAvailablePermissions(ctx, builder)
}

func (r *RolePermissionRepositoryImpl) FindByID(ctx context.Context, id int64) (*entities.RolePermission, error) {
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

	item, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entities.RolePermission])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (r *RolePermissionRepositoryImpl) Create(ctx context.Context, entity entities.RolePermission) (*entities.RolePermission, error) {
	values := map[string]any{
		"role_id":       entity.RoleId,
		"permission_id": entity.PermissionId,
		"effect":        entity.Effect,
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

	created, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entities.RolePermission])
	if err != nil {
		return nil, err
	}

	return &created, nil
}

func (r *RolePermissionRepositoryImpl) CreateBulk(ctx context.Context, items []entities.RolePermission) ([]entities.RolePermission, error) {
	if len(items) == 0 {
		return []entities.RolePermission{}, nil
	}

	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	builder := r.sb.Insert(tableName).
		Columns("role_id", "permission_id", "effect").
		Suffix("RETURNING " + columnList())

	for _, entity := range items {
		builder = builder.Values(entity.RoleId, entity.PermissionId, entity.Effect)
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

	created, err := pgx.CollectRows(rows, pgx.RowToStructByName[entities.RolePermission])
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return created, nil
}

func (r *RolePermissionRepositoryImpl) Update(ctx context.Context, id int64, data map[string]any) (*entities.RolePermission, error) {
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

	updated, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entities.RolePermission])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return &updated, nil
}

func (r *RolePermissionRepositoryImpl) ReplaceByRole(ctx context.Context, roleID int64, items []entities.RolePermission) ([]entities.RolePermission, error) {
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	deleteQuery, deleteArgs, err := r.sb.Delete(tableName).
		Where(sq.Eq{"role_id": roleID}).
		ToSql()
	if err != nil {
		return nil, err
	}
	if _, err := tx.Exec(ctx, deleteQuery, deleteArgs...); err != nil {
		return nil, err
	}

	if len(items) == 0 {
		if err := tx.Commit(ctx); err != nil {
			return nil, err
		}
		return []entities.RolePermission{}, nil
	}

	builder := r.sb.Insert(tableName).
		Columns("role_id", "permission_id", "effect").
		Suffix("RETURNING " + columnList())

	for _, entity := range items {
		builder = builder.Values(roleID, entity.PermissionId, entity.Effect)
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

	updated, err := pgx.CollectRows(rows, pgx.RowToStructByName[entities.RolePermission])
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return updated, nil
}

func (r *RolePermissionRepositoryImpl) Delete(ctx context.Context, id int64) error {
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

func (r *RolePermissionRepositoryImpl) detailBuilder(params shared.ListParams) sq.SelectBuilder {
	builder := r.sb.Select(
		"rp.id AS id",
		"rp.role_id AS role_id",
		"rp.permission_id AS permission_id",
		"rp.effect AS effect",
		"rp.created_at AS created_at",
		"p.code AS permission_code",
		"p.name AS permission_name",
		"m.id AS module_id",
		"m.code AS module_code",
		"m.name AS module_name",
		"a.id AS app_id",
		"a.code AS app_code",
		"a.name AS app_name",
	).
		From(tableName + " rp").
		Join("permissions p ON p.id = rp.permission_id").
		LeftJoin("modules m ON m.id = p.module_id").
		Join("apps a ON a.id = COALESCE(m.app_id, p.app_id)").
		OrderBy("rp.id DESC").
		Limit(params.Limit).
		Offset(params.Offset)

	if params.Search != "" {
		search := "%" + params.Search + "%"
		builder = builder.Where(sq.Or{
			sq.Expr("LOWER(p.code) LIKE ?", search),
			sq.Expr("LOWER(p.name) LIKE ?", search),
			sq.Expr("LOWER(m.code) LIKE ?", search),
			sq.Expr("LOWER(m.name) LIKE ?", search),
			sq.Expr("LOWER(a.code) LIKE ?", search),
			sq.Expr("LOWER(a.name) LIKE ?", search),
		})
	}

	return builder
}

func (r *RolePermissionRepositoryImpl) roleSummaryBuilder(params shared.ListParams) sq.SelectBuilder {
	builder := r.sb.Select(
		"ro.id AS role_id",
		"ro.code AS role_code",
		"ro.name AS role_name",
		"ro.description AS role_description",
		"ro.scope AS role_scope",
		"a.id AS app_id",
		"a.code AS app_code",
		"a.name AS app_name",
		"COUNT(rp.permission_id) AS permission_count",
	).
		From(tableName+" rp").
		Join("roles ro ON ro.id = rp.role_id").
		Join("permissions p ON p.id = rp.permission_id").
		LeftJoin("modules m ON m.id = p.module_id").
		Join("apps a ON a.id = COALESCE(m.app_id, p.app_id)").
		GroupBy("ro.id", "ro.code", "ro.name", "ro.description", "ro.scope", "a.id", "a.code", "a.name").
		OrderBy("ro.id DESC", "a.id DESC").
		Limit(params.Limit).
		Offset(params.Offset)

	if params.Search != "" {
		search := "%" + params.Search + "%"
		builder = builder.Where(sq.Or{
			sq.Expr("LOWER(ro.code) LIKE ?", search),
			sq.Expr("LOWER(ro.name) LIKE ?", search),
			sq.Expr("LOWER(a.code) LIKE ?", search),
			sq.Expr("LOWER(a.name) LIKE ?", search),
		})
	}

	return builder
}

func (r *RolePermissionRepositoryImpl) availablePermissionBuilder(params shared.ListParams) sq.SelectBuilder {
	builder := r.sb.Select(
		"p.id AS id",
		"p.app_id AS app_id",
		"p.module_id AS module_id",
		"p.action_id AS action_id",
		"p.code AS code",
		"p.name AS name",
		"p.description AS description",
		"p.risk_level AS risk_level",
		"p.is_system AS is_system",
		"p.status AS status",
		"p.created_at AS created_at",
		"p.updated_at AS updated_at",
		"m.name AS module_name",
		"m.code AS module_code",
	).
		From("permissions p").
		LeftJoin("modules m ON m.id = p.module_id").
		Join("apps a ON a.id = COALESCE(m.app_id, p.app_id)").
		OrderBy("m.name ASC NULLS LAST", "p.name ASC").
		Limit(params.Limit).
		Offset(params.Offset)

	if params.Search != "" {
		search := "%" + params.Search + "%"
		builder = builder.Where(sq.Or{
			sq.Expr("LOWER(p.code) LIKE ?", search),
			sq.Expr("LOWER(p.name) LIKE ?", search),
			sq.Expr("LOWER(m.code) LIKE ?", search),
			sq.Expr("LOWER(m.name) LIKE ?", search),
		})
	}

	return builder
}

func (r *RolePermissionRepositoryImpl) findDetails(ctx context.Context, builder sq.SelectBuilder) ([]entities.RolePermissionDetail, error) {
	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[entities.RolePermissionDetail])
}

func (r *RolePermissionRepositoryImpl) findRoleSummaries(ctx context.Context, builder sq.SelectBuilder) ([]entities.RolePermissionSummary, error) {
	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[entities.RolePermissionSummary])
}

func (r *RolePermissionRepositoryImpl) findAvailablePermissions(ctx context.Context, builder sq.SelectBuilder) ([]entities.AvailablePermissionRow, error) {
	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[entities.AvailablePermissionRow])
}

func columns() []string {
	return []string{
		"id",
		"role_id",
		"permission_id",
		"effect",
		"created_at",
	}
}

func columnList() string {
	return "id, role_id, permission_id, effect, created_at"
}

func canUpdateTimestamp() bool {
	for _, column := range columns() {
		if column == "updated_at" {
			return true
		}
	}
	return false
}
