package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	PublicHost         string
	Port               string
	Environment        string
	DBUrl              string
	DBAuthToken        string
	CookiesAuthSecret  string
	GithubClientID     string
	GithubClientSecret string
	GoogleClientID     string
	GoogleClientSecret string
)

func LoadConfig() {
	godotenv.Load()

	PublicHost = getEnv("PUBLIC_HOST", "localhost")
	Port = getEnv("PORT", "8080")
	Environment = getEnvOrError("ENVIRONMENT")
	DBUrl = getEnvOrError("TURSO_DATABASE_URL")
	DBAuthToken = getEnvOrError("TURSO_AUTH_TOKEN")
	CookiesAuthSecret = getEnv("COOKIES_AUTH_SECRET", "some-very-secret-key")
	GithubClientID = getEnvOrError("GITHUB_CLIENT_ID")
	GithubClientSecret = getEnvOrError("GITHUB_CLIENT_SECRET")
	GoogleClientID = getEnvOrError("GOOGLE_CLIENT_ID")
	GoogleClientSecret = getEnvOrError("GOOGLE_CLIENT_SECRET")

}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvOrError(key string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	panic(fmt.Sprintf("Environment variable %s is not set", key))
}

func getEnvAsInt(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		i, err := strconv.Atoi(value)
		if err != nil {
			return fallback
		}
		return i
	}
	return fallback
}

func getEnvAsBool(key string, fallback bool) bool {
	if value, ok := os.LookupEnv(key); ok {
		b, err := strconv.ParseBool(value)
		if err != nil {
			return fallback
		}
		return b
	}
	return fallback
}
