package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Environment string
	ServerPort  int
	Database    DatabaseConfig
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.SetDefault("environment", "development")
	viper.SetDefault("server_port", 8000)

	if err := viper.ReadInConfig(); err != nil {
		// this is if the config file doesnt exist, we'll rely on to env vars and defaults
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	// override with env vars
	viper.AutomaticEnv()

	// and then unmarshal into Config struct [we cant use them in their format yet]
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to decode into struct: %w", err)
	}

	return &cfg, nil
}
