package config

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv             string
	DatabaseURL        string
	Port               string
	TrustedProxies     []string
	AuthCookieSecure   bool
	AuthCookieDomain   string
	AuthCookieSameSite http.SameSite
}

func Load() (*Config, error) {
	var err error = godotenv.Load()

	if err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	databaseURL, err := buildDatabaseURL()
	if err != nil {
		return nil, err
	}

	appEnv := os.Getenv("APP_ENV")
	if appEnv == "" {
		appEnv = "local"
	}

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	authCookieSecure, err := optionalBoolEnv("AUTH_COOKIE_SECURE", appEnv == "production")
	if err != nil {
		return nil, err
	}

	authCookieSameSite, err := sameSiteEnv("AUTH_COOKIE_SAME_SITE", http.SameSiteLaxMode)
	if err != nil {
		return nil, err
	}

	var config *Config = &Config{
		AppEnv:             appEnv,
		DatabaseURL:        databaseURL,
		Port:               port,
		TrustedProxies:     commaSeparatedEnv("TRUSTED_PROXIES"),
		AuthCookieSecure:   authCookieSecure,
		AuthCookieDomain:   strings.TrimSpace(os.Getenv("AUTH_COOKIE_DOMAIN")),
		AuthCookieSameSite: authCookieSameSite,
	}

	return config, nil
}

func buildDatabaseURL() (string, error) {
	host, err := requiredEnvValues("DB_HOST")
	if err != nil {
		return "", err
	}

	port, err := requiredEnvValues("DB_PORT")
	if err != nil {
		return "", err
	}

	user, err := requiredEnvValues("DB_USER")
	if err != nil {
		return "", err
	}

	password, err := requiredEnvValues("DB_PASSWORD")
	if err != nil {
		return "", err
	}

	name, err := requiredEnvValues("DB_NAME")
	if err != nil {
		return "", err
	}

	sslMode := os.Getenv("DB_SSLMODE")
	if sslMode == "" {
		sslMode = "disable"
	}

	databaseURL := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(user, password),
		Host:   host + ":" + port,
		Path:   name,
	}

	query := databaseURL.Query()
	query.Set("sslmode", sslMode)
	databaseURL.RawQuery = query.Encode()

	return databaseURL.String(), nil
}

func requiredEnvValues(key string) (string, error) {
	value := os.Getenv(key)
	if value == "" {
		return "", fmt.Errorf("%s is required", key)
	}

	return value, nil
}

func commaSeparatedEnv(key string) []string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return nil
	}

	parts := strings.Split(value, ",")
	values := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			values = append(values, part)
		}
	}

	if len(values) == 0 {
		return nil
	}

	return values
}

func optionalBoolEnv(key string, defaultValue bool) (bool, error) {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return defaultValue, nil
	}

	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return false, fmt.Errorf("%s must be a boolean value", key)
	}

	return parsed, nil
}

func sameSiteEnv(key string, defaultValue http.SameSite) (http.SameSite, error) {
	value := strings.ToLower(strings.TrimSpace(os.Getenv(key)))
	if value == "" {
		return defaultValue, nil
	}

	switch value {
	case "default":
		return http.SameSiteDefaultMode, nil
	case "lax":
		return http.SameSiteLaxMode, nil
	case "strict":
		return http.SameSiteStrictMode, nil
	case "none":
		return http.SameSiteNoneMode, nil
	default:
		return http.SameSiteDefaultMode, fmt.Errorf("%s must be one of default, lax, strict, or none", key)
	}
}
