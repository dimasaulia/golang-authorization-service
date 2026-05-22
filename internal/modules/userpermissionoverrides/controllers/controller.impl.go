package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/open-suite/authorization/internal/entities"
	"github.com/open-suite/authorization/internal/modules/userpermissionoverrides/dto"
	"github.com/open-suite/authorization/internal/modules/userpermissionoverrides/repositories"
	"github.com/open-suite/authorization/internal/modules/userpermissionoverrides/services"
	"github.com/open-suite/authorization/internal/platform/logger"
	"github.com/open-suite/authorization/internal/shared"
	"github.com/open-suite/authorization/internal/shared/response"
)

type UserPermissionOverrideControllerImpl struct {
	UserPermissionOverrideService services.UserPermissionOverrideService
	response                      *response.Sender
	log                           *logger.LayerLogger
}

func NewUserPermissionOverrideController(service services.UserPermissionOverrideService, sender *response.Sender, appLogger *logger.Logger) UserPermissionOverrideController {
	return &UserPermissionOverrideControllerImpl{
		UserPermissionOverrideService: service,
		response:                      sender,
		log:                           appLogger.Layer("controller.user_permission_overrides"),
	}
}

func (c *UserPermissionOverrideControllerImpl) Find(w http.ResponseWriter, r *http.Request) {
	end := c.log.Start(r.Context(), "Find")

	params := shared.NewListParamsFromRequest(r)
	items, err := c.UserPermissionOverrideService.Find(r.Context(), params)
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusInternalServerError, "error.internal", nil)
		return
	}

	end(nil, "count", len(items))
	c.response.Success(w, r, http.StatusOK, "user_permission_overrides.find.success", dto.ListResponse[entities.UserPermissionOverride]{Items: items})
}

func (c *UserPermissionOverrideControllerImpl) FindByID(w http.ResponseWriter, r *http.Request) {
	end := c.log.Start(r.Context(), "FindByID")

	id, err := parseID(r)
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "user_permission_overrides.invalid_id", nil)
		return
	}

	item, err := c.UserPermissionOverrideService.FindByID(r.Context(), id)
	if errors.Is(err, repositories.ErrNotFound) {
		end(err)
		c.response.Error(w, r, http.StatusNotFound, "user_permission_overrides.not_found", nil)
		return
	}
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusInternalServerError, "error.internal", nil)
		return
	}

	end(nil)
	c.response.Success(w, r, http.StatusOK, "user_permission_overrides.find_by_id.success", item)
}

func (c *UserPermissionOverrideControllerImpl) Create(w http.ResponseWriter, r *http.Request) {
	end := c.log.Start(r.Context(), "Create")

	var request dto.CreateUserPermissionOverrideRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "user_permission_overrides.invalid_payload", nil)
		return
	}

	item, err := c.UserPermissionOverrideService.Create(r.Context(), request)
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "user_permission_overrides.create.failed", nil)
		return
	}

	end(nil)
	c.response.Success(w, r, http.StatusCreated, "user_permission_overrides.create.success", item)
}

func (c *UserPermissionOverrideControllerImpl) Update(w http.ResponseWriter, r *http.Request) {
	end := c.log.Start(r.Context(), "Update")

	id, err := parseID(r)
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "user_permission_overrides.invalid_id", nil)
		return
	}

	var request dto.UpdateUserPermissionOverrideRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "user_permission_overrides.invalid_payload", nil)
		return
	}

	item, err := c.UserPermissionOverrideService.Update(r.Context(), id, request)
	if errors.Is(err, repositories.ErrNotFound) {
		end(err)
		c.response.Error(w, r, http.StatusNotFound, "user_permission_overrides.not_found", nil)
		return
	}
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "user_permission_overrides.update.failed", nil)
		return
	}

	end(nil)
	c.response.Success(w, r, http.StatusOK, "user_permission_overrides.update.success", item)
}

func (c *UserPermissionOverrideControllerImpl) Delete(w http.ResponseWriter, r *http.Request) {
	end := c.log.Start(r.Context(), "Delete")

	id, err := parseID(r)
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "user_permission_overrides.invalid_id", nil)
		return
	}

	if err := c.UserPermissionOverrideService.Delete(r.Context(), id); err != nil {
		if errors.Is(err, repositories.ErrNotFound) {
			end(err)
			c.response.Error(w, r, http.StatusNotFound, "user_permission_overrides.not_found", nil)
			return
		}
		end(err)
		c.response.Error(w, r, http.StatusInternalServerError, "error.internal", nil)
		return
	}

	end(nil)
	c.response.Success(w, r, http.StatusOK, "user_permission_overrides.delete.success", nil)
}

func parseID(r *http.Request) (int64, error) {
	return strconv.ParseInt(r.PathValue("id"), 10, 64)
}
