package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/open-suite/authorization/internal/entities"
	"github.com/open-suite/authorization/internal/modules/organizations/dto"
	"github.com/open-suite/authorization/internal/modules/organizations/repositories"
	"github.com/open-suite/authorization/internal/modules/organizations/services"
	"github.com/open-suite/authorization/internal/platform/logger"
	"github.com/open-suite/authorization/internal/shared"
	"github.com/open-suite/authorization/internal/shared/response"
)

type OrganizationControllerImpl struct {
	OrganizationService services.OrganizationService
	response            *response.Sender
	log                 *logger.LayerLogger
}

func NewOrganizationController(service services.OrganizationService, sender *response.Sender, appLogger *logger.Logger) OrganizationController {
	return &OrganizationControllerImpl{
		OrganizationService: service,
		response:            sender,
		log:                 appLogger.Layer("controller.organizations"),
	}
}

func (c *OrganizationControllerImpl) Find(w http.ResponseWriter, r *http.Request) {
	end := c.log.Start(r.Context(), "Find")

	params := shared.NewListParamsFromRequest(r)
	items, err := c.OrganizationService.Find(r.Context(), params)
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusInternalServerError, "error.internal", nil)
		return
	}

	end(nil, "count", len(items))
	c.response.Success(w, r, http.StatusOK, "organizations.find.success", dto.ListResponse[entities.Organization]{Items: items})
}

func (c *OrganizationControllerImpl) FindByID(w http.ResponseWriter, r *http.Request) {
	end := c.log.Start(r.Context(), "FindByID")

	id, err := parseID(r)
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "organizations.invalid_id", nil)
		return
	}

	item, err := c.OrganizationService.FindByID(r.Context(), id)
	if errors.Is(err, repositories.ErrNotFound) {
		end(err)
		c.response.Error(w, r, http.StatusNotFound, "organizations.not_found", nil)
		return
	}
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusInternalServerError, "error.internal", nil)
		return
	}

	end(nil)
	c.response.Success(w, r, http.StatusOK, "organizations.find_by_id.success", item)
}

func (c *OrganizationControllerImpl) Create(w http.ResponseWriter, r *http.Request) {
	end := c.log.Start(r.Context(), "Create")

	var request dto.CreateOrganizationRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "organizations.invalid_payload", nil)
		return
	}

	item, err := c.OrganizationService.Create(r.Context(), request)
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "organizations.create.failed", nil)
		return
	}

	end(nil)
	c.response.Success(w, r, http.StatusCreated, "organizations.create.success", item)
}

func (c *OrganizationControllerImpl) Update(w http.ResponseWriter, r *http.Request) {
	end := c.log.Start(r.Context(), "Update")

	id, err := parseID(r)
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "organizations.invalid_id", nil)
		return
	}

	var request dto.UpdateOrganizationRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "organizations.invalid_payload", nil)
		return
	}

	item, err := c.OrganizationService.Update(r.Context(), id, request)
	if errors.Is(err, repositories.ErrNotFound) {
		end(err)
		c.response.Error(w, r, http.StatusNotFound, "organizations.not_found", nil)
		return
	}
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "organizations.update.failed", nil)
		return
	}

	end(nil)
	c.response.Success(w, r, http.StatusOK, "organizations.update.success", item)
}

func (c *OrganizationControllerImpl) Delete(w http.ResponseWriter, r *http.Request) {
	end := c.log.Start(r.Context(), "Delete")

	id, err := parseID(r)
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "organizations.invalid_id", nil)
		return
	}

	if err := c.OrganizationService.Delete(r.Context(), id); err != nil {
		if errors.Is(err, repositories.ErrNotFound) {
			end(err)
			c.response.Error(w, r, http.StatusNotFound, "organizations.not_found", nil)
			return
		}
		end(err)
		c.response.Error(w, r, http.StatusInternalServerError, "error.internal", nil)
		return
	}

	end(nil)
	c.response.Success(w, r, http.StatusOK, "organizations.delete.success", nil)
}

func parseID(r *http.Request) (int64, error) {
	return strconv.ParseInt(r.PathValue("id"), 10, 64)
}
