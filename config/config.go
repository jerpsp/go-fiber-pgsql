package config

import (
	"strings"
	"time"

	"github.com/jerpsp/go-fiber-beginner/pkg/database"
	"github.com/jerpsp/go-fiber-beginner/pkg/storage"
	"github.com/spf13/viper"
)

type (
	Server struct {
		ENV          string        `mapstructure:"ENV" validate:"required"`
		Port         int           `mapstructure:"SERVER_PORT" validate:"required"`
		Timeout      time.Duration `mapstructure:"TIMEOUT" validate:"required"`
		AllowOrigins string        `mapstructure:"ALLOW_ORIGINS" validate:"required"`
	}

	JWT struct {
		Secret          string        `mapstructure:"JWT_SECRET" validate:"required"`
		AccessTokenExp  time.Duration `mapstructure:"JWT_ACCESS_TOKEN_EXP" validate:"required"`
		RefreshTokenExp time.Duration `mapstructure:"JWT_REFRESH_TOKEN_EXP" validate:"required"`
	}

	Config struct {
		Server     *Server                  `mapstructure:"server" validate:"required"`
		PostgresDB *database.PostgresConfig `mapstructure:"postgresdb" validate:"required"`
		JWT        *JWT                     `mapstructure:"jwt" validate:"required"`
		Redis      *database.RedisConfig    `mapstructure:"redis" validate:"required"`
		AWS        *storage.AWSConfig       `mapstructure:"aws" validate:"required"`
	}
)

func InitConfig() *Config {
	var server Server
	var postgresDB database.PostgresConfig
	var jwt JWT
	var redis database.RedisConfig
	var aws storage.AWSConfig

	viper.SetConfigName("dev")
	viper.SetConfigType("env")
	viper.AddConfigPath("environments")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Try to read from config file, but don't panic if not found (for cloud environments)
	if err := viper.ReadInConfig(); err != nil {
		// Only panic if we're in development mode and config file is missing
		if strings.ToLower(viper.GetString("ENV")) == "development" {
			panic(err)
		}
	}

	if err := viper.Unmarshal(&server); err != nil {
		panic(err)
	}

	if err := viper.Unmarshal(&postgresDB); err != nil {
		panic(err)
	}

	if err := viper.Unmarshal(&jwt); err != nil {
		panic(err)
	}

	if err := viper.Unmarshal(&redis); err != nil {
		panic(err)
	}
	if err := viper.Unmarshal(&aws); err != nil {
		panic(err)
	}

	cfg := &Config{
		Server:     &server,
		PostgresDB: &postgresDB,
		JWT:        &jwt,
		Redis:      &redis,
		AWS:        &aws,
	}

	return cfg
}
