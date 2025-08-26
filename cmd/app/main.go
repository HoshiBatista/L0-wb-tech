package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"l0-wb-tech/internal/cache"
	"l0-wb-tech/internal/database"
	"l0-wb-tech/internal/handlers"
	"l0-wb-tech/internal/logger"
	"l0-wb-tech/internal/server"
	"l0-wb-tech/internal/streaming"
)

const (
	shutdownTimeout = time.Duration(5) * time.Second
	httpPort        = "8081"
	dbConnStr       = "postgres://admin:admin@localhost:5432/l0_database?sslmode=disable"
	kafkaTopic      = "orders"
)

func setupApplication(log *slog.Logger) (*server.Server, *streaming.Consumer) {
	db, err := database.New(dbConnStr, log)

	if err != nil {
		log.Error("Не удалось подключиться к БД", slog.Any("error", err))
		os.Exit(1)
	}

	if err := db.Migrate(context.Background()); err != nil {
		log.Error("Не удалось выполнить миграцию", slog.Any("error", err))
		os.Exit(1)
	}

	orderCache := cache.New()
	orders, err := db.GetAllOrders(context.Background())

	if err != nil {
		log.Error("Не удалось загрузить заказы из БД", slog.Any("error", err))
		os.Exit(1)
	}

	orderCache.Load(orders)
	log.Info("Кэш успешно восстановлен", slog.Int("orders_loaded", len(orders)))

	var kafkaBrokers = []string{"localhost:9092"}

	httpHandler := handlers.New(orderCache, db, log)
	httpServer := server.New(httpPort, httpHandler, log)
	kafkaConsumer := streaming.NewConsumer(kafkaBrokers, kafkaTopic, db, orderCache, log)

	return httpServer, kafkaConsumer
}

func main() {
	log := logger.New()
	log.Info("Сервис запускается...")

	httpServer, kafkaConsumer := setupApplication(log)

	mainCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go kafkaConsumer.Run(mainCtx)

	go func() {
		if err := httpServer.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("HTTP-сервер аварийно остановлен", slog.Any("error", err))
			stop()
		}
	}()

	log.Info("Сервис успешно запущен и готов к работе")
	<-mainCtx.Done()

	log.Info("Сервис останавливается...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := httpServer.Stop(shutdownCtx); err != nil {
		log.Error("Ошибка при остановке HTTP-сервера", slog.Any("error", err))
	}

	log.Info("Сервис успешно остановлен")
}
