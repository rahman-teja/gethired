package config

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/cors"
)

type Application struct {
	Port string
	Name string
}

type Database struct {
	Driver      string
	Dsn         string
	DBName      string
	MaxLifetime time.Duration
	MaxOpenConn int
	MaxIdleConn int
}

type Cors struct {
	Options cors.Options
}

type Config struct {
	Application Application
	Database    Database
	Cors        Cors
}

func Init() *Config {
	appConfig := new(Config)

	appConfig.ReadApplicationConfig()
	appConfig.ReadDatabaseConfig()
	appConfig.ReadCORSConfig()

	return appConfig
}

func (c *Config) mustLoadInt(s string, def int) (res int) {
	res, err := strconv.Atoi(s)
	if err != nil {
		return def
	}

	return res
}

func (c *Config) ReadApplicationConfig() {
	c.Application.Port = os.Getenv("APP_PORT")
	c.Application.Name = os.Getenv("APP_NAME")

	if c.Application.Port == "" {
		c.Application.Port = "3030"
	}

	if c.Application.Name == "" {
		c.Application.Name = "Todo Apps"
	}
}

func (c *Config) ReadDatabaseConfig() {
	maxopencon := c.mustLoadInt(os.Getenv("DB_MAX_OPEN_CONN"), 0)
	maxidlecon := c.mustLoadInt(os.Getenv("DB_MAX_IDLE_CONN"), 20)
	maxlifecon := c.mustLoadInt(os.Getenv("DB_MAX_LIFE_CONN"), 3)

	host := os.Getenv("MYSQL_HOST")
	port := os.Getenv("MYSQL_PORT")
	user := os.Getenv("MYSQL_USER")
	pass := os.Getenv("MYSQL_PASSWORD")
	db := os.Getenv("MYSQL_DBNAME")

	if port != "" {
		port = ":" + port
	}

	// "root:@tcp(localhost:3306)/brankas"
	c.Database.Dsn = fmt.Sprintf("%s:%s@tcp(%s%s)/%s", user, pass, host, port, db)
	c.Database.Driver = "mysql"
	c.Database.MaxOpenConn = maxopencon
	c.Database.MaxIdleConn = maxidlecon
	c.Database.MaxLifetime = time.Minute * time.Duration(maxlifecon)
}

func (c *Config) ReadCORSConfig() {
	var headers []string = []string{
		"Accept",
		"Authorization",
		"Content-Type",
		"X-CSRF-Token",
	}

	var origins []string = []string{"*"}
	var methods []string = []string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodDelete,
		http.MethodOptions,
	}

	maxAgeInSeconds := c.mustLoadInt(os.Getenv("CORS_MAX_AGE"), 60)

	c.Cors.Options = cors.Options{
		AllowedOrigins: origins,
		AllowedMethods: methods,
		AllowedHeaders: headers,
		MaxAge:         maxAgeInSeconds,
	}
}
