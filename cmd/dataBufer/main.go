package main

import (
	"KPI_Drive_test/internal/config"
	"KPI_Drive_test/internal/http-server/handlers/bufer"
	"KPI_Drive_test/internal/http-server/middleware/logger"
	"KPI_Drive_test/internal/stan"
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	ostan "github.com/nats-io/stan.go"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	envDev     = "dev"                // уровень логирования
	envProd    = "prod"               // уровень логирования
	configPath = "config/config.yaml" // путь до конфиг файла
)

func main() {
	// инициализация конфига
	cfg := config.MustLoad(configPath)

	// инициализация логгера
	log := SetupLogger(cfg.Env)

	log.Info("App started", slog.String("env", cfg.Env))
	log.Debug("Debugging started")

	// создаем клиент NATS
	client, err := stan.NewClient(
		cfg.Nats.ClusterId,
		cfg.Nats.ClientId,
		cfg.Nats.Url,
	)
	if err != nil {
		log.Error("failed to init connection to nats: ", err)
		os.Exit(1)
	}
	defer client.Close()

	// подписка на канал
	go func() {
		_, err = client.Subscribe("facts", func(m *ostan.Msg) {
			log.Info("New order to NATS!")
			stan.FactMessage(log, m)
		})
		if err != nil {
			log.Error("failed to subscribe to nats channel: ", err)
			os.Exit(1)
		}
	}()

	log.Info("nats successfully connected and listening")

	// инициализация роутера и настройка мидлвейров
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(logger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	// REST маршруты
	router.Post("/api/fact", bufer.GetFact(log, client))

	log.Info("starting server", slog.String("address", cfg.Address))

	// настройка и запуск сервера
	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("failed to start server: ", err)
		}
	}()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("shutdown error: ", err)
	} else {
		log.Info("server stopped gracefully")
	}
}

func SetupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)

	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
