package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/open-suite/authorization/internal/entities"
	"github.com/open-suite/authorization/internal/modules/users/dto"
	"github.com/open-suite/authorization/internal/modules/users/repositories"
	"github.com/open-suite/authorization/internal/modules/users/services"
	"github.com/open-suite/authorization/internal/platform/logger"
	"github.com/open-suite/authorization/internal/shared"
	"github.com/open-suite/authorization/internal/shared/response"
)

type UserControllerImpl struct {
	UserService services.UserService
	response    *response.Sender
	log         *logger.LayerLogger
}

func NewUserController(service services.UserService, sender *response.Sender, appLogger *logger.Logger) UserController {
	return &UserControllerImpl{
		UserService: service,
		response:    sender,
		log:         appLogger.Layer("controller.users"),
	}
}

func (c *UserControllerImpl) Find(w http.ResponseWriter, r *http.Request) {
	end := c.log.Start(r.Context(), "Find")

	params := shared.NewListParamsFromRequest(r)
	items, err := c.UserService.Find(r.Context(), params)
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusInternalServerError, "error.internal", nil)
		return
	}

	end(nil, "count", len(items))
	c.response.Success(w, r, http.StatusOK, "users.find.success", dto.ListResponse[entities.User]{Items: items})
}

func (c *UserControllerImpl) FindByID(w http.ResponseWriter, r *http.Request) {
	end := c.log.Start(r.Context(), "FindByID")

	id, err := parseID(r)
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "users.invalid_id", nil)
		return
	}

	item, err := c.UserService.FindByID(r.Context(), id)
	if errors.Is(err, repositories.ErrNotFound) {
		end(err)
		c.response.Error(w, r, http.StatusNotFound, "users.not_found", nil)
		return
	}
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusInternalServerError, "error.internal", nil)
		return
	}

	end(nil)
	c.response.Success(w, r, http.StatusOK, "users.find_by_id.success", item)
}

func (c *UserControllerImpl) Create(w http.ResponseWriter, r *http.Request) {
	end := c.log.Start(r.Context(), "Create")

	var request dto.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "users.invalid_payload", nil)
		return
	}

	item, err := c.UserService.Create(r.Context(), request)
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "users.create.failed", nil)
		return
	}

	end(nil)
	c.response.Success(w, r, http.StatusCreated, "users.create.success", item)
}

func (c *UserControllerImpl) Signup(w http.ResponseWriter, r *http.Request) {
	end := c.log.Start(r.Context(), "Signup")

	var request dto.SignupUserRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "users.invalid_payload", nil)
		return
	}

	item, err := c.UserService.Signup(r.Context(), request)
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "users.signup.failed", nil)
		return
	}

	end(nil)
	c.response.Success(w, r, http.StatusCreated, "users.signup.success", item)
}

func (c *UserControllerImpl) SignupWithGoogle(w http.ResponseWriter, r *http.Request) {
	end := c.log.Start(r.Context(), "SignupWithGoogle")

	var request dto.GoogleSignupRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "users.invalid_payload", nil)
		return
	}

	item, err := c.UserService.SignupWithGoogle(r.Context(), request)
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "users.signup_google.failed", nil)
		return
	}

	end(nil)
	c.response.Success(w, r, http.StatusCreated, "users.signup_google.success", item)
}

func (c *UserControllerImpl) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	end := c.log.Start(r.Context(), "VerifyEmail")

	request := dto.VerifyEmailRequest{Code: r.URL.Query().Get("code")}
	if request.Code == "" {
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			end(err)
			c.response.Error(w, r, http.StatusBadRequest, "users.invalid_payload", nil)
			return
		}
	}

	if err := c.UserService.VerifyEmail(r.Context(), request); err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "users.verify_email.failed", nil)
		return
	}

	end(nil)
	c.response.Success(w, r, http.StatusOK, "users.verify_email.success", nil)
}

func (c *UserControllerImpl) Update(w http.ResponseWriter, r *http.Request) {
	end := c.log.Start(r.Context(), "Update")

	id, err := parseID(r)
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "users.invalid_id", nil)
		return
	}

	var request dto.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "users.invalid_payload", nil)
		return
	}

	item, err := c.UserService.Update(r.Context(), id, request)
	if errors.Is(err, repositories.ErrNotFound) {
		end(err)
		c.response.Error(w, r, http.StatusNotFound, "users.not_found", nil)
		return
	}
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "users.update.failed", nil)
		return
	}

	end(nil)
	c.response.Success(w, r, http.StatusOK, "users.update.success", item)
}

func (c *UserControllerImpl) Delete(w http.ResponseWriter, r *http.Request) {
	end := c.log.Start(r.Context(), "Delete")

	id, err := parseID(r)
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "users.invalid_id", nil)
		return
	}

	if err := c.UserService.Delete(r.Context(), id); err != nil {
		if errors.Is(err, repositories.ErrNotFound) {
			end(err)
			c.response.Error(w, r, http.StatusNotFound, "users.not_found", nil)
			return
		}
		end(err)
		c.response.Error(w, r, http.StatusInternalServerError, "error.internal", nil)
		return
	}

	end(nil)
	c.response.Success(w, r, http.StatusOK, "users.delete.success", nil)
}

func parseID(r *http.Request) (int64, error) {
	return strconv.ParseInt(r.PathValue("id"), 10, 64)
}
