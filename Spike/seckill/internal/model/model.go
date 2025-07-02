
package model

import "time"

// Product represents a product that can be purchased.
type Product struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Stock int64  `json:"stock"`
}

// Order represents a customer's order.
type Order struct {
	ID        int64     `json:"id"`
	ProductID int64     `json:"product_id"`
	UserID    int64     `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}
