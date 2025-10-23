package publisher

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"obs-tools-usage/kafka/events"
)

// PaymentPublisher handles publishing payment events to Kafka
type PaymentPublisher struct {
	producer sarama.SyncProducer
	logger   *logrus.Logger
}

// NewPaymentPublisher creates a new payment publisher
func NewPaymentPublisher(brokers []string, logger *logrus.Logger) (*PaymentPublisher, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true
	config.Producer.Compression = sarama.CompressionSnappy

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka producer: %w", err)
	}

	return &PaymentPublisher{
		producer: producer,
		logger:   logger,
	}, nil
}

// PublishPaymentCompleted publishes a payment completed event
func (p *PaymentPublisher) PublishPaymentCompleted(ctx context.Context, event *events.PaymentCompletedEvent) error {
	event.EventID = uuid.New().String()
	event.EventType = events.PaymentCompletedEventType
	event.Timestamp = time.Now()

	message, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal payment completed event: %w", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: events.PaymentEventsTopic,
		Key:   sarama.StringEncoder(event.PaymentID),
		Value: sarama.ByteEncoder(message),
		Headers: []sarama.RecordHeader{
			{Key: []byte("event_type"), Value: []byte(event.EventType)},
			{Key: []byte("payment_id"), Value: []byte(event.PaymentID)},
			{Key: []byte("user_id"), Value: []byte(event.UserID)},
		},
	}

	partition, offset, err := p.producer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("failed to send payment completed event: %w", err)
	}

	p.logger.WithFields(logrus.Fields{
		"event_id":  event.EventID,
		"payment_id": event.PaymentID,
		"user_id":   event.UserID,
		"topic":     events.PaymentEventsTopic,
		"partition": partition,
		"offset":    offset,
	}).Info("Payment completed event published")

	return nil
}

// PublishPaymentFailed publishes a payment failed event
func (p *PaymentPublisher) PublishPaymentFailed(ctx context.Context, event *events.PaymentFailedEvent) error {
	event.EventID = uuid.New().String()
	event.EventType = events.PaymentFailedEventType
	event.Timestamp = time.Now()

	message, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal payment failed event: %w", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: events.PaymentEventsTopic,
		Key:   sarama.StringEncoder(event.PaymentID),
		Value: sarama.ByteEncoder(message),
		Headers: []sarama.RecordHeader{
			{Key: []byte("event_type"), Value: []byte(event.EventType)},
			{Key: []byte("payment_id"), Value: []byte(event.PaymentID)},
			{Key: []byte("user_id"), Value: []byte(event.UserID)},
		},
	}

	partition, offset, err := p.producer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("failed to send payment failed event: %w", err)
	}

	p.logger.WithFields(logrus.Fields{
		"event_id":  event.EventID,
		"payment_id": event.PaymentID,
		"user_id":   event.UserID,
		"topic":     events.PaymentEventsTopic,
		"partition": partition,
		"offset":    offset,
	}).Info("Payment failed event published")

	return nil
}

// PublishPaymentRefunded publishes a payment refunded event
func (p *PaymentPublisher) PublishPaymentRefunded(ctx context.Context, event *events.PaymentRefundedEvent) error {
	event.EventID = uuid.New().String()
	event.EventType = events.PaymentRefundedEventType
	event.Timestamp = time.Now()

	message, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal payment refunded event: %w", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: events.PaymentEventsTopic,
		Key:   sarama.StringEncoder(event.PaymentID),
		Value: sarama.ByteEncoder(message),
		Headers: []sarama.RecordHeader{
			{Key: []byte("event_type"), Value: []byte(event.EventType)},
			{Key: []byte("payment_id"), Value: []byte(event.PaymentID)},
			{Key: []byte("user_id"), Value: []byte(event.UserID)},
		},
	}

	partition, offset, err := p.producer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("failed to send payment refunded event: %w", err)
	}

	p.logger.WithFields(logrus.Fields{
		"event_id":  event.EventID,
		"payment_id": event.PaymentID,
		"user_id":   event.UserID,
		"topic":     events.PaymentEventsTopic,
		"partition": partition,
		"offset":    offset,
	}).Info("Payment refunded event published")

	return nil
}

// PublishStockUpdate publishes a stock update event
func (p *PaymentPublisher) PublishStockUpdate(ctx context.Context, event *events.StockUpdateEvent) error {
	event.EventID = uuid.New().String()
	event.EventType = events.StockUpdateEventType
	event.Timestamp = time.Now()

	message, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal stock update event: %w", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: events.StockEventsTopic,
		Key:   sarama.StringEncoder(fmt.Sprintf("%d", event.ProductID)),
		Value: sarama.ByteEncoder(message),
		Headers: []sarama.RecordHeader{
			{Key: []byte("event_type"), Value: []byte(event.EventType)},
			{Key: []byte("product_id"), Value: []byte(fmt.Sprintf("%d", event.ProductID))},
		},
	}

	partition, offset, err := p.producer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("failed to send stock update event: %w", err)
	}

	p.logger.WithFields(logrus.Fields{
		"event_id":   event.EventID,
		"product_id": event.ProductID,
		"quantity":   event.Quantity,
		"operation":  event.Operation,
		"topic":      events.StockEventsTopic,
		"partition":  partition,
		"offset":     offset,
	}).Info("Stock update event published")

	return nil
}

// PublishBasketCleared publishes a basket cleared event
func (p *PaymentPublisher) PublishBasketCleared(ctx context.Context, event *events.BasketClearedEvent) error {
	event.EventID = uuid.New().String()
	event.EventType = events.BasketClearedEventType
	event.Timestamp = time.Now()

	message, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal basket cleared event: %w", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: events.BasketEventsTopic,
		Key:   sarama.StringEncoder(event.UserID),
		Value: sarama.ByteEncoder(message),
		Headers: []sarama.RecordHeader{
			{Key: []byte("event_type"), Value: []byte(event.EventType)},
			{Key: []byte("user_id"), Value: []byte(event.UserID)},
			{Key: []byte("basket_id"), Value: []byte(event.BasketID)},
		},
	}

	partition, offset, err := p.producer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("failed to send basket cleared event: %w", err)
	}

	p.logger.WithFields(logrus.Fields{
		"event_id":  event.EventID,
		"user_id":   event.UserID,
		"basket_id": event.BasketID,
		"topic":     events.BasketEventsTopic,
		"partition": partition,
		"offset":    offset,
	}).Info("Basket cleared event published")

	return nil
}

// Close closes the publisher
func (p *PaymentPublisher) Close() error {
	return p.producer.Close()
}
