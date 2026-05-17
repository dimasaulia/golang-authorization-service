package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/open-suite/authorization/internal/entities"
	"github.com/open-suite/authorization/internal/modules/apppermissionmanifests/dto"
	"github.com/open-suite/authorization/internal/modules/apppermissionmanifests/repositories"
	"github.com/open-suite/authorization/internal/modules/apppermissionmanifests/services"
	"github.com/open-suite/authorization/internal/platform/logger"
	"github.com/open-suite/authorization/internal/shared/response"
)

type AppPermissionManifestControllerImpl struct {
	AppPermissionManifestService services.AppPermissionManifestService
	response                     *response.Sender
	log                          *logger.LayerLogger
}

func NewAppPermissionManifestController(service services.AppPermissionManifestService, sender *response.Sender, appLogger *logger.Logger) AppPermissionManifestController {
	return &AppPermissionManifestControllerImpl{
		AppPermissionManifestService: service,
		response:                     sender,
		log:                          appLogger.Layer("controller.app_permission_manifests"),
	}
}

func (c *AppPermissionManifestControllerImpl) Find(w http.ResponseWriter, r *http.Request) {
	end := c.log.Start(r.Context(), "Find")

	limit := parseUintQuery(r, "limit", 20)
	offset := parseUintQuery(r, "offset", 0)
	items, err := c.AppPermissionManifestService.Find(r.Context(), limit, offset)
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusInternalServerError, "error.internal", nil)
		return
	}

	end(nil, "count", len(items))
	c.response.Success(w, r, http.StatusOK, "app_permission_manifests.find.success", dto.ListResponse[entities.AppPermissionManifest]{Items: items})
}

func (c *AppPermissionManifestControllerImpl) FindByID(w http.ResponseWriter, r *http.Request) {
	end := c.log.Start(r.Context(), "FindByID")

	id, err := parseID(r)
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "app_permission_manifests.invalid_id", nil)
		return
	}

	item, err := c.AppPermissionManifestService.FindByID(r.Context(), id)
	if errors.Is(err, repositories.ErrNotFound) {
		end(err)
		c.response.Error(w, r, http.StatusNotFound, "app_permission_manifests.not_found", nil)
		return
	}
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusInternalServerError, "error.internal", nil)
		return
	}

	end(nil)
	c.response.Success(w, r, http.StatusOK, "app_permission_manifests.find_by_id.success", item)
}

func (c *AppPermissionManifestControllerImpl) Create(w http.ResponseWriter, r *http.Request) {
	end := c.log.Start(r.Context(), "Create")

	var request dto.CreateAppPermissionManifestRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "app_permission_manifests.invalid_payload", nil)
		return
	}

	item, err := c.AppPermissionManifestService.Create(r.Context(), request)
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "app_permission_manifests.create.failed", nil)
		return
	}

	end(nil)
	c.response.Success(w, r, http.StatusCreated, "app_permission_manifests.create.success", item)
}

func (c *AppPermissionManifestControllerImpl) Update(w http.ResponseWriter, r *http.Request) {
	end := c.log.Start(r.Context(), "Update")

	id, err := parseID(r)
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "app_permission_manifests.invalid_id", nil)
		return
	}

	var request dto.UpdateAppPermissionManifestRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "app_permission_manifests.invalid_payload", nil)
		return
	}

	item, err := c.AppPermissionManifestService.Update(r.Context(), id, request)
	if errors.Is(err, repositories.ErrNotFound) {
		end(err)
		c.response.Error(w, r, http.StatusNotFound, "app_permission_manifests.not_found", nil)
		return
	}
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "app_permission_manifests.update.failed", nil)
		return
	}

	end(nil)
	c.response.Success(w, r, http.StatusOK, "app_permission_manifests.update.success", item)
}

func (c *AppPermissionManifestControllerImpl) Delete(w http.ResponseWriter, r *http.Request) {
	end := c.log.Start(r.Context(), "Delete")

	id, err := parseID(r)
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "app_permission_manifests.invalid_id", nil)
		return
	}

	if err := c.AppPermissionManifestService.Delete(r.Context(), id); err != nil {
		if errors.Is(err, repositories.ErrNotFound) {
			end(err)
			c.response.Error(w, r, http.StatusNotFound, "app_permission_manifests.not_found", nil)
			return
		}
		end(err)
		c.response.Error(w, r, http.StatusInternalServerError, "error.internal", nil)
		return
	}

	end(nil)
	c.response.Success(w, r, http.StatusOK, "app_permission_manifests.delete.success", nil)
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
