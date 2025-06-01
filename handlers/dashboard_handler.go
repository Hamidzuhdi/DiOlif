package handlers

import (
	"database/sql"
	"html/template"
	"konveksi-app/repositories"
	"net/http"
)

type DashboardHandler struct {
	DB   *sql.DB
	Tmpl *template.Template
}

func (h *DashboardHandler) HandleDashboard(w http.ResponseWriter, r *http.Request) {
	stats, err := repositories.GetDashboardStats(h.DB)
	if err != nil {
		http.Error(w, "Gagal mengambil data dashboard", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"OverduePayments":     stats.OverduePayments,
		"OverdueTransactions": stats.OverdueTransactions,
		"AllPaymentsDone":     stats.Allpaymentsdone,
		"AllPaymentsPending":  stats.Allpaymentspending,
		"AllPaymentsCancelled": stats.Allpaymentscancelled,
		"ReminderBayar":    stats.ReminderPayments,
		"ReminderKerjakan": stats.ReminderTransactions,
	}

	h.Tmpl.ExecuteTemplate(w, "dashboard.html", data)
}
