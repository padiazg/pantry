package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// Config holds all application configuration.
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	LogLevel  string
	LogFormat string
}

// ServerConfig holds HTTP server configuration.
type ServerConfig struct {
	Port            int
	Host            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration
}

// DatabaseConfig holds database connection configuration.
type DatabaseConfig struct {
	URL string
}

// Load reads configuration from file and environment variables.
func Load(cfgFile string) (*Config, error) {
	setDefaults()

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName(".pantry")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		viper.AddConfigPath("$HOME")
	}

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	viper.SetEnvPrefix("PANTRY")
	viper.AutomaticEnv()

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	return &cfg, nil
}

func setDefaults() {
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.readtimeout", 15*time.Second)
	viper.SetDefault("server.writetimeout", 15*time.Second)
	viper.SetDefault("server.idletimeout", 60*time.Second)
	viper.SetDefault("server.shutdowntimeout", 30*time.Second)

	viper.SetDefault("database.url", "postgres://user:pass@localhost:5432/pantry?sslmode=disable")

	viper.SetDefault("loglevel", "info")
	viper.SetDefault("logformat", "json")
}
