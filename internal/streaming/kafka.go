package streaming

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"

	"l0-wb-tech/internal/cache"
	"l0-wb-tech/internal/database"
	"l0-wb-tech/internal/models"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader *kafka.Reader
	logger *slog.Logger
	db     *database.Storage
	cache  *cache.Cache
}

func NewConsumer(
	brokers []string,
	topic string,
	db *database.Storage,
	cache *cache.Cache,
	logger *slog.Logger,
) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		GroupID: "orders-group",
	})
	return &Consumer{
		reader: reader,
		logger: logger,
		db:     db,
		cache:  cache,
	}
}

func (c *Consumer) Run(ctx context.Context) {
	c.logger.Info("Kafka consumer запущен, слушает топик", slog.String("topic", c.reader.Config().Topic))

	for {
		msg, err := c.reader.FetchMessage(ctx)

		if err != nil {
			if errors.Is(err, context.Canceled) {
				break
			}
			c.logger.Error("Ошибка при получении сообщения из Kafka", slog.Any("error", err))
			continue
		}

		c.handleMessage(ctx, msg.Value)

		if err := c.reader.CommitMessages(ctx, msg); err != nil {
			c.logger.Error("Не удалось подтвердить сообщение в Kafka", slog.Any("error", err))
		}
	}

	c.logger.Info("Kafka consumer останавливается...")

	if err := c.reader.Close(); err != nil {
		c.logger.Error("Ошибка при закрытии Kafka reader", slog.Any("error", err))
	}
}

func (c *Consumer) handleMessage(ctx context.Context, msgValue []byte) {
	var order models.Order

	//nolint:musttag
	if err := json.Unmarshal(msgValue, &order); err != nil {
		c.logger.Error(
			"Ошибка парсинга JSON из Kafka, сообщение проигнорировано",
			slog.Any("error", err),
			slog.String("raw_message", string(msgValue)),
		)
		return
	}

	// TODO validation

	if err := c.db.SaveOrder(ctx, &order); err != nil {
		c.logger.Error(
			"Ошибка сохранения заказа в БД",
			slog.Any("error", err),
			slog.String("order_uid", order.OrderUID),
		)
		return
	}

	c.cache.Set(order)
	c.logger.Info("Заказ успешно обработан и сохранен", slog.String("order_uid", order.OrderUID))
}
