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

const tableName = "user_roles"

var ErrNotFound = errors.New("user role not found")

type UserRoleRepositoryImpl struct {
	db  *database.Database
	log *logger.LayerLogger
	sb  sq.StatementBuilderType
}

func NewUserRoleRepository(db *database.Database, appLogger *logger.Logger) UserRoleRepository {
	return &UserRoleRepositoryImpl{
		db:  db,
		log: appLogger.Layer("repository.user_roles"),
		sb:  sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *UserRoleRepositoryImpl) Find(ctx context.Context, params shared.ListParams) ([]entities.UserRole, error) {
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

	return pgx.CollectRows(rows, pgx.RowToStructByName[entities.UserRole])
}

func (r *UserRoleRepositoryImpl) FindByID(ctx context.Context, id int64) (*entities.UserRole, error) {
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

	item, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entities.UserRole])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (r *UserRoleRepositoryImpl) FindByUserAndRole(ctx context.Context, userID int64, roleID int64) (*entities.UserRole, error) {
	query, args, err := r.sb.Select(columns()...).
		From(tableName).
		Where(sq.Eq{"user_id": userID, "role_id": roleID}).
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

	item, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entities.UserRole])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (r *UserRoleRepositoryImpl) FindRoleIDsByUserID(ctx context.Context, userID int64) ([]int64, error) {
	query, args, err := r.sb.Select("role_id").
		From(tableName).
		Where(sq.Eq{"user_id": userID}).
		OrderBy("role_id ASC").
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []int64{}
	for rows.Next() {
		var roleID int64
		if err := rows.Scan(&roleID); err != nil {
			return nil, err
		}
		items = append(items, roleID)
	}
	return items, rows.Err()
}

func (r *UserRoleRepositoryImpl) FindAssignedRolesByUserID(ctx context.Context, userID int64) ([]entities.UserAssignedRole, error) {
	query, args, err := r.sb.Select(
		"r.id",
		"r.code",
		"r.name",
		"r.scope",
		"r.app_id",
		"a.code AS app_code",
		"a.name AS app_name",
	).
		From(tableName+" ur").
		Join("roles r ON r.id = ur.role_id").
		LeftJoin("apps a ON a.id = r.app_id").
		Where(sq.Eq{"ur.user_id": userID}).
		OrderBy("r.name ASC", "r.id ASC").
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[entities.UserAssignedRole])
}

func (r *UserRoleRepositoryImpl) Create(ctx context.Context, entity entities.UserRole) (*entities.UserRole, error) {
	values := map[string]any{
		"user_id":         entity.UserId,
		"role_id":         entity.RoleId,
		"app_id":          entity.AppId,
		"organization_id": entity.OrganizationId,
		"expires_at":      entity.ExpiresAt,
		"assigned_by":     entity.AssignedBy,
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

	created, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entities.UserRole])
	if err != nil {
		return nil, err
	}

	return &created, nil
}

func (r *UserRoleRepositoryImpl) ReplaceRolesForUser(ctx context.Context, userID int64, roleIDs []int64, organizationID *int64, assignedBy *int64) ([]entities.UserRole, error) {
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	deleteQuery, deleteArgs, err := r.sb.Delete(tableName).
		Where(sq.Eq{"user_id": userID}).
		ToSql()
	if err != nil {
		return nil, err
	}
	if _, err := tx.Exec(ctx, deleteQuery, deleteArgs...); err != nil {
		return nil, err
	}

	items := make([]entities.UserRole, 0, len(roleIDs))
	seen := map[int64]struct{}{}
	for _, roleID := range roleIDs {
		if roleID == 0 {
			continue
		}
		if _, exists := seen[roleID]; exists {
			continue
		}
		seen[roleID] = struct{}{}

		values := map[string]any{
			"user_id":         userID,
			"role_id":         roleID,
			"organization_id": organizationID,
			"assigned_by":     assignedBy,
		}
		query, args, err := r.sb.Insert(tableName).
			SetMap(values).
			Suffix("RETURNING " + columnList()).
			ToSql()
		if err != nil {
			return nil, err
		}

		rows, err := tx.Query(ctx, query, args...)
		if err != nil {
			return nil, err
		}
		created, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entities.UserRole])
		rows.Close()
		if err != nil {
			return nil, err
		}
		items = append(items, created)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *UserRoleRepositoryImpl) Update(ctx context.Context, id int64, data map[string]any) (*entities.UserRole, error) {
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

	updated, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entities.UserRole])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return &updated, nil
}

func (r *UserRoleRepositoryImpl) Delete(ctx context.Context, id int64) error {
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
		"role_id",
		"app_id",
		"organization_id",
		"expires_at",
		"assigned_by",
		"created_at",
	}
}

func columnList() string {
	return "id, user_id, role_id, app_id, organization_id, expires_at, assigned_by, created_at"
}

func canUpdateTimestamp() bool {
	for _, column := range columns() {
		if column == "updated_at" {
			return true
		}
	}
	return false
}
