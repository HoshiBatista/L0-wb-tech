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
)

func main() {
	log := logger.New()
	log.Info("Сервис запускается...")

	const connStr string = "postgres://myuser:mypassword@localhost:5432/orders_db?sslmode=disable"

	db, err := database.New(connStr, log)

	if err != nil {
		log.Error("Не удалось подключиться к базе данных", slog.Any("error", err))
		os.Exit(1)
	}

	if err := db.Migrate(context.Background()); err != nil {
		log.Error("Не удалось выполнить миграцию", slog.Any("error", err))
		os.Exit(1)
	}

	orderCache := cache.New()

	log.Info("Восстановление кэша из базы данных...")

	orders, err := db.GetAllOrders(context.Background())

	if err != nil {
		log.Error("Не удалось загрузить заказы из БД для восстановления кэша", slog.Any("error", err))
		os.Exit(1)
	}

	orderCache.Load(orders)

	log.Info("Кэш успешно восстановлен", slog.Int("загружено_заказов", len(orders)))

	httpHandler := handlers.New(orderCache, log)
	httpServer := server.New("8081", httpHandler, log)

	go func() {
		if err := httpServer.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("HTTP-сервер аварийно остановлен", slog.Any("error", err))
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Сервис останавливается...")

	const num_seconds int8 = 5

	shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Duration(num_seconds)*time.Second)

	defer cancel()

	if err := httpServer.Stop(shutdownCtx); err != nil {
		log.Error("Ошибка при остановке HTTP-сервера", slog.Any("error", err))
	}

	log.Info("Сервис успешно остановлен")
}
