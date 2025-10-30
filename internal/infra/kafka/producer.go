package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"vado_server/internal/config/port"
	"vado_server/internal/domain/chat"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type Producer struct {
	writer *kafka.Writer
	log    *zap.SugaredLogger
}

func NewProducer(topic string, log *zap.SugaredLogger) *Producer {
	brokers := []string{"localhost:" + port.Kafka}
	writer := &kafka.Writer{
		Addr:                   kafka.TCP(brokers...),
		Topic:                  topic,
		Balancer:               &kafka.LeastBytes{}, // Балансировщик для распределения сообщений по партициям (можно использовать другие: Hash, RoundRobin)
		AllowAutoTopicCreation: true,                // Авто создание топика
	}

	return &Producer{
		writer: writer,
		log:    log,
	}
}

// SendMessage — безопасная отправка сообщения с ретраями и логами
func (p *Producer) SendMessage(key, value []byte) error {
	msg := kafka.Message{
		Key:   key,
		Value: value,
		Time:  time.Now(),
	}

	err := p.writer.WriteMessages(context.Background(), msg)
	if err != nil {
		p.log.Errorw("Kafka write error", "key", string(key), "err", err)
		return err
	}

	p.log.Debugw("Kafka message sent", "key", string(key), "size", len(value))
	return nil
}

func (p *Producer) SendChatMessage(msg *chat.MessageLog) error {
	data, err := json.Marshal(msg)
	if err != nil {
		p.log.Errorw("Failed to marshal message", "err", err)
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(fmt.Sprintf("user_%d", msg.UserID)),
		Value: data,
		Time:  msg.Timestamp,
	})

	if err != nil {
		p.log.Errorw("Kafka write error", "err", err)
		return err
	}
	p.log.Debugw("Kafka message sent", "user_id", msg.UserID, "size", len(data))
	return nil
}

// Close — корректно закрывает соединение
func (p *Producer) Close() error {
	p.log.Info("Closing Kafka producer...")
	return p.writer.Close()
}
