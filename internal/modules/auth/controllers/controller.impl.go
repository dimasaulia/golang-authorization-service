package controllers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/open-suite/authorization/internal/modules/auth/services"
	"github.com/open-suite/authorization/internal/platform/logger"
	"github.com/open-suite/authorization/internal/shared/response"
)

type AuthControllerImpl struct {
	AuthService services.AuthService
	response    *response.Sender
	log         *logger.LayerLogger
}

func NewAuthController(service services.AuthService, sender *response.Sender, appLogger *logger.Logger) AuthController {
	return &AuthControllerImpl{
		AuthService: service,
		response:    sender,
		log:         appLogger.Layer("controller.auth"),
	}
}

func (c *AuthControllerImpl) GoogleRedirect(w http.ResponseWriter, r *http.Request) {
	end := c.log.Start(r.Context(), "GoogleRedirect")

	organizationID, err := parseOptionalInt64(r.URL.Query().Get("organization_id"))
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "auth.google_redirect.failed", nil)
		return
	}

	redirectURL, err := c.AuthService.GoogleRedirectURL(r.Context(), organizationID)
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "auth.google_redirect.failed", nil)
		return
	}

	end(nil)
	http.Redirect(w, r, redirectURL, http.StatusFound)
}

func (c *AuthControllerImpl) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	end := c.log.Start(r.Context(), "GoogleCallback")

	if oauthErr := r.URL.Query().Get("error"); oauthErr != "" {
		end(nil, "oauth_error", oauthErr)
		c.response.Error(w, r, http.StatusBadRequest, "auth.google_callback.failed", map[string]string{"error": oauthErr})
		return
	}

	result, redirectURL, err := c.AuthService.HandleGoogleCallback(
		r.Context(),
		r.URL.Query().Get("code"),
		r.URL.Query().Get("state"),
	)
	if err != nil {
		end(err)
		if redirectURL != "" {
			http.Redirect(w, r, redirectURL, http.StatusFound)
			return
		}
		c.response.Error(w, r, http.StatusBadRequest, "auth.google_callback.failed", nil)
		return
	}

	end(nil)
	if redirectURL != "" {
		http.Redirect(w, r, redirectURL, http.StatusFound)
		return
	}
	c.response.Success(w, r, http.StatusOK, "auth.google_callback.success", result)
}

func (c *AuthControllerImpl) UserApps(w http.ResponseWriter, r *http.Request) {
	userID, err := parseUserID(r)
	if err != nil {
		c.response.Error(w, r, http.StatusBadRequest, "auth.access.invalid_request", nil)
		return
	}
	result, err := c.AuthService.Apps(r.Context(), userID)
	if err != nil {
		c.response.Error(w, r, http.StatusInternalServerError, "auth.access.failed", nil)
		return
	}
	c.response.Success(w, r, http.StatusOK, "auth.access.apps.success", result)
}

func (c *AuthControllerImpl) AccessSummary(w http.ResponseWriter, r *http.Request) {
	userID, appCode, err := parseAccessPath(r)
	if err != nil {
		c.response.Error(w, r, http.StatusBadRequest, "auth.access.invalid_request", nil)
		return
	}
	result, err := c.AuthService.AccessSummary(r.Context(), userID, appCode)
	if err != nil {
		c.response.Error(w, r, http.StatusInternalServerError, "auth.access.failed", nil)
		return
	}
	c.response.Success(w, r, http.StatusOK, "auth.access.summary.success", result)
}

func (c *AuthControllerImpl) AccessMenus(w http.ResponseWriter, r *http.Request) {
	userID, appCode, err := parseAccessPath(r)
	if err != nil {
		c.response.Error(w, r, http.StatusBadRequest, "auth.access.invalid_request", nil)
		return
	}
	result, err := c.AuthService.Menus(r.Context(), userID, appCode)
	if err != nil {
		c.response.Error(w, r, http.StatusInternalServerError, "auth.access.failed", nil)
		return
	}
	c.response.Success(w, r, http.StatusOK, "auth.access.menus.success", map[string]any{"items": result})
}

func (c *AuthControllerImpl) AccessPermissions(w http.ResponseWriter, r *http.Request) {
	userID, appCode, err := parseAccessPath(r)
	if err != nil {
		c.response.Error(w, r, http.StatusBadRequest, "auth.access.invalid_request", nil)
		return
	}
	result, err := c.AuthService.Permissions(r.Context(), userID, appCode)
	if err != nil {
		c.response.Error(w, r, http.StatusInternalServerError, "auth.access.failed", nil)
		return
	}
	c.response.Success(w, r, http.StatusOK, "auth.access.permissions.success", map[string]any{"items": result})
}

func (c *AuthControllerImpl) AccessCheck(w http.ResponseWriter, r *http.Request) {
	userID, appCode, err := parseAccessPath(r)
	if err != nil {
		c.response.Error(w, r, http.StatusBadRequest, "auth.access.invalid_request", nil)
		return
	}
	permission := strings.TrimSpace(r.URL.Query().Get("permission"))
	if permission == "" {
		c.response.Error(w, r, http.StatusBadRequest, "auth.access.invalid_request", nil)
		return
	}
	result, err := c.AuthService.Check(r.Context(), userID, appCode, permission)
	if err != nil {
		c.response.Error(w, r, http.StatusInternalServerError, "auth.access.failed", nil)
		return
	}
	c.response.Success(w, r, http.StatusOK, "auth.access.check.success", result)
}

func (c *AuthControllerImpl) AccessToken(w http.ResponseWriter, r *http.Request) {
	userID, appCode, err := parseAccessPath(r)
	if err != nil {
		c.response.Error(w, r, http.StatusBadRequest, "auth.access.invalid_request", nil)
		return
	}
	result, err := c.AuthService.AccessToken(r.Context(), userID, appCode)
	if err != nil {
		c.response.Error(w, r, http.StatusInternalServerError, "auth.access.failed", nil)
		return
	}
	c.response.Success(w, r, http.StatusOK, "auth.access.token.success", result)
}

func parseOptionalInt64(value string) (int64, error) {
	if value == "" {
		return 0, nil
	}
	return strconv.ParseInt(value, 10, 64)
}

func parseUserID(r *http.Request) (int64, error) {
	return strconv.ParseInt(r.PathValue("user_id"), 10, 64)
}

func parseAccessPath(r *http.Request) (int64, string, error) {
	userID, err := parseUserID(r)
	if err != nil {
		return 0, "", err
	}
	appCode := strings.TrimSpace(r.PathValue("app"))
	if appCode == "" {
		return 0, "", errors.New("empty app")
	}
	return userID, appCode, nil
}
