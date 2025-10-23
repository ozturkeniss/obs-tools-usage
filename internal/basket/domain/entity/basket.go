package entity

import (
	"time"
)

// Basket represents a shopping basket
type Basket struct {
	ID        string            `json:"id" redis:"id"`
	UserID    string            `json:"user_id" redis:"user_id"`
	Items     []BasketItem      `json:"items" redis:"items"`
	Total     float64           `json:"total" redis:"total"`
	CreatedAt time.Time         `json:"created_at" redis:"created_at"`
	UpdatedAt time.Time         `json:"updated_at" redis:"updated_at"`
	ExpiresAt time.Time         `json:"expires_at" redis:"expires_at"`
	Metadata  map[string]string `json:"metadata,omitempty" redis:"metadata"`
}

// BasketItem represents an item in the basket
type BasketItem struct {
	ProductID int     `json:"product_id" redis:"product_id"`
	Name      string  `json:"name" redis:"name"`
	Price     float64 `json:"price" redis:"price"`
	Quantity  int     `json:"quantity" redis:"quantity"`
	Subtotal  float64 `json:"subtotal" redis:"subtotal"`
	Category  string  `json:"category,omitempty" redis:"category"`
}

// CalculateTotal calculates the total price of the basket
func (b *Basket) CalculateTotal() {
	total := 0.0
	for i := range b.Items {
		b.Items[i].Subtotal = b.Items[i].Price * float64(b.Items[i].Quantity)
		total += b.Items[i].Subtotal
	}
	b.Total = total
	b.UpdatedAt = time.Now()
}

// AddItem adds an item to the basket
func (b *Basket) AddItem(productID int, name string, price float64, quantity int, category string) {
	// Check if item already exists
	for i := range b.Items {
		if b.Items[i].ProductID == productID {
			b.Items[i].Quantity += quantity
			b.Items[i].Subtotal = b.Items[i].Price * float64(b.Items[i].Quantity)
			b.CalculateTotal()
			return
		}
	}

	// Add new item
	item := BasketItem{
		ProductID: productID,
		Name:      name,
		Price:     price,
		Quantity:  quantity,
		Category:  category,
	}
	item.Subtotal = item.Price * float64(item.Quantity)
	
	b.Items = append(b.Items, item)
	b.CalculateTotal()
}

// RemoveItem removes an item from the basket
func (b *Basket) RemoveItem(productID int) {
	for i := range b.Items {
		if b.Items[i].ProductID == productID {
			b.Items = append(b.Items[:i], b.Items[i+1:]...)
			b.CalculateTotal()
			return
		}
	}
}

// UpdateItemQuantity updates the quantity of an item
func (b *Basket) UpdateItemQuantity(productID int, quantity int) {
	for i := range b.Items {
		if b.Items[i].ProductID == productID {
			if quantity <= 0 {
				b.RemoveItem(productID)
			} else {
				b.Items[i].Quantity = quantity
				b.Items[i].Subtotal = b.Items[i].Price * float64(b.Items[i].Quantity)
				b.CalculateTotal()
			}
			return
		}
	}
}

// Clear removes all items from the basket
func (b *Basket) Clear() {
	b.Items = []BasketItem{}
	b.Total = 0.0
	b.UpdatedAt = time.Now()
}

// IsExpired checks if the basket is expired
func (b *Basket) IsExpired() bool {
	return time.Now().After(b.ExpiresAt)
}

// GetItemCount returns the total number of items in the basket
func (b *Basket) GetItemCount() int {
	count := 0
	for _, item := range b.Items {
		count += item.Quantity
	}
	return count
}
