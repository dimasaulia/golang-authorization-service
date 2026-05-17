package services

import (
	"context"

	"github.com/open-suite/authorization/internal/entities"
	"github.com/open-suite/authorization/internal/modules/auditlogs/dto"
	"github.com/open-suite/authorization/internal/modules/auditlogs/repositories"
	"github.com/open-suite/authorization/internal/platform/logger"
)

type AuditLogServiceImpl struct {
	AuditLogRepository repositories.AuditLogRepository
	log                *logger.LayerLogger
}

func NewAuditLogService(repository repositories.AuditLogRepository, appLogger *logger.Logger) AuditLogService {
	return &AuditLogServiceImpl{
		AuditLogRepository: repository,
		log:                appLogger.Layer("service.audit_logs"),
	}
}

func (s *AuditLogServiceImpl) Find(ctx context.Context, limit uint64, offset uint64) ([]entities.AuditLog, error) {
	end := s.log.Start(ctx, "Find")
	items, err := s.AuditLogRepository.Find(ctx, limit, offset)
	end(err, "count", len(items))
	return items, err
}

func (s *AuditLogServiceImpl) FindByID(ctx context.Context, id int64) (*entities.AuditLog, error) {
	end := s.log.Start(ctx, "FindByID", "id", id)
	item, err := s.AuditLogRepository.FindByID(ctx, id)
	end(err)
	return item, err
}

func (s *AuditLogServiceImpl) Create(ctx context.Context, request dto.CreateAuditLogRequest) (*entities.AuditLog, error) {
	end := s.log.Start(ctx, "Create")
	item, err := s.AuditLogRepository.Create(ctx, entities.AuditLog{
		OrganizationId: request.OrganizationId,
		AppId:          request.AppId,
		ActorUserId:    request.ActorUserId,
		TargetUserId:   request.TargetUserId,
		PermissionId:   request.PermissionId,
		Action:         request.Action,
		ResourceType:   request.ResourceType,
		ResourceId:     request.ResourceId,
		Result:         request.Result,
		MetadataJson:   request.MetadataJson,
		IpAddress:      request.IpAddress,
		UserAgent:      request.UserAgent,
	})
	end(err)
	return item, err
}

func (s *AuditLogServiceImpl) Update(ctx context.Context, id int64, request dto.UpdateAuditLogRequest) (*entities.AuditLog, error) {
	end := s.log.Start(ctx, "Update", "id", id)

	data := map[string]any{}
	if request.OrganizationId != nil {
		data["organization_id"] = *request.OrganizationId
	}
	if request.AppId != nil {
		data["app_id"] = *request.AppId
	}
	if request.ActorUserId != nil {
		data["actor_user_id"] = *request.ActorUserId
	}
	if request.TargetUserId != nil {
		data["target_user_id"] = *request.TargetUserId
	}
	if request.PermissionId != nil {
		data["permission_id"] = *request.PermissionId
	}
	if request.Action != nil {
		data["action"] = *request.Action
	}
	if request.ResourceType != nil {
		data["resource_type"] = *request.ResourceType
	}
	if request.ResourceId != nil {
		data["resource_id"] = *request.ResourceId
	}
	if request.Result != nil {
		data["result"] = *request.Result
	}
	if request.MetadataJson != nil {
		data["metadata_json"] = *request.MetadataJson
	}
	if request.IpAddress != nil {
		data["ip_address"] = *request.IpAddress
	}
	if request.UserAgent != nil {
		data["user_agent"] = *request.UserAgent
	}

	item, err := s.AuditLogRepository.Update(ctx, id, data)
	end(err)
	return item, err
}

func (s *AuditLogServiceImpl) Delete(ctx context.Context, id int64) error {
	end := s.log.Start(ctx, "Delete", "id", id)
	err := s.AuditLogRepository.Delete(ctx, id)
	end(err)
	return err
}
