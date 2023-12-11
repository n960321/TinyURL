package database

import (
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	gorm.DB
}

type Config struct {
	Type        *string `mapstructure:"type"`
	Host        *string `mapstructure:"host"`
	Port        *int    `mapstructure:"port"`
	User        *string `mapstructure:"user"`
	Password    *string `mapstructure:"password"`
	DBName      *string `mapstructure:"db_name"`
	SslMode     *string `mapstructure:"ssl_mode"`
	MigratePath *string `mapstructure:"migrate_path"`
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

func NewDatabase(config *Config) *Database {
	m, err := migrate.New(config.GetMigrateURL(), config.GetDatabaseURL())

	if err != nil {
		log.Panicf("connect db failed when migrate, err: %v", err)
	}

	if err := m.Up(); err != nil {
		if err != migrate.ErrNoChange {
			log.Panicf("Migrate failed, err %v", err)
		}
	}
	db, err := gorm.Open(postgres.Open(config.GetDSN()), &gorm.Config{})

	if err != nil {
		return nil
	}

	return &Database{*db}
}
