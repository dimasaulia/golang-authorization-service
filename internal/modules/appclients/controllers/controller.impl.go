package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/open-suite/authorization/internal/entities"
	"github.com/open-suite/authorization/internal/modules/appclients/dto"
	"github.com/open-suite/authorization/internal/modules/appclients/repositories"
	"github.com/open-suite/authorization/internal/modules/appclients/services"
	"github.com/open-suite/authorization/internal/platform/logger"
	"github.com/open-suite/authorization/internal/shared"
	"github.com/open-suite/authorization/internal/shared/response"
)

type AppClientControllerImpl struct {
	AppClientService services.AppClientService
	response         *response.Sender
	log              *logger.LayerLogger
}

func NewAppClientController(service services.AppClientService, sender *response.Sender, appLogger *logger.Logger) AppClientController {
	return &AppClientControllerImpl{
		AppClientService: service,
		response:         sender,
		log:              appLogger.Layer("controller.app_clients"),
	}
}

func (c *AppClientControllerImpl) Find(w http.ResponseWriter, r *http.Request) {
	end := c.log.Start(r.Context(), "Find")

	params := shared.NewListParamsFromRequest(r)
	items, err := c.AppClientService.Find(r.Context(), params)
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusInternalServerError, "error.internal", nil)
		return
	}

	end(nil, "count", len(items))
	c.response.Success(w, r, http.StatusOK, "app_clients.find.success", dto.ListResponse[entities.AppClient]{Items: items})
}

func (c *AppClientControllerImpl) FindByID(w http.ResponseWriter, r *http.Request) {
	end := c.log.Start(r.Context(), "FindByID")

	id, err := parseID(r)
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "app_clients.invalid_id", nil)
		return
	}

	item, err := c.AppClientService.FindByID(r.Context(), id)
	if errors.Is(err, repositories.ErrNotFound) {
		end(err)
		c.response.Error(w, r, http.StatusNotFound, "app_clients.not_found", nil)
		return
	}
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusInternalServerError, "error.internal", nil)
		return
	}

	end(nil)
	c.response.Success(w, r, http.StatusOK, "app_clients.find_by_id.success", item)
}

func (c *AppClientControllerImpl) Create(w http.ResponseWriter, r *http.Request) {
	end := c.log.Start(r.Context(), "Create")

	var request dto.CreateAppClientRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "app_clients.invalid_payload", nil)
		return
	}

	item, err := c.AppClientService.Create(r.Context(), request)
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "app_clients.create.failed", nil)
		return
	}

	end(nil)
	c.response.Success(w, r, http.StatusCreated, "app_clients.create.success", item)
}

func (c *AppClientControllerImpl) Update(w http.ResponseWriter, r *http.Request) {
	end := c.log.Start(r.Context(), "Update")

	id, err := parseID(r)
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "app_clients.invalid_id", nil)
		return
	}

	var request dto.UpdateAppClientRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "app_clients.invalid_payload", nil)
		return
	}

	item, err := c.AppClientService.Update(r.Context(), id, request)
	if errors.Is(err, repositories.ErrNotFound) {
		end(err)
		c.response.Error(w, r, http.StatusNotFound, "app_clients.not_found", nil)
		return
	}
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "app_clients.update.failed", nil)
		return
	}

	end(nil)
	c.response.Success(w, r, http.StatusOK, "app_clients.update.success", item)
}

func (c *AppClientControllerImpl) Delete(w http.ResponseWriter, r *http.Request) {
	end := c.log.Start(r.Context(), "Delete")

	id, err := parseID(r)
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "app_clients.invalid_id", nil)
		return
	}

	if err := c.AppClientService.Delete(r.Context(), id); err != nil {
		if errors.Is(err, repositories.ErrNotFound) {
			end(err)
			c.response.Error(w, r, http.StatusNotFound, "app_clients.not_found", nil)
			return
		}
		end(err)
		c.response.Error(w, r, http.StatusInternalServerError, "error.internal", nil)
		return
	}

	end(nil)
	c.response.Success(w, r, http.StatusOK, "app_clients.delete.success", nil)
}

func parseID(r *http.Request) (int64, error) {
	return strconv.ParseInt(r.PathValue("id"), 10, 64)
}
