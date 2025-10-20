package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"konveksi-app/models"
	"konveksi-app/repositories"

	"github.com/gorilla/mux"
)

type CustomerHandler struct {
	Repo *repositories.CustomerRepository
}

func (h *CustomerHandler) CreateCustomer(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Name     string                  `json:"name"`
        Type     string                  `json:"type"`
        Contact  string                  `json:"contact"`
        Address  string                  `json:"address"`
        Uniforms []models.CustomerUniform `json:"uniforms"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    customer := models.Customer{
        Name: req.Name,
        Type: req.Type,
        Contact: req.Contact,
        Address: req.Address,
    }
    if err := h.Repo.CreateWithUniforms(&customer, req.Uniforms); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(customer)
}

func (h *CustomerHandler) AddCustomerUniform(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    customerID, err := strconv.Atoi(params["id"])
    if err != nil {
        http.Error(w, "Invalid customer ID", http.StatusBadRequest)
        return
    }
    var u models.CustomerUniform
    if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    u.CustomerID = customerID
    if err := h.Repo.AddCustomerUniform(&u); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(u)
}

func (h *CustomerHandler) GetCustomer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	customer, err := h.Repo.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(customer)
}

func (h *CustomerHandler) GetAllCustomers(w http.ResponseWriter, r *http.Request) {
	customers, err := h.Repo.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(customers)
}

func (h *CustomerHandler) UpdateCustomer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var customer models.Customer
	if err := json.NewDecoder(r.Body).Decode(&customer); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	customer.ID = id

	if err := h.Repo.Update(&customer); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(customer)
}

func (h *CustomerHandler) GetCustomerUniform(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    id, err := strconv.Atoi(params["id"])
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }
    uniform, err := h.Repo.GetCustomerUniformByID(id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(uniform)
}
func (h *CustomerHandler) GetUniformsByCustomerID(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    id, err := strconv.Atoi(params["id"])
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }
    uniforms, err := h.Repo.GetUniformsByCustomerID(id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(uniforms)
}

func (h *CustomerHandler) DeleteCustomerUniform(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    id, err := strconv.Atoi(params["id"])
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }
    if err := h.Repo.DeleteCustomerUniform(id); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusNoContent)
}

// Update uniform + catat history
func (h *CustomerHandler) UpdateCustomerUniform(w http.ResponseWriter, r *http.Request) {
    idStr := mux.Vars(r)["id"]
    id, _ := strconv.Atoi(idStr)
    var req struct {
        UniformName string  `json:"uniform_name"`
        Size        string  `json:"size"`
        Price       float64 `json:"price"`
        Notes       string  `json:"notes"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }
    err := h.Repo.UpdateCustomerUniformWithHistory(id, req.UniformName, req.Size, req.Price, req.Notes)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusOK)
}

// Get history harga
func (h *CustomerHandler) GetUniformPriceHistory(w http.ResponseWriter, r *http.Request) {
    idStr := mux.Vars(r)["id"]
    id, _ := strconv.Atoi(idStr)
    history, err := h.Repo.GetUniformPriceHistory(id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(history)
}

func (h *CustomerHandler) DeleteCustomer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := h.Repo.Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}