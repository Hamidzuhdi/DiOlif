package repositories

import (
	"database/sql"
	"konveksi-app/models"
    _ "github.com/go-sql-driver/mysql"
)

type TransactionRepository struct {
	DB *sql.DB
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

func (r *TransactionRepository) GetAllStudentOrders() ([]struct {
    Transaction models.Transaksi
    Items       []models.StudentOrderItem
}, error) {
    query := `
        SELECT t.id, t.customer_id, c.name AS customer_name, 
            t.transaction_date, t.payment_date, t.status, 
            t.total_price, t.notes, t.created_at
        FROM transactions t
        JOIN customers c ON t.customer_id = c.id
        WHERE EXISTS (
            SELECT 1 FROM student_order_items soi WHERE soi.transaction_id = t.id
        )
        ORDER BY t.created_at DESC`
    rows, err := r.DB.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var results []struct {
        Transaction models.Transaksi
        Items       []models.StudentOrderItem
    }

    for rows.Next() {
        var t models.Transaksi
        err := rows.Scan(
            &t.ID, &t.CustomerID, &t.Customer_name,
            &t.Transaksidate, &t.Paymentdate, &t.Status,
            &t.Total, &t.Notes, &t.CreatedAt,
        )
        if err != nil {
            return nil, err
        }

        // Ambil student_order_items untuk transaksi ini
        itemsQuery := `
            SELECT id, customer_id, student_name, grade, transaction_id, uniform_name, size, quantity, unit_price, subtotal, notes, created_at
            FROM student_order_items WHERE transaction_id = ?`
        itemRows, err := r.DB.Query(itemsQuery, t.ID)
        if err != nil {
            return nil, err
        }
        var items []models.StudentOrderItem
        for itemRows.Next() {
            var item models.StudentOrderItem
            err := itemRows.Scan(
                &item.ID, &item.CustomerID, &item.StudentName, &item.Grade, &item.TransactionID,
                &item.UniformName, &item.Size, &item.Quantity, &item.UnitPrice, &item.Subtotal, &item.Notes, &item.CreatedAt,
            )
            if err != nil {
                itemRows.Close()
                return nil, err
            }
            items = append(items, item)
        }
        itemRows.Close()

        results = append(results, struct {
            Transaction models.Transaksi
            Items       []models.StudentOrderItem
        }{
            Transaction: t,
            Items:       items,
        })
    }
    return results, nil
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

func (r *TransactionRepository) Create(transaction *models.Transaksi) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}

	// Insert header transaksi
	res, err := tx.Exec(
		`INSERT INTO transactions 
		(customer_id, transaction_date, payment_date, status, total_price, notes) 
		VALUES (?, ?, ?, ?, ?, ?)`,
		transaction.CustomerID, transaction.Transaksidate,
		transaction.Paymentdate, transaction.Status,
		transaction.Total, transaction.Notes,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	transactionID, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return err
	}
	transaction.ID = int(transactionID)

	// Insert detail items
	for _, item := range transaction.Items {
		_, err := tx.Exec(
			`INSERT INTO order_items 
			(transaction_id, uniform_name, size, quantity, unit_price, notes) 
			VALUES (?, ?, ?, ?, ?, ?)`,
			transaction.ID, item.UniformName, item.Size,
			item.Quantity, item.UnitPrice, item.Notes,
		)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

// Untuk transaksi biasa (nama_murid & kelas kosong)
func (r *TransactionRepository) GetByIDNormal(id int) (*models.Transaksi, error) {
    var t models.Transaksi
    query := `
        SELECT t.id, t.customer_id, c.name AS customer_name, 
            t.transaction_date, t.payment_date, t.status, 
            t.total_price, t.notes, t.created_at
        FROM transactions t
        JOIN customers c ON t.customer_id = c.id
        WHERE t.id = ? AND (t.nama_murid IS NULL OR t.nama_murid = '') AND (t.kelas IS NULL OR t.kelas = '')`
    err := r.DB.QueryRow(query, id).Scan(
        &t.ID, &t.CustomerID, &t.Customer_name,
        &t.Transaksidate, &t.Paymentdate, &t.Status,
        &t.Total, &t.Notes, &t.CreatedAt,
    )
    if err != nil {
        return nil, err
    }
    // Ambil order items
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
    // Ambil student_order_items
    itemsQuery := `
        SELECT id, customer_id, student_name, grade, transaction_id, uniform_name, size, quantity, unit_price, subtotal, notes, created_at
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
            &item.UniformName, &item.Size, &item.Quantity, &item.UnitPrice, &item.Subtotal, &item.Notes, &item.CreatedAt,
        )
        if err != nil {
            return &t, nil, err
        }
        items = append(items, item)
    }
    return &t, items, nil
}

func (r *TransactionRepository) UpdateStatus(id int, status string) error {
	_, err := r.DB.Exec(
		"UPDATE transactions SET status = ? WHERE id = ?",
		status, id,
	)
	return err
}

// Ambil semua transaksi yang berasal dari order_items (transaksi biasa)
func (r *TransactionRepository) GetAllTransactionOrder() ([]models.Transaksi, error) {
    query := `
        SELECT t.id, t.customer_id, c.name AS customer_name, t.transaction_date, t.payment_date, t.status, t.total_price, t.notes, t.created_at
        FROM transactions t
        JOIN customers c ON t.customer_id = c.id
        WHERE EXISTS (
            SELECT 1 FROM order_items oi WHERE oi.transaction_id = t.id
        )
        ORDER BY t.id DESC`
    rows, err := r.DB.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var transactions []models.Transaksi
    for rows.Next() {
        var t models.Transaksi
        err := rows.Scan(
            &t.ID, &t.CustomerID, &t.Customer_name,
            &t.Transaksidate, &t.Paymentdate, &t.Status,
            &t.Total, &t.Notes, &t.CreatedAt,
        )
        if err != nil {
            return nil, err
        }
        transactions = append(transactions, t)
    }
    return transactions, nil
}

// Ambil semua transaksi yang berasal dari student_order_items (repeat order)
func (r *TransactionRepository) GetAllTransactionStudent() ([]models.Transaksi, error) {
    query := `
        SELECT t.id, t.customer_id, c.name AS customer_name, t.transaction_date, t.payment_date, t.status, t.total_price, t.notes, t.created_at
        FROM transactions t
        JOIN customers c ON t.customer_id = c.id
        WHERE EXISTS (
            SELECT 1 FROM student_order_items soi WHERE soi.transaction_id = t.id
        )
        ORDER BY t.id DESC`
    rows, err := r.DB.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var transactions []models.Transaksi
    for rows.Next() {
        var t models.Transaksi
        err := rows.Scan(
            &t.ID, &t.CustomerID, &t.Customer_name,
            &t.Transaksidate, &t.Paymentdate, &t.Status,
            &t.Total, &t.Notes, &t.CreatedAt,
        )
        if err != nil {
            return nil, err
        }
        transactions = append(transactions, t)
    }
    return transactions, nil
}

func (r *TransactionRepository) GetRemindTransactions() ([]models.Transaksi, error) {
    rows, err := r.DB.Query(`
        SELECT t.*, c.name AS customer_name
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
        &t.Total, &t.Notes, &t.CreatedAt,&t.UpdatedAt, &t.Customer_name, // <-- tambahkan ini
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
        SELECT t.*, c.name AS customer_name
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
        &t.Total, &t.Notes, &t.CreatedAt,&t.UpdatedAt, &t.Customer_name, // <-- tambahkan ini
    	)
        if err != nil {
            return nil, err
        }
        transactions = append(transactions, t)
    }
    return transactions, nil
}

type OrderItemUpdate struct {
    UniformName string
    Size        string
    Quantity    int
    UnitPrice   float64
    Notes       string
}

func (r *TransactionRepository) UpdateOrderItemsNormal(transactionID int, items []OrderItemUpdate) error {
    // Pastikan transaksi ini memang transaksi normal
    var count int
    err := r.DB.QueryRow(`SELECT COUNT(*) FROM transactions WHERE id = ? AND (nama_murid IS NULL OR nama_murid = '') AND (kelas IS NULL OR kelas = '')`, transactionID).Scan(&count)
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
    // Delete existing items
    _, err = tx.Exec("DELETE FROM order_items WHERE transaction_id = ?", transactionID)
    if err != nil {
        tx.Rollback()
        return err
    }
    // Insert new items & hitung total baru
    total := 0.0
    for _, item := range items {
        subtotal := float64(item.Quantity) * item.UnitPrice
        total += subtotal
        _, err := tx.Exec(
            `INSERT INTO order_items (transaction_id, uniform_name, size, quantity, unit_price, notes) VALUES (?, ?, ?, ?, ?, ?)`,
            transactionID, item.UniformName, item.Size, item.Quantity, item.UnitPrice, item.Notes, 
        )
        if err != nil {
            tx.Rollback()
            return err
        }
    }
    // Update total_price di tabel transactions
    _, err = tx.Exec("UPDATE transactions SET total_price = ? WHERE id = ?", total, transactionID)
    if err != nil {
        tx.Rollback()
        return err
    }
    return tx.Commit()
}

func (r *TransactionRepository) UpdateOrderItemsStudent(transactionID int, studentItems []models.StudentOrderItem) error {
    tx, err := r.DB.Begin()
    if err != nil {
        return err
    }
    _, err = tx.Exec("DELETE FROM student_order_items WHERE transaction_id = ?", transactionID)
    if err != nil {
        tx.Rollback()
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
            tx.Rollback()
            return err
        }
    }
    // Update total_price di tabel transactions
    _, err = tx.Exec("UPDATE transactions SET total_price = ? WHERE id = ?", total, transactionID)
    if err != nil {
        tx.Rollback()
        return err
    }
    return tx.Commit()
}

func (r *TransactionRepository) UpdateTransaction(id int, transactionDate, paymentDate, notes string) error {
    _, err := r.DB.Exec(
        "UPDATE transactions SET transaction_date = ?, payment_date = ?, notes = ? WHERE id = ?",
        transactionDate, paymentDate, notes, id,
    )
    return err
}

func (r *TransactionRepository) CreateStudentOrder(transaction *models.Transaksi, studentItems []models.StudentOrderItem) error {
    tx, err := r.DB.Begin()
    if err != nil {
        return err
    }

    // Hitung total
    total := 0.0
    for _, item := range studentItems {
        total += float64(item.Quantity) * item.UnitPrice
    }
    transaction.Total = total

    // Insert transaksi (header)
    res, err := tx.Exec(
        `INSERT INTO transactions 
        (customer_id, transaction_date, payment_date, status, total_price, notes) 
        VALUES (?, ?, ?, ?, ?, ?)`,
        transaction.CustomerID, transaction.Transaksidate,
        transaction.Paymentdate, transaction.Status,
        transaction.Total, transaction.Notes,
    )
    if err != nil {
        tx.Rollback()
        return err
    }
    transactionID, err := res.LastInsertId()
    if err != nil {
        tx.Rollback()
        return err
    }
    transaction.ID = int(transactionID)

    // Insert student_order_items
    for _, item := range studentItems {
        _, err := tx.Exec(
            `INSERT INTO student_order_items 
            (customer_id, student_name, grade, transaction_id, uniform_name, size, quantity, unit_price, notes) 
            VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
            item.CustomerID, item.StudentName, item.Grade, transaction.ID,
            item.UniformName, item.Size, item.Quantity, item.UnitPrice, item.Notes,
        )
        if err != nil {
            tx.Rollback()
            return err
        }
    }
    return tx.Commit()
}

// Ambil transaksi biasa (nama_murid & kelas kosong)
func (r *TransactionRepository) GetNormalTransactions() ([]models.Transaksi, error) {
    query := `
        SELECT t.id, t.customer_id, c.name AS customer_name, t.transaction_date, t.payment_date, t.status, t.total_price, t.notes, t.created_at
        FROM transactions t
        JOIN customers c ON t.customer_id = c.id
        WHERE (t.nama_murid IS NULL OR t.nama_murid = '')
          AND (t.kelas IS NULL OR t.kelas = '')
        ORDER BY t.created_at DESC`
    rows, err := r.DB.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var transactions []models.Transaksi
    for rows.Next() {
        var t models.Transaksi
        err := rows.Scan(
            &t.ID, &t.CustomerID, &t.Customer_name,
            &t.Transaksidate, &t.Paymentdate, &t.Status,
            &t.Total, &t.Notes, &t.CreatedAt,
        )
        if err != nil {
            return nil, err
        }
        transactions = append(transactions, t)
    }
    return transactions, nil
}