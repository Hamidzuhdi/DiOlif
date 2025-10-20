package handlers

import (
    "encoding/json"
    "fmt"
    "konveksi-app/models"
    "konveksi-app/repositories"
    "log"
    "net/http"
    "strconv"
    "regexp"
    "github.com/gorilla/mux"
    "os"
    "strings"
    "time"
)

type TransactionHandler struct {
    Repo *repositories.TransactionRepository
}

// Create normal transaction (item order)
func (h *TransactionHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
    log.Println("CreateTransaction handler called")
    
    var req struct {
        CustomerID      int    `json:"customer_id"`
        TransactionDate string `json:"transaction_date"`
        PaymentDate     string `json:"payment_date"`
        Status          string `json:"status"`
        Notes           string `json:"notes"`
        Items           []struct {
            UniformName string  `json:"uniform_name"`
            Size        string  `json:"size"`
            Quantity    int     `json:"quantity"`
            UnitPrice   float64 `json:"unit_price"`
            Notes       string  `json:"notes"`
        } `json:"items"`
    }

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        log.Printf("Error decoding request body: %v", err)
        http.Error(w, "Invalid JSON format", http.StatusBadRequest)
        return
    }

    log.Printf("Received transaction data: %+v", req)

    // Validate required fields
    if req.CustomerID == 0 {
        log.Println("Customer ID is missing")
        http.Error(w, "Customer ID is required", http.StatusBadRequest)
        return
    }

    if len(req.Items) == 0 {
        log.Println("No items provided")
        http.Error(w, "At least one item is required", http.StatusBadRequest)
        return
    }

    // Calculate total
    var total float64
    for _, item := range req.Items {
        total += item.UnitPrice * float64(item.Quantity)
    }

    // Create transaction object
    transaction := &models.Transaksi{
        CustomerID:    req.CustomerID,
        Transaksidate: req.TransactionDate,
        Paymentdate:   req.PaymentDate,
        Status:        req.Status,
        Total:         total,
        Notes:         req.Notes,
        Items:         make([]models.OrderItem, len(req.Items)),
    }

    // Map items
    for i, item := range req.Items {
        transaction.Items[i] = models.OrderItem{
            UniformName: item.UniformName,
            Size:        item.Size,
            Quantity:    item.Quantity,
            UnitPrice:   item.UnitPrice,
            Subtotal:    item.UnitPrice * float64(item.Quantity),
            Notes:       item.Notes,
        }
    }

    log.Printf("Creating transaction with total: %.2f", total)

    // Save to repository
    if err := h.Repo.Create(transaction); err != nil {
        log.Printf("Error creating transaction: %v", err)
        http.Error(w, "Failed to create transaction: "+err.Error(), http.StatusInternalServerError)
        return
    }

    log.Printf("Transaction created successfully with ID: %d", transaction.ID)

    // Return success response
    w.Header().Set("Content-Type", "application/json")
    response := map[string]interface{}{
        "id":      transaction.ID,
        "message": "Transaction created successfully",
        "total":   transaction.Total,
    }

    json.NewEncoder(w).Encode(response)
}

// Create student order transaction
func (h *TransactionHandler) CreateStudentOrder(w http.ResponseWriter, r *http.Request) {
    log.Println("CreateStudentOrder handler called")
    
    w.Header().Set("Content-Type", "application/json")
    
    var req struct {
        CustomerID      int    `json:"customer_id"`
        TransactionDate string `json:"transaction_date"`
        PaymentDate     string `json:"payment_date"`
        Status          string `json:"status"`
        Notes           string `json:"notes"`
        Items           []struct {
            StudentName string  `json:"student_name"`
            Grade       string  `json:"grade"`
            UniformName string  `json:"uniform_name"`
            Size        string  `json:"size"`
            Quantity    int     `json:"quantity"`
            UnitPrice   float64 `json:"unit_price"`
            Notes       string  `json:"notes"`
        } `json:"items"`
    }

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        log.Printf("Error decoding request: %v", err)
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    log.Printf("Received student order data: %+v", req)

    // Validate required fields
    if req.CustomerID == 0 {
        http.Error(w, "Customer ID is required", http.StatusBadRequest)
        return
    }

    if len(req.Items) == 0 {
        http.Error(w, "At least one item is required", http.StatusBadRequest)
        return
    }

    // Calculate total
    var total float64
    for _, item := range req.Items {
        total += item.UnitPrice * float64(item.Quantity)
    }

    // Create transaction object
    transaction := &models.Transaksi{
        CustomerID:    req.CustomerID,
        Transaksidate: req.TransactionDate,
        Paymentdate:   req.PaymentDate,
        Status:        req.Status,
        Total:         total,
        Notes:         req.Notes,
    }

    // Map student items
    var studentItems []models.StudentOrderItem
    for _, item := range req.Items {
        studentItems = append(studentItems, models.StudentOrderItem{
            CustomerID:  req.CustomerID,
            StudentName: item.StudentName,
            Grade:       item.Grade,
            UniformName: item.UniformName,
            Size:        item.Size,
            Quantity:    item.Quantity,
            UnitPrice:   item.UnitPrice,
            Subtotal:    item.UnitPrice * float64(item.Quantity),
            Notes:       item.Notes,
        })
    }

    log.Printf("Creating student order with total: %.2f and %d items", total, len(studentItems))

    // Save to repository
    if err := h.Repo.CreateStudentOrder(transaction, studentItems); err != nil {
        log.Printf("Error creating student order: %v", err)
        http.Error(w, "Failed to create transaction", http.StatusInternalServerError)
        return
    }

    log.Printf("Student order created successfully with ID: %d", transaction.ID)

    // Return success response
    response := map[string]interface{}{
        "id":      transaction.ID,
        "message": "Student order created successfully",
        "total":   transaction.Total,
    }

    json.NewEncoder(w).Encode(response)
}

// Get all transactions with filtering and pagination
func (h *TransactionHandler) GetAllTransactions(w http.ResponseWriter, r *http.Request) {
    results, err := h.Repo.GetAllTransactions()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    // Build response dengan field mapping yang explicit + HasStudentInfo
    var response []struct {
        ID              int     `json:"id"`
        CustomerID      int     `json:"customer_id"`
        CustomerName    string  `json:"customer_name"`
        TransactionDate string  `json:"transaction_date"`
        PaymentDate     string  `json:"payment_date"`
        Status          string  `json:"status"`
        TotalPrice      float64 `json:"total_price"`
        Notes           string  `json:"notes"`
        CreatedAt       string  `json:"created_at"`
        TransactionType string  `json:"transaction_type"`
        ItemCount       int     `json:"item_count"`
        HasStudentInfo  bool    `json:"has_student_info"`
    }
    
    for _, result := range results {
        itemCount, err := h.Repo.GetTransactionItemCount(result.ID, result.TransactionType)
        if err != nil {
            itemCount = 0
        }
        
        // Check if has student info (nama siswa dan kelas)
        hasStudentInfo := h.Repo.HasStudentInfo(result.ID)
        
        response = append(response, struct {
            ID              int     `json:"id"`
            CustomerID      int     `json:"customer_id"`
            CustomerName    string  `json:"customer_name"`
            TransactionDate string  `json:"transaction_date"`
            PaymentDate     string  `json:"payment_date"`
            Status          string  `json:"status"`
            TotalPrice      float64 `json:"total_price"`
            Notes           string  `json:"notes"`
            CreatedAt       string  `json:"created_at"`
            TransactionType string  `json:"transaction_type"`
            ItemCount       int     `json:"item_count"`
            HasStudentInfo  bool    `json:"has_student_info"`
        }{
            ID:              result.ID,
            CustomerID:      result.CustomerID,
            CustomerName:    result.Customer_name,
            TransactionDate: result.Transaksidate,
            PaymentDate:     result.Paymentdate,
            Status:          result.Status,
            TotalPrice:      result.Total,
            Notes:           result.Notes,
            CreatedAt:       result.CreatedAt,
            TransactionType: result.TransactionType,
            ItemCount:       itemCount,
            HasStudentInfo:  hasStudentInfo,
        })
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

// Get transactions by customer ID
func (h *TransactionHandler) GetCustomerTransactions(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    
    // Get customer ID from URL
    vars := mux.Vars(r)
    customerIDStr := vars["customerID"]
    customerID, err := strconv.Atoi(customerIDStr)
    if err != nil {
        http.Error(w, "Invalid customer ID", http.StatusBadRequest)
        return
    }
    
    // Get status filter from query param
    statusFilter := r.URL.Query().Get("status")
    
    // Get transactions from repository
    transactions, err := h.Repo.GetTransactionsByCustomerID(customerID, statusFilter)
    if err != nil {
        log.Printf("Error getting customer transactions: %v", err)
        http.Error(w, "Failed to get transactions", http.StatusInternalServerError)
        return
    }
    
    json.NewEncoder(w).Encode(transactions)
}

// Get single normal transaction by ID
func (h *TransactionHandler) GetByIDNormal(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    id, err := strconv.Atoi(params["id"])
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }
    
    trx, err := h.Repo.GetByIDNormal(id)
    if err != nil {
        http.Error(w, "Transaction not found", http.StatusNotFound)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(trx)
}

// Get single student order transaction by ID
func (h *TransactionHandler) GetByIDStudentOrder(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    id, err := strconv.Atoi(params["id"])
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }
    
    trx, items, err := h.Repo.GetByIDStudentOrder(id)
    if err != nil {
        http.Error(w, "Transaction not found", http.StatusNotFound)
        return
    }
    
    resp := struct {
        Transaction *models.Transaksi         `json:"transaction"`
        Items       []models.StudentOrderItem `json:"items"`
    }{
        Transaction: trx,
        Items:       items,
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(resp)
}

// Get customer uniforms by transaction ID
func (h *TransactionHandler) GetCustomerUniformsByTransactionID(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    transactionID, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid transaction ID", http.StatusBadRequest)
        return
    }

    uniforms, err := h.Repo.GetCustomerUniformsByTransactionID(transactionID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(uniforms)
}

// Update transaction status
func (h *TransactionHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    id, err := strconv.Atoi(params["id"])
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    var input struct {
        Status string `json:"status"`
    }
    if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    if err := h.Repo.UpdateStatus(id, input.Status); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}

// Update transaction status (alternative endpoint)
func (h *TransactionHandler) UpdateTransactionStatus(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    
    vars := mux.Vars(r)
    transactionIDStr := vars["transactionID"]
    transactionID, err := strconv.Atoi(transactionIDStr)
    if err != nil {
        http.Error(w, "Invalid transaction ID", http.StatusBadRequest)
        return
    }
    
    var req struct {
        Status string `json:"status"`
    }
    
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }
    
    // Validate status
    validStatuses := []string{"paid", "pending", "cancelled"}
    isValid := false
    for _, status := range validStatuses {
        if req.Status == status {
            isValid = true
            break
        }
    }
    
    if !isValid {
        http.Error(w, "Invalid status. Must be: paid, pending, or cancelled", http.StatusBadRequest)
        return
    }
    
    // Update status in repository
    if err := h.Repo.UpdateTransactionStatus(transactionID, req.Status); err != nil {
        log.Printf("Error updating transaction status: %v", err)
        http.Error(w, "Failed to update transaction status", http.StatusInternalServerError)
        return
    }
    
    response := map[string]interface{}{
        "message": "Transaction status updated successfully",
        "status":  req.Status,
    }
    
    json.NewEncoder(w).Encode(response)
}

// Update transaction header (for kelolatransaksi page)
func (h *TransactionHandler) UpdateTransactionHeader(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    id, err := strconv.Atoi(params["id"])
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    var req struct {
        TransactionDate string `json:"transaction_date"`
        PaymentDate     string `json:"payment_date"`
        Notes           string `json:"notes"`
        Status          string `json:"status,omitempty"` // Optional status update
    }

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }

    // Update transaction header
    err = h.Repo.UpdateTransaction(id, req.TransactionDate, req.PaymentDate, req.Notes)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Update status if provided
    if req.Status != "" {
        if err := h.Repo.UpdateStatus(id, req.Status); err != nil {
            log.Printf("Warning: Failed to update status: %v", err)
        }
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"message": "Transaction updated successfully"})
}

// Update transaction (legacy)
func (h *TransactionHandler) UpdateTransaction(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    id, err := strconv.Atoi(params["id"])
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }
    
    var req struct {
        TransactionDate string `json:"transaction_date"`
        PaymentDate     string `json:"payment_date"`
        Notes           string `json:"notes"`
    }
    
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }
    
    if err := h.Repo.UpdateTransaction(id, req.TransactionDate, req.PaymentDate, req.Notes); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.WriteHeader(http.StatusNoContent)
}

// Update bulk normal order items
func (h *TransactionHandler) UpdateOrderItemsNormal(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    id, err := strconv.Atoi(params["id"])
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }
    
    var req struct {
        Items []struct {
            UniformName string  `json:"uniform_name"`
            Size        string  `json:"size"`
            Quantity    int     `json:"quantity"`
            UnitPrice   float64 `json:"unit_price"`
            Notes       string  `json:"notes"`
        } `json:"items"`
    }
    
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }
    
    orderItems := make([]repositories.OrderItemUpdate, len(req.Items))
    for i, item := range req.Items {
        orderItems[i] = repositories.OrderItemUpdate{
            UniformName: item.UniformName,
            Size:        item.Size,
            Quantity:    item.Quantity,
            UnitPrice:   item.UnitPrice,
            Notes:       item.Notes,
        }
    }
    
    if err := h.Repo.UpdateOrderItemsNormal(id, orderItems); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.WriteHeader(http.StatusNoContent)
}

// Update bulk student order items
func (h *TransactionHandler) UpdateOrderItemsStudent(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    id, err := strconv.Atoi(params["id"])
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }
    
    var req struct {
        Items []struct {
            CustomerID  int    `json:"customer_id"`
            StudentName string `json:"student_name"`
            Grade       string `json:"grade"`
            UniformName string `json:"uniform_name"`
            Size        string `json:"size"`
            Quantity    int    `json:"quantity"`
            UnitPrice   float64 `json:"unit_price"`
            Notes       string `json:"notes"`
        } `json:"items"`
    }
    
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }
    
    var studentItems []models.StudentOrderItem
    for _, item := range req.Items {
        studentItems = append(studentItems, models.StudentOrderItem{
            CustomerID:  item.CustomerID,
            StudentName: item.StudentName,
            Grade:       item.Grade,
            UniformName: item.UniformName,
            Size:        item.Size,
            Quantity:    item.Quantity,
            UnitPrice:   item.UnitPrice,
            Notes:       item.Notes,
        })
    }
    
    if err := h.Repo.UpdateOrderItemsStudent(id, studentItems); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.WriteHeader(http.StatusNoContent)
}

// Update individual student order item (for detail page)
func (h *TransactionHandler) UpdateStudentOrderItem(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    id, err := strconv.Atoi(params["id"])
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    var req struct {
        StudentName   string  `json:"student_name"`
        Grade         string  `json:"grade"`
        UniformName   string  `json:"uniform_name"`
        Size          string  `json:"size"`
        Quantity      int     `json:"quantity"`
        UnitPrice     float64 `json:"unit_price"`
        Notes         string  `json:"notes"`
        TransactionID int     `json:"transaction_id"`
    }

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }

    // Update student order item
    err = h.Repo.UpdateStudentOrderItem(id, req.StudentName, req.Grade, req.UniformName, req.Size, req.Quantity, req.UnitPrice, req.Notes)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Recalculate transaction total if transaction_id provided
    if req.TransactionID > 0 {
        if err := h.Repo.RecalculateTransactionTotal(req.TransactionID); err != nil {
            log.Printf("Warning: Failed to recalculate transaction total: %v", err)
        }
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"message": "Student order item updated successfully"})
}

// Update individual normal order item (for detail page)  
func (h *TransactionHandler) UpdateNormalOrderItem(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    id, err := strconv.Atoi(params["id"])
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    var req struct {
        UniformName   string  `json:"uniform_name"`
        Size          string  `json:"size"`
        Quantity      int     `json:"quantity"`
        UnitPrice     float64 `json:"unit_price"`
        Notes         string  `json:"notes"`
        TransactionID int     `json:"transaction_id"`
    }

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }

    // Update normal order item
    err = h.Repo.UpdateNormalOrderItem(id, req.UniformName, req.Size, req.Quantity, req.UnitPrice, req.Notes)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Recalculate transaction total if transaction_id provided
    if req.TransactionID > 0 {
        if err := h.Repo.RecalculateTransactionTotal(req.TransactionID); err != nil {
            log.Printf("Warning: Failed to recalculate transaction total: %v", err)
        }
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"message": "Order item updated successfully"})
}

// Get overdue payment transactions
func (h *TransactionHandler) GetOverduePaymentTransactions(w http.ResponseWriter, r *http.Request) {
    transactions, err := h.Repo.GetOverduePaymentTransactions()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(transactions)
}

// Get overdue transactions
func (h *TransactionHandler) GetOverdueTransactions(w http.ResponseWriter, r *http.Request) {
    transactions, err := h.Repo.GetOverdueTransactions()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(transactions)
}

// Get remind transactions
func (h *TransactionHandler) GetRemindTransactions(w http.ResponseWriter, r *http.Request) {
    transactions, err := h.Repo.GetRemindTransactions()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(transactions)
}

// Get remind payments
func (h *TransactionHandler) GetRemindPayments(w http.ResponseWriter, r *http.Request) {
    transactions, err := h.Repo.GetRemindPayments()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(transactions)
}

func (h *TransactionHandler) PrintKuitansi(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    id, err := strconv.Atoi(params["id"])
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    trx, items, err := h.Repo.GetByIDStudentOrder(id)
    if err != nil {
        http.Error(w, "Transaction not found", http.StatusNotFound)
        return
    }

    // Get customer data for address and phone
    customer, err := h.Repo.GetCustomerByID(trx.CustomerID)
    var customerAddress, customerPhone string
    
    if err != nil {
        log.Printf("Error getting customer data for ID %d: %v", trx.CustomerID, err)
        customerAddress = "Alamat tidak tersedia"
        customerPhone = "No. Telp tidak tersedia"
    } else if customer == nil {
        log.Printf("Customer not found for ID %d", trx.CustomerID)
        customerAddress = "Alamat tidak tersedia"
        customerPhone = "No. Telp tidak tersedia"
    } else {
        // Check if address and contact are empty
        if customer.Address == "" || customer.Address == " " {
            customerAddress = "Alamat belum diisi"
        } else {
            customerAddress = customer.Address
        }
        
        if customer.Contact == "" || customer.Contact == " " {
            customerPhone = "No. Telp belum diisi"
        } else {
            customerPhone = customer.Contact
        }
        
        log.Printf("Customer data loaded - Address: '%s', Contact: '%s'", customerAddress, customerPhone)
    }

    // Read HTML template
    htmlTemplate, err := os.ReadFile("invoice.html")
    if err != nil {
        http.Error(w, "Template not found", http.StatusInternalServerError)
        return
    }

    htmlContent := string(htmlTemplate)

    // Replace placeholder values dengan data transaksi
    htmlContent = strings.ReplaceAll(htmlContent, "01.07.2022", formatDisplayDate(trx.Transaksidate))
    htmlContent = strings.ReplaceAll(htmlContent, "(Nama Sekolah)", trx.Customer_name)
    htmlContent = strings.ReplaceAll(htmlContent, "(Alamat)", customerAddress)
    htmlContent = strings.ReplaceAll(htmlContent, "(Notelp)", customerPhone)

    // Continue with rest of the function...
    // (Generate table rows, etc. - same as before)
    
    // Generate table rows untuk items
    var tableRows strings.Builder
    for _, item := range items {
        tableRows.WriteString(fmt.Sprintf(`
                    <tr class="tm_table_baseline">
                      <td class="tm_width_3 tm_primary_color">%s</td>
                      <td class="tm_width_4">%s</td>
                      <td class="tm_width_2">%s</td>
                      <td class="tm_width_1">%d</td>
                      <td class="tm_width_2 tm_text_right">%s</td>
                    </tr>`,
            item.StudentName,
            item.UniformName,
            item.Size,
            item.Quantity,
            formatCurrency(item.UnitPrice*float64(item.Quantity))))
    }

    // Replace table content using regex
    tableRegex := regexp.MustCompile(`(?s)<tbody>\s*<tr class="tm_table_baseline">.*?</tr>\s*</tbody>`)
    newTableBody := fmt.Sprintf(`<tbody>%s
                  </tbody>`, tableRows.String())
    htmlContent = tableRegex.ReplaceAllString(htmlContent, newTableBody)

    // Replace prices
    htmlContent = strings.ReplaceAll(htmlContent, "300.000", formatCurrency(trx.Total))

    // Replace status pembayaran
    statusText := getStatusText(trx.Status)
    htmlContent = strings.ReplaceAll(htmlContent, "(DP/LUNAS)", statusText)

    // Generate summary table
    summary := make(map[string]map[string]int)
    for _, item := range items {
        if summary[item.UniformName] == nil {
            summary[item.UniformName] = make(map[string]int)
        }
        summary[item.UniformName][item.Size] += item.Quantity
    }

    // Generate summary rows
    var summaryRows strings.Builder
    for name, sizes := range summary {
        for size, qty := range sizes {
            summaryRows.WriteString(fmt.Sprintf(`
                    <tr class="tm_table_baseline">
                      <td class="tm_width_4">%s</td>
                      <td class="tm_width_2">%s</td>
                      <td class="tm_width_1">%d</td>
                    </tr>`, name, size, qty))
        }
    }

    // Replace summary content using regex
    summaryIndex := strings.Index(htmlContent, "SUMMARY:")
    if summaryIndex != -1 {
        summarySection := htmlContent[summaryIndex:]
        summaryRegex := regexp.MustCompile(`(?s)<tbody>\s*<tr class="tm_table_baseline">.*?</tr>\s*</tbody>`)
        newSummaryBody := fmt.Sprintf(`<tbody>%s
                  </tbody>`, summaryRows.String())
        updatedSummarySection := summaryRegex.ReplaceAllString(summarySection, newSummaryBody)
        htmlContent = htmlContent[:summaryIndex] + updatedSummarySection
    }

    // Set header untuk HTML response
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    w.Write([]byte(htmlContent))
}

func (h *TransactionHandler) PrintKuitansibiasa(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    id, err := strconv.Atoi(params["id"])
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    trx, err := h.Repo.GetByIDNormal(id)
    if err != nil {
        http.Error(w, "Transaction not found", http.StatusNotFound)
        return
    }

    // Get customer data for address and phone
    customer, err := h.Repo.GetCustomerByID(trx.CustomerID)
    var customerAddress, customerPhone string
    
    if err != nil {
        log.Printf("Error getting customer data for ID %d: %v", trx.CustomerID, err)
        customerAddress = "Alamat tidak tersedia"
        customerPhone = "No. Telp tidak tersedia"
    } else if customer == nil {
        log.Printf("Customer not found for ID %d", trx.CustomerID)
        customerAddress = "Alamat tidak tersedia"
        customerPhone = "No. Telp tidak tersedia"
    } else {
        // Check if address and contact are empty
        if customer.Address == "" || customer.Address == " " {
            customerAddress = "Alamat belum diisi"
        } else {
            customerAddress = customer.Address
        }
        
        if customer.Contact == "" || customer.Contact == " " {
            customerPhone = "No. Telp belum diisi"
        } else {
            customerPhone = customer.Contact
        }
        
        log.Printf("Customer data loaded - Address: '%s', Contact: '%s'", customerAddress, customerPhone)
    }

    // Read HTML template
    htmlTemplate, err := os.ReadFile("invoicebiasa.html")
    if err != nil {
        http.Error(w, "Template not found", http.StatusInternalServerError)
        return
    }

    htmlContent := string(htmlTemplate)

    // Replace placeholder values
    htmlContent = strings.ReplaceAll(htmlContent, "01.07.2022", formatDisplayDate(trx.Transaksidate))
    htmlContent = strings.ReplaceAll(htmlContent, "(Nama Sekolah)", trx.Customer_name)
    htmlContent = strings.ReplaceAll(htmlContent, "(Alamat)", customerAddress)
    htmlContent = strings.ReplaceAll(htmlContent, "(Notelp)", customerPhone)

    // Generate table rows untuk pesanan
    if len(trx.Items) > 0 {
        var allTableRows strings.Builder
        for _, item := range trx.Items {
            allTableRows.WriteString(fmt.Sprintf(`
                    <tr class="tm_table_baseline">
                      <td class="tm_width_4">%s</td>
                      <td class="tm_width_2">%s</td>
                      <td class="tm_width_1">%d</td>
                      <td class="tm_width_2 tm_text_right">%s</td>
                    </tr>`, item.UniformName, item.Size, item.Quantity, formatCurrency(item.Subtotal)))
        }

        // Replace table content using regex
        tableRegex := regexp.MustCompile(`(?s)<tbody>\s*<tr class="tm_table_baseline">.*?</tr>\s*</tbody>`)
        newTableBody := fmt.Sprintf(`<tbody>%s
                  </tbody>`, allTableRows.String())
        htmlContent = tableRegex.ReplaceAllString(htmlContent, newTableBody)

        // Generate summary
        summary := make(map[string]map[string]int)
        for _, item := range trx.Items {
            if summary[item.UniformName] == nil {
                summary[item.UniformName] = make(map[string]int)
            }
            summary[item.UniformName][item.Size] += item.Quantity
        }

        // Generate summary rows
        var allSummaryRows strings.Builder
        for name, sizes := range summary {
            for size, qty := range sizes {
                allSummaryRows.WriteString(fmt.Sprintf(`
                    <tr class="tm_table_baseline">
                      <td class="tm_width_4">%s</td>
                      <td class="tm_width_2">%s</td>
                      <td class="tm_width_1">%d</td>
                    </tr>`, name, size, qty))
            }
        }

        // Replace summary using regex
        summaryIndex := strings.Index(htmlContent, "SUMMARY:")
        if summaryIndex != -1 {
            summarySection := htmlContent[summaryIndex:]
            summaryRegex := regexp.MustCompile(`(?s)<tbody>\s*<tr class="tm_table_baseline">.*?</tr>\s*</tbody>`)
            newSummaryBody := fmt.Sprintf(`<tbody>%s
                  </tbody>`, allSummaryRows.String())
            updatedSummarySection := summaryRegex.ReplaceAllString(summarySection, newSummaryBody)
            htmlContent = htmlContent[:summaryIndex] + updatedSummarySection
        }
    }

    // Replace grand total
    htmlContent = strings.ReplaceAll(htmlContent, "300.000", formatCurrency(trx.Total))

    // Replace status pembayaran
    statusText := getStatusText(trx.Status)
    htmlContent = strings.ReplaceAll(htmlContent, "(DP/LUNAS)", statusText)

    // Set header untuk HTML response
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    w.Write([]byte(htmlContent))
}

// Helper function untuk format tanggal display
func formatDisplayDate(dateString string) string {
    if dateString == "" {
        return time.Now().Format("02.01.2006")
    }
    
    // Parse different date formats
    layouts := []string{
        "2006-01-02",
        "2006-01-02 15:04:05",
        "02.01.2006",
        "02/01/2006",
    }
    
    for _, layout := range layouts {
        if t, err := time.Parse(layout, dateString); err == nil {
            return t.Format("02.01.2006")
        }
    }
    
    // If can't parse, return as is
    return dateString
}

// Helper function untuk format currency
func formatCurrency(amount float64) string {
    return fmt.Sprintf("%.0f", amount)
}

// Helper function untuk get status text
func getStatusText(status string) string {
    switch status {
    case "paid":
        return "LUNAS"
    case "pending":
        return "DP"
    case "cancelled":
        return "DIBATALKAN"
    default:
        return "DP"
    }
}
