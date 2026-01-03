package handlers

import (
	"encoding/json"
	"net/http"

	"quizzler/database"
	"quizzler/middleware"
	"quizzler/models"

	"golang.org/x/crypto/bcrypt"
)

func Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		http.Error(w, `{"error": "Email and password are required"}`, http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, `{"error": "Failed to hash password"}`, http.StatusInternalServerError)
		return
	}

	// Make the first user admin for funsies
	var userCount int
	database.DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&userCount)
	isAdmin := userCount == 0

	result, err := database.DB.Exec("INSERT INTO users (email, password, admin) VALUES (?, ?, ?)", req.Email, string(hashedPassword), isAdmin)
	if err != nil {
		http.Error(w, `{"error": "Email already exists"}`, http.StatusConflict)
		return
	}

	userID, _ := result.LastInsertId()

	token, err := middleware.GenerateToken(int(userID))
	if err != nil {
		http.Error(w, `{"error": "Failed to generate token"}`, http.StatusInternalServerError)
		return
	}

	user := models.User{
		ID:    int(userID),
		Email: req.Email,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.AuthResponse{
		Token: token,
		User:  user,
	})
}

func Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		http.Error(w, `{"error": "Email and password are required"}`, http.StatusBadRequest)
		return
	}

	var user models.User
	err := database.DB.QueryRow("SELECT id, email, password FROM users WHERE email = ?", req.Email).
		Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		http.Error(w, `{"error": "Invalid credentials"}`, http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		http.Error(w, `{"error": "Invalid credentials"}`, http.StatusUnauthorized)
		return
	}

	token, err := middleware.GenerateToken(user.ID)
	if err != nil {
		http.Error(w, `{"error": "Failed to generate token"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.AuthResponse{
		Token: token,
		User:  user,
	})
}
