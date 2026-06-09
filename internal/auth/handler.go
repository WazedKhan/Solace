package auth

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

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

func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	limit, err := strconv.Atoi(queryParams.Get("limit"))
	if err != nil {
		limit = 10
	}

	offset, err := strconv.Atoi(queryParams.Get("offset"))
	if err != nil {
		offset = 0
	}

	query := GetUserQuery{
		Limit:  limit,
		Offset: offset,
		Search: queryParams.Get("search"),
	}

	users, err := h.service.GetUsers(r.Context(), query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var res []UserResponse
	for _, user := range users {
		res = append(res, UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		})
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}
