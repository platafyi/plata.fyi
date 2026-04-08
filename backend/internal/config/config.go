package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	DBURL           string
	TurnstileSecret string
	IPHMACSecret    string
	AppURL          string
	APIURL          string
	Port            int
}

func Load() (*Config, error) {
	cfg := &Config{
		DBURL:           getEnv("DB_URL", "postgres://platafyi:platafyi@localhost:5432/platafyi?sslmode=disable"),
		TurnstileSecret: os.Getenv("TURNSTILE_SECRET"),
		IPHMACSecret:    os.Getenv("IP_HMAC_SECRET"),
		AppURL:          getEnv("APP_URL", "http://localhost:3000"),
		APIURL:          getEnv("API_URL", "http://localhost:8080"),
	}

	portStr := getEnv("PORT", "8080")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("invalid PORT %q: %w", portStr, err)
	}
	cfg.Port = port

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
