package config

import (
	"errors"
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Port        int    `mapstructure:"PORT"`
	DatabaseURL string `mapstructure:"DATABASE_URL"`
	// Storage              StorageConfig
}

// type StorageConfig struct {
// 	LogoPath string `mapstructure:"STORAGE_LOGO_PATH"`
// 	BaseURL  string `mapstructure:"STORAGE_BASE_URL"`
// }

func LoadConfig() (*Config, error) {
	viper.AutomaticEnv()

	// Explicitly bind environment variables to config keys
	keys := []string{
		"PORT",
		"DATABASE_URL",
		// "STORAGE_LOGO_PATH",
		// "STORAGE_BASE_URL",
	}

	for _, key := range keys {
		if err := viper.BindEnv(key); err != nil {
			return nil, fmt.Errorf("failed to bind environment variable %s: %w", key, err)
		}
	}

	// Debug: Print all environment variables
	// for _, e := range os.Environ() {
	// pair := strings.SplitN(e, "=", 2)
	// log.Printf("ENV: %s=%s\n", pair[0], pair[1])
	// }

	// Optionally read from a config file if specified by an environment variable
	configFile := viper.GetString("CONFIG_FILE")
	if configFile != "" {
		viper.SetConfigFile(configFile)
		if err := viper.ReadInConfig(); err != nil {
			return nil, err
		}
		log.Println("Using config file:", configFile)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	// Load from .env local
	if cfg == (Config{}) {
		viper.AutomaticEnv()
		viper.SetConfigFile(".env")

		if err := viper.ReadInConfig(); err != nil {
			return nil, err
		}

		if err := viper.Unmarshal(&cfg); err != nil {
			return nil, err
		}
	}

	// Debug: Print the loaded configuration
	// log.Printf("Loaded Config: %+v\n", cfg)

	// if cfg.Storage.LogoPath == "" {
	// 	cfg.Storage.LogoPath = "/root/storage/logo" // Default path
	// }
	// if cfg.Storage.BaseURL == "" {
	// 	cfg.Storage.BaseURL = "http://localhost:8080/public" // Default URL
	// }

	// Validate required config values
	if err := validateConfig(cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func validateConfig(cfg Config) error {
	missingEnvVars := []string{}

	if cfg.DatabaseURL == "" {
		missingEnvVars = append(missingEnvVars, "DATABASE_URL")
	}
	// Validate storage configuration
	// if cfg.Storage.LogoPath == "" {
	// 	missingEnvVars = append(missingEnvVars, "STORAGE_LOGO_PATH")
	// }
	// Ensure the logo directory exists and is writable
	// if err := os.MkdirAll(cfg.Storage.LogoPath, 0755); err != nil {
	// 	return fmt.Errorf("failed to create logo directory: %w", err)
	// }
	if len(missingEnvVars) > 0 {
		return errors.New(fmt.Sprintf("Missing required environment variables: %v", missingEnvVars))
	}

	return nil
}
