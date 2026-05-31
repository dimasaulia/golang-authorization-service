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

const tableName = "team_members"

var ErrNotFound = errors.New("team member not found")

type TeamMemberRepositoryImpl struct {
	db  *database.Database
	log *logger.LayerLogger
	sb  sq.StatementBuilderType
}

func NewTeamMemberRepository(db *database.Database, appLogger *logger.Logger) TeamMemberRepository {
	return &TeamMemberRepositoryImpl{
		db:  db,
		log: appLogger.Layer("repository.team_members"),
		sb:  sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *TeamMemberRepositoryImpl) Find(ctx context.Context, params shared.ListParams) ([]entities.TeamMember, error) {
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

	return pgx.CollectRows(rows, pgx.RowToStructByName[entities.TeamMember])
}

func (r *TeamMemberRepositoryImpl) FindByID(ctx context.Context, id int64) (*entities.TeamMember, error) {
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

	item, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entities.TeamMember])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (r *TeamMemberRepositoryImpl) FindTeamIDsByUserID(ctx context.Context, userID int64) ([]int64, error) {
	query, args, err := r.sb.Select("team_id").
		From(tableName).
		Where(sq.Eq{"user_id": userID}).
		OrderBy("team_id ASC").
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
		var teamID int64
		if err := rows.Scan(&teamID); err != nil {
			return nil, err
		}
		items = append(items, teamID)
	}
	return items, rows.Err()
}

func (r *TeamMemberRepositoryImpl) FindAssignedTeamsByUserID(ctx context.Context, userID int64) ([]entities.UserAssignedTeam, error) {
	query, args, err := r.sb.Select("t.id", "t.code", "t.name").
		From(tableName+" tm").
		Join("teams t ON t.id = tm.team_id").
		Where(sq.Eq{"tm.user_id": userID}).
		OrderBy("t.name ASC", "t.id ASC").
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[entities.UserAssignedTeam])
}

func (r *TeamMemberRepositoryImpl) Create(ctx context.Context, entity entities.TeamMember) (*entities.TeamMember, error) {
	values := map[string]any{
		"team_id":      entity.TeamId,
		"user_id":      entity.UserId,
		"role_in_team": entity.RoleInTeam,
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

	created, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entities.TeamMember])
	if err != nil {
		return nil, err
	}

	return &created, nil
}

func (r *TeamMemberRepositoryImpl) ReplaceTeamsForUser(ctx context.Context, userID int64, teamIDs []int64) ([]entities.TeamMember, error) {
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

	items := make([]entities.TeamMember, 0, len(teamIDs))
	seen := map[int64]struct{}{}
	for _, teamID := range teamIDs {
		if teamID == 0 {
			continue
		}
		if _, exists := seen[teamID]; exists {
			continue
		}
		seen[teamID] = struct{}{}

		query, args, err := r.sb.Insert(tableName).
			SetMap(map[string]any{
				"team_id": teamID,
				"user_id": userID,
			}).
			Suffix("RETURNING " + columnList()).
			ToSql()
		if err != nil {
			return nil, err
		}

		rows, err := tx.Query(ctx, query, args...)
		if err != nil {
			return nil, err
		}
		created, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entities.TeamMember])
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

func (r *TeamMemberRepositoryImpl) Update(ctx context.Context, id int64, data map[string]any) (*entities.TeamMember, error) {
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

	updated, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entities.TeamMember])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return &updated, nil
}

func (r *TeamMemberRepositoryImpl) Delete(ctx context.Context, id int64) error {
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
		"team_id",
		"user_id",
		"role_in_team",
		"created_at",
	}
}

func columnList() string {
	return "id, team_id, user_id, role_in_team, created_at"
}

func canUpdateTimestamp() bool {
	for _, column := range columns() {
		if column == "updated_at" {
			return true
		}
	}
	return false
}
