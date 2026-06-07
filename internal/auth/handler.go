package auth

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/WazedKhan/Solace/internal/auth/utils"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	password := "wazed"
	hashed_password, err := utils.HashPassword(password)
	if err != nil {
		log.Println("failed to hash the password!", err)
	}
	fmt.Fprintln(w, hashed_password)
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.service.Register(r.Context(), req)
	if err != nil {
		switch err {
		case ErrEmailAlreadyExists:
			http.Error(w, "email already exists", http.StatusConflict)
			return
		case ErrInvalidInput:
			http.Error(w, "invalid input", http.StatusBadRequest)
			return
		default:
			http.Error(w, "internal server error", http.StatusInternalServerError)
			log.Panicln(err)
			return
		}

	}

	res := RegisterResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}
