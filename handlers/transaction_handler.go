package handlers

import (
	"encoding/json"
	"konveksi-app/models"
	"konveksi-app/repositories"
	"net/http"
	"strconv"
    "time"

    "github.com/gorilla/mux"
    "github.com/jung-kurt/gofpdf"
    "bytes"
)

type TransactionHandler struct {
	Repo *repositories.TransactionRepository
}
func (h *TransactionHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var request struct {
		CustomerID      int     `json:"customer_id"`
		TransactionDate string  `json:"transaction_date"`
		PaymentDate     string  `json:"payment_date"`
		Status          string  `json:"status"`
		Notes           string  `json:"notes"`
		Items           []struct {
			UniformName string     `json:"uniform_name"`
			Size          string  `json:"size"`
			Quantity      int     `json:"quantity"`
			UnitPrice     float64 `json:"unit_price"`
			Notes         string  `json:"notes"`
		} `json:"items"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, `{"error": "Invalid JSON: `+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

 // Validasi tanggal (optional, biar tetap string)
    layout := "2006-01-02"
    if request.TransactionDate != "" {
        if _, err := time.Parse(layout, request.TransactionDate); err != nil {
            http.Error(w, `{"error": "Invalid transaction_date: `+err.Error()+`"}`, http.StatusBadRequest)
            return
        }
    }
    if request.PaymentDate != "" {
        if _, err := time.Parse(layout, request.PaymentDate); err != nil {
            http.Error(w, `{"error": "Invalid payment_date: `+err.Error()+`"}`, http.StatusBadRequest)
            return
        }
    }

	// Buat objek transaksi dari input
	transaction := models.Transaksi{
		CustomerID:    request.CustomerID,
		Transaksidate: request.TransactionDate,
		Paymentdate:   request.PaymentDate,
		Status:        request.Status,
		Notes:         request.Notes,
	}

	// Hitung total dan salin item
	transaction.Total = 0
	for _, item := range request.Items {
		subtotal := item.UnitPrice * float64(item.Quantity)
		transaction.Items = append(transaction.Items, models.OrderItem{
			UniformName: item.UniformName,
			Size:      item.Size,
			Quantity:  item.Quantity,
			UnitPrice: item.UnitPrice,
			Notes:     item.Notes,
			Subtotal:  subtotal,
		})
		transaction.Total += subtotal
	}

	// Simpan ke database
	if err := h.Repo.Create(&transaction); err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	// Berhasil
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(transaction)
}

// func (h *TransactionHandler) GetTransaction(w http.ResponseWriter, r *http.Request) {
//     params := mux.Vars(r)
//     id, err := strconv.Atoi(params["id"])
//     if err != nil {
//         http.Error(w, "Invalid ID", http.StatusBadRequest)
//         return
//     }
//     trx, err := h.Repo.GetByID(id)
//     if err != nil {
//         http.Error(w, "Transaction not found", http.StatusNotFound)
//         return
//     }
//     w.Header().Set("Content-Type", "application/json")
//     json.NewEncoder(w).Encode(trx)
// }

func (h *TransactionHandler) GetAllStudentOrders(w http.ResponseWriter, r *http.Request) {
    results, err := h.Repo.GetAllStudentOrders()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(results)
}

// Untuk transaksi dari order_items (transaksi biasa)
func (h *TransactionHandler) GetAllTransactionOrder(w http.ResponseWriter, r *http.Request) {
    transactions, err := h.Repo.GetAllTransactionOrder()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(transactions)
}

// Untuk transaksi dari student_order_items (repeat order)
func (h *TransactionHandler) GetAllTransactionStudent(w http.ResponseWriter, r *http.Request) {
    transactions, err := h.Repo.GetAllTransactionStudent()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(transactions)
}

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

func (h *TransactionHandler) GetOverduePaymentTransactions(w http.ResponseWriter, r *http.Request) {
    transactions, err := h.Repo.GetOverduePaymentTransactions()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(transactions)
}

func (h *TransactionHandler) GetOverdueTransactions(w http.ResponseWriter, r *http.Request) {
    transactions, err := h.Repo.GetOverdueTransactions()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(transactions)
}

func (h *TransactionHandler) GetRemindTransactions(w http.ResponseWriter, r *http.Request) {
    transactions, err := h.Repo.GetRemindTransactions()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(transactions)
}

func (h *TransactionHandler) GetRemindPayments(w http.ResponseWriter, r *http.Request) {
    transactions, err := h.Repo.GetRemindPayments()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(transactions)
}

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
func (h *TransactionHandler) PrintKuitansi(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    id, err := strconv.Atoi(params["id"])
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }
    // Ambil data transaksi dan items
    trx, items, err := h.Repo.GetByIDStudentOrder(id)
    if err != nil {
        http.Error(w, "Transaction not found", http.StatusNotFound)
        return
    }
    
    // Generate PDF
    pdf := gofpdf.New("P", "mm", "A4", "")
    pdf.AddPage()
    pdf.SetFont("Arial", "B", 16)
    pdf.Cell(40, 10, "Kuitansi Pembayaran")
    pdf.Ln(12)
    
    pdf.SetFont("Arial", "", 12)
    // Detail transaksi
    pdf.Cell(40, 10, "ID Transaksi: "+strconv.Itoa(trx.ID))
    pdf.Ln(8)
    pdf.Cell(40, 10, "Pelanggan: "+trx.Customer_name)
    pdf.Ln(8)
    pdf.Cell(40, 10, "Tanggal: "+trx.Transaksidate)
    pdf.Ln(8)
    pdf.Cell(40, 10, "Status: "+trx.Status)
    pdf.Ln(8)
    pdf.Cell(40, 10, "Total: Rp "+strconv.FormatFloat(trx.Total, 'f', 0, 64))
    pdf.Ln(12)
    
    // Summary seragam (hitung dari items)
    summary := make(map[string]int)
    for _, item := range items {
        summary[item.UniformName] += item.Quantity
    }
    
    pdf.Cell(40, 10, "Summary Seragam:")
    pdf.Ln(8)
    for uniformName, qty := range summary {
        pdf.Cell(0, 8, "- "+uniformName+": "+strconv.Itoa(qty)+" pcs")
        pdf.Ln(8)
    }
    pdf.Ln(4)
    
    // Daftar items detail
    pdf.Cell(40, 10, "Detail Siswa:")
    pdf.Ln(8)
    for _, item := range items {
        pdf.Cell(0, 8, "- "+item.StudentName+" ("+item.Grade+") - "+item.UniformName+" "+item.Size+" x"+strconv.Itoa(item.Quantity))
        pdf.Ln(8)
    }
    
    var buf bytes.Buffer
    err = pdf.Output(&buf)
    if err != nil {
        http.Error(w, "Failed to generate PDF", http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/pdf")
    w.Header().Set("Content-Disposition", "inline; filename=kuitansi.pdf")
    w.Write(buf.Bytes())
}

func (h *TransactionHandler) PrintKuitansibiasa(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    id, err := strconv.Atoi(params["id"])
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }
    // Ambil data transaksi dan items
    trx, err := h.Repo.GetByIDNormal(id)
    if err != nil {
        http.Error(w, "Transaction not found", http.StatusNotFound)
        return
    }
    
    // Generate PDF
    pdf := gofpdf.New("P", "mm", "A4", "")
    pdf.AddPage()
    pdf.SetFont("Arial", "B", 16)
    pdf.Cell(40, 10, "Kuitansi Pembayaran")
    pdf.Ln(12)
    
    pdf.SetFont("Arial", "", 12)
    // Detail transaksi
    pdf.Cell(40, 10, "ID Transaksi: "+strconv.Itoa(trx.ID))
    pdf.Ln(8)
    pdf.Cell(40, 10, "Pelanggan: "+trx.Customer_name)
    pdf.Ln(8)
    pdf.Cell(40, 10, "Tanggal: "+trx.Transaksidate)
    pdf.Ln(8)
    pdf.Cell(40, 10, "Status: "+trx.Status)
    pdf.Ln(8)
    pdf.Cell(40, 10, "Total: Rp "+strconv.FormatFloat(trx.Total, 'f', 0, 64))
    pdf.Ln(12)
    
    // Summary seragam (hitung dari order_items)
    summary := make(map[string]int)
    for _, item := range trx.Items {
        summary[item.UniformName] += item.Quantity
    }
    
    pdf.Cell(40, 10, "Summary Seragam:")
    pdf.Ln(8)
    for uniformName, qty := range summary {
        pdf.Cell(0, 8, "- "+uniformName+": "+strconv.Itoa(qty)+" pcs")
        pdf.Ln(8)
    }
    pdf.Ln(4)
    
    // Daftar items detail
    pdf.Cell(40, 10, "Detail Items:")
    pdf.Ln(8)
    for _, item := range trx.Items {
        pdf.Cell(0, 8, "- "+item.UniformName+" "+item.Size+" x"+strconv.Itoa(item.Quantity)+" = Rp "+strconv.FormatFloat(item.Subtotal, 'f', 0, 64))
        pdf.Ln(8)
    }
    
    var buf bytes.Buffer
    err = pdf.Output(&buf)
    if err != nil {
        http.Error(w, "Failed to generate PDF", http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/pdf")
    w.Header().Set("Content-Disposition", "inline; filename=kuitansi.pdf")
    w.Write(buf.Bytes())
}

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

func (h *TransactionHandler) CreateStudentOrder(w http.ResponseWriter, r *http.Request) {
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
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }
    transaction := models.Transaksi{
        CustomerID:    req.CustomerID,
        Transaksidate: req.TransactionDate,
        Paymentdate:   req.PaymentDate,
        Status:        req.Status,
        Notes:         req.Notes,
    }
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
            Notes:       item.Notes,
        })
    }
    if err := h.Repo.CreateStudentOrder(&transaction, studentItems); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(transaction)
}

// API untuk transaksi biasa
func (h *TransactionHandler) GetNormalTransactions(w http.ResponseWriter, r *http.Request) {
    transactions, err := h.Repo.GetNormalTransactions()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(transactions)
}