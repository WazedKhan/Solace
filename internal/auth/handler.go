package auth

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
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
			log.Println(err)
			return
		}
	}

	res := RegisterResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Println(err)
	}
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

	res := make([]UserResponse, 0)
	for _, user := range users {
		res = append(res, UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		})
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Println(err)
	}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	res, err := h.service.Login(r.Context(), req)
	if err != nil {
		switch err {
		case ErrInvalidCredentials:
			http.Error(w, "email or password didn't match", http.StatusUnauthorized)
			return
		case ErrInvalidInput:
			http.Error(w, "invalid input", http.StatusBadRequest)
			return
		default:
			http.Error(w, "internal server error", http.StatusInternalServerError)
			log.Println(err)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Println(err)
	}
}
