package repositories

import (
	"database/sql"
	"fmt"
	"konveksi-app/models"
)

type CustomerRepository struct {
	DB *sql.DB
}

func (r *CustomerRepository) CreateWithUniforms(customer *models.Customer, uniforms []models.CustomerUniform) error {
    tx, err := r.DB.Begin()
    if err != nil {
        return err
    }
    res, err := tx.Exec(
        "INSERT INTO customers (name, type, contact, address) VALUES (?, ?, ?, ?)",
        customer.Name, customer.Type, customer.Contact, customer.Address,
    )
    if err != nil {
        tx.Rollback()
        return err
    }
    customerID, _ := res.LastInsertId()
    customer.ID = int(customerID)

    // Validasi unik kombinasi uniform_name + size
    unique := map[string]bool{}
    for _, u := range uniforms {
        key := u.UniformName + "|" + u.Size
        if unique[key] {
            tx.Rollback()
            return fmt.Errorf("ukuran '%s' untuk seragam '%s' sudah ada", u.Size, u.UniformName)
        }
        unique[key] = true
        _, err := tx.Exec(
            "INSERT INTO customer_uniforms (customer_id, uniform_name, size, price, notes) VALUES (?, ?, ?, ?, ?)",
            customer.ID, u.UniformName, u.Size, u.Price, u.Notes,
        )
        if err != nil {
            tx.Rollback()
            return err
        }
    }
    return tx.Commit()
}

func (r *CustomerRepository) AddCustomerUniform(u *models.CustomerUniform) error {
    // Validasi unik kombinasi uniform_name + size untuk customer ini
    var count int
    err := r.DB.QueryRow(
        `SELECT COUNT(*) FROM customer_uniforms WHERE customer_id = ? AND uniform_name = ? AND size = ?`,
        u.CustomerID, u.UniformName, u.Size,
    ).Scan(&count)
    if err != nil {
        return err
    }
    if count > 0 {
        return fmt.Errorf("ukuran '%s' untuk seragam '%s' sudah ada", u.Size, u.UniformName)
    }
    _, err = r.DB.Exec(
        `INSERT INTO customer_uniforms (customer_id, uniform_name, size, price, notes) VALUES (?, ?, ?, ?, ?)`,
        u.CustomerID, u.UniformName, u.Size, u.Price, u.Notes,
    )
    return err
}

func (r *CustomerRepository) GetByID(id int) (*models.Customer, error) {
	query := `
		SELECT id, name, type, contact, address, created_at
		FROM customers WHERE id = ?`

	var customer models.Customer
	err := r.DB.QueryRow(query, id).Scan(
		&customer.ID, &customer.Name, &customer.Type,
		&customer.Contact, &customer.Address,
		&customer.CreatedAt,
	)

	return &customer, err
}

func (r *CustomerRepository) GetUniformsByCustomerID(customerID int) ([]models.CustomerUniform, error) {
    rows, err := r.DB.Query(
        "SELECT id, customer_id, uniform_name, size, price, notes, created_at FROM customer_uniforms WHERE customer_id = ?",
        customerID,
    )
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    var uniforms []models.CustomerUniform
    for rows.Next() {
        var u models.CustomerUniform
        if err := rows.Scan(&u.ID, &u.CustomerID, &u.UniformName, &u.Size, &u.Price, &u.Notes, &u.CreatedAt); err != nil {
            return nil, err
        }
        uniforms = append(uniforms, u)
    }
    return uniforms, nil
}

func (r *CustomerRepository) GetAll() ([]models.Customer, error) {
    query := `
        SELECT id, name, type, contact, address, created_at
        FROM customers ORDER BY name`  // Koma setelah created_at dihapus

    rows, err := r.DB.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var customers []models.Customer
    for rows.Next() {
        var c models.Customer
        err := rows.Scan(
            &c.ID, &c.Name, &c.Type,
            &c.Contact, &c.Address,
            &c.CreatedAt,
        )
        if err != nil {
            return nil, err
        }
        customers = append(customers, c)
    }

    return customers, nil
}

func (r *CustomerRepository) Update(customer *models.Customer) error {
    query := `
        UPDATE customers 
        SET name = ?, type = ?, contact = ?, address = ? 
        WHERE id = ?`  // Koma setelah address dihapus

    _, err := r.DB.Exec(query,
        customer.Name, customer.Type,
        customer.Contact, customer.Address,
        customer.ID,
    )

    return err
}

func (r *CustomerRepository) GetCustomerUniformByID(id int) (*models.CustomerUniform, error) {
    var u models.CustomerUniform
    err := r.DB.QueryRow(
        "SELECT id, customer_id, uniform_name, size, price, notes, created_at FROM customer_uniforms WHERE id = ?",
        id,
    ).Scan(&u.ID, &u.CustomerID, &u.UniformName, &u.Size, &u.Price, &u.Notes, &u.CreatedAt)
    if err != nil {
        return nil, err
    }
    return &u, nil
}

func (r *CustomerRepository) UpdateCustomerUniformWithHistory(id int, uniformName, size string, price float64, notes string) error {
    tx, err := r.DB.Begin()
    if err != nil {
        return err
    }

    // Ambil harga lama
    var oldPrice float64
    err = tx.QueryRow("SELECT price FROM customer_uniforms WHERE id = ?", id).Scan(&oldPrice)
    if err != nil {
        tx.Rollback()
        return err
    }

    // Jika harga berubah, catat ke history
    if oldPrice != price {
        _, err = tx.Exec("INSERT INTO customer_uniform_price_history (customer_uniform_id, old_price) VALUES (?, ?)", id, oldPrice)
        if err != nil {
            tx.Rollback()
            return err
        }
    }

    // Update uniform
    _, err = tx.Exec("UPDATE customer_uniforms SET uniform_name=?, size=?, price=?, notes=? WHERE id=?", uniformName, size, price, notes, id)
    if err != nil {
        tx.Rollback()
        return err
    }

    return tx.Commit()
}

// Ambil riwayat harga
func (r *CustomerRepository) GetUniformPriceHistory(uniformID int) ([]struct {
    OldPrice  float64
    ChangedAt string
}, error) {
    rows, err := r.DB.Query("SELECT old_price, changed_at FROM customer_uniform_price_history WHERE customer_uniform_id = ? ORDER BY changed_at ASC", uniformID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    var history []struct {
        OldPrice  float64
        ChangedAt string
    }
    for rows.Next() {
        var h struct {
            OldPrice  float64
            ChangedAt string
        }
        if err := rows.Scan(&h.OldPrice, &h.ChangedAt); err != nil {
            return nil, err
        }
        history = append(history, h)
    }
    return history, nil
}

func (r *CustomerRepository) DeleteCustomerUniform(id int) error {
    _, err := r.DB.Exec("DELETE FROM customer_uniforms WHERE id = ?", id)
    return err
}
func (r *CustomerRepository) Delete(id int) error {
	_, err := r.DB.Exec("DELETE FROM customers WHERE id = ?", id)
	return err
}