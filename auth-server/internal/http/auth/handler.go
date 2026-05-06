package auth

import (
	"encoding/json"
	"net/http"

	"github.com/leandrogutierrez148/acomm/auth-server/internal/models"
	"github.com/leandrogutierrez148/acomm/auth-server/internal/repositories"
	"github.com/leandrogutierrez148/acomm/auth-server/internal/security"
)

type AuthHandler struct {
	repo *repositories.UsersRepository
}

func NewAuthHandler(repo *repositories.UsersRepository) *AuthHandler {
	return &AuthHandler{repo: repo}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	existingUser, err := h.repo.GetUserByEmail(req.Email)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if existingUser != nil {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	hashedPassword, err := security.HashPassword(req.Password)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	user := &models.User{
		Email:        req.Email,
		PasswordHash: &hashedPassword,
		UserType:     models.RoleCustomer, // Default to customer
	}

	if err := h.repo.CreateUser(user); err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.repo.GetUserByEmail(req.Email)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if user == nil || user.PasswordHash == nil || !security.CheckPasswordHash(req.Password, *user.PasswordHash) {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	token, err := security.GenerateJWT(user)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	resp := models.AuthResponse{
		Token: token,
	}
	resp.User.ID = user.ID.String()
	resp.User.Email = user.Email
	resp.User.UserType = string(user.UserType)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
