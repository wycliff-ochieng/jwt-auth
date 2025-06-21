package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	model "github.com/wycliff-ochieng/models"
	"github.com/wycliff-ochieng/service"
)

type AuthHandler struct {
	l     *log.Logger
	UServ *service.UserService
}

type RegisterReq struct {
	ID        int64
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type LoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	User         *model.UserResponse
	AccessToken  string
	RefreshToken string
}

func NewAuthHandle(l *log.Logger, UServ *service.UserService) *AuthHandler {
	return &AuthHandler{l: l, UServ: UServ}
}

func (h *AuthHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var req RegisterReq

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "unable to create user", http.StatusInternalServerError)
		return
	}

	if req.Firstname == "" || req.Lastname == "" || req.Email == "" || req.Password == "" {
		http.Error(w, "user could not be created", http.StatusExpectationFailed)
		return
	}

	//register user
	user, err := h.UServ.Register(req.ID, req.Firstname, req.Lastname, req.Email, req.Password)
	if err == service.ErrEmailExists {
		http.Error(w, "email already exists ", http.StatusConflict)
		return
	}
	if err != nil {
		h.l.Printf("registration failed %v", err)
		http.Error(w, "fialed to register", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&user)

}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {

	var req LoginReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "issue with your request", http.StatusInternalServerError)
		return
	}

	if req.Email == "" || req.Password == "" {
		http.Error(w, "input is required", http.StatusExpectationFailed)
		return
	}

	//authenticate user
	tokens, user, err := h.UServ.Login(req.Email, req.Password)
	if err == service.ErrInvalidPassword || err == service.ErrUserNotFound {
		http.Error(w, "failed to be logged in", http.StatusUnauthorized)
		h.l.Printf("Why fail: %v", err)
		return
	}
	if err != nil {
		http.Error(w, "something unexpected occurred", http.StatusInternalServerError)
		h.l.Printf("Login error : %v", err)
		return
	}

	//return response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(AuthResponse{
		User:         user,
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	})

}
