package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/open-suite/authorization/internal/entities"
	"github.com/open-suite/authorization/internal/modules/accesscacheversions/dto"
	"github.com/open-suite/authorization/internal/modules/accesscacheversions/repositories"
	"github.com/open-suite/authorization/internal/modules/accesscacheversions/services"
	"github.com/open-suite/authorization/internal/platform/logger"
	"github.com/open-suite/authorization/internal/shared"
	"github.com/open-suite/authorization/internal/shared/response"
)

type AccessCacheVersionControllerImpl struct {
	AccessCacheVersionService services.AccessCacheVersionService
	response                  *response.Sender
	log                       *logger.LayerLogger
}

func NewAccessCacheVersionController(service services.AccessCacheVersionService, sender *response.Sender, appLogger *logger.Logger) AccessCacheVersionController {
	return &AccessCacheVersionControllerImpl{
		AccessCacheVersionService: service,
		response:                  sender,
		log:                       appLogger.Layer("controller.access_cache_versions"),
	}
}

func (c *AccessCacheVersionControllerImpl) Find(w http.ResponseWriter, r *http.Request) {
	end := c.log.Start(r.Context(), "Find")

	params := shared.NewListParamsFromRequest(r)
	items, err := c.AccessCacheVersionService.Find(r.Context(), params)
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusInternalServerError, "error.internal", nil)
		return
	}

	end(nil, "count", len(items))
	c.response.Success(w, r, http.StatusOK, "access_cache_versions.find.success", dto.ListResponse[entities.AccessCacheVersion]{Items: items})
}

func (c *AccessCacheVersionControllerImpl) FindByID(w http.ResponseWriter, r *http.Request) {
	end := c.log.Start(r.Context(), "FindByID")

	id, err := parseID(r)
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "access_cache_versions.invalid_id", nil)
		return
	}

	item, err := c.AccessCacheVersionService.FindByID(r.Context(), id)
	if errors.Is(err, repositories.ErrNotFound) {
		end(err)
		c.response.Error(w, r, http.StatusNotFound, "access_cache_versions.not_found", nil)
		return
	}
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusInternalServerError, "error.internal", nil)
		return
	}

	end(nil)
	c.response.Success(w, r, http.StatusOK, "access_cache_versions.find_by_id.success", item)
}

func (c *AccessCacheVersionControllerImpl) Create(w http.ResponseWriter, r *http.Request) {
	end := c.log.Start(r.Context(), "Create")

	var request dto.CreateAccessCacheVersionRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "access_cache_versions.invalid_payload", nil)
		return
	}

	item, err := c.AccessCacheVersionService.Create(r.Context(), request)
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "access_cache_versions.create.failed", nil)
		return
	}

	end(nil)
	c.response.Success(w, r, http.StatusCreated, "access_cache_versions.create.success", item)
}

func (c *AccessCacheVersionControllerImpl) Update(w http.ResponseWriter, r *http.Request) {
	end := c.log.Start(r.Context(), "Update")

	id, err := parseID(r)
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "access_cache_versions.invalid_id", nil)
		return
	}

	var request dto.UpdateAccessCacheVersionRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "access_cache_versions.invalid_payload", nil)
		return
	}

	item, err := c.AccessCacheVersionService.Update(r.Context(), id, request)
	if errors.Is(err, repositories.ErrNotFound) {
		end(err)
		c.response.Error(w, r, http.StatusNotFound, "access_cache_versions.not_found", nil)
		return
	}
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "access_cache_versions.update.failed", nil)
		return
	}

	end(nil)
	c.response.Success(w, r, http.StatusOK, "access_cache_versions.update.success", item)
}

func (c *AccessCacheVersionControllerImpl) Delete(w http.ResponseWriter, r *http.Request) {
	end := c.log.Start(r.Context(), "Delete")

	id, err := parseID(r)
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "access_cache_versions.invalid_id", nil)
		return
	}

	if err := c.AccessCacheVersionService.Delete(r.Context(), id); err != nil {
		if errors.Is(err, repositories.ErrNotFound) {
			end(err)
			c.response.Error(w, r, http.StatusNotFound, "access_cache_versions.not_found", nil)
			return
		}
		end(err)
		c.response.Error(w, r, http.StatusInternalServerError, "error.internal", nil)
		return
	}

	end(nil)
	c.response.Success(w, r, http.StatusOK, "access_cache_versions.delete.success", nil)
}

func parseID(r *http.Request) (int64, error) {
	return strconv.ParseInt(r.PathValue("id"), 10, 64)
}
