package config

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/omikkel/restore-vercel-deployments/internal/logger"
)

// Default values for configuration
const (
	defaultAPIURL          = "https://vercel.com/api"
	defaultLogLevel        = "info"
	defaultRestoreCooldown = 250
)

// Config holds all configuration values for the application
type Config struct {
	LogLevel        int
	APIURL          string
	APIToken        string
	RestoreCooldown time.Duration
}

// Load reads configuration from environment variables and .env file
func Load() (*Config, error) {
	// Load .env file if it exists (ignores error if file doesn't exist)
	_ = godotenv.Load()

	apiToken := getEnv("VERCEL_API_TOKEN", "")
	if apiToken == "" {
		return nil, errors.New("VERCEL_API_TOKEN environment variable is required")
	}

	return &Config{
		LogLevel:        parseLogLevel(getEnv("LOG_LEVEL", defaultLogLevel)),
		APIURL:          getEnv("VERCEL_API_URL", defaultAPIURL),
		APIToken:        apiToken,
		RestoreCooldown: parseRestoreCooldown(getEnv("RESTORE_COOLDOWN_MS", strconv.Itoa(defaultRestoreCooldown))),
	}, nil
}

// getEnv returns the value of an environment variable or a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// parseLogLevel converts a string log level to logger level constant
func parseLogLevel(level string) int {
	switch strings.ToLower(level) {
	case "debug":
		return logger.LevelDebug
	case "info":
		return logger.LevelInfo
	case "error":
		return logger.LevelError
	case "disabled":
		return logger.LevelDisabled
	default:
		return logger.LevelInfo
	}
}

// parseRestoreCooldown parses the cooldown string to a duration
func parseRestoreCooldown(cooldownStr string) time.Duration {
	cooldown, err := strconv.Atoi(cooldownStr)
	if err != nil {
		cooldown = defaultRestoreCooldown
	}
	return time.Duration(cooldown) * time.Millisecond
}
