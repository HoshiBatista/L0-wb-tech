package database

import (
	"context"
	"fmt"
	"log/slog"

	"l0-wb-tech/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Storage struct {
	db     *gorm.DB
	logger *slog.Logger
}

func New(connStr string, logger *slog.Logger) (*Storage, error) {
	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})

	if err != nil {
		return nil, fmt.Errorf("не удалось подключиться к базе данных: %w", err)
	}

	return &Storage{db: db, logger: logger}, nil
}

func (s *Storage) Migrate(ctx context.Context) error {
	s.logger.Info("Запуск миграции базы данных...")

	err := s.db.WithContext(ctx).AutoMigrate(&models.Order{}, &models.Delivery{}, &models.Payment{}, &models.Item{})

	if err != nil {
		s.logger.Error("Ошибка миграции базы данных", slog.Any("error", err))
		return err
	}

	s.logger.Info("Миграция базы данных успешно завершена")

	return nil
}

func (s *Storage) SaveOrder(ctx context.Context, order *models.Order) error {
	result := s.db.WithContext(ctx).Create(order)

	if result.Error != nil {
		return fmt.Errorf("ошибка при сохранении заказа: %w", result.Error)
	}

	return nil
}

func (s *Storage) GetAllOrders(ctx context.Context) ([]models.Order, error) {
	var orders []models.Order
	result := s.db.WithContext(ctx).Preload("Delivery").Preload("Payment").Preload("Items").Find(&orders)

	if result.Error != nil {
		return nil, fmt.Errorf("ошибка при загрузке заказов из БД: %w", result.Error)
	}

	return orders, nil
}
