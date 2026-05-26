package controllers

import (
	"fmt"
	"net/http"
	"strconv"

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

	fmt.Println("URL => ")
	fmt.Println(r.URL)

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

func parseOptionalInt64(value string) (int64, error) {
	if value == "" {
		return 0, nil
	}
	return strconv.ParseInt(value, 10, 64)
}
