package repositories

import (
	"database/sql"
)

type DashboardStats struct {
	OverduePayments     int
	OverdueTransactions int
	Allpaymentsdone     int
	Allpaymentspending   int
	Allpaymentscancelled int
	ReminderPayments     int
	ReminderTransactions int
}

func GetDashboardStats(db *sql.DB) (DashboardStats, error) {
	var stats DashboardStats

	err := db.QueryRow(`
		SELECT COUNT(*) FROM transactions
		WHERE status != 'paid'
		AND payment_date IS NOT NULL
		AND payment_date < CURDATE()
	`).Scan(&stats.OverduePayments)
	if err != nil {
		return stats, err
	}

	err = db.QueryRow(`
		SELECT COUNT(*) FROM transactions
		WHERE status = 'pending'
		AND transaction_date IS NOT NULL
		AND transaction_date < CURDATE()
	`).Scan(&stats.OverdueTransactions)
	if err != nil {
		return stats, err
	}

	err = db.QueryRow(`
		SELECT COUNT(*) AS total_paid_transactions
		FROM transactions
		WHERE status = 'paid'
		`).Scan(&stats.Allpaymentsdone)
	if err != nil {
		return stats, err
	}

	err = db.QueryRow(`
		SELECT COUNT(*) AS total_pending_transactions
		FROM transactions
		WHERE status = 'pending'
		`).Scan(&stats.Allpaymentspending)
	if err != nil {
		return stats, err
	}

	err = db.QueryRow(`
		SELECT COUNT(*) AS total_cancelled_transactions
		FROM transactions
		WHERE status = 'cancelled'
		`).Scan(&stats.Allpaymentscancelled)
	if err != nil {
		return stats, err
	}
	err = db.QueryRow(`
		SELECT COUNT(*) AS overdue_reminder_count
		FROM transactions
		WHERE status != 'paid'
		AND payment_date IS NOT NULL
		AND DATE_SUB(payment_date, INTERVAL 2 DAY) <= CURDATE()
		AND CURDATE() < payment_date
	`).Scan(&stats.ReminderPayments)
	if err != nil {
		return stats, err
	}

	err = db.QueryRow(`
		SELECT COUNT(*) AS task_reminder_count
		FROM transactions
		WHERE status = 'pending'
		AND transaction_date IS NOT NULL
		AND DATE_SUB(transaction_date, INTERVAL 2 DAY) <= CURDATE()
		AND CURDATE() < transaction_date
	`).Scan(&stats.ReminderTransactions)
	if err != nil {
		return stats, err
	}

	return stats, nil
}


