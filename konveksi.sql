CREATE DATABASE konveksi_bude;

USE konveksi_bude;

CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(100) NOT NULL,
    password varchar(12) not null,
    contact VARCHAR(100),
    address TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- Tabel Customer
CREATE TABLE customers (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    type ENUM('TK', 'SD', 'SMP', 'Kelompok Tadarus', 'Lainnya') NOT NULL,
    contact VARCHAR(100),
    address TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE customer_uniforms (
    id INT AUTO_INCREMENT PRIMARY KEY,
    customer_id INT NOT NULL,
    uniform_name VARCHAR(100) NOT NULL,
    size VARCHAR(20) NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (customer_id) REFERENCES customers(id),
    UNIQUE KEY uniq_cust_uniform_size (customer_id, uniform_name, size)
);

CREATE TABLE customer_uniform_price_history (
    id INT AUTO_INCREMENT PRIMARY KEY,
    customer_uniform_id INT NOT NULL,
    old_price DECIMAL(10,2) NOT NULL,
    changed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (customer_uniform_id) REFERENCES customer_uniforms(id) ON DELETE CASCADE
);

-- Transaksi
CREATE TABLE transactions (
    id INT AUTO_INCREMENT PRIMARY KEY,
    customer_id INT NOT NULL,
    transaction_date DATE NOT NULL,
    payment_date DATE,
    status ENUM('pending', 'paid', 'cancelled') DEFAULT 'pending',
    total_price DECIMAL(12,2),
    notes TEXT,
    nama_murid varchar(30) ,
    kelas varchar(10),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (customer_id) REFERENCES customers(id)
);

CREATE TABLE student_order_items (
    id INT AUTO_INCREMENT PRIMARY KEY,
    customer_id INT NOT NULL,         -- SD mana
    student_name VARCHAR(100) NOT NULL,
    grade VARCHAR(10),                -- Kelas
    transaction_id INT NOT NULL,
    uniform_name VARCHAR(100) NOT NULL,
    size VARCHAR(20) NOT NULL,
    quantity INT NOT NULL DEFAULT 1,
    unit_price DECIMAL(10,2) NOT NULL,
    subtotal DECIMAL(12,2) GENERATED ALWAYS AS (quantity * unit_price) STORED,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE,
    FOREIGN KEY (transaction_id) REFERENCES transactions(id) ON DELETE CASCADE
);

-- Detail Pesanan
CREATE TABLE order_items (
    id INT AUTO_INCREMENT PRIMARY KEY,
    transaction_id INT NOT NULL,
    uniform_name VARCHAR(100) NOT NULL,
    size VARCHAR(20) NOT NULL,
    quantity INT NOT NULL,
    unit_price DECIMAL(10,2) NOT NULL,
    subtotal DECIMAL(12,2) GENERATED ALWAYS AS (quantity * unit_price) STORED,
    notes TEXT,
    FOREIGN KEY (transaction_id) REFERENCES transactions(id) ON DELETE CASCADE
);
SET GLOBAL sql_mode='NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION';


DROP VIEW IF EXISTS view_prices;

SELECT COUNT(*) AS overdue_payment_count
FROM transactions
WHERE status != 'paid'
  AND payment_date IS NOT NULL
  AND payment_date < CURDATE();


SELECT COUNT(*)
FROM transactions
WHERE status = 'pending'
  AND transaction_date IS NOT NULL
  AND transaction_date < CURDATE();


CREATE OR REPLACE VIEW view_prices AS
SELECT 
    cu.id,
    c.name AS customer_name,
    cu.uniform_name,
    cu.size,
    cu.price,
    cu.notes
FROM customer_uniforms cu
JOIN customers c ON cu.customer_id = c.id;


SELECT COUNT(*) AS total_paid_transactions
FROM transactions
WHERE status = 'paid';


SELECT t.*, c.name AS customer_name
FROM transactions t
JOIN customers c ON t.customer_id = c.id
WHERE t.status != 'paid'
  AND t.payment_date IS NOT NULL
  AND t.payment_date < CURDATE();


SELECT t.*, c.name AS customer_name
FROM transactions t
JOIN customers c ON t.customer_id = c.id
WHERE t.status = 'pending'
  AND t.transaction_date IS NOT NULL
  AND t.transaction_date < CURDATE();


SELECT COUNT(*) AS overdue_reminder_count
FROM transactions
WHERE status != 'paid'
  AND payment_date IS NOT NULL
  AND DATE_SUB(payment_date, INTERVAL 2 DAY) <= CURDATE()
  AND CURDATE() < payment_date;


SELECT COUNT(*) AS task_reminder_count
FROM transactions
WHERE status = 'pending'
  AND transaction_date IS NOT NULL
  AND DATE_SUB(transaction_date, INTERVAL 2 DAY) <= CURDATE()
  AND CURDATE() < transaction_date;


SELECT t.*, c.name AS customer_name
FROM transactions t
JOIN customers c ON t.customer_id = c.id
WHERE t.status != 'paid'
  AND t.payment_date IS NOT NULL
  AND CURDATE() >= DATE_SUB(t.payment_date, INTERVAL 2 DAY)
  AND CURDATE() < t.payment_date;


SELECT t.*, c.name AS customer_name
FROM transactions t
JOIN customers c ON t.customer_id = c.id
WHERE t.status != 'paid'
  AND t.transaction_date IS NOT NULL
  AND CURDATE() >= DATE_SUB(t.transaction_date, INTERVAL 2 DAY)
  AND CURDATE() < t.transaction_date;
  
  
SELECT t.*
FROM transactions t
WHERE EXISTS (
  SELECT 1 FROM order_items oi WHERE oi.transaction_id = t.id
)
ORDER BY t.id DESC;

SELECT t.*
FROM transactions t
WHERE EXISTS (
  SELECT 1 FROM student_order_items soi WHERE soi.transaction_id = t.id
)
ORDER BY t.id DESC;




