package models


type StudentOrderItem struct {
    ID            int       `json:"id"`
    CustomerID    int       `json:"customer_id"`
    StudentName   string    `json:"student_name"`
    Grade         string    `json:"grade"`
    TransactionID int       `json:"transaction_id"`
    UniformName   string    `json:"uniform_name"`
    Size          string    `json:"size"`
    Quantity      int       `json:"quantity"`
    UnitPrice     float64   `json:"unit_price"`
    Subtotal      float64   `json:"subtotal"`
    Notes         string    `json:"notes"`
    CreatedAt     string `json:"created_at"`
}