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

const tableName = "user_identities"

var ErrNotFound = errors.New("user identity not found")

type UserIdentityRepositoryImpl struct {
	db  *database.Database
	log *logger.LayerLogger
	sb  sq.StatementBuilderType
}

func NewUserIdentityRepository(db *database.Database, appLogger *logger.Logger) UserIdentityRepository {
	return &UserIdentityRepositoryImpl{
		db:  db,
		log: appLogger.Layer("repository.user_identities"),
		sb:  sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *UserIdentityRepositoryImpl) Find(ctx context.Context, limit uint64, offset uint64) ([]entities.UserIdentity, error) {
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

	return pgx.CollectRows(rows, pgx.RowToStructByName[entities.UserIdentity])
}

func (r *UserIdentityRepositoryImpl) FindByID(ctx context.Context, id int64) (*entities.UserIdentity, error) {
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

	item, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entities.UserIdentity])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (r *UserIdentityRepositoryImpl) Create(ctx context.Context, entity entities.UserIdentity) (*entities.UserIdentity, error) {
	values := map[string]any{
		"user_id":          entity.UserId,
		"provider":         entity.Provider,
		"provider_user_id": entity.ProviderUserId,
		"username":         entity.Username,
		"email":            entity.Email,
		"is_primary":       entity.IsPrimary,
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

	created, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entities.UserIdentity])
	if err != nil {
		return nil, err
	}

	return &created, nil
}

func (r *UserIdentityRepositoryImpl) Update(ctx context.Context, id int64, data map[string]any) (*entities.UserIdentity, error) {
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

	updated, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entities.UserIdentity])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return &updated, nil
}

func (r *UserIdentityRepositoryImpl) Delete(ctx context.Context, id int64) error {
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
		"provider",
		"provider_user_id",
		"username",
		"email",
		"is_primary",
		"created_at",
	}
}

func columnList() string {
	return "id, user_id, provider, provider_user_id, username, email, is_primary, created_at"
}

func canUpdateTimestamp() bool {
	for _, column := range columns() {
		if column == "updated_at" {
			return true
		}
	}
	return false
}
