package main

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/rahman-teja/gethired/internal/httphandler"
	log "github.com/sirupsen/logrus"
	"gitlab.com/rteja-library3/rdecoder"
)

func InitHttpHandler(db *sql.DB, opt cors.Options) http.Handler {
	log.Info("Initiate http handler")
	defer log.Info("Http handler ready")

	httpProperty := httphandler.HTTPHandlerProperty{
		DB:             db,
		DefaultDecoder: rdecoder.NewJSONEncoder(),
	}

	r := chi.NewRouter()

	// A good base middleware stack
	// r.Use(middleware.RequestID)
	// r.Use(middleware.RealIP)
	// r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.New(opt).Handler)
	r.Use(middleware.Compress(5, "gzip"))

	// Health Checking...
	r.Get("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("SUCCESS"))
	})

	r.Mount("/activity-groups", httphandler.NewActivityHttpHandler(httpProperty))
	r.Mount("/todo-items", httphandler.NewToDoHttpHandler(httpProperty))

	chi.Walk(r, func(method, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		route = strings.Replace(route, "/*/", "/", -1)
		log.Printf("[%s] %s", method, route)
		return nil
	})

	return r
}
