package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/open-suite/authorization/internal/entities"
	"github.com/open-suite/authorization/internal/modules/permissions/dto"
	"github.com/open-suite/authorization/internal/modules/permissions/repositories"
	"github.com/open-suite/authorization/internal/modules/permissions/services"
	"github.com/open-suite/authorization/internal/platform/logger"
	"github.com/open-suite/authorization/internal/shared/response"
)

type PermissionControllerImpl struct {
	PermissionService services.PermissionService
	response          *response.Sender
	log               *logger.LayerLogger
}

func NewPermissionController(service services.PermissionService, sender *response.Sender, appLogger *logger.Logger) PermissionController {
	return &PermissionControllerImpl{
		PermissionService: service,
		response:          sender,
		log:               appLogger.Layer("controller.permissions"),
	}
}

func (c *PermissionControllerImpl) Find(w http.ResponseWriter, r *http.Request) {
	end := c.log.Start(r.Context(), "Find")

	limit := parseUintQuery(r, "limit", 20)
	offset := parseUintQuery(r, "offset", 0)
	items, err := c.PermissionService.Find(r.Context(), limit, offset)
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusInternalServerError, "error.internal", nil)
		return
	}

	end(nil, "count", len(items))
	c.response.Success(w, r, http.StatusOK, "permissions.find.success", dto.ListResponse[entities.Permission]{Items: items})
}

func (c *PermissionControllerImpl) FindByID(w http.ResponseWriter, r *http.Request) {
	end := c.log.Start(r.Context(), "FindByID")

	id, err := parseID(r)
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "permissions.invalid_id", nil)
		return
	}

	item, err := c.PermissionService.FindByID(r.Context(), id)
	if errors.Is(err, repositories.ErrNotFound) {
		end(err)
		c.response.Error(w, r, http.StatusNotFound, "permissions.not_found", nil)
		return
	}
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusInternalServerError, "error.internal", nil)
		return
	}

	end(nil)
	c.response.Success(w, r, http.StatusOK, "permissions.find_by_id.success", item)
}

func (c *PermissionControllerImpl) Create(w http.ResponseWriter, r *http.Request) {
	end := c.log.Start(r.Context(), "Create")

	var request dto.CreatePermissionRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "permissions.invalid_payload", nil)
		return
	}

	item, err := c.PermissionService.Create(r.Context(), request)
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "permissions.create.failed", nil)
		return
	}

	end(nil)
	c.response.Success(w, r, http.StatusCreated, "permissions.create.success", item)
}

func (c *PermissionControllerImpl) Update(w http.ResponseWriter, r *http.Request) {
	end := c.log.Start(r.Context(), "Update")

	id, err := parseID(r)
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "permissions.invalid_id", nil)
		return
	}

	var request dto.UpdatePermissionRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "permissions.invalid_payload", nil)
		return
	}

	item, err := c.PermissionService.Update(r.Context(), id, request)
	if errors.Is(err, repositories.ErrNotFound) {
		end(err)
		c.response.Error(w, r, http.StatusNotFound, "permissions.not_found", nil)
		return
	}
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "permissions.update.failed", nil)
		return
	}

	end(nil)
	c.response.Success(w, r, http.StatusOK, "permissions.update.success", item)
}

func (c *PermissionControllerImpl) Delete(w http.ResponseWriter, r *http.Request) {
	end := c.log.Start(r.Context(), "Delete")

	id, err := parseID(r)
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "permissions.invalid_id", nil)
		return
	}

	if err := c.PermissionService.Delete(r.Context(), id); err != nil {
		if errors.Is(err, repositories.ErrNotFound) {
			end(err)
			c.response.Error(w, r, http.StatusNotFound, "permissions.not_found", nil)
			return
		}
		end(err)
		c.response.Error(w, r, http.StatusInternalServerError, "error.internal", nil)
		return
	}

	end(nil)
	c.response.Success(w, r, http.StatusOK, "permissions.delete.success", nil)
}

func parseID(r *http.Request) (int64, error) {
	return strconv.ParseInt(r.PathValue("id"), 10, 64)
}

func parseUintQuery(r *http.Request, key string, fallback uint64) uint64 {
	value := r.URL.Query().Get(key)
	if value == "" {
		return fallback
	}

	parsed, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return fallback
	}

	return parsed
}
