package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/gommon/log"
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

var Envs = getEnvironment()

func getEnvironment() Environment {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	log.Info("Loading environment variables...")

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
	log.Infof("%s not found in environment, defaulting to: %s", key, fallback)
	return fallback
}

func getBoolEnv(key string, fallback bool) bool {
	if value := getEnv(key, ""); value != "" {
		return value == "true"
	}
	log.Infof("%s not found in environment, defaulting to: %v", key, fallback)
	return fallback
}
