package main

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all configurable values for the application.
type Config struct {
// rewrote this part
	Port     int
	LogLevel string
}

// LoadConfig reads environment variables (optionally from a .env file) and
// returns a populated Config struct. It aborts the program if required values
// are missing or malformed.
func LoadConfig() *Config {
	// Load .env file if present; ignore error because the file may be absent.
	_ = godotenv.Load()

	portStr := getEnv("PORT", "8080")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("invalid PORT value: %v", err)
	}

	logLevel := getEnv("LOG_LEVEL", "info")

	return &Config{
		Port:     port,
		LogLevel: logLevel,
	}
}

// getEnv fetches a key from the environment or returns a fallback.
func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}