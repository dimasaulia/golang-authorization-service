package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/open-suite/authorization/internal/entities"
	"github.com/open-suite/authorization/internal/modules/menus/dto"
	"github.com/open-suite/authorization/internal/modules/menus/repositories"
	"github.com/open-suite/authorization/internal/modules/menus/services"
	"github.com/open-suite/authorization/internal/platform/logger"
	"github.com/open-suite/authorization/internal/shared"
	"github.com/open-suite/authorization/internal/shared/response"
)

type MenuControllerImpl struct {
	MenuService services.MenuService
	response    *response.Sender
	log         *logger.LayerLogger
}

func NewMenuController(service services.MenuService, sender *response.Sender, appLogger *logger.Logger) MenuController {
	return &MenuControllerImpl{
		MenuService: service,
		response:    sender,
		log:         appLogger.Layer("controller.menus"),
	}
}

func (c *MenuControllerImpl) Find(w http.ResponseWriter, r *http.Request) {
	end := c.log.Start(r.Context(), "Find")

	params := shared.NewListParamsFromRequest(r)
	items, err := c.MenuService.Find(r.Context(), params)
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusInternalServerError, "error.internal", nil)
		return
	}

	end(nil, "count", len(items))
	c.response.Success(w, r, http.StatusOK, "menus.find.success", dto.ListResponse[entities.Menu]{Items: items})
}

func (c *MenuControllerImpl) FindByID(w http.ResponseWriter, r *http.Request) {
	end := c.log.Start(r.Context(), "FindByID")

	id, err := parseID(r)
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "menus.invalid_id", nil)
		return
	}

	item, err := c.MenuService.FindByID(r.Context(), id)
	if errors.Is(err, repositories.ErrNotFound) {
		end(err)
		c.response.Error(w, r, http.StatusNotFound, "menus.not_found", nil)
		return
	}
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusInternalServerError, "error.internal", nil)
		return
	}

	end(nil)
	c.response.Success(w, r, http.StatusOK, "menus.find_by_id.success", item)
}

func (c *MenuControllerImpl) Create(w http.ResponseWriter, r *http.Request) {
	end := c.log.Start(r.Context(), "Create")

	var request dto.CreateMenuRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "menus.invalid_payload", nil)
		return
	}

	item, err := c.MenuService.Create(r.Context(), request)
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "menus.create.failed", nil)
		return
	}

	end(nil)
	c.response.Success(w, r, http.StatusCreated, "menus.create.success", item)
}

func (c *MenuControllerImpl) Update(w http.ResponseWriter, r *http.Request) {
	end := c.log.Start(r.Context(), "Update")

	id, err := parseID(r)
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "menus.invalid_id", nil)
		return
	}

	var request dto.UpdateMenuRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "menus.invalid_payload", nil)
		return
	}

	item, err := c.MenuService.Update(r.Context(), id, request)
	if errors.Is(err, repositories.ErrNotFound) {
		end(err)
		c.response.Error(w, r, http.StatusNotFound, "menus.not_found", nil)
		return
	}
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "menus.update.failed", nil)
		return
	}

	end(nil)
	c.response.Success(w, r, http.StatusOK, "menus.update.success", item)
}

func (c *MenuControllerImpl) Delete(w http.ResponseWriter, r *http.Request) {
	end := c.log.Start(r.Context(), "Delete")

	id, err := parseID(r)
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "menus.invalid_id", nil)
		return
	}

	if err := c.MenuService.Delete(r.Context(), id); err != nil {
		if errors.Is(err, repositories.ErrNotFound) {
			end(err)
			c.response.Error(w, r, http.StatusNotFound, "menus.not_found", nil)
			return
		}
		end(err)
		c.response.Error(w, r, http.StatusInternalServerError, "error.internal", nil)
		return
	}

	end(nil)
	c.response.Success(w, r, http.StatusOK, "menus.delete.success", nil)
}

func parseID(r *http.Request) (int64, error) {
	return strconv.ParseInt(r.PathValue("id"), 10, 64)
}
