package config

import (
	"fmt"
	"os"

	"strconv"

	"github.com/joho/godotenv"
	"github.com/pelletier/go-toml/v2"
)

// Path to the default configuration TOML file
const defaultConfigFilePath = "config.toml"

var config Config

func init() {
	err := LoadConfig()
	if err != nil {
		panic(err)
	}
}

func GetConfig() Config {
	return config
}

func LoadConfig() error {
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("cannot load ENV config: '%w'", err)
	}

	config, err = newConfigWithDefaults()
	if err != nil {
		return fmt.Errorf("cannot set config with default values: '%w'", err)
	}

	err = updateConfigFromEnv(&config)
	if err != nil {
		return fmt.Errorf("cannot update config from ENV: '%w'", err)
	}

	return nil
}

func newConfigWithDefaults() (Config, error) {
	var defaultConfig Config

	data, err := os.ReadFile(defaultConfigFilePath)
	if err != nil {
		return Config{}, fmt.Errorf("could not read default config file '%s': %w", defaultConfigFilePath, err)
	}

	err = toml.Unmarshal(data, &defaultConfig)
	if err != nil {
		return Config{}, fmt.Errorf("error unmarshalling default config from '%s': %w", defaultConfigFilePath, err)
	}
	return defaultConfig, nil
}

func updateConfigFromEnv(cfg *Config) error {
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
