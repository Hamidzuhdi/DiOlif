package models


type Transaksi struct {
	ID            int       `json:"id" gorm:"primaryKey;autoIncrement"`
	CustomerID    int       `json:"customer_id"`
	Transaksidate string `json:"transaction_date"`
	Paymentdate   string `json:"payment_date"`
	Status        string    `json:"status"` // paid, unpaid
	Total         float64   `json:"total_price"`
	Notes         string    `json:"notes"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
	Customer_name string      `json:"customer_name"`
	Items         []OrderItem `json:"items" gorm:"foreignKey:TransaksiID"`
}
