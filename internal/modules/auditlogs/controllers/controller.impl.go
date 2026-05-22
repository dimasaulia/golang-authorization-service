package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/open-suite/authorization/internal/entities"
	"github.com/open-suite/authorization/internal/modules/auditlogs/dto"
	"github.com/open-suite/authorization/internal/modules/auditlogs/repositories"
	"github.com/open-suite/authorization/internal/modules/auditlogs/services"
	"github.com/open-suite/authorization/internal/platform/logger"
	"github.com/open-suite/authorization/internal/shared"
	"github.com/open-suite/authorization/internal/shared/response"
)

type AuditLogControllerImpl struct {
	AuditLogService services.AuditLogService
	response        *response.Sender
	log             *logger.LayerLogger
}

func NewAuditLogController(service services.AuditLogService, sender *response.Sender, appLogger *logger.Logger) AuditLogController {
	return &AuditLogControllerImpl{
		AuditLogService: service,
		response:        sender,
		log:             appLogger.Layer("controller.audit_logs"),
	}
}

func (c *AuditLogControllerImpl) Find(w http.ResponseWriter, r *http.Request) {
	end := c.log.Start(r.Context(), "Find")

	params := shared.NewListParamsFromRequest(r)
	items, err := c.AuditLogService.Find(r.Context(), params)
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusInternalServerError, "error.internal", nil)
		return
	}

	end(nil, "count", len(items))
	c.response.Success(w, r, http.StatusOK, "audit_logs.find.success", dto.ListResponse[entities.AuditLog]{Items: items})
}

func (c *AuditLogControllerImpl) FindByID(w http.ResponseWriter, r *http.Request) {
	end := c.log.Start(r.Context(), "FindByID")

	id, err := parseID(r)
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "audit_logs.invalid_id", nil)
		return
	}

	item, err := c.AuditLogService.FindByID(r.Context(), id)
	if errors.Is(err, repositories.ErrNotFound) {
		end(err)
		c.response.Error(w, r, http.StatusNotFound, "audit_logs.not_found", nil)
		return
	}
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusInternalServerError, "error.internal", nil)
		return
	}

	end(nil)
	c.response.Success(w, r, http.StatusOK, "audit_logs.find_by_id.success", item)
}

func (c *AuditLogControllerImpl) Create(w http.ResponseWriter, r *http.Request) {
	end := c.log.Start(r.Context(), "Create")

	var request dto.CreateAuditLogRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "audit_logs.invalid_payload", nil)
		return
	}

	item, err := c.AuditLogService.Create(r.Context(), request)
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "audit_logs.create.failed", nil)
		return
	}

	end(nil)
	c.response.Success(w, r, http.StatusCreated, "audit_logs.create.success", item)
}

func (c *AuditLogControllerImpl) Update(w http.ResponseWriter, r *http.Request) {
	end := c.log.Start(r.Context(), "Update")

	id, err := parseID(r)
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "audit_logs.invalid_id", nil)
		return
	}

	var request dto.UpdateAuditLogRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "audit_logs.invalid_payload", nil)
		return
	}

	item, err := c.AuditLogService.Update(r.Context(), id, request)
	if errors.Is(err, repositories.ErrNotFound) {
		end(err)
		c.response.Error(w, r, http.StatusNotFound, "audit_logs.not_found", nil)
		return
	}
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "audit_logs.update.failed", nil)
		return
	}

	end(nil)
	c.response.Success(w, r, http.StatusOK, "audit_logs.update.success", item)
}

func (c *AuditLogControllerImpl) Delete(w http.ResponseWriter, r *http.Request) {
	end := c.log.Start(r.Context(), "Delete")

	id, err := parseID(r)
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "audit_logs.invalid_id", nil)
		return
	}

	if err := c.AuditLogService.Delete(r.Context(), id); err != nil {
		if errors.Is(err, repositories.ErrNotFound) {
			end(err)
			c.response.Error(w, r, http.StatusNotFound, "audit_logs.not_found", nil)
			return
		}
		end(err)
		c.response.Error(w, r, http.StatusInternalServerError, "error.internal", nil)
		return
	}

	end(nil)
	c.response.Success(w, r, http.StatusOK, "audit_logs.delete.success", nil)
}

func parseID(r *http.Request) (int64, error) {
	return strconv.ParseInt(r.PathValue("id"), 10, 64)
}
