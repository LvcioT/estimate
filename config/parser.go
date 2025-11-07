package config

import (
	"fmt"
	"os"

	"strconv"

	"github.com/joho/godotenv"
	"github.com/pelletier/go-toml/v2"
)

// Path to the default configuration TOML file
const defaultConfigFilePath = "config/default.toml"

var config AppConfig

func GetConfig() AppConfig {
	return config
}

func LoadConfig() error {
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("cannot load config: '%w'", err)
	}

	return nil
}

func newConfigWithDefaults() (AppConfig, error) {
	var defaultConfig AppConfig

	data, err := os.ReadFile(defaultConfigFilePath)
	if err != nil {
		return AppConfig{}, fmt.Errorf("could not read default config file '%s': %w", defaultConfigFilePath, err)
	}

	err = toml.Unmarshal(data, &defaultConfig)
	if err != nil {
		return AppConfig{}, fmt.Errorf("error unmarshalling default config from '%s': %w", defaultConfigFilePath, err)
	}
	return defaultConfig, nil
}

func updateConfigFromEnv(cfg *AppConfig) error {
	if portStr := os.Getenv("GIN_PORT"); portStr != "" {
		port, err := strconv.Atoi(portStr)
		if err != nil {
			return fmt.Errorf("invalid GIN_PORT value: %w", err)
		}
		cfg.Gin.Port = port
	}

	if debugStr := os.Getenv("GIN_DEBUG"); debugStr != "" {
		debug, err := strconv.ParseBool(debugStr)
		if err != nil {
			return fmt.Errorf("invalid GIN_DEBUG value: %w", err)
		}
		cfg.Gin.Debug = debug
	}

	if file := os.Getenv("SQLITE_FILE"); file != "" {
		cfg.Sqlite.File = file
	}

	return nil
}
