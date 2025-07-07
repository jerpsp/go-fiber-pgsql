package config

import (
	"os"
	"strings"
	"time"

	"github.com/jerpsp/go-fiber-beginner/pkg/database"
	"github.com/spf13/viper"
)

type (
	Server struct {
		ENV          string        `mapstructure:"ENV" validate:"required"`
		Port         int           `mapstructure:"SERVER_PORT" validate:"required"`
		Timeout      time.Duration `mapstructure:"TIMEOUT" validate:"required"`
		AllowOrigins string        `mapstructure:"ALLOW_ORIGINS" validate:"required"`
	}

	Config struct {
		Server     *Server                  `mapstructure:"server" validate:"required"`
		PostgresDB *database.PostgresConfig `mapstructure:"postgresdb" validate:"required"`
	}
)

func InitConfig() *Config {
	var server Server
	var postgresDB database.PostgresConfig

	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "dev"
	}

	viper.SetConfigName(env)
	viper.SetConfigType("env")
	viper.AddConfigPath("environments")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	if err := viper.Unmarshal(&server); err != nil {
		panic(err)
	}

	if err := viper.Unmarshal(&postgresDB); err != nil {
		panic(err)
	}

	cfg := &Config{
		Server:     &server,
		PostgresDB: &postgresDB,
	}

	return cfg
}
