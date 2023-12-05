package database

import (
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

func NewDatabase() *Database {

	m, err := migrate.New("file://deployment/database/migrations", "postgres://postgres:admin@localhost:5432/postgres?sslmode=disable")

	if err != nil {
		log.Panicf("connect db failed when migrate, err: %v", err)
	}

	if err := m.Up(); err != nil {
		log.Panicf("Migrate failed, err %v", err)
	}
	dsn := "host=localhost user=postgres password=admin dbname=postgres port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil
	}

	return &Database{*db}
}
