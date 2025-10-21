package product

import (
	"time"

	"github.com/sirupsen/logrus"
)

// BusinessEvent represents a business event
type BusinessEvent struct {
	EventType    string                 `json:"event_type"`
	EventName    string                 `json:"event_name"`
	Timestamp    time.Time              `json:"timestamp"`
	UserID       string                 `json:"user_id,omitempty"`
	SessionID    string                 `json:"session_id,omitempty"`
	ProductID    int                    `json:"product_id,omitempty"`
	ProductName  string                 `json:"product_name,omitempty"`
	Category     string                 `json:"category,omitempty"`
	Price        float64                `json:"price,omitempty"`
	Stock        int                    `json:"stock,omitempty"`
	OldValue     interface{}            `json:"old_value,omitempty"`
	NewValue     interface{}            `json:"new_value,omitempty"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
}

// LogBusinessEvent logs a business event
func LogBusinessEvent(logger *logrus.Entry, event BusinessEvent) {
	// Prepare event fields
	eventFields := map[string]interface{}{
		"business_event": true,
		"event_type":     event.EventType,
		"event_name":     event.EventName,
		"timestamp":      event.Timestamp,
		"user_id":        event.UserID,
		"session_id":     event.SessionID,
		"product_id":     event.ProductID,
		"product_name":   event.ProductName,
		"category":       event.Category,
		"price":          event.Price,
		"stock":          event.Stock,
		"old_value":      event.OldValue,
		"new_value":      event.NewValue,
		"metadata":       event.Metadata,
	}
	
	// Mask sensitive data in event fields
	maskedEventFields := MaskFields(eventFields)
	
	logger.WithFields(maskedEventFields).Info("Business event occurred")
}

// LogProductCreated logs when a product is created
func LogProductCreated(logger *logrus.Entry, product Product, userID string) {
	event := BusinessEvent{
		EventType:   "product",
		EventName:   "product_created",
		Timestamp:   time.Now(),
		UserID:      userID,
		ProductID:   product.ID,
		ProductName: product.Name,
		Category:    product.Category,
		Price:       product.Price,
		Stock:       product.Stock,
		Metadata: map[string]interface{}{
			"description": product.Description,
			"created_at":  product.CreatedAt,
		},
	}
	
	LogBusinessEvent(logger, event)
}

// LogProductUpdated logs when a product is updated
func LogProductUpdated(logger *logrus.Entry, oldProduct, newProduct Product, userID string) {
	event := BusinessEvent{
		EventType:   "product",
		EventName:   "product_updated",
		Timestamp:   time.Now(),
		UserID:      userID,
		ProductID:   newProduct.ID,
		ProductName: newProduct.Name,
		Category:    newProduct.Category,
		Price:       newProduct.Price,
		Stock:       newProduct.Stock,
		OldValue: map[string]interface{}{
			"name":        oldProduct.Name,
			"description": oldProduct.Description,
			"price":       oldProduct.Price,
			"stock":       oldProduct.Stock,
			"category":    oldProduct.Category,
		},
		NewValue: map[string]interface{}{
			"name":        newProduct.Name,
			"description": newProduct.Description,
			"price":       newProduct.Price,
			"stock":       newProduct.Stock,
			"category":    newProduct.Category,
		},
		Metadata: map[string]interface{}{
			"updated_at": newProduct.UpdatedAt,
		},
	}
	
	LogBusinessEvent(logger, event)
}

// LogProductDeleted logs when a product is deleted
func LogProductDeleted(logger *logrus.Entry, product Product, userID string) {
	event := BusinessEvent{
		EventType:   "product",
		EventName:   "product_deleted",
		Timestamp:   time.Now(),
		UserID:      userID,
		ProductID:   product.ID,
		ProductName: product.Name,
		Category:    product.Category,
		Price:       product.Price,
		Stock:       product.Stock,
		Metadata: map[string]interface{}{
			"deleted_at": time.Now(),
		},
	}
	
	LogBusinessEvent(logger, event)
}

// LogLowStockAlert logs when stock is low
func LogLowStockAlert(logger *logrus.Entry, product Product, threshold int) {
	event := BusinessEvent{
		EventType:   "inventory",
		EventName:   "low_stock_alert",
		Timestamp:   time.Now(),
		ProductID:   product.ID,
		ProductName: product.Name,
		Category:    product.Category,
		Stock:       product.Stock,
		Metadata: map[string]interface{}{
			"threshold": threshold,
			"alert_type": "low_stock",
		},
	}
	
	LogBusinessEvent(logger, event)
}

// LogHighValueProduct logs when a high-value product is accessed
func LogHighValueProduct(logger *logrus.Entry, product Product, userID string) {
	event := BusinessEvent{
		EventType:   "product",
		EventName:   "high_value_product_accessed",
		Timestamp:   time.Now(),
		UserID:      userID,
		ProductID:   product.ID,
		ProductName: product.Name,
		Category:    product.Category,
		Price:       product.Price,
		Metadata: map[string]interface{}{
			"value_threshold": 1000.0,
			"access_type":     "view",
		},
	}
	
	LogBusinessEvent(logger, event)
}
