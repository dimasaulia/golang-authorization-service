package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/open-suite/authorization/internal/entities"
	"github.com/open-suite/authorization/internal/modules/modules/dto"
	"github.com/open-suite/authorization/internal/modules/modules/repositories"
	"github.com/open-suite/authorization/internal/modules/modules/services"
	"github.com/open-suite/authorization/internal/platform/logger"
	"github.com/open-suite/authorization/internal/shared"
	"github.com/open-suite/authorization/internal/shared/response"
)

type ModuleControllerImpl struct {
	ModuleService services.ModuleService
	response      *response.Sender
	log           *logger.LayerLogger
}

func NewModuleController(service services.ModuleService, sender *response.Sender, appLogger *logger.Logger) ModuleController {
	return &ModuleControllerImpl{
		ModuleService: service,
		response:      sender,
		log:           appLogger.Layer("controller.modules"),
	}
}

func (c *ModuleControllerImpl) Find(w http.ResponseWriter, r *http.Request) {
	end := c.log.Start(r.Context(), "Find")

	params := shared.NewListParamsFromRequest(r)
	items, err := c.ModuleService.Find(r.Context(), params)
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusInternalServerError, "error.internal", nil)
		return
	}

	end(nil, "count", len(items))
	c.response.Success(w, r, http.StatusOK, "modules.find.success", dto.ListResponse[entities.Module]{Items: items})
}

func (c *ModuleControllerImpl) FindByUnique(w http.ResponseWriter, r *http.Request) {
	end := c.log.Start(r.Context(), "FindByUnique")

	identifier := strings.TrimSpace(r.PathValue("id"))
	if identifier == "" {
		err := errors.New("empty module identifier")
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "modules.empty_id", nil)
		return
	}

	item, err := c.ModuleService.FindByUnique(r.Context(), identifier)
	if errors.Is(err, repositories.ErrNotFound) {
		end(err)
		c.response.Error(w, r, http.StatusNotFound, "modules.not_found", nil)
		return
	}
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusInternalServerError, "error.internal", nil)
		return
	}

	end(nil)
	c.response.Success(w, r, http.StatusOK, "modules.find_by_unique.success", item)
}

func (c *ModuleControllerImpl) FindByApp(w http.ResponseWriter, r *http.Request) {
	end := c.log.Start(r.Context(), "FindByApp")

	appIdentifier := strings.TrimSpace(r.PathValue("app"))
	if appIdentifier == "" {
		err := errors.New("empty app identifier")
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "apps.empty_id", nil)
		return
	}

	params := shared.NewListParamsFromRequest(r)
	items, err := c.ModuleService.FindByApp(r.Context(), appIdentifier, params)
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusInternalServerError, "error.internal", nil)
		return
	}

	end(nil, "count", len(items))
	c.response.Success(w, r, http.StatusOK, "modules.find_by_app.success", dto.ListResponse[entities.Module]{Items: items})
}

func (c *ModuleControllerImpl) Create(w http.ResponseWriter, r *http.Request) {
	end := c.log.Start(r.Context(), "Create")

	var request dto.CreateModuleRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "modules.invalid_payload", nil)
		return
	}

	item, err := c.ModuleService.Create(r.Context(), request)
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "modules.create.failed", nil)
		return
	}

	end(nil)
	c.response.Success(w, r, http.StatusCreated, "modules.create.success", item)
}

func (c *ModuleControllerImpl) Update(w http.ResponseWriter, r *http.Request) {
	end := c.log.Start(r.Context(), "Update")

	id, err := parseID(r)
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "modules.invalid_id", nil)
		return
	}

	var request dto.UpdateModuleRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "modules.invalid_payload", nil)
		return
	}

	item, err := c.ModuleService.Update(r.Context(), id, request)
	if errors.Is(err, repositories.ErrNotFound) {
		end(err)
		c.response.Error(w, r, http.StatusNotFound, "modules.not_found", nil)
		return
	}
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "modules.update.failed", nil)
		return
	}

	end(nil)
	c.response.Success(w, r, http.StatusOK, "modules.update.success", item)
}

func (c *ModuleControllerImpl) Delete(w http.ResponseWriter, r *http.Request) {
	end := c.log.Start(r.Context(), "Delete")

	id, err := parseID(r)
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "modules.invalid_id", nil)
		return
	}

	if err := c.ModuleService.Delete(r.Context(), id); err != nil {
		if errors.Is(err, repositories.ErrNotFound) {
			end(err)
			c.response.Error(w, r, http.StatusNotFound, "modules.not_found", nil)
			return
		}
		end(err)
		c.response.Error(w, r, http.StatusInternalServerError, "error.internal", nil)
		return
	}

	end(nil)
	c.response.Success(w, r, http.StatusOK, "modules.delete.success", nil)
}

func parseID(r *http.Request) (int64, error) {
	return strconv.ParseInt(r.PathValue("id"), 10, 64)
}
