package controllers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/open-suite/authorization/internal/modules/auth/dto"
	"github.com/open-suite/authorization/internal/modules/auth/services"
	"github.com/open-suite/authorization/internal/platform/logger"
	"github.com/open-suite/authorization/internal/shared/requestctx"
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

func (c *AuthControllerImpl) JWKS(w http.ResponseWriter, r *http.Request) {
	end := c.log.Start(r.Context(), "JWKS")
	jwks, err := c.AuthService.JWKS(r.Context())
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "auth.jwks.failed", nil)
		return
	}

	end(nil)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(jwks)
}

func (c *AuthControllerImpl) Login(w http.ResponseWriter, r *http.Request) {
	end := c.log.Start(r.Context(), "Login")

	var request dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "auth.login.invalid_payload", nil)
		return
	}

	result, err := c.AuthService.Login(r.Context(), request)
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusUnauthorized, "auth.login.failed", nil)
		return
	}
	if request.SetCookie {
		setSessionCookies(w, r, result)
	}

	end(nil)
	c.response.Success(w, r, http.StatusOK, "auth.login.success", result)
}

func (c *AuthControllerImpl) Refresh(w http.ResponseWriter, r *http.Request) {
	end := c.log.Start(r.Context(), "Refresh")

	var request dto.RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		if !errors.Is(err, io.EOF) {
			end(err)
			c.response.Error(w, r, http.StatusBadRequest, "auth.refresh.invalid_payload", nil)
			return
		}
	}
	if request.RefreshToken == "" {
		if cookie, err := r.Cookie("refresh_token"); err == nil {
			request.RefreshToken = cookie.Value
		}
	}

	result, err := c.AuthService.Refresh(r.Context(), request)
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusUnauthorized, "auth.refresh.failed", nil)
		return
	}
	if request.SetCookie {
		setSessionCookies(w, r, result)
	}

	end(nil)
	c.response.Success(w, r, http.StatusOK, "auth.refresh.success", result)
}

func (c *AuthControllerImpl) Logout(w http.ResponseWriter, r *http.Request) {
	end := c.log.Start(r.Context(), "Logout")

	var request dto.LogoutRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		if !errors.Is(err, io.EOF) {
			end(err)
			c.response.Error(w, r, http.StatusBadRequest, "auth.logout.invalid_payload", nil)
			return
		}
	}
	if request.RefreshToken == "" {
		if cookie, err := r.Cookie("refresh_token"); err == nil {
			request.RefreshToken = cookie.Value
		}
	}

	if err := c.AuthService.Logout(r.Context(), request); err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "auth.logout.failed", nil)
		return
	}
	clearSessionCookies(w, r)

	end(nil)
	c.response.Success(w, r, http.StatusOK, "auth.logout.success", nil)
}

func (c *AuthControllerImpl) KeycloakRedirect(w http.ResponseWriter, r *http.Request) {
	end := c.log.Start(r.Context(), "KeycloakRedirect")

	redirectURL, err := c.AuthService.KeycloakRedirectURL(
		r.Context(),
		keycloakAppCallbackURL(r),
		strings.TrimSpace(r.URL.Query().Get("prompt")),
	)
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "auth.keycloak_redirect.failed", nil)
		return
	}

	end(nil)
	http.Redirect(w, r, redirectURL, http.StatusFound)
}

func (c *AuthControllerImpl) KeycloakCallback(w http.ResponseWriter, r *http.Request) {
	end := c.log.Start(r.Context(), "KeycloakCallback")

	if oauthErr := r.URL.Query().Get("error"); oauthErr != "" {
		end(nil, "oauth_error", oauthErr)
		redirectURL, err := c.AuthService.HandleKeycloakErrorCallback(r.Context(), r.URL.Query().Get("state"), oauthErr)
		if err != nil {
			end(err)
		}
		if redirectURL != "" {
			http.Redirect(w, r, redirectURL, http.StatusFound)
			return
		}
		c.response.Error(w, r, http.StatusBadRequest, "auth.keycloak_callback.failed", map[string]string{"error": oauthErr})
		return
	}

	result, redirectURL, err := c.AuthService.HandleKeycloakCallback(
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
		c.response.Error(w, r, http.StatusBadRequest, "auth.keycloak_callback.failed", nil)
		return
	}

	end(nil)
	if redirectURL != "" {
		http.Redirect(w, r, redirectURL, http.StatusFound)
		return
	}
	c.response.Success(w, r, http.StatusOK, "auth.keycloak_callback.success", result)
}

func (c *AuthControllerImpl) KeycloakExchange(w http.ResponseWriter, r *http.Request) {
	end := c.log.Start(r.Context(), "KeycloakExchange")

	var request dto.KeycloakExchangeRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "auth.keycloak_exchange.invalid_payload", nil)
		return
	}

	result, err := c.AuthService.ExchangeKeycloakCallbackCode(r.Context(), request.Code)
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusUnauthorized, "auth.keycloak_exchange.failed", nil)
		return
	}

	end(nil)
	c.response.Success(w, r, http.StatusOK, "auth.keycloak_exchange.success", result)
}

func keycloakAppCallbackURL(r *http.Request) string {
	callbackURL := strings.TrimSpace(r.URL.Query().Get("callback_url"))
	if callbackURL == "" {
		callbackURL = strings.TrimSpace(r.URL.Query().Get("callback"))
	}
	return callbackURL
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

func (c *AuthControllerImpl) CurrentUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := currentUserID(r)
	if !ok {
		c.response.Error(w, r, http.StatusUnauthorized, "auth.unauthorized", nil)
		return
	}
	result, err := c.AuthService.CurrentUserAccess(r.Context(), userID)
	if err != nil {
		if errors.Is(err, services.ErrInvalidSession) {
			c.response.Error(w, r, http.StatusUnauthorized, "auth.unauthorized", nil)
			return
		}
		c.response.Error(w, r, http.StatusInternalServerError, "auth.current_user.failed", nil)
		return
	}
	c.response.Success(w, r, http.StatusOK, "auth.current_user.success", result)
}

func (c *AuthControllerImpl) UpdateCurrentUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := currentUserID(r)
	if !ok {
		c.response.Error(w, r, http.StatusUnauthorized, "auth.unauthorized", nil)
		return
	}

	var request dto.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		c.response.Error(w, r, http.StatusBadRequest, "auth.update_user.invalid_payload", nil)
		return
	}

	result, err := c.AuthService.UpdateUser(r.Context(), userID, request)
	if err != nil {
		c.response.Error(w, r, http.StatusBadRequest, "auth.update_user.failed", nil)
		return
	}
	c.response.Success(w, r, http.StatusOK, "auth.update_user.success", result)
}

func (c *AuthControllerImpl) CurrentUserApps(w http.ResponseWriter, r *http.Request) {
	userID, ok := currentUserID(r)
	if !ok {
		c.response.Error(w, r, http.StatusUnauthorized, "auth.unauthorized", nil)
		return
	}
	result, err := c.AuthService.Apps(r.Context(), userID)
	if err != nil {
		c.response.Error(w, r, http.StatusInternalServerError, "auth.access.failed", nil)
		return
	}
	c.response.Success(w, r, http.StatusOK, "auth.access.apps.success", result)
}

func (c *AuthControllerImpl) CurrentUserAccessSummary(w http.ResponseWriter, r *http.Request) {
	userID, appCode, err := parseCurrentUserAccessPath(r)
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

func (c *AuthControllerImpl) CurrentUserAccessMenus(w http.ResponseWriter, r *http.Request) {
	userID, appCode, err := parseCurrentUserAccessPath(r)
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

func (c *AuthControllerImpl) CurrentUserAccessPermissions(w http.ResponseWriter, r *http.Request) {
	userID, appCode, err := parseCurrentUserAccessPath(r)
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

func (c *AuthControllerImpl) CurrentUserAccessCheck(w http.ResponseWriter, r *http.Request) {
	userID, appCode, err := parseCurrentUserAccessPath(r)
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

func (c *AuthControllerImpl) CurrentUserAccessToken(w http.ResponseWriter, r *http.Request) {
	userID, appCode, err := parseCurrentUserAccessPath(r)
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

func (c *AuthControllerImpl) UpdateUser(w http.ResponseWriter, r *http.Request) {
	userID, err := parseUserID(r)
	if err != nil {
		c.response.Error(w, r, http.StatusBadRequest, "auth.update_user.invalid_payload", nil)
		return
	}

	var request dto.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		c.response.Error(w, r, http.StatusBadRequest, "auth.update_user.invalid_payload", nil)
		return
	}

	result, err := c.AuthService.UpdateUser(r.Context(), userID, request)
	if err != nil {
		c.response.Error(w, r, http.StatusBadRequest, "auth.update_user.failed", nil)
		return
	}
	c.response.Success(w, r, http.StatusOK, "auth.update_user.success", result)
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

func currentUserID(r *http.Request) (int64, bool) {
	userID := requestctx.UserID(r.Context())
	return userID, userID > 0
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

func parseCurrentUserAccessPath(r *http.Request) (int64, string, error) {
	userID, ok := currentUserID(r)
	if !ok {
		return 0, "", errors.New("empty user")
	}
	appCode := strings.TrimSpace(r.PathValue("app"))
	if appCode == "" {
		return 0, "", errors.New("empty app")
	}
	return userID, appCode, nil
}

func setSessionCookies(w http.ResponseWriter, r *http.Request, session *dto.SessionResponse) {
	secure := r.TLS != nil || strings.EqualFold(r.Header.Get("X-Forwarded-Proto"), "https")
	setCookie(w, "access_token", session.AccessToken, int(session.ExpiresIn), secure)
	if session.RefreshToken != "" {
		setCookie(w, "refresh_token", session.RefreshToken, int(session.RefreshExpiresIn), secure)
	}
	if session.IDToken != "" {
		setCookie(w, "id_token", session.IDToken, int(session.ExpiresIn), secure)
	}
}

func clearSessionCookies(w http.ResponseWriter, r *http.Request) {
	secure := r.TLS != nil || strings.EqualFold(r.Header.Get("X-Forwarded-Proto"), "https")
	setCookie(w, "access_token", "", -1, secure)
	setCookie(w, "refresh_token", "", -1, secure)
	setCookie(w, "id_token", "", -1, secure)
}

func setCookie(w http.ResponseWriter, name string, value string, maxAge int, secure bool) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		MaxAge:   maxAge,
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
	})
}
