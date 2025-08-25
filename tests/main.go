package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/segmentio/kafka-go"
)

func main() {
	byteValue, err := ioutil.ReadFile("../model.json")

	if err != nil {
		log.Fatalf("Ошибка чтения test_order.json: %v", err)
	}

	var data map[string]interface{}

	if err := json.Unmarshal(byteValue, &data); err != nil {
		log.Fatalf("Ошибка парсинга JSON: %v", err)
	}

	orderUID, ok := data["order_uid"].(string)

	if !ok || orderUID == "" {
		log.Fatal("Ключ 'order_uid' не найден или пуст в test_order.json")
	}

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{"localhost:9092"},
		Topic:    "orders",
		Balancer: &kafka.LeastBytes{},
	})

	defer writer.Close()

	fmt.Printf("Отправляем заказ с UID: %s\n", orderUID)

	err = writer.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte(orderUID), 
			Value: byteValue,
		},
	)

	if err != nil {
		log.Fatalf("Ошибка отправки сообщения в Kafka: %v", err)
	}

	fmt.Println("Сообщение успешно отправлено!")
}