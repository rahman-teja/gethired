package main

import (
	"database/sql"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func InitMigration(db *sql.DB, dbname string) error {
	args := os.Args
	isUp := true

	if len(args) > 1 {
		isUp = args[1] == "up"
	}

	log.Println("start migrate")
	defer log.Println("Finish migrate")

	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migration",
		dbname,
		driver,
	)
	if err != nil {
		return err
	}

	if isUp {
		log.Println("Migrate up")
		log.Error(m.Up())
	} else {
		log.Println("Migrate down")
		m.Down()
	}

	return nil
}
