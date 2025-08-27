package streaming

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"time"

	"l0-wb-tech/internal/cache"
	"l0-wb-tech/internal/database"
	"l0-wb-tech/internal/models"

	"github.com/segmentio/kafka-go"
	"gorm.io/gorm"
)

type Consumer struct {
	reader *kafka.Reader
	logger *slog.Logger
	db     *database.Storage
	cache  *cache.Cache
}

const maxWaitTime = 10 * time.Second

func NewConsumer(
	brokers []string,
	topic string,
	db *database.Storage,
	cache *cache.Cache,
	logger *slog.Logger,
) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        brokers,
		Topic:          topic,
		GroupID:        "orders-processor-group-final",
		CommitInterval: 0,
		StartOffset:    kafka.FirstOffset,
		MaxWait:        maxWaitTime,
	})
	return &Consumer{
		reader: reader,
		logger: logger,
		db:     db,
		cache:  cache,
	}
}

func (c *Consumer) Run(ctx context.Context) {
	c.logger.Info("Запуск Kafka потребителя", "topic", c.reader.Config().Topic)
	defer func() {
		if err := c.reader.Close(); err != nil {
			c.logger.Error("Ошибка при закрытии Kafka reader", "error", err)
		}
		c.logger.Info("Kafka потребитель остановлен")
	}()

	var jsonBuffer strings.Builder
	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				c.logger.Info("Получен сигнал завершения работы")
				break
			}
			c.logger.Error("Ошибка при чтении сообщения из Kafka", "error", err)
			time.Sleep(1 * time.Second)
			continue
		}

		c.logger.Debug("Получено новое сообщение",
			"offset", msg.Offset,
			"partition", msg.Partition,
			"key", string(msg.Key))

		jsonBuffer.Write(msg.Value)

		var order models.Order
		bufferContent := jsonBuffer.String()

		if err := json.Unmarshal([]byte(bufferContent), &order); err != nil {
			if errors.Is(err, io.ErrUnexpectedEOF) ||
				strings.Contains(err.Error(), "unexpected end of JSON") {
				c.logger.Debug("Накопление JSON, сообщение еще не полное")
				continue
			}
			c.logger.Error("Ошибка парсинга JSON сообщения",
				"error", err,
				"message", bufferContent)

			jsonBuffer.Reset()
			c.commitMessage(ctx, msg)

			continue
		}

		c.logger.Debug("JSON успешно собран из нескольких сообщений")

		if err := c.processOrder(ctx, order); err != nil {
			c.logger.Error("Ошибка обработки заказа",
				"error", err,
				"order_uid", order.OrderUID)
			continue
		}

		jsonBuffer.Reset()

		c.commitMessage(ctx, msg)
	}
}

func (c *Consumer) processOrder(ctx context.Context, order models.Order) error {
	_, err := c.db.GetOrderByUID(ctx, order.OrderUID)

	if err == nil {
		c.logger.Warn("Заказ уже существует в БД", "order_uid", order.OrderUID)
		return nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("ошибка проверки дубликата: %w", err)
	}

	if err := c.db.SaveOrder(ctx, &order); err != nil {
		return fmt.Errorf("ошибка сохранения в БД: %w", err)
	}

	c.cache.Set(order)

	c.logger.Info("Заказ успешно обработан",
		"order_uid", order.OrderUID,
		"track_number", order.TrackNumber)

	return nil
}

func (c *Consumer) commitMessage(ctx context.Context, msg kafka.Message) {
	if err := c.reader.CommitMessages(ctx, msg); err != nil {
		c.logger.Error("Не удалось подтвердить сообщение (сделать commit)", "offset", msg.Offset, "error", err)
	}
}
