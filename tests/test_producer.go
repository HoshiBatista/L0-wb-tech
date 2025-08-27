package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
)

func main() {
	content, err := os.ReadFile("model.json")

	if err != nil {
		log.Fatal("Ошибка чтения файла:", err)
	}

	var order map[string]interface{}

	if err := json.Unmarshal(content, &order); err != nil {
		log.Fatal("Ошибка парсинга JSON:", err)
	}

	const num = 10

	w := &kafka.Writer{
		Addr:         kafka.TCP("localhost:9092"),
		Topic:        "orders",
		BatchTimeout: time.Duration(num) * time.Millisecond,
	}

	err = w.WriteMessages(context.Background(),
		kafka.Message{
			Value: content,
		},
	)

	if err != nil {
		log.Fatal("Ошибка отправки сообщения:", err)
	}

	fmt.Println("Сообщение отправлено в Kafka")

	if err := w.Close(); err != nil {
		log.Fatal("Ошибка при закрытии writer:", err)
	}
}
