package database

import (
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct {
	gorm.DB
}

type Config struct {
	Type         *string `mapstructure:"type"`
	Host         *string `mapstructure:"host"`
	Port         *int    `mapstructure:"port"`
	User         *string `mapstructure:"user"`
	Password     *string `mapstructure:"password"`
	DBName       *string `mapstructure:"db_name"`
	SslMode      *string `mapstructure:"ssl_mode"`
	MigratePath  *string `mapstructure:"migrate_path"`
	MaxIdleConns *int    `mapstructure:"max_idle_conns"`
	MaxOpenConns *int    `mapstructure:"max_open_conns"`
}

func (c *Config) GetMigrateURL() string {
	return fmt.Sprintf("file://%s", *c.MigratePath)
}

func (c *Config) GetDatabaseURL() string {
	return fmt.Sprintf("%s://%s:%s@%s:%d/%s?sslmode=%s",
		*c.Type,
		*c.User,
		*c.Password,
		*c.Host,
		*c.Port,
		*c.DBName,
		*c.SslMode)
}

func (c *Config) GetDSN() string {
	if c.Type == nil {
		return ""
	}
	if *c.Type == "postgres" {
		return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
			*c.Host,
			*c.User,
			*c.Password,
			*c.DBName,
			*c.Port,
			*c.SslMode)
	}
	return ""
}

func NewDatabase(config *Config, isLocal bool) *Database {
	m, err := migrate.New(config.GetMigrateURL(), config.GetDatabaseURL())

	if err != nil {
		log.Panic().Err(err).Msgf("Connect Database Failed When Migrate")
	}

	if err := m.Up(); err != nil {
		if err != migrate.ErrNoChange {
			log.Panic().Err(err).Msgf("Migrate Up failed")
		}
	}
	db, err := gorm.Open(postgres.Open(config.GetDSN()), &gorm.Config{
		Logger: New(logger.Config{
			SlowThreshold: 200 * time.Millisecond,
			Colorful:      isLocal,
			LogLevel:      logger.Info,
		}),
	})

	if err != nil {
		log.Panic().Err(err).Msgf("Connect to Database failed")
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(*config.MaxIdleConns)
	sqlDB.SetMaxOpenConns(*config.MaxOpenConns)

	log.Info().Msgf("Connect to Database [%v] Successful!", config.GetDSN())

	return &Database{*db}
}
