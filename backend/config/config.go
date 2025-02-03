package config

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Environment struct {
	Production         bool
	PublicHost         string
	Port               string
	PrometheusPort     string
	FrontendUrl        string
	DBConnectionString string
	SessionSecret      string
	HMAC               string
	GoogleID           string
	GoogleSecret       string
	GoogleCallback     string
}

var Env = getEnvironment()

func getEnvironment() Environment {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		os.Exit(1)
	}

	slog.Info("Loading environment variables.")

	return Environment{
		Production:         getBoolEnv("PRODUCTION", false),
		Port:               getEnv("PORT", ":8080"),
		PrometheusPort:     getEnv("PROMETHEUS_PORT", ":8081"),
		FrontendUrl:        getEnv("FRONTEND_URL", "http://localhost:5173"),
		DBConnectionString: getEnv("DB_CONNECTION_STRING", ""),
		SessionSecret:      getEnv("SESSION_SECRET", ""),
		HMAC:               getEnv("HMAC", ""),
		GoogleID:           getEnv("GOOGLE_ID", ""),
		GoogleSecret:       getEnv("GOOGLE_SECRET", ""),
		GoogleCallback:     getEnv("GOOGLE_CALLBACK", ""),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	slog.Info(fmt.Sprintf("%s not found in environment, defaulting to: %s", key, fallback))
	return fallback
}

func getBoolEnv(key string, fallback bool) bool {
	if value := getEnv(key, ""); value != "" {
		return value == "true"
	}
	slog.Info(fmt.Sprintf("%s not found in environment, defaulting to: %s", key, strconv.FormatBool(fallback)))
	return fallback
}
