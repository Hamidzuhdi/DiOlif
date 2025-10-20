package handlers

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "konveksi-app/repositories"
    "log"
    "net/http"
    "time"
)

type DashboardHandler struct {
    DB *sql.DB
}

// GetDashboardStats - API endpoint untuk mendapatkan statistik dashboard
func (h *DashboardHandler) GetDashboardStats(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    
    // Log request
    log.Printf("Dashboard stats requested from: %s", r.RemoteAddr)
    
    stats, err := repositories.GetDashboardStats(h.DB)
    if err != nil {
        log.Printf("Error getting dashboard stats: %v", err)
        w.WriteHeader(http.StatusInternalServerError)
        errorResponse := map[string]interface{}{
            "error": "Gagal mengambil data dashboard",
            "message": err.Error(),
        }
        json.NewEncoder(w).Encode(errorResponse)
        return
    }

    // Structure response untuk API
    response := map[string]interface{}{
        "success": true,
        "data": map[string]interface{}{
            "total_customers":     0, // Will be populated by frontend
            "total_transactions":  stats.Allpaymentsdone + stats.Allpaymentspending + stats.Allpaymentscancelled,
            "pending_orders":      stats.Allpaymentspending,
            "paid_orders":         stats.Allpaymentsdone,
            "cancelled_orders":    stats.Allpaymentscancelled,
            "overdue_payments":    stats.OverduePayments,
            "overdue_transactions": stats.OverdueTransactions,
            "reminder_payments":   stats.ReminderPayments,
            "reminder_transactions": stats.ReminderTransactions,
            "total_revenue":       0, // Will be calculated by frontend
        },
    }

    log.Printf("Dashboard stats response: %+v", response)
    json.NewEncoder(w).Encode(response)
}

// GetNotifications - API endpoint untuk mendapatkan notifikasi
func (h *DashboardHandler) GetNotifications(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    
    // Log request
    log.Printf("Notifications requested from: %s", r.RemoteAddr)
    
    stats, err := repositories.GetDashboardStats(h.DB)
    if err != nil {
        log.Printf("Error getting notifications: %v", err)
        w.WriteHeader(http.StatusInternalServerError)
        errorResponse := map[string]interface{}{
            "error": "Gagal mengambil notifikasi",
            "message": err.Error(),
        }
        json.NewEncoder(w).Encode(errorResponse)
        return
    }

    notifications := []map[string]interface{}{}

    // Overdue Payments Notification
    if stats.OverduePayments > 0 {
        notifications = append(notifications, map[string]interface{}{
            "type":        "overdue_payment",
            "title":       "Pembayaran Terlambat",
            "message":     fmt.Sprintf("%d transaksi dengan pembayaran yang sudah terlambat", stats.OverduePayments),
            "count":       stats.OverduePayments,
            "icon":        "ph-currency-circle-dollar",
            "color":       "danger",
            "action_url":  "/kelolatransaksi?filter=overdue_payment",
            "created_at":  time.Now().Format("2006-01-02 15:04:05"),
        })
    }

    // Overdue Transactions Notification
    if stats.OverdueTransactions > 0 {
        notifications = append(notifications, map[string]interface{}{
            "type":        "overdue_transaction",
            "title":       "Transaksi Terlambat",
            "message":     fmt.Sprintf("%d transaksi pending yang sudah melewati target tanggal", stats.OverdueTransactions),
            "count":       stats.OverdueTransactions,
            "icon":        "ph-clock",
            "color":       "warning",
            "action_url":  "/kelolatransaksi?filter=overdue_transaction",
            "created_at":  time.Now().Format("2006-01-02 15:04:05"),
        })
    }

    // Reminder Payments Notification
    if stats.ReminderPayments > 0 {
        notifications = append(notifications, map[string]interface{}{
            "type":        "reminder_payment",
            "title":       "Reminder Pembayaran",
            "message":     fmt.Sprintf("%d transaksi mendekati deadline pembayaran", stats.ReminderPayments),
            "count":       stats.ReminderPayments,
            "icon":        "ph-bell",
            "color":       "info",
            "action_url":  "/kelolatransaksi?filter=reminder_payment",
            "created_at":  time.Now().Format("2006-01-02 15:04:05"),
        })
    }

    // Reminder Transactions Notification
    if stats.ReminderTransactions > 0 {
        notifications = append(notifications, map[string]interface{}{
            "type":        "reminder_transaction",
            "title":       "Reminder Pengerjaan",
            "message":     fmt.Sprintf("%d transaksi mendekati deadline pengerjaan", stats.ReminderTransactions),
            "count":       stats.ReminderTransactions,
            "icon":        "ph-wrench",
            "color":       "primary",
            "action_url":  "/kelolatransaksi?filter=reminder_transaction",
            "created_at":  time.Now().Format("2006-01-02 15:04:05"),
        })
    }

    response := map[string]interface{}{
        "success": true,
        "data": map[string]interface{}{
            "notifications": notifications,
            "total_count":   len(notifications),
            "has_notifications": len(notifications) > 0,
        },
    }

    log.Printf("Notifications response: %+v", response)
    json.NewEncoder(w).Encode(response)
}