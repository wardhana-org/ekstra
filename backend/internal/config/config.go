package config

import (
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	Port        string
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

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	var config *Config = &Config{
		DatabaseURL: databaseURL,
		Port:        port,
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
