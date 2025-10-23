package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
	"obs-tools-usage/kafka/events"
)

// PaymentEventHandler handles payment events
type PaymentEventHandler interface {
	HandlePaymentCompleted(ctx context.Context, event *events.PaymentCompletedEvent) error
	HandlePaymentFailed(ctx context.Context, event *events.PaymentFailedEvent) error
	HandlePaymentRefunded(ctx context.Context, event *events.PaymentRefundedEvent) error
	HandleStockUpdate(ctx context.Context, event *events.StockUpdateEvent) error
	HandleBasketCleared(ctx context.Context, event *events.BasketClearedEvent) error
}

// PaymentConsumer handles consuming payment events from Kafka
type PaymentConsumer struct {
	consumerGroup sarama.ConsumerGroup
	handler       PaymentEventHandler
	logger        *logrus.Logger
	topics        []string
}

// NewPaymentConsumer creates a new payment consumer
func NewPaymentConsumer(
	brokers []string,
	groupID string,
	handler PaymentEventHandler,
	logger *logrus.Logger,
) (*PaymentConsumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Consumer.Group.Session.Timeout = 10 * time.Second
	config.Consumer.Group.Heartbeat.Interval = 3 * time.Second

	consumerGroup, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer group: %w", err)
	}

	return &PaymentConsumer{
		consumerGroup: consumerGroup,
		handler:       handler,
		logger:        logger,
		topics: []string{
			events.PaymentEventsTopic,
			events.StockEventsTopic,
			events.BasketEventsTopic,
		},
	}, nil
}

// Start starts consuming messages
func (c *PaymentConsumer) Start(ctx context.Context) error {
	c.logger.Info("Starting payment consumer...")

	for {
		select {
		case <-ctx.Done():
			c.logger.Info("Payment consumer stopped")
			return ctx.Err()
		default:
			err := c.consumerGroup.Consume(ctx, c.topics, c)
			if err != nil {
				c.logger.WithError(err).Error("Error consuming messages")
				time.Sleep(5 * time.Second)
			}
		}
	}
}

// Stop stops the consumer
func (c *PaymentConsumer) Stop() error {
	return c.consumerGroup.Close()
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (c *PaymentConsumer) Setup(sarama.ConsumerGroupSession) error {
	c.logger.Info("Consumer group session started")
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (c *PaymentConsumer) Cleanup(sarama.ConsumerGroupSession) error {
	c.logger.Info("Consumer group session ended")
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages()
func (c *PaymentConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message := <-claim.Messages():
			if message == nil {
				return nil
			}

			c.logger.WithFields(logrus.Fields{
				"topic":     message.Topic,
				"partition": message.Partition,
				"offset":    message.Offset,
				"key":       string(message.Key),
			}).Debug("Received message")

			if err := c.processMessage(context.Background(), message); err != nil {
				c.logger.WithError(err).WithFields(logrus.Fields{
					"topic":     message.Topic,
					"partition": message.Partition,
					"offset":    message.Offset,
				}).Error("Failed to process message")
			}

			session.MarkMessage(message, "")

		case <-session.Context().Done():
			return nil
		}
	}
}

// processMessage processes a single message
func (c *PaymentConsumer) processMessage(ctx context.Context, message *sarama.ConsumerMessage) error {
	// Get event type from headers
	var eventType string
	for _, header := range message.Headers {
		if string(header.Key) == "event_type" {
			eventType = string(header.Value)
			break
		}
	}

	if eventType == "" {
		return fmt.Errorf("event type not found in message headers")
	}

	switch eventType {
	case events.PaymentCompletedEventType:
		var event events.PaymentCompletedEvent
		if err := json.Unmarshal(message.Value, &event); err != nil {
			return fmt.Errorf("failed to unmarshal payment completed event: %w", err)
		}
		return c.handler.HandlePaymentCompleted(ctx, &event)

	case events.PaymentFailedEventType:
		var event events.PaymentFailedEvent
		if err := json.Unmarshal(message.Value, &event); err != nil {
			return fmt.Errorf("failed to unmarshal payment failed event: %w", err)
		}
		return c.handler.HandlePaymentFailed(ctx, &event)

	case events.PaymentRefundedEventType:
		var event events.PaymentRefundedEvent
		if err := json.Unmarshal(message.Value, &event); err != nil {
			return fmt.Errorf("failed to unmarshal payment refunded event: %w", err)
		}
		return c.handler.HandlePaymentRefunded(ctx, &event)

	case events.StockUpdateEventType:
		var event events.StockUpdateEvent
		if err := json.Unmarshal(message.Value, &event); err != nil {
			return fmt.Errorf("failed to unmarshal stock update event: %w", err)
		}
		return c.handler.HandleStockUpdate(ctx, &event)

	case events.BasketClearedEventType:
		var event events.BasketClearedEvent
		if err := json.Unmarshal(message.Value, &event); err != nil {
			return fmt.Errorf("failed to unmarshal basket cleared event: %w", err)
		}
		return c.handler.HandleBasketCleared(ctx, &event)

	default:
		c.logger.WithField("event_type", eventType).Warn("Unknown event type")
		return nil
	}
}
