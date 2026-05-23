package services

import (
	"context"
	"strconv"
	"strings"

	"github.com/open-suite/authorization/internal/entities"
	"github.com/open-suite/authorization/internal/modules/rolepermissions/dto"
	"github.com/open-suite/authorization/internal/modules/rolepermissions/repositories"
	"github.com/open-suite/authorization/internal/platform/logger"
	"github.com/open-suite/authorization/internal/shared"
)

type RolePermissionServiceImpl struct {
	RolePermissionRepository repositories.RolePermissionRepository
	log                      *logger.LayerLogger
}

func NewRolePermissionService(repository repositories.RolePermissionRepository, appLogger *logger.Logger) RolePermissionService {
	return &RolePermissionServiceImpl{
		RolePermissionRepository: repository,
		log:                      appLogger.Layer("service.role_permissions"),
	}
}

func (s *RolePermissionServiceImpl) Find(ctx context.Context, params shared.ListParams) ([]entities.RolePermission, error) {
	end := s.log.Start(ctx, "Find")
	items, err := s.RolePermissionRepository.Find(ctx, params)
	end(err, "count", len(items))
	return items, err
}

func (s *RolePermissionServiceImpl) FindByApp(ctx context.Context, appUnique string, params shared.ListParams) ([]entities.RolePermissionDetail, error) {
	appUnique = strings.TrimSpace(appUnique)
	end := s.log.Start(ctx, "FindByApp", "app_unique", appUnique)

	id, err := strconv.ParseInt(appUnique, 10, 64)
	if err == nil {
		items, err := s.RolePermissionRepository.FindByAppID(ctx, id, params)
		end(err, "count", len(items))
		return items, err
	}

	items, err := s.RolePermissionRepository.FindByAppCode(ctx, appUnique, params)
	end(err, "count", len(items))
	return items, err
}

func (s *RolePermissionServiceImpl) FindByRole(ctx context.Context, roleUnique string, params shared.ListParams) ([]entities.RolePermissionDetail, error) {
	roleUnique = strings.TrimSpace(roleUnique)
	end := s.log.Start(ctx, "FindByRole", "role_unique", roleUnique)

	roleID, err := s.resolveRoleID(ctx, roleUnique)
	if err != nil {
		end(err)
		return nil, err
	}

	items, err := s.RolePermissionRepository.FindByRole(ctx, roleID, params)
	end(err, "count", len(items))
	return items, err
}

func (s *RolePermissionServiceImpl) FindRoleSummaries(ctx context.Context, params shared.ListParams) ([]entities.RolePermissionSummary, error) {
	end := s.log.Start(ctx, "FindRoleSummaries")
	items, err := s.RolePermissionRepository.FindRoleSummaries(ctx, params)
	end(err, "count", len(items))
	return items, err
}

func (s *RolePermissionServiceImpl) FindRoleSummariesByApp(ctx context.Context, appUnique string, params shared.ListParams) ([]entities.RolePermissionSummary, error) {
	appUnique = strings.TrimSpace(appUnique)
	end := s.log.Start(ctx, "FindRoleSummariesByApp", "app_unique", appUnique)

	id, err := strconv.ParseInt(appUnique, 10, 64)
	if err == nil {
		items, err := s.RolePermissionRepository.FindRoleSummariesByAppID(ctx, id, params)
		end(err, "count", len(items))
		return items, err
	}

	items, err := s.RolePermissionRepository.FindRoleSummariesByAppCode(ctx, appUnique, params)
	end(err, "count", len(items))
	return items, err
}

func (s *RolePermissionServiceImpl) FindAvailablePermissionsByApp(ctx context.Context, appUnique string, params shared.ListParams) ([]entities.AvailablePermissionModule, error) {
	appUnique = strings.TrimSpace(appUnique)
	end := s.log.Start(ctx, "FindAvailablePermissionsByApp", "app_unique", appUnique)

	var (
		rows []entities.AvailablePermissionRow
		err  error
	)

	id, parseErr := strconv.ParseInt(appUnique, 10, 64)
	if parseErr == nil {
		rows, err = s.RolePermissionRepository.FindAvailablePermissionsByAppID(ctx, id, params)
	} else {
		rows, err = s.RolePermissionRepository.FindAvailablePermissionsByAppCode(ctx, appUnique, params)
	}
	if err != nil {
		end(err)
		return nil, err
	}

	items := groupAvailablePermissions(rows)
	end(nil, "count", len(items))
	return items, nil
}

func (s *RolePermissionServiceImpl) FindByID(ctx context.Context, id int64) (*entities.RolePermission, error) {
	end := s.log.Start(ctx, "FindByID", "id", id)
	item, err := s.RolePermissionRepository.FindByID(ctx, id)
	end(err)
	return item, err
}

func (s *RolePermissionServiceImpl) Create(ctx context.Context, request dto.CreateRolePermissionRequest) (*entities.RolePermission, error) {
	end := s.log.Start(ctx, "Create")
	item, err := s.RolePermissionRepository.Create(ctx, entities.RolePermission{
		RoleId:       request.RoleId,
		PermissionId: request.PermissionId,
		Effect:       defaultEffect(request.Effect),
	})
	end(err)
	return item, err
}

func (s *RolePermissionServiceImpl) CreateBulk(ctx context.Context, request []dto.CreateBulkRolePermissionRequest) ([]entities.RolePermission, error) {
	end := s.log.Start(ctx, "CreateBulk", "count", len(request))

	items := make([]entities.RolePermission, 0)
	seen := map[[2]int64]struct{}{}
	for _, rolePermissions := range request {
		effect := defaultEffect(rolePermissions.Effect)
		for _, permissionID := range rolePermissions.PermissionIds {
			key := [2]int64{rolePermissions.RoleId, permissionID}
			if _, exists := seen[key]; exists {
				continue
			}
			seen[key] = struct{}{}
			items = append(items, entities.RolePermission{
				RoleId:       rolePermissions.RoleId,
				PermissionId: permissionID,
				Effect:       effect,
			})
		}
	}

	created, err := s.RolePermissionRepository.CreateBulk(ctx, items)
	end(err, "count", len(created))
	return created, err
}

func (s *RolePermissionServiceImpl) Update(ctx context.Context, id int64, request dto.UpdateRolePermissionRequest) (*entities.RolePermission, error) {
	end := s.log.Start(ctx, "Update", "id", id)

	data := map[string]any{}
	if request.RoleId != nil {
		data["role_id"] = *request.RoleId
	}
	if request.PermissionId != nil {
		data["permission_id"] = *request.PermissionId
	}
	if request.Effect != nil {
		data["effect"] = *request.Effect
	}

	item, err := s.RolePermissionRepository.Update(ctx, id, data)
	end(err)
	return item, err
}

func (s *RolePermissionServiceImpl) UpdateByRole(ctx context.Context, roleUnique string, request dto.UpdateRolePermissionByRoleRequest) ([]entities.RolePermission, error) {
	roleUnique = strings.TrimSpace(roleUnique)
	end := s.log.Start(ctx, "UpdateByRole", "role_unique", roleUnique)

	roleID, err := s.resolveRoleID(ctx, roleUnique)
	if err != nil {
		end(err)
		return nil, err
	}

	effect := defaultEffect(request.Effect)
	items := make([]entities.RolePermission, 0, len(request.PermissionIds))
	seen := map[int64]struct{}{}
	for _, permissionID := range request.PermissionIds {
		if _, exists := seen[permissionID]; exists {
			continue
		}
		seen[permissionID] = struct{}{}
		items = append(items, entities.RolePermission{
			RoleId:       roleID,
			PermissionId: permissionID,
			Effect:       effect,
		})
	}

	updated, err := s.RolePermissionRepository.ReplaceByRole(ctx, roleID, items)
	end(err, "count", len(updated))
	return updated, err
}

func (s *RolePermissionServiceImpl) Delete(ctx context.Context, id int64) error {
	end := s.log.Start(ctx, "Delete", "id", id)
	err := s.RolePermissionRepository.Delete(ctx, id)
	end(err)
	return err
}

func defaultEffect(effect string) string {
	effect = strings.TrimSpace(effect)
	if effect == "" {
		return "allow"
	}
	return effect
}

func (s *RolePermissionServiceImpl) resolveRoleID(ctx context.Context, roleUnique string) (int64, error) {
	id, err := strconv.ParseInt(roleUnique, 10, 64)
	if err == nil {
		return id, nil
	}

	return s.RolePermissionRepository.FindRoleIDByCode(ctx, roleUnique)
}

func groupAvailablePermissions(rows []entities.AvailablePermissionRow) []entities.AvailablePermissionModule {
	items := make([]entities.AvailablePermissionModule, 0)
	indexes := map[int64]int{}
	withoutModuleIndex := -1

	for _, row := range rows {
		permission := entities.Permission{
			ID:          row.ID,
			AppId:       row.AppId,
			ModuleId:    row.ModuleId,
			ActionId:    row.ActionId,
			Code:        row.Code,
			Name:        row.Name,
			Description: row.Description,
			RiskLevel:   row.RiskLevel,
			IsSystem:    row.IsSystem,
			Status:      row.Status,
			CreatedAt:   row.CreatedAt,
			UpdatedAt:   row.UpdatedAt,
		}

		if row.ModuleId == nil {
			if withoutModuleIndex == -1 {
				items = append(items, entities.AvailablePermissionModule{
					ModuleName:  row.ModuleName,
					ModuleCode:  row.ModuleCode,
					Permissions: []entities.Permission{},
				})
				withoutModuleIndex = len(items) - 1
			}
			items[withoutModuleIndex].Permissions = append(items[withoutModuleIndex].Permissions, permission)
			continue
		}

		index, exists := indexes[*row.ModuleId]
		if !exists {
			items = append(items, entities.AvailablePermissionModule{
				ModuleId:    row.ModuleId,
				ModuleName:  row.ModuleName,
				ModuleCode:  row.ModuleCode,
				Permissions: []entities.Permission{},
			})
			index = len(items) - 1
			indexes[*row.ModuleId] = index
		}

		items[index].Permissions = append(items[index].Permissions, permission)
	}

	return items
}
