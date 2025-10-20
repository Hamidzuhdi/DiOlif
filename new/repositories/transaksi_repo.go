package repositories

import (
    "database/sql"
    "konveksi-app/models"
    _ "github.com/go-sql-driver/mysql"
    "fmt"
    "log"
)

type TransactionRepository struct {
    DB *sql.DB
}

type OrderItemUpdate struct {
    UniformName string
    Size        string
    Quantity    int
    UnitPrice   float64
    Notes       string
}

func (r *TransactionRepository) Create(transaction *models.Transaksi) error {
    log.Printf("Starting transaction creation for customer: %d", transaction.CustomerID)
    
    tx, err := r.DB.Begin()
    if err != nil {
        log.Printf("Error starting transaction: %v", err)
        return fmt.Errorf("failed to start transaction: %v", err)
    }
    defer func() {
        if err != nil {
            if rollbackErr := tx.Rollback(); rollbackErr != nil {
                log.Printf("Error rolling back transaction: %v", rollbackErr)
            }
        }
    }()

    result, err := tx.Exec(`
        INSERT INTO transactions 
        (customer_id, transaction_date, payment_date, status, total_price, notes) 
        VALUES (?, ?, ?, ?, ?, ?)`,
        transaction.CustomerID, 
        transaction.Transaksidate,
        transaction.Paymentdate, 
        transaction.Status,
        transaction.Total, 
        transaction.Notes,
    )
    if err != nil {
        log.Printf("Error inserting transaction: %v", err)
        return fmt.Errorf("failed to insert transaction: %v", err)
    }

    transactionID, err := result.LastInsertId()
    if err != nil {
        log.Printf("Error getting transaction ID: %v", err)
        return fmt.Errorf("failed to get transaction ID: %v", err)
    }
    transaction.ID = int(transactionID)

    log.Printf("Transaction inserted with ID: %d", transaction.ID)

    for i, item := range transaction.Items {
        log.Printf("Inserting item %d: %s size %s qty %d price %.2f", 
            i+1, item.UniformName, item.Size, item.Quantity, item.UnitPrice)
        
        _, err := tx.Exec(`
            INSERT INTO order_items 
            (transaction_id, uniform_name, size, quantity, unit_price, notes) 
            VALUES (?, ?, ?, ?, ?, ?)`,
            transaction.ID, 
            item.UniformName, 
            item.Size,
            item.Quantity, 
            item.UnitPrice, 
            item.Notes,
        )
        if err != nil {
            log.Printf("Error inserting order item %d: %v", i+1, err)
            return fmt.Errorf("failed to insert order item %d: %v", i+1, err)
        }
    }

    log.Printf("All %d items inserted successfully", len(transaction.Items))

    if err = tx.Commit(); err != nil {
        log.Printf("Error committing transaction: %v", err)
        return fmt.Errorf("failed to commit transaction: %v", err)
    }

    log.Printf("Transaction %d committed successfully", transaction.ID)
    return nil
}

func (r *TransactionRepository) CreateStudentOrder(transaction *models.Transaksi, studentItems []models.StudentOrderItem) error {
    log.Printf("Starting student order creation for customer: %d", transaction.CustomerID)
    
    tx, err := r.DB.Begin()
    if err != nil {
        log.Printf("Error starting transaction: %v", err)
        return err
    }
    defer func() {
        if err != nil {
            tx.Rollback()
            log.Printf("Transaction rolled back due to error: %v", err)
        }
    }()

    var total float64
    for _, item := range studentItems {
        total += item.UnitPrice * float64(item.Quantity)
    }
    transaction.Total = total

    result, err := tx.Exec(`
        INSERT INTO transactions 
        (customer_id, transaction_date, payment_date, status, total_price, notes) 
        VALUES (?, ?, ?, ?, ?, ?)`,
        transaction.CustomerID, 
        transaction.Transaksidate,
        transaction.Paymentdate, 
        transaction.Status,
        transaction.Total, 
        transaction.Notes,
    )
    if err != nil {
        log.Printf("Error inserting transaction: %v", err)
        return err
    }

    transactionID, err := result.LastInsertId()
    if err != nil {
        log.Printf("Error getting transaction ID: %v", err)
        return err
    }
    transaction.ID = int(transactionID)

    log.Printf("Transaction inserted with ID: %d", transaction.ID)

    for i, item := range studentItems {
        item.TransactionID = transaction.ID
        
        _, err := tx.Exec(`
            INSERT INTO student_order_items 
            (customer_id, student_name, grade, transaction_id, uniform_name, size, quantity, unit_price, notes) 
            VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
            item.CustomerID,
            item.StudentName,
            item.Grade,
            item.TransactionID,
            item.UniformName,
            item.Size,
            item.Quantity,
            item.UnitPrice,
            item.Notes,
        )
        if err != nil {
            log.Printf("Error inserting student item %d: %v", i+1, err)
            return err
        }
        
        log.Printf("Student item %d inserted: %s - %s (%s) x%d", 
            i+1, item.StudentName, item.UniformName, item.Size, item.Quantity)
    }

    log.Printf("All %d student items inserted successfully", len(studentItems))

    if err = tx.Commit(); err != nil {
        log.Printf("Error committing transaction: %v", err)
        return err
    }

    log.Printf("Student order %d committed successfully", transaction.ID)
    return nil
}

func (r *TransactionRepository) GetAllTransactions() ([]struct {
    models.Transaksi
    TransactionType string
}, error) {
    query := `
        SELECT t.id, t.customer_id, c.name AS customer_name, 
               t.transaction_date, t.payment_date, t.status, 
               COALESCE(t.total_price, 0) AS total_price, 
               COALESCE(t.notes, '') AS notes, 
               t.created_at,
               CASE 
                   WHEN EXISTS(SELECT 1 FROM student_order_items soi WHERE soi.transaction_id = t.id) THEN 'student_order'
                   WHEN EXISTS(SELECT 1 FROM order_items oi WHERE oi.transaction_id = t.id) THEN 'item_order'
                   ELSE 'unknown'
               END AS transaction_type
        FROM transactions t
        JOIN customers c ON t.customer_id = c.id
        ORDER BY t.created_at DESC`
    
    rows, err := r.DB.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var results []struct {
        models.Transaksi
        TransactionType string
    }

    for rows.Next() {
        var result struct {
            models.Transaksi
            TransactionType string
        }
        
        err := rows.Scan(
            &result.ID, 
            &result.CustomerID, 
            &result.Customer_name,
            &result.Transaksidate, 
            &result.Paymentdate, 
            &result.Status,
            &result.Total,
            &result.Notes, 
            &result.CreatedAt,
            &result.TransactionType,
        )
        if err != nil {
            return nil, err
        }
        results = append(results, result)
    }
    return results, nil
}

func (r *TransactionRepository) GetByIDNormal(id int) (*models.Transaksi, error) {
    var t models.Transaksi
    query := `
        SELECT t.id, t.customer_id, c.name AS customer_name, 
            t.transaction_date, t.payment_date, t.status, 
            t.total_price, t.notes, t.created_at
        FROM transactions t
        JOIN customers c ON t.customer_id = c.id
        WHERE t.id = ?`
    err := r.DB.QueryRow(query, id).Scan(
        &t.ID, &t.CustomerID, &t.Customer_name,
        &t.Transaksidate, &t.Paymentdate, &t.Status,
        &t.Total, &t.Notes, &t.CreatedAt,
    )
    if err != nil {
        return nil, err
    }
    
    itemsQuery := `
        SELECT id, uniform_name, size, quantity, unit_price, subtotal, notes
        FROM order_items WHERE transaction_id = ?`
    rows, err := r.DB.Query(itemsQuery, id)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    for rows.Next() {
        var item models.OrderItem
        err := rows.Scan(
            &item.ID, &item.UniformName, &item.Size,
            &item.Quantity, &item.UnitPrice,
            &item.Subtotal, &item.Notes,
        )
        if err != nil {
            return nil, err
        }
        t.Items = append(t.Items, item)
    }
    return &t, nil
}

func (r *TransactionRepository) GetByIDStudentOrder(id int) (*models.Transaksi, []models.StudentOrderItem, error) {
    var t models.Transaksi
    query := `
        SELECT t.id, t.customer_id, c.name AS customer_name, 
            t.transaction_date, t.payment_date, t.status, 
            t.total_price, t.notes, t.created_at
        FROM transactions t
        JOIN customers c ON t.customer_id = c.id
        WHERE t.id = ?`
    err := r.DB.QueryRow(query, id).Scan(
        &t.ID, &t.CustomerID, &t.Customer_name,
        &t.Transaksidate, &t.Paymentdate, &t.Status,
        &t.Total, &t.Notes, &t.CreatedAt,
    )
    if err != nil {
        return nil, nil, err
    }
    
    itemsQuery := `
        SELECT id, customer_id, student_name, grade, transaction_id, 
               uniform_name, size, quantity, unit_price, 
               (quantity * unit_price) as subtotal, notes, created_at
        FROM student_order_items WHERE transaction_id = ?`
    rows, err := r.DB.Query(itemsQuery, id)
    if err != nil {
        return &t, nil, err
    }
    defer rows.Close()
    
    var items []models.StudentOrderItem
    for rows.Next() {
        var item models.StudentOrderItem
        err := rows.Scan(
            &item.ID, &item.CustomerID, &item.StudentName, &item.Grade, &item.TransactionID,
            &item.UniformName, &item.Size, &item.Quantity, &item.UnitPrice, 
            &item.Subtotal, &item.Notes, &item.CreatedAt,
        )
        if err != nil {
            return &t, nil, err
        }
        items = append(items, item)
    }
    return &t, items, nil
}

func (r *TransactionRepository) GetCustomerUniformsByTransactionID(transactionID int) ([]models.CustomerUniform, error) {
    query := `
        SELECT cu.id, cu.customer_id, cu.uniform_name, cu.size, 
               cu.price, cu.created_at
        FROM customer_uniforms cu
        JOIN transactions t ON cu.customer_id = t.customer_id
        WHERE t.id = ?
        ORDER BY cu.uniform_name, cu.size
    `
    
    rows, err := r.DB.Query(query, transactionID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var uniforms []models.CustomerUniform
    for rows.Next() {
        var uniform models.CustomerUniform
        err := rows.Scan(
            &uniform.ID,
            &uniform.CustomerID,
            &uniform.UniformName,
            &uniform.Size,
            &uniform.Price,
            &uniform.CreatedAt,
        )
        if err != nil {
            return nil, err
        }
        uniforms = append(uniforms, uniform)
    }
    
    return uniforms, nil
}

func (r *TransactionRepository) UpdateStatus(id int, status string) error {
    _, err := r.DB.Exec(
        "UPDATE transactions SET status = ? WHERE id = ?",
        status, id,
    )
    return err
}

func (r *TransactionRepository) UpdateTransaction(id int, transactionDate, paymentDate, notes string) error {
    _, err := r.DB.Exec(
        "UPDATE transactions SET transaction_date = ?, payment_date = ?, notes = ? WHERE id = ?",
        transactionDate, paymentDate, notes, id,
    )
    return err
}

func (r *TransactionRepository) UpdateOrderItemsNormal(transactionID int, items []OrderItemUpdate) error {
    var count int
    err := r.DB.QueryRow(`SELECT COUNT(*) FROM transactions WHERE id = ?`, transactionID).Scan(&count)
    if err != nil || count == 0 {
        return sql.ErrNoRows
    }
    return r.UpdateOrderItems(transactionID, items)
}

func (r *TransactionRepository) UpdateOrderItems(transactionID int, items []OrderItemUpdate) error {
    tx, err := r.DB.Begin()
    if err != nil {
        return err
    }
    defer func() {
        if err != nil {
            tx.Rollback()
        }
    }()
    
    _, err = tx.Exec("DELETE FROM order_items WHERE transaction_id = ?", transactionID)
    if err != nil {
        return err
    }
    
    total := 0.0
    for _, item := range items {
        subtotal := float64(item.Quantity) * item.UnitPrice
        total += subtotal
        _, err := tx.Exec(
            `INSERT INTO order_items (transaction_id, uniform_name, size, quantity, unit_price, notes) VALUES (?, ?, ?, ?, ?, ?)`,
            transactionID, item.UniformName, item.Size, item.Quantity, item.UnitPrice, item.Notes, 
        )
        if err != nil {
            return err
        }
    }
    
    _, err = tx.Exec("UPDATE transactions SET total_price = ? WHERE id = ?", total, transactionID)
    if err != nil {
        return err
    }
    return tx.Commit()
}

func (r *TransactionRepository) UpdateOrderItemsStudent(transactionID int, studentItems []models.StudentOrderItem) error {
    tx, err := r.DB.Begin()
    if err != nil {
        return err
    }
    defer func() {
        if err != nil {
            tx.Rollback()
        }
    }()
    
    _, err = tx.Exec("DELETE FROM student_order_items WHERE transaction_id = ?", transactionID)
    if err != nil {
        return err
    }
    
    total := 0.0
    for _, item := range studentItems {
        subtotal := float64(item.Quantity) * item.UnitPrice
        total += subtotal
        _, err := tx.Exec(
            `INSERT INTO student_order_items 
            (customer_id, student_name, grade, transaction_id, uniform_name, size, quantity, unit_price, notes) 
            VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
            item.CustomerID, item.StudentName, item.Grade, transactionID,
            item.UniformName, item.Size, item.Quantity, item.UnitPrice, item.Notes,
        )
        if err != nil {
            return err
        }
    }
    
    _, err = tx.Exec("UPDATE transactions SET total_price = ? WHERE id = ?", total, transactionID)
    if err != nil {
        return err
    }
    return tx.Commit()
}

func (r *TransactionRepository) UpdateStudentOrderItem(itemID int, studentName, grade, uniformName, size string, quantity int, unitPrice float64, notes string) error {
    _, err := r.DB.Exec(
        `UPDATE student_order_items 
         SET student_name = ?, grade = ?, uniform_name = ?, size = ?, quantity = ?, unit_price = ?, notes = ?
         WHERE id = ?`,
        studentName, grade, uniformName, size, quantity, unitPrice, notes, itemID,
    )
    return err
}

func (r *TransactionRepository) UpdateNormalOrderItem(itemID int, uniformName, size string, quantity int, unitPrice float64, notes string) error {
    _, err := r.DB.Exec(
        `UPDATE order_items 
         SET uniform_name = ?, size = ?, quantity = ?, unit_price = ?, notes = ?
         WHERE id = ?`,
        uniformName, size, quantity, unitPrice, notes, itemID,
    )
    return err
}

func (r *TransactionRepository) RecalculateTransactionTotal(transactionID int) error {
    tx, err := r.DB.Begin()
    if err != nil {
        return err
    }
    defer func() {
        if err != nil {
            tx.Rollback()
        }
    }()

    var studentTotal, normalTotal float64
    
    err = tx.QueryRow(
        "SELECT COALESCE(SUM(quantity * unit_price), 0) FROM student_order_items WHERE transaction_id = ?",
        transactionID,
    ).Scan(&studentTotal)
    if err != nil {
        return err
    }
    
    err = tx.QueryRow(
        "SELECT COALESCE(SUM(quantity * unit_price), 0) FROM order_items WHERE transaction_id = ?",
        transactionID,
    ).Scan(&normalTotal)
    if err != nil {
        return err
    }
    
    total := studentTotal + normalTotal
    
    _, err = tx.Exec(
        "UPDATE transactions SET total_price = ? WHERE id = ?",
        total, transactionID,
    )
    if err != nil {
        return err
    }
    
    return tx.Commit()
}

func (r *TransactionRepository) GetOverduePaymentTransactions() ([]models.Transaksi, error) {
    rows, err := r.DB.Query(`
        SELECT t.id, t.customer_id, c.name AS customer_name, t.transaction_date, t.payment_date, t.status, t.total_price, t.notes, t.created_at
        FROM transactions t
        JOIN customers c ON t.customer_id = c.id
        WHERE t.status != 'paid'
          AND t.payment_date IS NOT NULL
          AND t.payment_date < CURDATE()
    `)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var transactions []models.Transaksi
    for rows.Next() {
        var t models.Transaksi
        err := rows.Scan(&t.ID, &t.CustomerID, &t.Customer_name, &t.Transaksidate, &t.Paymentdate, &t.Status, &t.Total, &t.Notes, &t.CreatedAt)
        if err != nil {
            return nil, err
        }
        transactions = append(transactions, t)
    }
    return transactions, nil
}

func (r *TransactionRepository) GetOverdueTransactions() ([]models.Transaksi, error) {
    rows, err := r.DB.Query(`
        SELECT t.id, t.customer_id, c.name AS customer_name, t.transaction_date, t.payment_date, t.status, t.total_price, t.notes, t.created_at
        FROM transactions t
        JOIN customers c ON t.customer_id = c.id
        WHERE t.status = 'pending'
          AND t.transaction_date IS NOT NULL
          AND t.transaction_date < CURDATE()
    `)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var transactions []models.Transaksi
    for rows.Next() {
        var t models.Transaksi
        err := rows.Scan(&t.ID, &t.CustomerID, &t.Customer_name, &t.Transaksidate, &t.Paymentdate, &t.Status, &t.Total, &t.Notes, &t.CreatedAt)
        if err != nil {
            return nil, err
        }
        transactions = append(transactions, t)
    }
    return transactions, nil
}

func (r *TransactionRepository) GetRemindTransactions() ([]models.Transaksi, error) {
    rows, err := r.DB.Query(`
        SELECT t.id, t.customer_id, t.transaction_date, t.payment_date, t.status,
               t.total_price, t.notes, t.created_at, t.updated_at, c.name AS customer_name
        FROM transactions t
        JOIN customers c ON t.customer_id = c.id
        WHERE t.status != 'paid'
        AND t.transaction_date IS NOT NULL
        AND CURDATE() >= DATE_SUB(t.transaction_date, INTERVAL 2 DAY)
        AND CURDATE() < t.transaction_date
    `)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var transactions []models.Transaksi
    for rows.Next() {
        var t models.Transaksi
        err := rows.Scan(
            &t.ID, &t.CustomerID, &t.Transaksidate, &t.Paymentdate, &t.Status,
            &t.Total, &t.Notes, &t.CreatedAt, &t.UpdatedAt, &t.Customer_name,
        )
        if err != nil {
            return nil, err
        }
        transactions = append(transactions, t)
    }
    return transactions, nil
}

func (r *TransactionRepository) GetRemindPayments() ([]models.Transaksi, error) {
    rows, err := r.DB.Query(`
        SELECT t.id, t.customer_id, t.transaction_date, t.payment_date, t.status,
               t.total_price, t.notes, t.created_at, t.updated_at, c.name AS customer_name
        FROM transactions t
        JOIN customers c ON t.customer_id = c.id
        WHERE t.status != 'paid'
        AND t.payment_date IS NOT NULL
        AND CURDATE() >= DATE_SUB(t.payment_date, INTERVAL 2 DAY)
        AND CURDATE() < t.payment_date
    `)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var transactions []models.Transaksi
    for rows.Next() {
        var t models.Transaksi
        err := rows.Scan(
            &t.ID, &t.CustomerID, &t.Transaksidate, &t.Paymentdate, &t.Status,
            &t.Total, &t.Notes, &t.CreatedAt, &t.UpdatedAt, &t.Customer_name,
        )
        if err != nil {
            return nil, err
        }
        transactions = append(transactions, t)
    }
    return transactions, nil
}

func (r *TransactionRepository) GetTransactionsByCustomerID(customerID int, statusFilter string) ([]struct {
    ID              int     `json:"id"`
    CustomerID      int     `json:"customer_id"`
    CustomerName    string  `json:"customer_name"`
    TransactionDate string  `json:"transaction_date"`
    PaymentDate     string  `json:"payment_date"`
    Status          string  `json:"status"`
    TotalPrice      float64 `json:"total_price"`
    Notes           string  `json:"notes"`
    CreatedAt       string  `json:"created_at"`
    ItemCount       int     `json:"item_count"`
    HasStudentInfo  bool    `json:"has_student_info"`
}, error) {
    baseQuery := `
        SELECT 
            t.id, 
            t.customer_id,
            c.name as customer_name,
            t.transaction_date,
            t.payment_date,
            t.status,
            t.total_price,
            COALESCE(t.notes, '') as notes,
            t.created_at
        FROM transactions t
        LEFT JOIN customers c ON t.customer_id = c.id
        WHERE t.customer_id = ?`
    
    args := []interface{}{customerID}
    
    if statusFilter != "" {
        baseQuery += " AND t.status = ?"
        args = append(args, statusFilter)
    }
    
    baseQuery += " ORDER BY t.transaction_date DESC"
    
    rows, err := r.DB.Query(baseQuery, args...)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var transactions []struct {
        ID              int     `json:"id"`
        CustomerID      int     `json:"customer_id"`
        CustomerName    string  `json:"customer_name"`
        TransactionDate string  `json:"transaction_date"`
        PaymentDate     string  `json:"payment_date"`
        Status          string  `json:"status"`
        TotalPrice      float64 `json:"total_price"`
        Notes           string  `json:"notes"`
        CreatedAt       string  `json:"created_at"`
        ItemCount       int     `json:"item_count"`
        HasStudentInfo  bool    `json:"has_student_info"`
    }
    
    for rows.Next() {
        var t struct {
            ID              int     `json:"id"`
            CustomerID      int     `json:"customer_id"`
            CustomerName    string  `json:"customer_name"`
            TransactionDate string  `json:"transaction_date"`
            PaymentDate     string  `json:"payment_date"`
            Status          string  `json:"status"`
            TotalPrice      float64 `json:"total_price"`
            Notes           string  `json:"notes"`
            CreatedAt       string  `json:"created_at"`
            ItemCount       int     `json:"item_count"`
            HasStudentInfo  bool    `json:"has_student_info"`
        }
        
        err := rows.Scan(
            &t.ID,
            &t.CustomerID,
            &t.CustomerName,
            &t.TransactionDate,
            &t.PaymentDate,
            &t.Status,
            &t.TotalPrice,
            &t.Notes,
            &t.CreatedAt,
        )
        if err != nil {
            return nil, err
        }
        
        t.ItemCount = r.getTransactionItemCountByID(t.ID)
        t.HasStudentInfo = r.HasStudentInfo(t.ID)
        
        transactions = append(transactions, t)
    }
    
    return transactions, nil
}

func (r *TransactionRepository) UpdateTransactionStatus(transactionID int, status string) error {
    _, err := r.DB.Exec(
        "UPDATE transactions SET status = ?, updated_at = NOW() WHERE id = ?",
        status, transactionID,
    )
    return err
}

func (r *TransactionRepository) HasStudentInfo(transactionID int) bool {
    var count int
    
    query := `SELECT COUNT(*) FROM student_order_items 
              WHERE transaction_id = ? 
              AND student_name IS NOT NULL 
              AND student_name != '' 
              AND grade IS NOT NULL 
              AND grade != ''`
    
    err := r.DB.QueryRow(query, transactionID).Scan(&count)
    if err != nil {
        return false
    }
    
    return count > 0
}

func (r *TransactionRepository) GetTransactionItemCount(transactionID int, transactionType string) (int, error) {
    var query string
    if transactionType == "student_order" {
        query = "SELECT COUNT(*) FROM student_order_items WHERE transaction_id = ?"
    } else {
        query = "SELECT COUNT(*) FROM order_items WHERE transaction_id = ?"
    }
    
    var count int
    err := r.DB.QueryRow(query, transactionID).Scan(&count)
    return count, err
}

func (r *TransactionRepository) getTransactionItemCountByID(transactionID int) int {
    var studentCount int
    err := r.DB.QueryRow("SELECT COUNT(*) FROM student_order_items WHERE transaction_id = ?", transactionID).Scan(&studentCount)
    if err == nil && studentCount > 0 {
        return studentCount
    }
    
    var orderCount int
    err = r.DB.QueryRow("SELECT COUNT(*) FROM order_items WHERE transaction_id = ?", transactionID).Scan(&orderCount)
    if err == nil {
        return orderCount
    }
    
    return 0
}

func (r *TransactionRepository) GetCustomerByID(customerID int) (*models.Customer, error) {
    query := `SELECT id, name, address, contact FROM customers WHERE id = ?`
    
    var customer models.Customer
    err := r.DB.QueryRow(query, customerID).Scan(
        &customer.ID, 
        &customer.Name, 
        &customer.Address, 
        &customer.Contact,
    )
    if err != nil {
        log.Printf("Error getting customer by ID %d: %v", customerID, err)
        return nil, err
    }
    
    return &customer, nil
}