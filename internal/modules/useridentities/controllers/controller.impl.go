package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/open-suite/authorization/internal/entities"
	"github.com/open-suite/authorization/internal/modules/useridentities/dto"
	"github.com/open-suite/authorization/internal/modules/useridentities/repositories"
	"github.com/open-suite/authorization/internal/modules/useridentities/services"
	"github.com/open-suite/authorization/internal/platform/logger"
	"github.com/open-suite/authorization/internal/shared/response"
)

type UserIdentityControllerImpl struct {
	UserIdentityService services.UserIdentityService
	response            *response.Sender
	log                 *logger.LayerLogger
}

func NewUserIdentityController(service services.UserIdentityService, sender *response.Sender, appLogger *logger.Logger) UserIdentityController {
	return &UserIdentityControllerImpl{
		UserIdentityService: service,
		response:            sender,
		log:                 appLogger.Layer("controller.user_identities"),
	}
}

func (c *UserIdentityControllerImpl) Find(w http.ResponseWriter, r *http.Request) {
	end := c.log.Start(r.Context(), "Find")

	limit := parseUintQuery(r, "limit", 20)
	offset := parseUintQuery(r, "offset", 0)
	items, err := c.UserIdentityService.Find(r.Context(), limit, offset)
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusInternalServerError, "error.internal", nil)
		return
	}

	end(nil, "count", len(items))
	c.response.Success(w, r, http.StatusOK, "user_identities.find.success", dto.ListResponse[entities.UserIdentity]{Items: items})
}

func (c *UserIdentityControllerImpl) FindByID(w http.ResponseWriter, r *http.Request) {
	end := c.log.Start(r.Context(), "FindByID")

	id, err := parseID(r)
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "user_identities.invalid_id", nil)
		return
	}

	item, err := c.UserIdentityService.FindByID(r.Context(), id)
	if errors.Is(err, repositories.ErrNotFound) {
		end(err)
		c.response.Error(w, r, http.StatusNotFound, "user_identities.not_found", nil)
		return
	}
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusInternalServerError, "error.internal", nil)
		return
	}

	end(nil)
	c.response.Success(w, r, http.StatusOK, "user_identities.find_by_id.success", item)
}

func (c *UserIdentityControllerImpl) Create(w http.ResponseWriter, r *http.Request) {
	end := c.log.Start(r.Context(), "Create")

	var request dto.CreateUserIdentityRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "user_identities.invalid_payload", nil)
		return
	}

	item, err := c.UserIdentityService.Create(r.Context(), request)
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "user_identities.create.failed", nil)
		return
	}

	end(nil)
	c.response.Success(w, r, http.StatusCreated, "user_identities.create.success", item)
}

func (c *UserIdentityControllerImpl) Update(w http.ResponseWriter, r *http.Request) {
	end := c.log.Start(r.Context(), "Update")

	id, err := parseID(r)
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "user_identities.invalid_id", nil)
		return
	}

	var request dto.UpdateUserIdentityRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "user_identities.invalid_payload", nil)
		return
	}

	item, err := c.UserIdentityService.Update(r.Context(), id, request)
	if errors.Is(err, repositories.ErrNotFound) {
		end(err)
		c.response.Error(w, r, http.StatusNotFound, "user_identities.not_found", nil)
		return
	}
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "user_identities.update.failed", nil)
		return
	}

	end(nil)
	c.response.Success(w, r, http.StatusOK, "user_identities.update.success", item)
}

func (c *UserIdentityControllerImpl) Delete(w http.ResponseWriter, r *http.Request) {
	end := c.log.Start(r.Context(), "Delete")

	id, err := parseID(r)
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "user_identities.invalid_id", nil)
		return
	}

	if err := c.UserIdentityService.Delete(r.Context(), id); err != nil {
		if errors.Is(err, repositories.ErrNotFound) {
			end(err)
			c.response.Error(w, r, http.StatusNotFound, "user_identities.not_found", nil)
			return
		}
		end(err)
		c.response.Error(w, r, http.StatusInternalServerError, "error.internal", nil)
		return
	}

	end(nil)
	c.response.Success(w, r, http.StatusOK, "user_identities.delete.success", nil)
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
