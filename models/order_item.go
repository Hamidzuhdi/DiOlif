package models

import "time"

type OrderItem struct {
	ID          int       `json:"id"`
	UniformName   string  `json:"uniform_name"`
	Size        string    `json:"size"`
	Quantity    int       `json:"quantity"`
	UnitPrice   float64   `json:"unit_price"`
	Subtotal    float64   `json:"subtotal"`
	Notes       string    `json:"notes"`
	CreatedAt   time.Time `json:"created_at"`
}
