package kafka

import (
	"context"
	"fmt"
	"time"
	"vado_server/internal/app"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type Consumer struct {
	reader *kafka.Reader
	log    *zap.SugaredLogger
}

func NewConsumer(ctx *app.Context) *Consumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{fmt.Sprintf("localhost:%s", ctx.Cfg.KafkaPort)},
		Topic:          "chat",
		GroupID:        "chat-group",
		StartOffset:    kafka.FirstOffset,
		MinBytes:       1, // Читаем сразу, даже если сообщение одно
		MaxBytes:       10e6,
		CommitInterval: 0, // Вручную коммитим
	})

	return &Consumer{reader: r, log: ctx.Log}
}

func (c *Consumer) Run(ctx context.Context, handle func(key, value []byte) error) error {
	c.log.Infow("Kafka consumer started",
		"topic", c.reader.Config().Topic,
		"group", c.reader.Config().GroupID)

	for {
		m, err := c.reader.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				c.log.Infow("Kafka consumer stopped")
				return nil
			}
			c.log.Errorw("Kafka read error", "err", err)
			time.Sleep(2 * time.Second)
			continue
		}

		if errHandle := handle(m.Key, m.Value); errHandle != nil {
			c.log.Errorw("Handler error", "offset", m.Offset, "err", err)
			continue
		}

		c.log.Debugw("Processed message",
			"partition", m.Partition,
			"offset", m.Offset,
			"key", string(m.Key),
			"value", string(m.Value))

		// Ручной коммит после успешной обработки
		if err := c.reader.CommitMessages(ctx, m); err != nil {
			c.log.Errorw("Commit error", "err", err)
		} else {
			c.log.Debugw("Offset committed", "offset", m.Offset)
		}
	}
}

func (c *Consumer) Close() error {
	c.log.Info("Kafka consumer closing...")
	return c.reader.Close()
}
