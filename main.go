package main

import (
    "database/sql"
    "konveksi-app/handlers"
    "konveksi-app/repositories"
    "log"
    "net/http"

    _ "github.com/go-sql-driver/mysql"
    "github.com/gorilla/mux"
)

// sessions is a simple in-memory session store (for demo purposes only)
// var sessions = make(map[string]struct{})

func authMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Skip authentication untuk login page dan static assets
        if r.URL.Path == "/" || r.URL.Path == "/login" || 
           r.URL.Path == "/api/auth/login" || r.URL.Path == "/api/auth/logout" ||
           (len(r.URL.Path) >= 8 && r.URL.Path[:8] == "/assets/") {
            next.ServeHTTP(w, r)
            return
        }

        // Check session/authentication
        session, err := r.Cookie("session")
        if err != nil || session.Value == "" {
            http.Redirect(w, r, "/login", http.StatusSeeOther)
            return
        }
        
        // Validate session exists menggunakan handler's session store
        if !handlers.ValidateSession(session.Value) {
            log.Printf("Invalid session: %s for path: %s", session.Value, r.URL.Path)
            http.Redirect(w, r, "/login", http.StatusSeeOther)
            return
        }

        log.Printf("Valid session found for path: %s", r.URL.Path)
        next.ServeHTTP(w, r)
    })
}

func main() {
    // Initialize database connection
    db, err := sql.Open("mysql", "root@tcp(localhost:3306)/konveksi_bude")
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }
    defer db.Close()

    // Test database connection
    if err = db.Ping(); err != nil {
        log.Fatal("Failed to ping database:", err)
    }
    log.Println("Successfully connected to database")

    // Initialize repositories
    customerRepo := &repositories.CustomerRepository{DB: db}
    transactionRepo := &repositories.TransactionRepository{DB: db}
    userRepo := &repositories.UserRepository{DB: db} // Tambah user repo

    // Initialize handlers
    customerHandler := &handlers.CustomerHandler{Repo: customerRepo}
    transactionHandler := &handlers.TransactionHandler{Repo: transactionRepo}
    dashboardHandler := &handlers.DashboardHandler{DB: db}
    userHandler := &handlers.UserHandler{Repo: userRepo, DB: db} // Tambah user handler

    // Setup router
    r := mux.NewRouter()

    // Static files (CSS, JS, images)
    r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("assets/"))))

    // ==============================================
    // AUTHENTICATION ROUTES (tanpa middleware)
    // ==============================================

    // Root route - ke login
    r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "login.html")
    }).Methods("GET")

    // Login page
    r.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "login.html")
    }).Methods("GET")

    // Login API
    r.HandleFunc("/api/auth/login", userHandler.LoginAPI).Methods("POST")
    r.HandleFunc("/api/auth/logout", userHandler.LogoutAPI).Methods("POST")

    // ==============================================
    // PROTECTED ROUTES (dengan middleware)
    // ==============================================
    protected := r.PathPrefix("").Subrouter()
    protected.Use(authMiddleware)

    // Dashboard page
    protected.HandleFunc("/dashboard", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "index.html")
    }).Methods("GET")

    // Dashboard API routes
    protected.HandleFunc("/api/dashboard/stats", dashboardHandler.GetDashboardStats).Methods("GET")
    protected.HandleFunc("/api/dashboard/notifications", dashboardHandler.GetNotifications).Methods("GET")

    // Customer routes
    protected.HandleFunc("/kelolapelanggan", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "kelolapelanggan.html")
    }).Methods("GET")

    protected.HandleFunc("/api/customers", customerHandler.GetAllCustomers).Methods("GET")
    protected.HandleFunc("/api/customers", customerHandler.CreateCustomer).Methods("POST")
    protected.HandleFunc("/api/customers/{id}", customerHandler.GetCustomer).Methods("GET")
    protected.HandleFunc("/api/customers/{id}", customerHandler.UpdateCustomer).Methods("PUT")
    protected.HandleFunc("/api/customers/{id}", customerHandler.DeleteCustomer).Methods("DELETE")

    protected.HandleFunc("/tambahpelanggan", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "tambahpelanggan.html")
    }).Methods("GET")

    protected.HandleFunc("/edit-customer/{id}", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "editcust.html")
    }).Methods("GET")

    protected.HandleFunc("/detailpelanggan/{id}", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "detailpelanggan.html")
    }).Methods("GET")

    // Customer uniform routes
    protected.HandleFunc("/api/customers/{id}/uniforms", customerHandler.GetUniformsByCustomerID).Methods("GET")
    protected.HandleFunc("/api/customers/{id}/uniforms", customerHandler.AddCustomerUniform).Methods("POST")
    protected.HandleFunc("/api/customer-uniforms/{id}", customerHandler.GetCustomerUniform).Methods("GET")
    protected.HandleFunc("/api/customer-uniforms/{id}", customerHandler.UpdateCustomerUniform).Methods("PUT")
    protected.HandleFunc("/api/customer-uniforms/{id}", customerHandler.DeleteCustomerUniform).Methods("DELETE")
    protected.HandleFunc("/api/customer-uniforms/{id}/price-history", customerHandler.GetUniformPriceHistory).Methods("GET")

    // Transaction routes
    protected.HandleFunc("/kelolatransaksi", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "kelolatransaksi.html")
    }).Methods("GET")

    protected.HandleFunc("/formpesanan", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "formpesanan.html")
    }).Methods("GET")

    protected.HandleFunc("/formpesananperitem", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "formpesananperitem.html")
    }).Methods("GET")

    protected.HandleFunc("/detailpesananperitem", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "detailpesananperitem.html")
    }).Methods("GET")

    protected.HandleFunc("/detailpesanan", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "detailpesanan.html")
    }).Methods("GET")

    protected.HandleFunc("/detailpesananperitem/{id}", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "detailpesananperitem.html")
    }).Methods("GET")

    protected.HandleFunc("/detailpesanan/{id}", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "detailpesanan.html")
    }).Methods("GET")

    // Transaction API routes
    protected.HandleFunc("/api/transactions/all", transactionHandler.GetAllTransactions).Methods("GET")
    protected.HandleFunc("/api/transactions/normal/{id}", transactionHandler.GetByIDNormal).Methods("GET")
    protected.HandleFunc("/api/transactions/student/{id}", transactionHandler.GetByIDStudentOrder).Methods("GET")
    protected.HandleFunc("/api/transactions/{id}/customer-uniforms", transactionHandler.GetCustomerUniformsByTransactionID).Methods("GET")
    protected.HandleFunc("/api/transactions", transactionHandler.CreateTransaction).Methods("POST")
    protected.HandleFunc("/api/transactions/student", transactionHandler.CreateStudentOrder).Methods("POST")
    protected.HandleFunc("/api/transactions/{id}/status", transactionHandler.UpdateStatus).Methods("PUT")
    protected.HandleFunc("/api/customers/{customerID}/transactions", transactionHandler.GetCustomerTransactions).Methods("GET")
    protected.HandleFunc("/api/transactions/{transactionID}/status", transactionHandler.UpdateTransactionStatus).Methods("PUT")
    protected.HandleFunc("/api/transactions/{id}/print-kuitansi", transactionHandler.PrintKuitansi).Methods("GET")
    protected.HandleFunc("/api/transactions/{id}/print-kuitansi-biasa", transactionHandler.PrintKuitansibiasa).Methods("GET")
    protected.HandleFunc("/api/customers/list", customerHandler.GetAllCustomers).Methods("GET")
    protected.HandleFunc("/api/student-order-items/{id}", transactionHandler.UpdateStudentOrderItem).Methods("PUT")
    protected.HandleFunc("/api/order-items/{id}", transactionHandler.UpdateNormalOrderItem).Methods("PUT")
    protected.HandleFunc("/api/transactions/{id}/header", transactionHandler.UpdateTransactionHeader).Methods("PUT")
    protected.HandleFunc("/api/transactions/normal/{id}/items", transactionHandler.UpdateOrderItemsNormal).Methods("PUT")
    protected.HandleFunc("/api/transactions/student/{id}/items", transactionHandler.UpdateOrderItemsStudent).Methods("PUT")

    // CORS middleware
    r.Use(func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            w.Header().Set("Access-Control-Allow-Origin", "*")
            w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
            w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

            if r.Method == "OPTIONS" {
                w.WriteHeader(http.StatusOK)
                return
            }

            next.ServeHTTP(w, r)
        })
    })

    // Start server
    log.Println("==============================================")
    log.Println("🚀 DiOlif Fashion Management System")
    log.Println("📍 Server running on: http://localhost:8080")
    log.Println("🔑 Login: http://localhost:8080/")
    log.Println("🏠 Dashboard: http://localhost:8080/dashboard")
    log.Println("👥 Pelanggan: http://localhost:8080/kelolapelanggan")
    log.Println("🛒 Transaksi: http://localhost:8080/kelolatransaksi")
    log.Println("==============================================")

    log.Fatal(http.ListenAndServe(":8080", r))
}