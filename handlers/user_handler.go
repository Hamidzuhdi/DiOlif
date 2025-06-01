package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"konveksi-app/models"
	"konveksi-app/repositories"
	"database/sql"
	"strings"

	"github.com/gorilla/mux"
)

type UserHandler struct {
	Repo *repositories.UserRepository
	DB   *sql.DB
}

func (h *UserHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Form tidak valid", http.StatusBadRequest)
		return
	}

	username := strings.TrimSpace(r.FormValue("username"))
	password := strings.TrimSpace(r.FormValue("password"))

	// Validasi awal
	if username == "" || password == "" {
		http.Error(w, "Username dan password wajib diisi", http.StatusBadRequest)
		return
	}

	// Ambil user dari database
	var dbPassword string
	err := h.DB.QueryRow("SELECT password FROM users WHERE username = ?", username).Scan(&dbPassword)
	if err == sql.ErrNoRows {
		http.Error(w, "Username tidak ditemukan", http.StatusUnauthorized)
		return
	} else if err != nil {
		http.Error(w, "Terjadi kesalahan pada server", http.StatusInternalServerError)
		return
	}

	// Cocokkan password
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if password != dbPassword {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`<script>alert("Password salah"); window.history.back();</script>`))
		return
	}

	// Login berhasil
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.Repo.Create(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	user, err := h.Repo.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.Repo.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user.ID = id

	if err := h.Repo.Update(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
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
