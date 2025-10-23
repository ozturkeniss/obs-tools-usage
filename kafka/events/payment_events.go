package events

import (
	"time"
)

// PaymentCompletedEvent represents a payment completion event
type PaymentCompletedEvent struct {
	EventID     string                 `json:"event_id"`
	EventType   string                 `json:"event_type"`
	Timestamp   time.Time              `json:"timestamp"`
	PaymentID   string                 `json:"payment_id"`
	UserID      string                 `json:"user_id"`
	BasketID    string                 `json:"basket_id"`
	Amount      float64                `json:"amount"`
	Currency    string                 `json:"currency"`
	Items       []PaymentItemEvent     `json:"items"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// PaymentItemEvent represents a payment item in the event
type PaymentItemEvent struct {
	ProductID int     `json:"product_id"`
	Name      string  `json:"name"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
	Subtotal  float64 `json:"subtotal"`
	Category  string  `json:"category"`
}

// PaymentFailedEvent represents a payment failure event
type PaymentFailedEvent struct {
	EventID     string                 `json:"event_id"`
	EventType   string                 `json:"event_type"`
	Timestamp   time.Time              `json:"timestamp"`
	PaymentID   string                 `json:"payment_id"`
	UserID      string                 `json:"user_id"`
	BasketID    string                 `json:"basket_id"`
	Amount      float64                `json:"amount"`
	Currency    string                 `json:"currency"`
	Reason      string                 `json:"reason"`
	ErrorCode   string                 `json:"error_code"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// PaymentRefundedEvent represents a payment refund event
type PaymentRefundedEvent struct {
	EventID     string                 `json:"event_id"`
	EventType   string                 `json:"event_type"`
	Timestamp   time.Time              `json:"timestamp"`
	PaymentID   string                 `json:"payment_id"`
	UserID      string                 `json:"user_id"`
	Amount      float64                `json:"amount"`
	Currency    string                 `json:"currency"`
	Reason      string                 `json:"reason"`
	RefundID    string                 `json:"refund_id"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// StockUpdateEvent represents a stock update event
type StockUpdateEvent struct {
	EventID     string                 `json:"event_id"`
	EventType   string                 `json:"event_type"`
	Timestamp   time.Time              `json:"timestamp"`
	ProductID   int                    `json:"product_id"`
	Quantity    int                    `json:"quantity"`
	Operation   string                 `json:"operation"` // "decrease" or "increase"
	Reason      string                 `json:"reason"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// BasketClearedEvent represents a basket clearing event
type BasketClearedEvent struct {
	EventID     string                 `json:"event_id"`
	EventType   string                 `json:"event_type"`
	Timestamp   time.Time              `json:"timestamp"`
	UserID      string                 `json:"user_id"`
	BasketID    string                 `json:"basket_id"`
	Reason      string                 `json:"reason"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// Event types
const (
	PaymentCompletedEventType = "payment.completed"
	PaymentFailedEventType    = "payment.failed"
	PaymentRefundedEventType  = "payment.refunded"
	StockUpdateEventType      = "stock.updated"
	BasketClearedEventType    = "basket.cleared"
)

// Kafka topics
const (
	PaymentEventsTopic = "payment-events"
	StockEventsTopic   = "stock-events"
	BasketEventsTopic  = "basket-events"
)
