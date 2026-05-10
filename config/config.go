package config

import "os"

// Config holds all runtime configuration values for the application.
type Config struct {
	// ServerAddress is the TCP address the HTTP server listens on (e.g. ":8080").
	ServerAddress string
}

// Load reads configuration from environment variables and returns the populated Config.
// If an expected variable is not set, a sensible default is used.
func Load() *Config {
	return &Config{
		ServerAddress: getEnv("SERVER_ADDRESS", ":8080"),
	}
}

func getEnv(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}
