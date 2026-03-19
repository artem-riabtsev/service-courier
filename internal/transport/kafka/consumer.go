package kafka

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"service-courier/internal/event"

	"github.com/IBM/sarama"
)

type Consumer struct {
	consumer      sarama.ConsumerGroup
	processor     *event.EventProcessor
	topic         string
	consumerGroup string
	wg            sync.WaitGroup
}

func NewConsumer(brokers []string, topic, consumerGroup string, processor *event.EventProcessor) (*Consumer, error) {
	config := sarama.NewConfig()
	config.Version = sarama.V2_8_0_0
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	config.Consumer.Return.Errors = true
	config.Consumer.MaxProcessingTime = 10 * time.Second

	consumer, err := sarama.NewConsumerGroup(brokers, consumerGroup, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer group: %w", err)
	}

	return &Consumer{
		consumer:      consumer,
		processor:     processor,
		topic:         topic,
		consumerGroup: consumerGroup,
	}, nil
}

func (c *Consumer) Start(ctx context.Context) {
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()

		for {
			select {
			case <-ctx.Done():
				slog.Info("Kafka consumer stopping")
				return
			default:
				handler := &consumerGroupHandler{
					processor: c.processor,
					ready:     make(chan bool),
				}

				err := c.consumer.Consume(ctx, []string{c.topic}, handler)
				if err != nil {
					slog.Error("Error consuming from Kafka", "error", err)
					time.Sleep(5 * time.Second)
				}

				if ctx.Err() != nil {
					return
				}

				handler.ready = make(chan bool)
			}
		}
	}()

	slog.Info("Kafka consumer started", "topic", c.topic, "group", c.consumerGroup)
}

func (c *Consumer) Stop() error {
	c.wg.Wait()

	if err := c.consumer.Close(); err != nil {
		return fmt.Errorf("failed to close consumer: %w", err)
	}

	slog.Info("Kafka consumer stopped")
	return nil
}

type consumerGroupHandler struct {
	processor *event.EventProcessor
	ready     chan bool
}

func (h *consumerGroupHandler) Setup(session sarama.ConsumerGroupSession) error {
	close(h.ready)
	slog.Debug("Kafka consumer setup", "member_id", session.MemberID())
	return nil
}

func (h *consumerGroupHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	slog.Debug("Kafka consumer cleanup", "member_id", session.MemberID())
	return nil
}

func (h *consumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				slog.Debug("Message channel closed", "partition", claim.Partition())
				return nil
			}

			slog.Debug("Received message", "topic", message.Topic,
				"partition", message.Partition, "offset", message.Offset)

			ctx := context.Background()
			if err := h.processor.ProcessMessage(ctx, message.Value); err != nil {
				slog.Error("Failed to process message", "error", err)
				continue
			}

			session.MarkMessage(message, "")

		case <-session.Context().Done():
			return nil
		}
	}
}
