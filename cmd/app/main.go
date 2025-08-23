package main

import (
	"context"
	"log/slog"
	"os"

	"l0-wb-tech/internal/database"
	"l0-wb-tech/internal/logger"
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

	log.Info("Сервис успешно запущен и готов к работе")

	select {}
}
