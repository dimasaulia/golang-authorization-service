package services

import (
	"context"

	"github.com/open-suite/authorization/internal/entities"
	"github.com/open-suite/authorization/internal/modules/userroles/dto"
	"github.com/open-suite/authorization/internal/modules/userroles/repositories"
	"github.com/open-suite/authorization/internal/platform/logger"
	"github.com/open-suite/authorization/internal/shared"
)

type UserRoleServiceImpl struct {
	UserRoleRepository repositories.UserRoleRepository
	log                *logger.LayerLogger
}

func NewUserRoleService(repository repositories.UserRoleRepository, appLogger *logger.Logger) UserRoleService {
	return &UserRoleServiceImpl{
		UserRoleRepository: repository,
		log:                appLogger.Layer("service.user_roles"),
	}
}

func (s *UserRoleServiceImpl) Find(ctx context.Context, params shared.ListParams) ([]entities.UserRole, error) {
	end := s.log.Start(ctx, "Find")
	items, err := s.UserRoleRepository.Find(ctx, params)
	end(err, "count", len(items))
	return items, err
}

func (s *UserRoleServiceImpl) FindByID(ctx context.Context, id int64) (*entities.UserRole, error) {
	end := s.log.Start(ctx, "FindByID", "id", id)
	item, err := s.UserRoleRepository.FindByID(ctx, id)
	end(err)
	return item, err
}

func (s *UserRoleServiceImpl) Create(ctx context.Context, request dto.CreateUserRoleRequest) (*entities.UserRole, error) {
	end := s.log.Start(ctx, "Create")
	item, err := s.UserRoleRepository.Create(ctx, entities.UserRole{
		UserId:         request.UserId,
		RoleId:         request.RoleId,
		AppId:          request.AppId,
		OrganizationId: request.OrganizationId,
		ExpiresAt:      request.ExpiresAt,
		AssignedBy:     request.AssignedBy,
	})
	end(err)
	return item, err
}

func (s *UserRoleServiceImpl) AssignRolesToUser(ctx context.Context, userID int64, roleIDs []int64, organizationID *int64, assignedBy *int64) ([]entities.UserRole, error) {
	end := s.log.Start(ctx, "AssignRolesToUser", "user_id", userID, "count", len(roleIDs))
	items := make([]entities.UserRole, 0, len(roleIDs))
	for _, roleID := range roleIDs {
		if roleID == 0 {
			continue
		}
		item, err := s.UserRoleRepository.Create(ctx, entities.UserRole{
			UserId:         userID,
			RoleId:         roleID,
			OrganizationId: organizationID,
			AssignedBy:     assignedBy,
		})
		if err != nil {
			end(err)
			return nil, err
		}
		items = append(items, *item)
	}
	end(nil, "assigned", len(items))
	return items, nil
}

func (s *UserRoleServiceImpl) FindRoleIDsByUserID(ctx context.Context, userID int64) ([]int64, error) {
	end := s.log.Start(ctx, "FindRoleIDsByUserID", "user_id", userID)
	items, err := s.UserRoleRepository.FindRoleIDsByUserID(ctx, userID)
	end(err, "count", len(items))
	return items, err
}

func (s *UserRoleServiceImpl) FindAssignedRolesByUserID(ctx context.Context, userID int64) ([]entities.UserAssignedRole, error) {
	end := s.log.Start(ctx, "FindAssignedRolesByUserID", "user_id", userID)
	items, err := s.UserRoleRepository.FindAssignedRolesByUserID(ctx, userID)
	end(err, "count", len(items))
	return items, err
}

func (s *UserRoleServiceImpl) ReplaceRolesForUser(ctx context.Context, userID int64, roleIDs []int64, organizationID *int64, assignedBy *int64) ([]entities.UserRole, error) {
	end := s.log.Start(ctx, "ReplaceRolesForUser", "user_id", userID, "count", len(roleIDs))
	items, err := s.UserRoleRepository.ReplaceRolesForUser(ctx, userID, roleIDs, organizationID, assignedBy)
	end(err, "assigned", len(items))
	return items, err
}

func (s *UserRoleServiceImpl) Update(ctx context.Context, id int64, request dto.UpdateUserRoleRequest) (*entities.UserRole, error) {
	end := s.log.Start(ctx, "Update", "id", id)

	data := map[string]any{}
	if request.UserId != nil {
		data["user_id"] = *request.UserId
	}
	if request.RoleId != nil {
		data["role_id"] = *request.RoleId
	}
	if request.AppId != nil {
		data["app_id"] = *request.AppId
	}
	if request.OrganizationId != nil {
		data["organization_id"] = *request.OrganizationId
	}
	if request.ExpiresAt != nil {
		data["expires_at"] = *request.ExpiresAt
	}
	if request.AssignedBy != nil {
		data["assigned_by"] = *request.AssignedBy
	}

	item, err := s.UserRoleRepository.Update(ctx, id, data)
	end(err)
	return item, err
}

func (s *UserRoleServiceImpl) Delete(ctx context.Context, id int64) error {
	end := s.log.Start(ctx, "Delete", "id", id)
	err := s.UserRoleRepository.Delete(ctx, id)
	end(err)
	return err
}
