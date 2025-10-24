package events

// Notification-specific event types
const (
	// User Events
	UserRegisteredEventType     = "user_registered"
	UserLoggedInEventType       = "user_logged_in"
	UserLoggedOutEventType      = "user_logged_out"
	UserProfileUpdatedEventType = "user_profile_updated"
	UserPasswordChangedEventType = "user_password_changed"
	
	// Product Events
	ProductCreatedEventType     = "product_created"
	ProductUpdatedEventType     = "product_updated"
	ProductDeletedEventType     = "product_deleted"
	ProductViewedEventType      = "product_viewed"
	ProductAddedToWishlistEventType = "product_added_to_wishlist"
	ProductRemovedFromWishlistEventType = "product_removed_from_wishlist"
	
	// Basket Events
	BasketCreatedEventType      = "basket_created"
	BasketUpdatedEventType      = "basket_updated"
	BasketItemAddedEventType    = "basket_item_added"
	BasketItemRemovedEventType  = "basket_item_removed"
	BasketItemUpdatedEventType  = "basket_item_updated"
	BasketAbandonedEventType    = "basket_abandoned"
	BasketRecoveredEventType    = "basket_recovered"
	
	// Order Events
	OrderCreatedEventType       = "order_created"
	OrderConfirmedEventType     = "order_confirmed"
	OrderShippedEventType       = "order_shipped"
	OrderDeliveredEventType     = "order_delivered"
	OrderCancelledEventType     = "order_cancelled"
	OrderReturnedEventType      = "order_returned"
	
	// Payment Events (already defined in payment_events.go)
	// PaymentCompletedEventType, PaymentFailedEventType, etc.
	
	// Inventory Events
	StockLowEventType           = "stock_low"
	StockOutEventType           = "stock_out"
	StockRestockedEventType     = "stock_restocked"
	
	// System Events
	SystemMaintenanceEventType  = "system_maintenance"
	SystemUpdateEventType       = "system_update"
	SystemAlertEventType        = "system_alert"
	
	// Marketing Events
	PromotionCreatedEventType   = "promotion_created"
	PromotionExpiredEventType   = "promotion_expired"
	NewsletterSentEventType     = "newsletter_sent"
	CampaignLaunchedEventType   = "campaign_launched"
)

// UserRegisteredEvent represents a user registration event
type UserRegisteredEvent struct {
	EventID   string `json:"event_id"`
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Timestamp string `json:"timestamp"`
}

// UserLoggedInEvent represents a user login event
type UserLoggedInEvent struct {
	EventID   string `json:"event_id"`
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	IPAddress string `json:"ip_address"`
	UserAgent string `json:"user_agent"`
	Timestamp string `json:"timestamp"`
}

// ProductCreatedEvent represents a product creation event
type ProductCreatedEvent struct {
	EventID     string `json:"event_id"`
	ProductID   int    `json:"product_id"`
	Name        string `json:"name"`
	Category    string `json:"category"`
	Price       float64 `json:"price"`
	Stock       int    `json:"stock"`
	CreatedBy   string `json:"created_by"`
	Timestamp   string `json:"timestamp"`
}

// ProductViewedEvent represents a product view event
type ProductViewedEvent struct {
	EventID   string `json:"event_id"`
	ProductID int    `json:"product_id"`
	UserID    string `json:"user_id"`
	SessionID string `json:"session_id"`
	Timestamp string `json:"timestamp"`
}

// BasketItemAddedEvent represents a basket item addition event
type BasketItemAddedEvent struct {
	EventID     string `json:"event_id"`
	UserID      string `json:"user_id"`
	BasketID    string `json:"basket_id"`
	ProductID   int    `json:"product_id"`
	ProductName string `json:"product_name"`
	Quantity    int    `json:"quantity"`
	Price       float64 `json:"price"`
	Timestamp   string `json:"timestamp"`
}

// BasketAbandonedEvent represents a basket abandonment event
type BasketAbandonedEvent struct {
	EventID     string `json:"event_id"`
	UserID      string `json:"user_id"`
	BasketID    string `json:"basket_id"`
	ItemCount   int    `json:"item_count"`
	TotalValue  float64 `json:"total_value"`
	AbandonedAt string `json:"abandoned_at"`
	Timestamp   string `json:"timestamp"`
}

// OrderCreatedEvent represents an order creation event
type OrderCreatedEvent struct {
	EventID     string `json:"event_id"`
	OrderID     string `json:"order_id"`
	UserID      string `json:"user_id"`
	TotalAmount float64 `json:"total_amount"`
	Currency    string `json:"currency"`
	ItemCount   int    `json:"item_count"`
	Timestamp   string `json:"timestamp"`
}

// OrderShippedEvent represents an order shipment event
type OrderShippedEvent struct {
	EventID       string `json:"event_id"`
	OrderID       string `json:"order_id"`
	UserID        string `json:"user_id"`
	TrackingNumber string `json:"tracking_number"`
	Carrier       string `json:"carrier"`
	EstimatedDelivery string `json:"estimated_delivery"`
	Timestamp     string `json:"timestamp"`
}

// StockLowEvent represents a low stock event
type StockLowEvent struct {
	EventID     string `json:"event_id"`
	ProductID   int    `json:"product_id"`
	ProductName string `json:"product_name"`
	CurrentStock int   `json:"current_stock"`
	Threshold   int    `json:"threshold"`
	Timestamp   string `json:"timestamp"`
}

// StockOutEvent represents a stock out event
type StockOutEvent struct {
	EventID     string `json:"event_id"`
	ProductID   int    `json:"product_id"`
	ProductName string `json:"product_name"`
	Timestamp   string `json:"timestamp"`
}

// SystemMaintenanceEvent represents a system maintenance event
type SystemMaintenanceEvent struct {
	EventID     string `json:"event_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
	Severity    string `json:"severity"`
	Timestamp   string `json:"timestamp"`
}

// PromotionCreatedEvent represents a promotion creation event
type PromotionCreatedEvent struct {
	EventID     string  `json:"event_id"`
	PromotionID string  `json:"promotion_id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Discount    float64 `json:"discount"`
	StartDate   string  `json:"start_date"`
	EndDate     string  `json:"end_date"`
	Timestamp   string  `json:"timestamp"`
}
