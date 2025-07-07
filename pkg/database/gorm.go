package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type PostgresConfig struct {
	Host     string `mapstructure:"POSTGRESDB_HOST" validate:"required"`
	Port     int    `mapstructure:"POSTGRESDB_PORT" validate:"required"`
	User     string `mapstructure:"POSTGRESDB_USERNAME" validate:"required"`
	Password string `mapstructure:"POSTGRESDB_PASSWORD" validate:"required"`
	DBName   string `mapstructure:"POSTGRESDB_NAME" validate:"required"`
	SSLMode  string `mapstructure:"POSTGRESDB_SSL_MODE" validate:"required"`
	Schema   string `mapstructure:"POSTGRESDB_SCHEMA" validate:"required"`
}

type GormDB struct {
	DB     *gorm.DB
	Config *PostgresConfig
}

func NewGormDB(cfg *PostgresConfig) *GormDB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s search_path=%s",
		cfg.Host,
		cfg.User,
		cfg.Password,
		cfg.DBName,
		cfg.Port,
		cfg.SSLMode,
		cfg.Schema,
	)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // Enable color
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		panic(err)
	}

	return &GormDB{DB: db, Config: cfg}
}

func NewGormDBWithoutDB(cfg *PostgresConfig) *GormDB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s port=%d sslmode=%s search_path=%s",
		cfg.Host,
		cfg.User,
		cfg.Password,
		cfg.Port,
		cfg.SSLMode,
		cfg.Schema,
	)

	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		panic(err)
	}

	return &GormDB{DB: db, Config: cfg}
}

func (s *GormDB) CreateDB() error {
	if tx := s.DB.Exec(fmt.Sprintf("CREATE DATABASE %s;", s.Config.DBName)); tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (s *GormDB) DropDB() error {
	if tx := s.DB.Exec(fmt.Sprintf("DROP DATABASE %s WITH (FORCE);", s.Config.DBName)); tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (s *GormDB) AddExtension() error {
	if tx := s.DB.Exec("CREATE EXTENSION \"uuid-ossp\";"); tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (s *GormDB) Disconnect() error {
	sqlDB, err := s.DB.DB()
	if err != nil {
		return err
	}

	if err := sqlDB.Close(); err != nil {
		return err
	}

	return nil
}
