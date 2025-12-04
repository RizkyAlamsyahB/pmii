package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	Database DatabaseConfig
	JWT      JWTConfig
	Server   ServerConfig
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret          string
	ExpirationHours int
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port           string
	AllowedOrigins string
	Environment    string
}

// Load loads configuration from .env file using Viper
func Load() (*Config, error) {
	// Set config file
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")

	// Auto bind environment variables
	viper.AutomaticEnv()

	// Read config file (optional - akan fallback ke env vars jika file tidak ada)
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: .env file not found, using environment variables: %v", err)
	}

	// Validate required configs
	requiredKeys := []string{
		"DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME",
		"JWT_SECRET", "JWT_EXPIRATION_HOURS",
		"PORT", "ENV", "ALLOWED_ORIGINS",
	}
	for _, key := range requiredKeys {
		if !viper.IsSet(key) || viper.GetString(key) == "" {
			return nil, fmt.Errorf("required configuration %s is not set", key)
		}
	}

	config := &Config{
		Database: DatabaseConfig{
			Host:     viper.GetString("DB_HOST"),
			Port:     viper.GetString("DB_PORT"),
			User:     viper.GetString("DB_USER"),
			Password: viper.GetString("DB_PASSWORD"),
			DBName:   viper.GetString("DB_NAME"),
		},
		JWT: JWTConfig{
			Secret:          viper.GetString("JWT_SECRET"),
			ExpirationHours: viper.GetInt("JWT_EXPIRATION_HOURS"),
		},
		Server: ServerConfig{
			Port:           viper.GetString("PORT"),
			AllowedOrigins: viper.GetString("ALLOWED_ORIGINS"),
			Environment:    viper.GetString("ENV"),
		},
	}

	log.Println("✅ Configuration loaded successfully")
	return config, nil
}
