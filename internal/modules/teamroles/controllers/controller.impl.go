package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/open-suite/authorization/internal/entities"
	"github.com/open-suite/authorization/internal/modules/teamroles/dto"
	"github.com/open-suite/authorization/internal/modules/teamroles/repositories"
	"github.com/open-suite/authorization/internal/modules/teamroles/services"
	"github.com/open-suite/authorization/internal/platform/logger"
	"github.com/open-suite/authorization/internal/shared/response"
)

type TeamRoleControllerImpl struct {
	TeamRoleService services.TeamRoleService
	response        *response.Sender
	log             *logger.LayerLogger
}

func NewTeamRoleController(service services.TeamRoleService, sender *response.Sender, appLogger *logger.Logger) TeamRoleController {
	return &TeamRoleControllerImpl{
		TeamRoleService: service,
		response:        sender,
		log:             appLogger.Layer("controller.team_roles"),
	}
}

func (c *TeamRoleControllerImpl) Find(w http.ResponseWriter, r *http.Request) {
	end := c.log.Start(r.Context(), "Find")

	limit := parseUintQuery(r, "limit", 20)
	offset := parseUintQuery(r, "offset", 0)
	items, err := c.TeamRoleService.Find(r.Context(), limit, offset)
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusInternalServerError, "error.internal", nil)
		return
	}

	end(nil, "count", len(items))
	c.response.Success(w, r, http.StatusOK, "team_roles.find.success", dto.ListResponse[entities.TeamRole]{Items: items})
}

func (c *TeamRoleControllerImpl) FindByID(w http.ResponseWriter, r *http.Request) {
	end := c.log.Start(r.Context(), "FindByID")

	id, err := parseID(r)
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "team_roles.invalid_id", nil)
		return
	}

	item, err := c.TeamRoleService.FindByID(r.Context(), id)
	if errors.Is(err, repositories.ErrNotFound) {
		end(err)
		c.response.Error(w, r, http.StatusNotFound, "team_roles.not_found", nil)
		return
	}
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusInternalServerError, "error.internal", nil)
		return
	}

	end(nil)
	c.response.Success(w, r, http.StatusOK, "team_roles.find_by_id.success", item)
}

func (c *TeamRoleControllerImpl) Create(w http.ResponseWriter, r *http.Request) {
	end := c.log.Start(r.Context(), "Create")

	var request dto.CreateTeamRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "team_roles.invalid_payload", nil)
		return
	}

	item, err := c.TeamRoleService.Create(r.Context(), request)
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "team_roles.create.failed", nil)
		return
	}

	end(nil)
	c.response.Success(w, r, http.StatusCreated, "team_roles.create.success", item)
}

func (c *TeamRoleControllerImpl) Update(w http.ResponseWriter, r *http.Request) {
	end := c.log.Start(r.Context(), "Update")

	id, err := parseID(r)
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "team_roles.invalid_id", nil)
		return
	}

	var request dto.UpdateTeamRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "team_roles.invalid_payload", nil)
		return
	}

	item, err := c.TeamRoleService.Update(r.Context(), id, request)
	if errors.Is(err, repositories.ErrNotFound) {
		end(err)
		c.response.Error(w, r, http.StatusNotFound, "team_roles.not_found", nil)
		return
	}
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "team_roles.update.failed", nil)
		return
	}

	end(nil)
	c.response.Success(w, r, http.StatusOK, "team_roles.update.success", item)
}

func (c *TeamRoleControllerImpl) Delete(w http.ResponseWriter, r *http.Request) {
	end := c.log.Start(r.Context(), "Delete")

	id, err := parseID(r)
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "team_roles.invalid_id", nil)
		return
	}

	if err := c.TeamRoleService.Delete(r.Context(), id); err != nil {
		if errors.Is(err, repositories.ErrNotFound) {
			end(err)
			c.response.Error(w, r, http.StatusNotFound, "team_roles.not_found", nil)
			return
		}
		end(err)
		c.response.Error(w, r, http.StatusInternalServerError, "error.internal", nil)
		return
	}

	end(nil)
	c.response.Success(w, r, http.StatusOK, "team_roles.delete.success", nil)
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
