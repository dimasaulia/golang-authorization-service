package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/open-suite/authorization/internal/entities"
	"github.com/open-suite/authorization/internal/modules/rolepermissions/dto"
	"github.com/open-suite/authorization/internal/modules/rolepermissions/repositories"
	"github.com/open-suite/authorization/internal/modules/rolepermissions/services"
	"github.com/open-suite/authorization/internal/platform/logger"
	"github.com/open-suite/authorization/internal/shared/response"
)

type RolePermissionControllerImpl struct {
	RolePermissionService services.RolePermissionService
	response              *response.Sender
	log                   *logger.LayerLogger
}

func NewRolePermissionController(service services.RolePermissionService, sender *response.Sender, appLogger *logger.Logger) RolePermissionController {
	return &RolePermissionControllerImpl{
		RolePermissionService: service,
		response:              sender,
		log:                   appLogger.Layer("controller.role_permissions"),
	}
}

func (c *RolePermissionControllerImpl) Find(w http.ResponseWriter, r *http.Request) {
	end := c.log.Start(r.Context(), "Find")

	limit := parseUintQuery(r, "limit", 20)
	offset := parseUintQuery(r, "offset", 0)
	items, err := c.RolePermissionService.Find(r.Context(), limit, offset)
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusInternalServerError, "error.internal", nil)
		return
	}

	end(nil, "count", len(items))
	c.response.Success(w, r, http.StatusOK, "role_permissions.find.success", dto.ListResponse[entities.RolePermission]{Items: items})
}

func (c *RolePermissionControllerImpl) FindByID(w http.ResponseWriter, r *http.Request) {
	end := c.log.Start(r.Context(), "FindByID")

	id, err := parseID(r)
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "role_permissions.invalid_id", nil)
		return
	}

	item, err := c.RolePermissionService.FindByID(r.Context(), id)
	if errors.Is(err, repositories.ErrNotFound) {
		end(err)
		c.response.Error(w, r, http.StatusNotFound, "role_permissions.not_found", nil)
		return
	}
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusInternalServerError, "error.internal", nil)
		return
	}

	end(nil)
	c.response.Success(w, r, http.StatusOK, "role_permissions.find_by_id.success", item)
}

func (c *RolePermissionControllerImpl) Create(w http.ResponseWriter, r *http.Request) {
	end := c.log.Start(r.Context(), "Create")

	var request dto.CreateRolePermissionRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "role_permissions.invalid_payload", nil)
		return
	}

	item, err := c.RolePermissionService.Create(r.Context(), request)
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "role_permissions.create.failed", nil)
		return
	}

	end(nil)
	c.response.Success(w, r, http.StatusCreated, "role_permissions.create.success", item)
}

func (c *RolePermissionControllerImpl) Update(w http.ResponseWriter, r *http.Request) {
	end := c.log.Start(r.Context(), "Update")

	id, err := parseID(r)
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "role_permissions.invalid_id", nil)
		return
	}

	var request dto.UpdateRolePermissionRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "role_permissions.invalid_payload", nil)
		return
	}

	item, err := c.RolePermissionService.Update(r.Context(), id, request)
	if errors.Is(err, repositories.ErrNotFound) {
		end(err)
		c.response.Error(w, r, http.StatusNotFound, "role_permissions.not_found", nil)
		return
	}
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "role_permissions.update.failed", nil)
		return
	}

	end(nil)
	c.response.Success(w, r, http.StatusOK, "role_permissions.update.success", item)
}

func (c *RolePermissionControllerImpl) Delete(w http.ResponseWriter, r *http.Request) {
	end := c.log.Start(r.Context(), "Delete")

	id, err := parseID(r)
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "role_permissions.invalid_id", nil)
		return
	}

	if err := c.RolePermissionService.Delete(r.Context(), id); err != nil {
		if errors.Is(err, repositories.ErrNotFound) {
			end(err)
			c.response.Error(w, r, http.StatusNotFound, "role_permissions.not_found", nil)
			return
		}
		end(err)
		c.response.Error(w, r, http.StatusInternalServerError, "error.internal", nil)
		return
	}

	end(nil)
	c.response.Success(w, r, http.StatusOK, "role_permissions.delete.success", nil)
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
