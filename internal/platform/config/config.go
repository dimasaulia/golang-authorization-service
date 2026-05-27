package config

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	App      AppConfig
	Logger   LoggerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	CORS     CORSConfig
	Mail     MailConfig
	OAuth    OAuthConfig
	Keycloak KeycloakConfig
	FreeIPA  FreeIPAConfig
	Authz    AuthzConfig
}

type AppConfig struct {
	Port      string
	PublicURL string
}

type CORSConfig struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	AllowCredentials bool
	MaxAge           int
}

type LoggerConfig struct {
	Level  string
	LogDir string
}

type DatabaseConfig struct {
	URL             string
	MaxOpenConns    int32
	MinConns        int32
	MaxConnLifetime time.Duration
	MaxConnIdleTime time.Duration
}

type RedisConfig struct {
	Addr     string
	Username string
	Password string
	DB       int
}

type MailConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
	Enabled  bool
}

type OAuthConfig struct {
	Google GoogleOAuthConfig
}

type GoogleOAuthConfig struct {
	ClientID              string
	ClientSecret          string
	RedirectURL           string
	Scopes                []string
	DefaultOrganizationID int64
	SuccessRedirectURL    string
	FailureRedirectURL    string
}

type KeycloakConfig struct {
	Enabled      bool
	BaseURL      string
	Realm        string
	ClientID     string
	ClientSecret string
}

type FreeIPAConfig struct {
	Enabled            bool
	BaseURL            string
	Username           string
	Password           string
	InsecureSkipVerify bool
}

type AuthzConfig struct {
	CacheTTL    time.Duration
	TokenSecret string
	TokenTTL    time.Duration
}

func Load() Config {
	_ = godotenv.Load()

	return Config{
		App: AppConfig{
			Port:      env("APP_PORT", "8080"),
			PublicURL: strings.TrimRight(env("APP_PUBLIC_URL", "http://localhost:8080"), "/"),
		},
		Logger: LoggerConfig{
			Level:  env("LOG_LEVEL", "info"),
			LogDir: env("LOG_DIR", ".logs"),
		},
		Database: DatabaseConfig{
			URL:             env("DATABASE_URL", "postgres://postgres:ldvgKEuDrPYuKULN@127.0.0.1:5432/authorization_service?sslmode=disable"),
			MaxOpenConns:    int32Env("DATABASE_MAX_OPEN_CONNS", 10),
			MinConns:        int32Env("DATABASE_MIN_CONNS", 1),
			MaxConnLifetime: durationEnv("DATABASE_MAX_CONN_LIFETIME", time.Hour),
			MaxConnIdleTime: durationEnv("DATABASE_MAX_CONN_IDLE_TIME", 30*time.Minute),
		},
		Redis: RedisConfig{
			Addr:     env("REDIS_ADDR", "dimasaulia.com:6379"),
			Username: env("REDIS_USERNAME", ""),
			Password: env("REDIS_PASSWORD", ""),
			DB:       intEnv("REDIS_DB", 0),
		},
		CORS: CORSConfig{
			AllowedOrigins:   splitEnv("CORS_ALLOWED_ORIGINS", []string{"*"}),
			AllowedMethods:   splitEnv("CORS_ALLOWED_METHODS", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"}),
			AllowedHeaders:   splitEnv("CORS_ALLOWED_HEADERS", []string{"Content-Type", "Authorization"}),
			AllowCredentials: boolEnv("CORS_ALLOW_CREDENTIALS", true),
			MaxAge:           intEnv("CORS_MAX_AGE", 86400),
		},
		Mail: MailConfig{
			Host:     env("MAIL_HOST", ""),
			Port:     intEnv("MAIL_PORT", 587),
			Username: env("MAIL_USERNAME", ""),
			Password: env("MAIL_PASSWORD", ""),
			From:     env("MAIL_FROM", "noreply@localhost"),
			Enabled:  boolEnv("MAIL_ENABLED", false),
		},
		OAuth: OAuthConfig{
			Google: GoogleOAuthConfig{
				ClientID:              env("GOOGLE_CLIENT_ID", ""),
				ClientSecret:          env("GOOGLE_CLIENT_SECRET", ""),
				RedirectURL:           env("GOOGLE_REDIRECT_URL", env("APP_PUBLIC_URL", "http://localhost:8080")+"/api/v1/auth/google/callback"),
				Scopes:                splitEnv("GOOGLE_SCOPES", []string{"openid", "email", "profile"}),
				DefaultOrganizationID: int64Env("GOOGLE_DEFAULT_ORGANIZATION_ID", 0),
				SuccessRedirectURL:    env("GOOGLE_SUCCESS_REDIRECT_URL", ""),
				FailureRedirectURL:    env("GOOGLE_FAILURE_REDIRECT_URL", ""),
			},
		},
		Keycloak: KeycloakConfig{
			Enabled:      boolEnv("KEYCLOAK_ENABLED", false),
			BaseURL:      strings.TrimRight(env("KEYCLOAK_BASE_URL", ""), "/"),
			Realm:        env("KEYCLOAK_REALM", ""),
			ClientID:     env("KEYCLOAK_ADMIN_CLIENT_ID", ""),
			ClientSecret: env("KEYCLOAK_ADMIN_CLIENT_SECRET", ""),
		},
		FreeIPA: FreeIPAConfig{
			Enabled:            boolEnv("FREEIPA_ENABLED", false),
			BaseURL:            strings.TrimRight(env("FREEIPA_BASE_URL", ""), "/"),
			Username:           env("FREEIPA_USERNAME", ""),
			Password:           env("FREEIPA_PASSWORD", ""),
			InsecureSkipVerify: boolEnv("FREEIPA_INSECURE_SKIP_VERIFY", false),
		},
		Authz: AuthzConfig{
			CacheTTL:    durationEnv("AUTHZ_CACHE_TTL", 5*time.Minute),
			TokenSecret: env("AUTHZ_TOKEN_SECRET", "change-me"),
			TokenTTL:    durationEnv("AUTHZ_TOKEN_TTL", 15*time.Minute),
		},
	}
}

func env(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}

func intEnv(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}

	return parsed
}

func int32Env(key string, fallback int32) int32 {
	return int32(intEnv(key, int(fallback)))
}

func int64Env(key string, fallback int64) int64 {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return fallback
	}

	return parsed
}

func durationEnv(key string, fallback time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := time.ParseDuration(value)
	if err != nil {
		return fallback
	}

	return parsed
}

func splitEnv(key string, fallback []string) []string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parts := strings.Split(value, ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}

func boolEnv(key string, fallback bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}
	return parsed
}
