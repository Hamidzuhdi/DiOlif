package handlers

import (
    "encoding/json"
    "net/http"
  
    "time"
    "crypto/rand"
    "encoding/hex"
    "konveksi-app/models"
    "konveksi-app/repositories"
    "database/sql"
    "strings"
    "log"

)

type UserHandler struct {
    Repo *repositories.UserRepository
    DB   *sql.DB
}

// Simple in-memory session store (untuk production gunakan Redis/database)
var sessions = make(map[string]string)

func generateSessionID() string {
    bytes := make([]byte, 16)
    rand.Read(bytes)
    return hex.EncodeToString(bytes)
}

// Export functions untuk digunakan di main.go
func ValidateSession(sessionID string) bool {
    _, exists := sessions[sessionID]
    log.Printf("Validating session %s: exists=%v", sessionID, exists)
    return exists
}

func GetSessions() map[string]string {
    return sessions
}

func (h *UserHandler) LoginAPI(w http.ResponseWriter, r *http.Request) {
    log.Println("LoginAPI called")
    
    var loginRequest struct {
        Username string `json:"username"`
        Password string `json:"password"`
    }

    if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
        log.Printf("Error decoding JSON: %v", err)
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }

    username := strings.TrimSpace(loginRequest.Username)
    password := strings.TrimSpace(loginRequest.Password)
    
    log.Printf("Login attempt - Username: '%s', Password: '%s'", username, password)

    if username == "" || password == "" {
        log.Println("Empty username or password")
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{
            "success": "false",
            "message": "Username dan password wajib diisi",
        })
        return
    }

    // Query user dari database
    var user models.User
    query := "SELECT id, username, password, contact, address, created_at FROM users WHERE username = ?"
    log.Printf("Executing query: %s with username: %s", query, username)
    
    err := h.DB.QueryRow(query, username).Scan(
        &user.ID, &user.Username, &user.Password, &user.Contact, &user.Address, &user.CreatedAt,
    )

    if err == sql.ErrNoRows {
        log.Printf("User not found: %s", username)
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusUnauthorized)
        json.NewEncoder(w).Encode(map[string]string{
            "success": "false",
            "message": "Username tidak ditemukan",
        })
        return
    } else if err != nil {
        log.Printf("Database error: %v", err)
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string{
            "success": "false",
            "message": "Terjadi kesalahan pada server",
        })
        return
    }

    log.Printf("User found - ID: %d, Username: %s, DB Password: '%s', Input Password: '%s'", 
        user.ID, user.Username, user.Password, password)

    // Verify password (plain text comparison)
    if password != user.Password {
        log.Printf("Password mismatch - Expected: '%s', Got: '%s'", user.Password, password)
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusUnauthorized)
        json.NewEncoder(w).Encode(map[string]string{
            "success": "false",
            "message": "Password salah",
        })
        return
    }

    log.Println("Password verified successfully")

    // Generate session
    sessionID := generateSessionID()
    sessions[sessionID] = username
    
    log.Printf("Session created - ID: %s, Username: %s", sessionID, username)

    // Set cookie
    cookie := &http.Cookie{
        Name:     "session",
        Value:    sessionID,
        Path:     "/",
        HttpOnly: true,
        Secure:   false, // Set ke true untuk HTTPS
        SameSite: http.SameSiteLaxMode,
        Expires:  time.Now().Add(24 * time.Hour), // 24 jam
    }
    http.SetCookie(w, cookie)

    log.Println("Cookie set successfully")

    // Response sukses
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "success": "true",
        "message": "Login berhasil",
        "user": map[string]interface{}{
            "id":       user.ID,
            "username": user.Username,
            "contact":  user.Contact,
            "address":  user.Address,
        },
    })
    
    log.Println("Login response sent successfully")
}

// Rest of the methods remain the same...
func (h *UserHandler) LogoutAPI(w http.ResponseWriter, r *http.Request) {
    // Get session cookie
    cookie, err := r.Cookie("session")
    if err == nil {
        // Remove session dari store
        delete(sessions, cookie.Value)
    }

    // Clear cookie
    clearCookie := &http.Cookie{
        Name:     "session",
        Value:    "",
        Path:     "/",
        HttpOnly: true,
        Expires:  time.Unix(0, 0), // Set expired
    }
    http.SetCookie(w, clearCookie)

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
        "success": "true",
        "message": "Logout berhasil",
    })
}