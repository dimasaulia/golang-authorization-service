package freeipa

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	"github.com/open-suite/authorization/internal/platform/config"
	"github.com/open-suite/authorization/internal/platform/logger"
)

var (
	ErrDisabled      = errors.New("freeipa integration is disabled")
	ErrNotConfigured = errors.New("freeipa integration is not configured")
)

type Client interface {
	Enabled() bool
	CreateUser(ctx context.Context, input CreateUserInput) (string, error)
	UpdateUser(ctx context.Context, uid string, input UpdateUserInput) (string, error)
	DeleteUser(ctx context.Context, uid string) error
}

type FreeIPA struct {
	cfg        config.FreeIPAConfig
	log        *logger.LayerLogger
	httpClient *http.Client
}

type CreateUserInput struct {
	Username    string
	Email       string
	DisplayName string
	Password    string
}

type UpdateUserInput struct {
	Username    string
	Email       string
	DisplayName string
	Password    string
}

type rpcRequest struct {
	Method string `json:"method"`
	Params []any  `json:"params"`
	ID     int64  `json:"id"`
}

type rpcResponse struct {
	Result any       `json:"result"`
	Error  *rpcError `json:"error"`
}

type rpcError struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
	Name    string `json:"name"`
}

func New(cfg config.Config, appLogger *logger.Logger) Client {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	if cfg.FreeIPA.InsecureSkipVerify {
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	jar, _ := cookiejar.New(nil)
	return &FreeIPA{
		cfg: cfg.FreeIPA,
		log: appLogger.Layer("platform.freeipa"),
		httpClient: &http.Client{
			Timeout:   15 * time.Second,
			Transport: transport,
			Jar:       jar,
		},
	}
}

func (f *FreeIPA) Enabled() bool {
	return f.cfg.Enabled
}

func (f *FreeIPA) CreateUser(ctx context.Context, input CreateUserInput) (string, error) {
	end := f.log.Start(ctx, "CreateUser")
	if !f.cfg.Enabled {
		end(ErrDisabled)
		return "", ErrDisabled
	}
	if err := f.validateConfig(); err != nil {
		end(err)
		return "", err
	}
	if err := f.login(ctx); err != nil {
		end(err)
		return "", err
	}

	uid := normalizeUID(input.Username)
	givenName, surname := splitDisplayName(input.DisplayName)
	options := map[string]any{
		"givenname": givenName,
		"sn":        surname,
		"cn":        strings.TrimSpace(input.DisplayName),
		"mail":      strings.ToLower(strings.TrimSpace(input.Email)),
		"version":   "2.251",
	}
	if strings.TrimSpace(input.Password) == "" {
		options["random"] = true
	} else {
		options["userpassword"] = input.Password
	}
	if strings.TrimSpace(input.DisplayName) == "" {
		options["cn"] = uid
	}

	if err := f.rpc(ctx, "user_add", []any{[]string{uid}, options}, nil); err != nil {
		if isDuplicateError(err) {
			end(nil, "uid", uid, "duplicate", true)
			return uid, nil
		}
		end(err)
		return "", err
	}

	end(nil, "uid", uid)
	return uid, nil
}

func (f *FreeIPA) UpdateUser(ctx context.Context, uid string, input UpdateUserInput) (string, error) {
	end := f.log.Start(ctx, "UpdateUser", "uid", uid)
	if !f.cfg.Enabled {
		end(ErrDisabled)
		return "", ErrDisabled
	}
	uid = normalizeUID(uid)
	if uid == "" {
		end(nil)
		return "", nil
	}
	if err := f.validateConfig(); err != nil {
		end(err)
		return "", err
	}
	if err := f.login(ctx); err != nil {
		end(err)
		return "", err
	}

	nextUID := uid
	options := map[string]any{"version": "2.251"}
	if strings.TrimSpace(input.Username) != "" {
		nextUID = normalizeUID(input.Username)
		if nextUID != "" && nextUID != uid {
			options["rename"] = nextUID
		}
	}
	if strings.TrimSpace(input.Email) != "" {
		options["mail"] = strings.ToLower(strings.TrimSpace(input.Email))
	}
	if strings.TrimSpace(input.DisplayName) != "" {
		givenName, surname := splitDisplayName(input.DisplayName)
		options["givenname"] = givenName
		options["sn"] = surname
		options["cn"] = strings.TrimSpace(input.DisplayName)
	}
	if input.Password != "" {
		options["userpassword"] = input.Password
	}

	if len(options) == 1 {
		end(nil, "uid", uid)
		return uid, nil
	}

	if err := f.rpc(ctx, "user_mod", []any{[]string{uid}, options}, nil); err != nil {
		end(err)
		return "", err
	}

	end(nil, "uid", nextUID)
	return nextUID, nil
}

func (f *FreeIPA) DeleteUser(ctx context.Context, uid string) error {
	end := f.log.Start(ctx, "DeleteUser", "uid", uid)
	if !f.cfg.Enabled {
		end(ErrDisabled)
		return ErrDisabled
	}
	uid = normalizeUID(uid)
	if uid == "" {
		end(nil)
		return nil
	}
	if err := f.validateConfig(); err != nil {
		end(err)
		return err
	}
	if err := f.login(ctx); err != nil {
		end(err)
		return err
	}

	err := f.rpc(ctx, "user_del", []any{[]string{uid}, map[string]any{"version": "2.251"}}, nil)
	if err != nil && !isNotFoundError(err) {
		end(err)
		return err
	}

	end(nil)
	return nil
}

func (f *FreeIPA) validateConfig() error {
	if f.cfg.BaseURL == "" || f.cfg.Username == "" || f.cfg.Password == "" {
		return ErrNotConfigured
	}
	return nil
}

func (f *FreeIPA) login(ctx context.Context) error {
	values := url.Values{}
	values.Set("user", f.cfg.Username)
	values.Set("password", f.cfg.Password)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, f.cfg.BaseURL+"/ipa/session/login_password", strings.NewReader(values.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "text/plain")
	req.Header.Set("Referer", f.cfg.BaseURL+"/ipa")

	resp, err := f.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
		return fmt.Errorf("freeipa login failed: status=%d body=%s", resp.StatusCode, string(body))
	}
	return nil
}

func (f *FreeIPA) rpc(ctx context.Context, method string, params []any, target any) error {
	payload, err := json.Marshal(rpcRequest{
		Method: method,
		Params: params,
		ID:     time.Now().UnixNano(),
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, f.cfg.BaseURL+"/ipa/session/json", bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Referer", f.cfg.BaseURL+"/ipa")

	resp, err := f.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return err
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("freeipa rpc failed: status=%d body=%s", resp.StatusCode, string(body))
	}

	var rpcResp rpcResponse
	if err := json.Unmarshal(body, &rpcResp); err != nil {
		return err
	}
	if rpcResp.Error != nil {
		return fmt.Errorf("freeipa rpc error: code=%d name=%s message=%s", rpcResp.Error.Code, rpcResp.Error.Name, rpcResp.Error.Message)
	}
	if target == nil {
		return nil
	}

	result, err := json.Marshal(rpcResp.Result)
	if err != nil {
		return err
	}
	return json.Unmarshal(result, target)
}

func normalizeUID(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	value = strings.ReplaceAll(value, " ", ".")
	return value
}

func splitDisplayName(displayName string) (string, string) {
	parts := strings.Fields(displayName)
	if len(parts) == 0 {
		return "User", "User"
	}
	if len(parts) == 1 {
		return parts[0], parts[0]
	}
	return parts[0], strings.Join(parts[1:], " ")
}

func isDuplicateError(err error) bool {
	message := strings.ToLower(err.Error())
	return strings.Contains(message, "already exists") || strings.Contains(message, "duplicate")
}

func isNotFoundError(err error) bool {
	message := strings.ToLower(err.Error())
	return strings.Contains(message, "not found") || strings.Contains(message, "no such entry")
}
