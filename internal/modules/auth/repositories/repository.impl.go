package repositories

import (
	"context"
	"sort"

	"github.com/jackc/pgx/v5"

	"github.com/open-suite/authorization/internal/modules/auth/dto"
	"github.com/open-suite/authorization/internal/platform/database"
	"github.com/open-suite/authorization/internal/platform/logger"
)

type AccessRepositoryImpl struct {
	db  *database.Database
	log *logger.LayerLogger
}

type permissionRow struct {
	Code string `db:"code"`
}

type menuRow struct {
	Code               string  `db:"code"`
	Path               string  `db:"path"`
	RequiredPermission *string `db:"required_permission"`
}

type appRow struct {
	Code string `db:"code"`
	Name string `db:"name"`
}

func NewAccessRepository(db *database.Database, appLogger *logger.Logger) AccessRepository {
	return &AccessRepositoryImpl{
		db:  db,
		log: appLogger.Layer("repository.auth_access"),
	}
}

func (r *AccessRepositoryImpl) FindAccessSummary(ctx context.Context, userID int64, appCode string) (*dto.AccessSummaryResponse, error) {
	permissions, err := r.findPermissions(ctx, userID, appCode)
	if err != nil {
		return nil, err
	}
	menus, err := r.findMenus(ctx, userID, appCode)
	if err != nil {
		return nil, err
	}

	items := make([]string, 0, len(permissions))
	for _, permission := range permissions {
		items = append(items, permission.Code)
	}
	sort.Strings(items)

	menuItems := make([]dto.AccessibleMenu, 0, len(menus))
	for _, menu := range menus {
		menuItems = append(menuItems, dto.AccessibleMenu{
			Code:               menu.Code,
			Path:               menu.Path,
			RequiredPermission: menu.RequiredPermission,
		})
	}

	return &dto.AccessSummaryResponse{
		App:         appCode,
		Menus:       menuItems,
		Permissions: items,
	}, nil
}

func (r *AccessRepositoryImpl) FindApps(ctx context.Context, userID int64) ([]dto.UserAppAccess, error) {
	query := `
WITH role_sources AS (
    SELECT ur.role_id
    FROM user_roles ur
    JOIN roles r ON r.id = ur.role_id AND r.status = 'active'
    WHERE ur.user_id = $1
      AND (ur.expires_at IS NULL OR ur.expires_at > NOW())
    UNION
    SELECT tr.role_id
    FROM team_members tm
    JOIN team_roles tr ON tr.team_id = tm.team_id
    JOIN roles r ON r.id = tr.role_id AND r.status = 'active'
    WHERE tm.user_id = $1
),
role_permissions_effective AS (
    SELECT p.app_id, rp.effect
    FROM role_sources rs
    JOIN role_permissions rp ON rp.role_id = rs.role_id
    JOIN permissions p ON p.id = rp.permission_id AND p.status = 'active'
),
override_permissions AS (
    SELECT p.app_id, upo.effect
    FROM user_permission_overrides upo
    JOIN permissions p ON p.id = upo.permission_id AND p.status = 'active'
    WHERE upo.user_id = $1
      AND (upo.expires_at IS NULL OR upo.expires_at > NOW())
),
effective_app_ids AS (
    SELECT app_id FROM role_permissions_effective WHERE effect = 'allow'
    UNION
    SELECT app_id FROM override_permissions WHERE effect = 'allow'
    EXCEPT
    SELECT app_id FROM role_permissions_effective WHERE effect = 'deny'
    EXCEPT
    SELECT app_id FROM override_permissions WHERE effect = 'deny'
)
SELECT DISTINCT a.code, a.name
FROM apps a
JOIN effective_app_ids e ON e.app_id = a.id
WHERE a.status = 'active'
ORDER BY a.code ASC`

	rows, err := r.db.Pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	apps, err := pgx.CollectRows(rows, pgx.RowToStructByName[appRow])
	if err != nil {
		return nil, err
	}
	items := make([]dto.UserAppAccess, 0, len(apps))
	for _, app := range apps {
		items = append(items, dto.UserAppAccess{Code: app.Code, Name: app.Name})
	}
	return items, nil
}

func (r *AccessRepositoryImpl) findPermissions(ctx context.Context, userID int64, appCode string) ([]permissionRow, error) {
	query := `
WITH role_sources AS (
    SELECT ur.role_id
    FROM user_roles ur
    JOIN roles r ON r.id = ur.role_id AND r.status = 'active'
    WHERE ur.user_id = $1
      AND (ur.expires_at IS NULL OR ur.expires_at > NOW())
    UNION
    SELECT tr.role_id
    FROM team_members tm
    JOIN team_roles tr ON tr.team_id = tm.team_id
    JOIN roles r ON r.id = tr.role_id AND r.status = 'active'
    WHERE tm.user_id = $1
),
role_permissions_effective AS (
    SELECT p.code, rp.effect
    FROM role_sources rs
    JOIN role_permissions rp ON rp.role_id = rs.role_id
    JOIN permissions p ON p.id = rp.permission_id AND p.status = 'active'
    JOIN apps a ON a.id = p.app_id AND a.status = 'active'
    WHERE LOWER(a.code) = LOWER($2)
),
override_permissions AS (
    SELECT p.code, upo.effect
    FROM user_permission_overrides upo
    JOIN permissions p ON p.id = upo.permission_id AND p.status = 'active'
    JOIN apps a ON a.id = p.app_id AND a.status = 'active'
    WHERE upo.user_id = $1
      AND LOWER(a.code) = LOWER($2)
      AND (upo.expires_at IS NULL OR upo.expires_at > NOW())
),
effective_permissions AS (
    SELECT code FROM role_permissions_effective WHERE effect = 'allow'
    UNION
    SELECT code FROM override_permissions WHERE effect = 'allow'
    EXCEPT
    SELECT code FROM role_permissions_effective WHERE effect = 'deny'
    EXCEPT
    SELECT code FROM override_permissions WHERE effect = 'deny'
)
SELECT code
FROM effective_permissions
ORDER BY code ASC`

	rows, err := r.db.Pool.Query(ctx, query, userID, appCode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[permissionRow])
}

func (r *AccessRepositoryImpl) findMenus(ctx context.Context, userID int64, appCode string) ([]menuRow, error) {
	query := `
WITH role_sources AS (
    SELECT ur.role_id
    FROM user_roles ur
    JOIN roles r ON r.id = ur.role_id AND r.status = 'active'
    WHERE ur.user_id = $1
      AND (ur.expires_at IS NULL OR ur.expires_at > NOW())
    UNION
    SELECT tr.role_id
    FROM team_members tm
    JOIN team_roles tr ON tr.team_id = tm.team_id
    JOIN roles r ON r.id = tr.role_id AND r.status = 'active'
    WHERE tm.user_id = $1
),
role_permissions_effective AS (
    SELECT p.code, rp.effect
    FROM role_sources rs
    JOIN role_permissions rp ON rp.role_id = rs.role_id
    JOIN permissions p ON p.id = rp.permission_id AND p.status = 'active'
    JOIN apps a ON a.id = p.app_id AND a.status = 'active'
    WHERE LOWER(a.code) = LOWER($2)
),
override_permissions AS (
    SELECT p.code, upo.effect
    FROM user_permission_overrides upo
    JOIN permissions p ON p.id = upo.permission_id AND p.status = 'active'
    JOIN apps a ON a.id = p.app_id AND a.status = 'active'
    WHERE upo.user_id = $1
      AND LOWER(a.code) = LOWER($2)
      AND (upo.expires_at IS NULL OR upo.expires_at > NOW())
),
effective_permissions AS (
    SELECT code FROM role_permissions_effective WHERE effect = 'allow'
    UNION
    SELECT code FROM override_permissions WHERE effect = 'allow'
    EXCEPT
    SELECT code FROM role_permissions_effective WHERE effect = 'deny'
    EXCEPT
    SELECT code FROM override_permissions WHERE effect = 'deny'
)
SELECT m.code, m.route_path AS path, p.code AS required_permission
FROM menus m
JOIN apps a ON a.id = m.app_id AND a.status = 'active'
LEFT JOIN permissions p ON p.id = m.required_permission_id
WHERE LOWER(a.code) = LOWER($2)
  AND m.status = 'active'
  AND (m.required_permission_id IS NULL OR p.code IN (SELECT code FROM effective_permissions))
ORDER BY m.sort_order ASC, m.code ASC`

	rows, err := r.db.Pool.Query(ctx, query, userID, appCode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[menuRow])
}
