package main

import (
	"database/sql"
	"konveksi-app/handlers"
	"konveksi-app/repositories"
	"log"
	"net/http"
	"html/template"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {
	// Initialize database connection
	db, err := sql.Open("mysql", "root@tcp(localhost:3306)/konveksi_bude")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Test database connection
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	// Initialize repositories
	customerRepo := &repositories.CustomerRepository{DB: db}
	// uniformRepo := &repositories.UniformRepository{DB: db}
	// priceRepo := &repositories.PriceRepository{DB: db}
	transactionRepo := &repositories.TransactionRepository{DB: db}
	userRepo := &repositories.UserRepository{DB: db}

	// Initialize handlers
	customerHandler := &handlers.CustomerHandler{Repo: customerRepo}
	// uniformHandler := &handlers.UniformHandler{Repo: uniformRepo}
	// priceHandler := &handlers.PriceHandler{Repo: priceRepo}
	transactionHandler := &handlers.TransactionHandler{Repo: transactionRepo}
	userHandler := &handlers.UserHandler{Repo: userRepo, DB: db}

	tmpl := template.Must(template.ParseGlob("templates/*.html"))

	// Setup router
	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "templates/index.html")
	})
	r.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "templates/index.html")
	})

	r.HandleFunc("/login", userHandler.LoginHandler).Methods("POST")

	// Dashboard setelah login
	
	dashboardHandler := handlers.DashboardHandler{DB: db, Tmpl: tmpl}
	r.HandleFunc("/dashboard", dashboardHandler.HandleDashboard).Methods("GET")

	// Overdue Payment
	r.HandleFunc("/overdue-payment", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/overduepayment.html")
	}).Methods("GET")
	r.HandleFunc("/api/overdue-payment", transactionHandler.GetOverduePaymentTransactions).Methods("GET")

	// Overdue Transaksi
	r.HandleFunc("/overdue-transaksi", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/overduetransaksi.html")
	}).Methods("GET")
	r.HandleFunc("/api/overdue-transaksi", transactionHandler.GetOverdueTransactions).Methods("GET")
	r.HandleFunc("/edit-detail-transaksi/{id}", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/editdetailtransaksi.html")
	}).Methods("GET")
	// Untuk transaksi biasa
	r.HandleFunc("/api/transactions/normal/{id}", transactionHandler.GetByIDNormal).Methods("GET")
	r.HandleFunc("/api/transactions/normal/{id}/order-items", transactionHandler.UpdateOrderItemsNormal).Methods("PUT")

	// Untuk repeat order
	// r.HandleFunc("/api/transactions/repeat/{id}", transactionHandler.GetByIDRepeatOrder).Methods("GET")
	// r.HandleFunc("/api/transactions/repeat/{id}/order-items", transactionHandler.UpdateOrderItemsRepeat).Methods("PUT")

	r.HandleFunc("/api/users", userHandler.CreateUser).Methods("POST")

	// Customer routes
	r.HandleFunc("/api/customers", customerHandler.CreateCustomer).Methods("POST")
	r.HandleFunc("/add-customer", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/addcust.html")
	}).Methods("GET")	
	r.HandleFunc("/api/customers", customerHandler.GetAllCustomers).Methods("GET")
	r.HandleFunc("/customers", func(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "templates/customers.html")
	}).Methods("GET")
	r.HandleFunc("/api/customers/{id}", customerHandler.GetCustomer).Methods("GET")
	r.HandleFunc("/api/customers/{id}", customerHandler.UpdateCustomer).Methods("PUT")
	r.HandleFunc("/edit-customer/{id}", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/editcust.html")
	}).Methods("GET")
	r.HandleFunc("/api/customers/{id}", customerHandler.DeleteCustomer).Methods("DELETE")
	// Untuk update customer uniform
	r.HandleFunc("/api/customer-uniforms/{id}", customerHandler.GetCustomerUniform).Methods("GET")
	r.HandleFunc("/api/customer-uniforms/{id}", customerHandler.UpdateCustomerUniform).Methods("PUT")
	r.HandleFunc("/api/customer-uniforms/{id}/price-history", customerHandler.GetUniformPriceHistory).Methods("GET")
	r.HandleFunc("/historyharga/{id}", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/historyharga.html")
	}).Methods("GET")
	r.HandleFunc("/api/transactions/{id}/kuitansi", transactionHandler.PrintKuitansi).Methods("GET")
	r.HandleFunc("/api/transactions/normal/{id}/kuitansi", transactionHandler.PrintKuitansibiasa).Methods("GET")
	r.HandleFunc("/edit-custuniform/{id}", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/editcustuniform.html")
	}).Methods("GET")
	// Halaman detail seragam per customer
	r.HandleFunc("/detailuniformcust/{id}", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/detailuniformcust.html")
	}).Methods("GET")

	// API: Get all uniforms milik customer
	r.HandleFunc("/api/customers/{id}/uniforms", customerHandler.GetUniformsByCustomerID).Methods("GET")

	// API: Delete customer_uniform
	r.HandleFunc("/api/customer-uniforms/{id}", customerHandler.DeleteCustomerUniform).Methods("DELETE")
	r.HandleFunc("/api/customers/{id}/uniforms", customerHandler.AddCustomerUniform).Methods("POST")
	r.HandleFunc("/add-custuniform/{id}", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/addcustuniform.html")
	}).Methods("GET")


	r.HandleFunc("/repeat-order", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/repeat_order.html")
	}).Methods("GET")

		// Halaman tambah repeat order
	r.HandleFunc("/repeat-order-form", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/repeat_order_form.html")
	}).Methods("GET")
	// API repeat order
	// r.HandleFunc("/api/repeat-orders", transactionHandler.CreateRepeatOrder).Methods("POST")


	// Transaction routes
	r.HandleFunc("/api/transactions", transactionHandler.CreateTransaction).Methods("POST")
	r.HandleFunc("/transaction-form", func(w http.ResponseWriter, r *http.Request) {
    	http.ServeFile(w, r, "templates/transaction_form.html")
	})
	r.HandleFunc("/api/transactions/student", transactionHandler.GetAllTransactionStudent).Methods("GET")
	r.HandleFunc("/api/transactions/order", transactionHandler.GetAllTransactionOrder).Methods("GET")
	r.HandleFunc("/transactions", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/transaksi.html")
	}).Methods("GET")
	// Untuk detail repeat order
	r.HandleFunc("/transaction-detail-repeat/{id}", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/detail_repeatorder.html")
	}).Methods("GET")
	r.HandleFunc("/api/transactions/normal", transactionHandler.GetNormalTransactions).Methods("GET")
	r.HandleFunc("/transaction-detail/{id}", func(w http.ResponseWriter, r *http.Request) {
    	http.ServeFile(w, r, "templates/detail_transaksi.html")
	}).Methods("GET")
	r.HandleFunc("/api/transactions/{id}/status", transactionHandler.UpdateStatus).Methods("PATCH")
	r.HandleFunc("/api/transactions/{id}/status", transactionHandler.UpdateStatus).Methods("PUT")
	r.HandleFunc("/remind-transaksi", func(w http.ResponseWriter, r *http.Request) {
    	http.ServeFile(w, r, "templates/remindtransaksi.html")
	}).Methods("GET")
	r.HandleFunc("/api/remind-transaksi", transactionHandler.GetRemindTransactions).Methods("GET")
	r.HandleFunc("/remind-payment", func(w http.ResponseWriter, r *http.Request) {
    	http.ServeFile(w, r, "templates/remindpayment.html")
	}).Methods("GET")
	r.HandleFunc("/editrepeatorder/{id}", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/editrepeatorder.html")
	}).Methods("GET")
	r.HandleFunc("/api/remind-payment", transactionHandler.GetRemindPayments).Methods("GET")
	r.HandleFunc("/edit-transaksi/{id}", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/edittransaksi.html")
	}).Methods("GET")
	r.HandleFunc("/api/transactions/{id}", transactionHandler.UpdateTransaction).Methods("PUT")
	r.HandleFunc("/editdetailtransaksi/{id}", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/editdetailtransaksi.html")
	}).Methods("GET")

	r.HandleFunc("/editdetailrepeat/{id}", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/editdetailrepeat.html")
	}).Methods("GET")


	r.HandleFunc("/repeat-order", func(w http.ResponseWriter, r *http.Request) {
    	http.ServeFile(w, r, "templates/repeat_order.html")
	}).Methods("GET")
	// r.HandleFunc("/api/transactions/repeat", transactionHandler.GetRepeatOrders).Methods("GET")

	r.HandleFunc("/api/transactions/student", transactionHandler.CreateStudentOrder).Methods("POST")
	r.HandleFunc("/api/transactions/student/{id}", transactionHandler.GetByIDStudentOrder).Methods("GET")
	r.HandleFunc("/api/transactions/student/{id}/order-items", transactionHandler.UpdateOrderItemsStudent).Methods("PUT")
	// Handler di transaction_handler.go
	r.HandleFunc("/api/transactions/student", transactionHandler.GetAllStudentOrders).Methods("GET")
	// Start server
	log.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}