package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

var db *sql.DB

const (
	migrationsDir = "file://./migrations"
)

func InitDb() (*sql.DB) {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading env file")
	}
	
	var err error
	db, err := sql.Open("postgres", os.Getenv("DB_URL"))

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
    m, err := migrate.New(migrationsDir, os.Getenv("DB_URL"))
    if err != nil {
        return err
    }

    if err := m.Up(); err != nil && err != migrate.ErrNoChange {
        return err
    }

    return nil
}