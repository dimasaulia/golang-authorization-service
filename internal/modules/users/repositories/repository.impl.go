package repositories

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/open-suite/authorization/internal/entities"
	"github.com/open-suite/authorization/internal/platform/database"
	"github.com/open-suite/authorization/internal/platform/logger"
	"github.com/open-suite/authorization/internal/shared"
)

const tableName = "users"
const userCredentialsTableName = "user_credentials"
const userIdentitiesTableName = "user_identities"
const userVerificationCodesTableName = "user_verification_codes"

var ErrNotFound = errors.New("user not found")

type UserRepositoryImpl struct {
	db  *database.Database
	log *logger.LayerLogger
	sb  sq.StatementBuilderType
}

func NewUserRepository(db *database.Database, appLogger *logger.Logger) UserRepository {
	return &UserRepositoryImpl{
		db:  db,
		log: appLogger.Layer("repository.users"),
		sb:  sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *UserRepositoryImpl) Find(ctx context.Context, params shared.ListParams) ([]entities.User, error) {
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

	return pgx.CollectRows(rows, pgx.RowToStructByName[entities.User])
}

func (r *UserRepositoryImpl) FindByID(ctx context.Context, id int64) (*entities.User, error) {
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

	item, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entities.User])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (r *UserRepositoryImpl) FindByEmail(ctx context.Context, email string) (*entities.User, error) {
	query, args, err := r.sb.Select(columns()...).
		From(tableName).
		Where(sq.Eq{"email": email}).
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

	item, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entities.User])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (r *UserRepositoryImpl) FindByIdentity(ctx context.Context, provider string, providerUserID string) (*entities.User, error) {
	userColumns := make([]string, 0, len(columns()))
	for _, column := range columns() {
		userColumns = append(userColumns, "u."+column)
	}

	query, args, err := r.sb.Select(userColumns...).
		From(tableName + " u").
		Join(userIdentitiesTableName + " ui ON ui.user_id = u.id").
		Where(sq.Eq{"ui.provider": provider, "ui.provider_user_id": providerUserID}).
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

	item, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entities.User])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (r *UserRepositoryImpl) FindIdentitiesByUserID(ctx context.Context, userID int64) ([]entities.UserIdentity, error) {
	query, args, err := r.sb.Select(identityColumns()...).
		From(userIdentitiesTableName).
		Where(sq.Eq{"user_id": userID}).
		OrderBy("is_primary DESC", "provider ASC").
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

func (r *UserRepositoryImpl) Create(ctx context.Context, input CreateUserInput) (*entities.User, error) {
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	values := map[string]any{
		"organization_id": input.User.OrganizationId,
		"username":        input.User.Username,
		"email":           input.User.Email,
		"display_name":    input.User.DisplayName,
		"type":            input.User.Type,
		"status":          input.User.Status,
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
	defer rows.Close()

	created, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entities.User])
	if err != nil {
		return nil, err
	}

	if input.PasswordHash != "" {
		if err := r.createCredential(ctx, tx, created.ID, input.PasswordHash, input.MustChangePassword); err != nil {
			return nil, err
		}
	}

	if input.Identity != nil {
		identity := *input.Identity
		identity.UserId = created.ID
		if err := r.createIdentity(ctx, tx, identity); err != nil {
			return nil, err
		}
	}

	if input.VerificationCode != nil {
		if err := r.createVerificationCode(ctx, tx, created.ID, *input.VerificationCode); err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &created, nil
}

func (r *UserRepositoryImpl) LinkIdentity(ctx context.Context, identity entities.UserIdentity) error {
	return r.createIdentity(ctx, r.db.Pool, identity)
}

func (r *UserRepositoryImpl) UpdateCredential(ctx context.Context, userID int64, passwordHash string, mustChangePassword bool) error {
	if passwordHash == "" {
		return nil
	}

	query, args, err := r.sb.Insert(userCredentialsTableName).
		SetMap(map[string]any{
			"user_id":              userID,
			"password_hash":        passwordHash,
			"must_change_password": mustChangePassword,
			"password_changed_at":  time.Now(),
		}).
		Suffix("ON CONFLICT (user_id) DO UPDATE SET password_hash = EXCLUDED.password_hash, must_change_password = EXCLUDED.must_change_password, password_changed_at = EXCLUDED.password_changed_at, updated_at = NOW()").
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.Pool.Exec(ctx, query, args...)
	return err
}

func (r *UserRepositoryImpl) UpdateIdentityProfile(ctx context.Context, userID int64, provider string, providerUserID string, username string, email string) error {
	data := map[string]any{}
	if strings.TrimSpace(providerUserID) != "" {
		data["provider_user_id"] = providerUserID
	}
	if strings.TrimSpace(username) != "" {
		data["username"] = username
	}
	if strings.TrimSpace(email) != "" {
		data["email"] = email
	}
	if len(data) == 0 {
		return nil
	}

	query, args, err := r.sb.Update(userIdentitiesTableName).
		SetMap(data).
		Where(sq.Eq{"user_id": userID, "provider": provider}).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.Pool.Exec(ctx, query, args...)
	return err
}

func (r *UserRepositoryImpl) Update(ctx context.Context, id int64, data map[string]any) (*entities.User, error) {
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

	updated, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entities.User])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return &updated, nil
}

func (r *UserRepositoryImpl) Delete(ctx context.Context, id int64) error {
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

func (r *UserRepositoryImpl) CreateVerificationCode(ctx context.Context, userID int64, input CreateVerificationCodeInput) error {
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query, args, err := r.sb.Update(userVerificationCodesTableName).
		Set("used_at", time.Now()).
		Where(sq.Eq{"user_id": userID, "purpose": input.Purpose}).
		Where(sq.Eq{"used_at": nil}).
		ToSql()
	if err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, query, args...); err != nil {
		return err
	}

	if err := r.createVerificationCode(ctx, tx, userID, input); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *UserRepositoryImpl) FindVerificationCode(ctx context.Context, purpose string, codeHash string) (*entities.UserVerificationCode, error) {
	println("FindVerificationCode SQL => START")

	query, args, err := r.sb.Select(
		"id",
		"user_id",
		"purpose",
		"code_hash",
		"expires_at",
		"used_at",
		"created_at",
	).
		From(userVerificationCodesTableName).
		Where(sq.Eq{"purpose": purpose, "code_hash": codeHash}).
		Where(sq.Eq{"used_at": nil}).
		// Where(sq.Gt{"expires_at": time.Now()}).
		OrderBy("id DESC").
		Limit(1).
		ToSql()
	if err != nil {
		println("FindVerificationCode SQL => ERROR")
		return nil, err
	}
	println("FindVerificationCode SQL =>", query)
	println("FindVerificationCode ARG purpose =>", purpose)
	println("FindVerificationCode ARG code_hash_prefix =>", (codeHash))
	println("FindVerificationCode ARG count =>", len(args))
	for index, arg := range args {
		println("FindVerificationCode ARG", index+1, "=>", fmt.Sprint(arg))
	}

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		println("FindVerificationCode Query error =>", err.Error())
		return nil, err
	}
	defer rows.Close()

	item, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entities.UserVerificationCode])
	if errors.Is(err, pgx.ErrNoRows) {
		r.debugVerificationCodeMiss(ctx, purpose, codeHash)
		println("FindVerificationCode CollectOneRow error =>", ErrNotFound.Error())
		return nil, ErrNotFound
	}
	if err != nil {
		println("FindVerificationCode CollectOneRow error =>", err.Error())
		return nil, err
	}
	println("FindVerificationCode Found => id:", item.ID, "user_id:", item.UserID, "purpose:", item.Purpose)

	return &item, nil
}

func (r *UserRepositoryImpl) debugVerificationCodeMiss(ctx context.Context, purpose string, codeHash string) {
	println("FindVerificationCode MISS => requested_purpose:", purpose)
	println("FindVerificationCode MISS => requested_code_hash:", codeHash)
	println("FindVerificationCode MISS => now:", time.Now().Format(time.RFC3339Nano))

	query, args, err := r.sb.Select(
		"id",
		"user_id",
		"purpose",
		"expires_at",
		"used_at",
		"created_at",
	).
		From(userVerificationCodesTableName).
		Where(sq.Eq{"code_hash": codeHash}).
		OrderBy("id DESC").
		Limit(5).
		ToSql()
	if err != nil {
		println("FindVerificationCode MISS debug query build error =>", err.Error())
		return
	}

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		println("FindVerificationCode MISS debug query error =>", err.Error())
		return
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		count++
		var (
			id        int64
			userID    int64
			rowPurpose string
			expiresAt time.Time
			usedAt    *time.Time
			createdAt *time.Time
		)
		if err := rows.Scan(&id, &userID, &rowPurpose, &expiresAt, &usedAt, &createdAt); err != nil {
			println("FindVerificationCode MISS debug scan error =>", err.Error())
			return
		}

		usedAtText := "<nil>"
		if usedAt != nil {
			usedAtText = usedAt.Format(time.RFC3339Nano)
		}
		createdAtText := "<nil>"
		if createdAt != nil {
			createdAtText = createdAt.Format(time.RFC3339Nano)
		}
		println("FindVerificationCode MISS row => id:", id, "user_id:", userID, "purpose:", rowPurpose, "expires_at:", expiresAt.Format(time.RFC3339Nano), "used_at:", usedAtText, "created_at:", createdAtText)
	}
	if err := rows.Err(); err != nil {
		println("FindVerificationCode MISS debug rows error =>", err.Error())
		return
	}
	if count == 0 {
		println("FindVerificationCode MISS => no rows found with requested code_hash")
	}
}

func (r *UserRepositoryImpl) UseVerificationCode(ctx context.Context, codeID int64) error {
	query, args, err := r.sb.Update(userVerificationCodesTableName).
		Set("used_at", time.Now()).
		Where(sq.Eq{"id": codeID}).
		Where(sq.Eq{"used_at": nil}).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return err
	}

	var updatedID int64
	if err := r.db.Pool.QueryRow(ctx, query, args...).Scan(&updatedID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}
		return err
	}

	return nil
}

func (r *UserRepositoryImpl) MarkEmailVerified(ctx context.Context, userID int64) error {
	query, args, err := r.sb.Update(tableName).
		Set("email_verified_at", time.Now()).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": userID}).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return err
	}

	var updatedID int64
	if err := r.db.Pool.QueryRow(ctx, query, args...).Scan(&updatedID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}
		return err
	}

	return nil
}

type txExecutor interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
}

func (r *UserRepositoryImpl) createCredential(ctx context.Context, tx txExecutor, userID int64, passwordHash string, mustChangePassword bool) error {
	values := map[string]any{
		"user_id":              userID,
		"password_hash":        passwordHash,
		"must_change_password": mustChangePassword,
		"password_changed_at":  time.Now(),
	}

	query, args, err := r.sb.Insert(userCredentialsTableName).SetMap(values).ToSql()
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, query, args...)
	return err
}

func (r *UserRepositoryImpl) createIdentity(ctx context.Context, tx txExecutor, identity entities.UserIdentity) error {
	values := map[string]any{
		"user_id":          identity.UserId,
		"provider":         identity.Provider,
		"provider_user_id": identity.ProviderUserId,
		"username":         identity.Username,
		"email":            identity.Email,
		"is_primary":       identity.IsPrimary,
	}

	query, args, err := r.sb.Insert(userIdentitiesTableName).SetMap(values).ToSql()
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, query, args...)
	return err
}

func (r *UserRepositoryImpl) createVerificationCode(ctx context.Context, tx txExecutor, userID int64, input CreateVerificationCodeInput) error {
	values := map[string]any{
		"user_id":    userID,
		"purpose":    input.Purpose,
		"code_hash":  input.CodeHash,
		"expires_at": input.ExpiresAt,
	}

	query, args, err := r.sb.Insert(userVerificationCodesTableName).SetMap(values).ToSql()
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, query, args...)
	return err
}

func columns() []string {
	return []string{
		"id",
		"organization_id",
		"username",
		"email",
		"display_name",
		"type",
		"status",
		"email_verified_at",
		"created_at",
		"updated_at",
	}
}

func columnList() string {
	return "id, organization_id, username, email, display_name, type, status, email_verified_at, created_at, updated_at"
}

func identityColumns() []string {
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

func canUpdateTimestamp() bool {
	for _, column := range columns() {
		if column == "updated_at" {
			return true
		}
	}
	return false
}

func maskHashPrefix(value string) string {
	if len(value) <= 12 {
		return value
	}
	return value[:12] + "..."
}
