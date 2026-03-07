package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Environment string         `mapstructure:"environment"`
	ServerPort  int            `mapstructure:"server_port"`
	Database    DatabaseConfig `mapstructure:"database"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
}

func Load() (*Config, error) {
	v := viper.New()

	// allow env to override config keys like host, port ...
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("./config")

	// defaults
	v.SetDefault("environment", "development")
	v.SetDefault("server_port", 8000)

	// read base config
	if err := v.ReadInConfig(); err != nil {
		// this is if the config file doesnt exist, we'll rely on to env vars and defaults
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	// merge with environment specific config
	env := v.GetString("environment")
	if env == "" {
		env = "development"
	}
	v.SetConfigName("config." + env)
	if err := v.MergeInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading env config file: %w", err)
		}
	}

	// and then unmarshal into Config struct [we cant use them in their format yet]
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to decode into struct: %w", err)

	}

	return &cfg, nil
}

func (c *Config) Validate() error {
	if c.ServerPort <= 0 || c.ServerPort > 65535 {
		return fmt.Errorf("invalid server port: %d", c.ServerPort)
	}

	if c.Database.Host == "" {
		return fmt.Errorf("db host is required")
	}

	if c.Database.Port <= 0 || c.Database.Port > 65535 {
		return fmt.Errorf("invalid db port: %d", c.Database.Port)
	}

	if c.Database.User == "" {
		return fmt.Errorf("db user is required")
	}

	if c.Database.DBName == "" {
		return fmt.Errorf("db name is required")
	}

	validEnvs := map[string]bool{"development": true, "staging": true, "production": true}
	if !validEnvs[c.Environment] {
		return fmt.Errorf("invalid environment: %s (must be development, staging or production)", c.Environment)
	}

	return nil
}
