package main

import (
	"database/sql"

	"github.com/rahman-teja/gethired/internal/config"
	log "github.com/sirupsen/logrus"
)

func InitDatabase(cfg config.Database) *sql.DB {
	log.Infof("Try to open connection [%s]:%s", cfg.Driver, cfg.Dsn)
	db, err := sql.Open(cfg.Driver, cfg.Dsn)
	if err != nil {
		log.Fatal("Failed to open a DB connection: ", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping DB: ", err)
	}

	db.SetConnMaxLifetime(cfg.MaxLifetime)
	db.SetMaxOpenConns(cfg.MaxOpenConn)
	db.SetMaxIdleConns(cfg.MaxIdleConn)
	log.Infof("cfg.MaxLifetime %.0f | cfg.MaxOpenConn %d | cfg.MaxIdleConn %d", cfg.MaxLifetime.Minutes(), cfg.MaxOpenConn, cfg.MaxIdleConn)
	log.Infof("Successfully open connection to [%s]:%s", cfg.Driver, cfg.Dsn)

	return db
}
