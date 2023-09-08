package db

import (
	"database/sql"
	"fmt"
	"log"
	"github.com/golang-migrate/migrate/v4"
    _ "github.com/golang-migrate/migrate/v4/database/postgres"
    _ "github.com/golang-migrate/migrate/v4/source/file"
)

var db *sql.DB

const (
	dbConnectionString = "postgresql://postgres:admin@localhost:5432/goTest"
	migrationsDir = "file://./migrations"
)

func InitDb() (*sql.DB) {
	var err error
	db, err := sql.Open("postgres", dbConnectionString)

	if err != nil {
		log.Fatal(err)
	}


	err = db.Ping()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Db connected successfully")

	if err := applyMigrations(); err != nil {
        log.Fatal(err)
    }

	return db

}

func applyMigrations() error {
    m, err := migrate.New(migrationsDir, dbConnectionString)
    if err != nil {
        return err
    }

    if err := m.Up(); err != nil && err != migrate.ErrNoChange {
        return err
    }

    return nil
}