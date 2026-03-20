package config

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type JWTConfig struct {
	Secret          string        `mapstructure:"secret"`
	AccessTokenExp  time.Duration `mapstructure:"access_token_exp"`
	RefreshTokenExp time.Duration `mapstructure:"refresh_token_exp"`
}

type Config struct {
	Environment string         `mapstructure:"environment"`
	ServerPort  int            `mapstructure:"server_port"`
	Database    DatabaseConfig `mapstructure:"database"`
	JWT         JWTConfig
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

	v.AddConfigPath(".")
	v.AddConfigPath("./internal/config")

	// defaults
	v.SetDefault("environment", "development")
	v.SetDefault("server_port", 8000)

	/// read base config
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("Warning: base config file not found, using defaults and environment variables")
		} else {
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

// TODO: update validate to check jwt secret
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

// OpenDB creates a gorm db conn using the config
func (c *Config) OpenDB() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", c.Database.Host, c.Database.Port, c.Database.User, c.Database.Password, c.Database.DBName, c.Database.SSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		// i'll add gorm loggin here
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %w", err)
	}

	// get the sql.DB obj to config the conn pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(25)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	return db, nil
}
