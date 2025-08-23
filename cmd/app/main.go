package main

import (
	"fmt"
	"l0-wb-tech/internal/logger"
)

func main() {
	log := logger.New()
	log.Info("Запуск сервиса...")
	log.Debug("Режим отладки включен")
	fmt.Println("Start project")
	log.Info("Сервис успешно запущен и готов к работе")
}