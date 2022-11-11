package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/joho/godotenv/autoload" // for development
	"github.com/rahman-teja/gethired/internal/config"
	"github.com/sirupsen/logrus"
	"gitlab.com/rteja-library3/rserver"
)

func main() {
	cfg := config.Init()

	logrus.Infof("Service %s run on port %s", cfg.Application.Name, cfg.Application.Port)

	db := InitDatabase(cfg.Database)

	err := InitMigration(db, cfg.Database.DBName)
	if err != nil {
		db.Close()

		log.Fatal(err)
	}

	// put handler last
	handler := InitHttpHandler(db, cfg.Cors.Options)

	// RUN SERVER
	server := rserver.NewServer(
		logrus.New(),
		rserver.
			NewOptions().
			SetHandler(handler).
			SetPort(cfg.Application.Port),
	)
	server.Start()

	csignal := make(chan os.Signal, 1)
	signal.Notify(csignal, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	// WAIT FOR IT
	<-csignal

	db.Close()
	server.Close()

	logrus.Infof("Service %s run on port %s successfully stopped", cfg.Application.Name, cfg.Application.Port)
}
