package models



type OrderItem struct {
	ID          int       `json:"id"`
	TransactionID int 	 `json:"transaction_id"`
	UniformName   string  `json:"uniform_name"`
	Size        string    `json:"size"`
	Quantity    int       `json:"quantity"`
	UnitPrice   float64   `json:"unit_price"`
	Subtotal    float64   `json:"subtotal"`
	Notes       string    `json:"notes"`
	// CreatedAt   string    `json:"created_at"`
	// UpdatedAt   string    `json:"updated_at"`
}
