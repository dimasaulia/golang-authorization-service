package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/open-suite/authorization/internal/entities"
	"github.com/open-suite/authorization/internal/modules/actions/dto"
	"github.com/open-suite/authorization/internal/modules/actions/repositories"
	"github.com/open-suite/authorization/internal/modules/actions/services"
	"github.com/open-suite/authorization/internal/platform/logger"
	"github.com/open-suite/authorization/internal/shared"
	"github.com/open-suite/authorization/internal/shared/response"
)

type ActionControllerImpl struct {
	ActionService services.ActionService
	response      *response.Sender
	log           *logger.LayerLogger
}

func NewActionController(service services.ActionService, sender *response.Sender, appLogger *logger.Logger) ActionController {
	return &ActionControllerImpl{
		ActionService: service,
		response:      sender,
		log:           appLogger.Layer("controller.actions"),
	}
}

func (c *ActionControllerImpl) Find(w http.ResponseWriter, r *http.Request) {
	end := c.log.Start(r.Context(), "Find")

	params := shared.NewListParamsFromRequest(r)
	items, err := c.ActionService.Find(r.Context(), params)
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusInternalServerError, "error.internal", nil)
		return
	}

	end(nil, "count", len(items))
	c.response.Success(w, r, http.StatusOK, "actions.find.success", dto.ListResponse[entities.Action]{Items: items})
}

func (c *ActionControllerImpl) FindByUnique(w http.ResponseWriter, r *http.Request) {
	end := c.log.Start(r.Context(), "FindByUnique")

	identifier := strings.TrimSpace(r.PathValue("id"))
	if identifier == "" {
		err := errors.New("empty action identifier")
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "actions.empty_id", nil)
		return
	}

	item, err := c.ActionService.FindByUnique(r.Context(), identifier)
	if errors.Is(err, repositories.ErrNotFound) {
		end(err)
		c.response.Error(w, r, http.StatusNotFound, "actions.not_found", nil)
		return
	}
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusInternalServerError, "error.internal", nil)
		return
	}

	end(nil)
	c.response.Success(w, r, http.StatusOK, "actions.find_by_unique.success", item)
}

func (c *ActionControllerImpl) Create(w http.ResponseWriter, r *http.Request) {
	end := c.log.Start(r.Context(), "Create")

	var request dto.CreateActionRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "actions.invalid_payload", nil)
		return
	}

	item, err := c.ActionService.Create(r.Context(), request)
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "actions.create.failed", nil)
		return
	}

	end(nil)
	c.response.Success(w, r, http.StatusCreated, "actions.create.success", item)
}

func (c *ActionControllerImpl) Update(w http.ResponseWriter, r *http.Request) {
	end := c.log.Start(r.Context(), "Update")

	id, err := parseID(r)
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "actions.invalid_id", nil)
		return
	}

	var request dto.UpdateActionRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "actions.invalid_payload", nil)
		return
	}

	item, err := c.ActionService.Update(r.Context(), id, request)
	if errors.Is(err, repositories.ErrNotFound) {
		end(err)
		c.response.Error(w, r, http.StatusNotFound, "actions.not_found", nil)
		return
	}
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "actions.update.failed", nil)
		return
	}

	end(nil)
	c.response.Success(w, r, http.StatusOK, "actions.update.success", item)
}

func (c *ActionControllerImpl) Delete(w http.ResponseWriter, r *http.Request) {
	end := c.log.Start(r.Context(), "Delete")

	id, err := parseID(r)
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "actions.invalid_id", nil)
		return
	}

	if err := c.ActionService.Delete(r.Context(), id); err != nil {
		if errors.Is(err, repositories.ErrNotFound) {
			end(err)
			c.response.Error(w, r, http.StatusNotFound, "actions.not_found", nil)
			return
		}
		end(err)
		c.response.Error(w, r, http.StatusInternalServerError, "error.internal", nil)
		return
	}

	end(nil)
	c.response.Success(w, r, http.StatusOK, "actions.delete.success", nil)
}

func parseID(r *http.Request) (int64, error) {
	return strconv.ParseInt(r.PathValue("id"), 10, 64)
}
