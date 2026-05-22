package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/open-suite/authorization/internal/entities"
	"github.com/open-suite/authorization/internal/modules/teammembers/dto"
	"github.com/open-suite/authorization/internal/modules/teammembers/repositories"
	"github.com/open-suite/authorization/internal/modules/teammembers/services"
	"github.com/open-suite/authorization/internal/platform/logger"
	"github.com/open-suite/authorization/internal/shared"
	"github.com/open-suite/authorization/internal/shared/response"
)

type TeamMemberControllerImpl struct {
	TeamMemberService services.TeamMemberService
	response          *response.Sender
	log               *logger.LayerLogger
}

func NewTeamMemberController(service services.TeamMemberService, sender *response.Sender, appLogger *logger.Logger) TeamMemberController {
	return &TeamMemberControllerImpl{
		TeamMemberService: service,
		response:          sender,
		log:               appLogger.Layer("controller.team_members"),
	}
}

func (c *TeamMemberControllerImpl) Find(w http.ResponseWriter, r *http.Request) {
	end := c.log.Start(r.Context(), "Find")

	params := shared.NewListParamsFromRequest(r)
	items, err := c.TeamMemberService.Find(r.Context(), params)
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusInternalServerError, "error.internal", nil)
		return
	}

	end(nil, "count", len(items))
	c.response.Success(w, r, http.StatusOK, "team_members.find.success", dto.ListResponse[entities.TeamMember]{Items: items})
}

func (c *TeamMemberControllerImpl) FindByID(w http.ResponseWriter, r *http.Request) {
	end := c.log.Start(r.Context(), "FindByID")

	id, err := parseID(r)
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "team_members.invalid_id", nil)
		return
	}

	item, err := c.TeamMemberService.FindByID(r.Context(), id)
	if errors.Is(err, repositories.ErrNotFound) {
		end(err)
		c.response.Error(w, r, http.StatusNotFound, "team_members.not_found", nil)
		return
	}
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusInternalServerError, "error.internal", nil)
		return
	}

	end(nil)
	c.response.Success(w, r, http.StatusOK, "team_members.find_by_id.success", item)
}

func (c *TeamMemberControllerImpl) Create(w http.ResponseWriter, r *http.Request) {
	end := c.log.Start(r.Context(), "Create")

	var request dto.CreateTeamMemberRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "team_members.invalid_payload", nil)
		return
	}

	item, err := c.TeamMemberService.Create(r.Context(), request)
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "team_members.create.failed", nil)
		return
	}

	end(nil)
	c.response.Success(w, r, http.StatusCreated, "team_members.create.success", item)
}

func (c *TeamMemberControllerImpl) Update(w http.ResponseWriter, r *http.Request) {
	end := c.log.Start(r.Context(), "Update")

	id, err := parseID(r)
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "team_members.invalid_id", nil)
		return
	}

	var request dto.UpdateTeamMemberRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "team_members.invalid_payload", nil)
		return
	}

	item, err := c.TeamMemberService.Update(r.Context(), id, request)
	if errors.Is(err, repositories.ErrNotFound) {
		end(err)
		c.response.Error(w, r, http.StatusNotFound, "team_members.not_found", nil)
		return
	}
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "team_members.update.failed", nil)
		return
	}

	end(nil)
	c.response.Success(w, r, http.StatusOK, "team_members.update.success", item)
}

func (c *TeamMemberControllerImpl) Delete(w http.ResponseWriter, r *http.Request) {
	end := c.log.Start(r.Context(), "Delete")

	id, err := parseID(r)
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "team_members.invalid_id", nil)
		return
	}

	if err := c.TeamMemberService.Delete(r.Context(), id); err != nil {
		if errors.Is(err, repositories.ErrNotFound) {
			end(err)
			c.response.Error(w, r, http.StatusNotFound, "team_members.not_found", nil)
			return
		}
		end(err)
		c.response.Error(w, r, http.StatusInternalServerError, "error.internal", nil)
		return
	}

	end(nil)
	c.response.Success(w, r, http.StatusOK, "team_members.delete.success", nil)
}

func parseID(r *http.Request) (int64, error) {
	return strconv.ParseInt(r.PathValue("id"), 10, 64)
}
