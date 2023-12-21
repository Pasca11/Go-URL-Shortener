package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"url-shortener/internal/config"
	"url-shortener/internal/http-server/handlers/redirect"
	"url-shortener/internal/http-server/handlers/url/delete"
	"url-shortener/internal/http-server/handlers/url/save"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/storage/postgres"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	cfg := config.Require()

	log := SetupLogger(cfg.Env)
	log.Info("App is  started", slog.String("env", cfg.Env))

	storage, err := postgres.New(cfg.DB.Host, cfg.DB.Port, cfg.DB.Username, cfg.DB.Password, cfg.DB.DBName)
	if err != nil {
		log.Error("failed to connect DB:", sl.Err(err))
		return
	}

	_ = storage

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Route("/url", func(r chi.Router) {
		r.Use(middleware.BasicAuth("url-shortener", map[string]string{
			cfg.Username: cfg.Password,
		}))

		r.Post("/", save.New(log, storage))
		// TODO add DELETE
		r.Delete("/", delete.New(log, storage))
	})

	router.Get("/{alias}", redirect.New(log, storage))

	log.Info("Server started", slog.String("Adress", cfg.Adress))

	server := http.Server{
		Addr:         cfg.Adress,
		Handler:      router,
		ReadTimeout:  cfg.HttpServer.Timeout,
		WriteTimeout: cfg.HttpServer.Timeout,
		IdleTimeout:  cfg.HttpServer.IdleTimeout,
	}

	defer storage.Close()

	if err := server.ListenAndServe(); err != nil {
		log.Error("cant start server")
		fmt.Println(err)
	}

}

func SetupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case "local":
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case "dev":
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case "prod":
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}
	return log
}
