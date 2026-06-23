package habit

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	jwt_token "github.com/WazedKhan/Solace/internal/auth/token"
	"github.com/WazedKhan/Solace/internal/pagination"
	"github.com/WazedKhan/Solace/internal/utils"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateHabit(w http.ResponseWriter, r *http.Request) {
	userId, ok := jwt_token.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}
	var req HabitRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	res, err := h.service.CreateHabit(r.Context(), req, userId)
	if err != nil {
		switch err {
		case ErrFailedWriting:
			http.Error(w, "failed to create habit", http.StatusInternalServerError)
			return
		case utils.ErrInvalidInput:
			http.Error(w, "invalid input", http.StatusBadRequest)
			return
		default:
			http.Error(w, "internal server error", http.StatusInternalServerError)
			log.Println(err)
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Println(err)
	}
}

func (h *Handler) GetHabits(w http.ResponseWriter, r *http.Request) {
	userId, ok := jwt_token.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}

	query := r.URL.Query()

	offset, err := strconv.Atoi(query.Get("offset"))
	if err != nil || offset < 0 {
		offset = 0
	}

	limit, err := strconv.Atoi(query.Get("limit"))
	if err != nil || limit < 1 {
		limit = 10
	}

	params := pagination.QueryParams{
		Limit:  limit,
		Offset: offset,
	}

	res, err := h.service.GetHabitsByUserID(r.Context(), userId, params)
	if err != nil {
		switch err {
		case ErrNotFound:
			http.Error(w, "no habit found", http.StatusBadRequest)
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

func (h *Handler) CheckIn(w http.ResponseWriter, r *http.Request) {
	userId, ok := jwt_token.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}

	habitID := r.PathValue("id")
	if habitID == "" {
		http.Error(w, "missing habit id", http.StatusBadRequest)
		return
	}

	crStreak, err := h.service.CheckIn(r.Context(), userId, habitID)
	if err != nil {
		switch err {
		case ErrNotFound:
			http.Error(w, "no habit found", http.StatusNotFound)
			return
		case ErrAlreadyChecked:
			http.Error(w, "already checked in", http.StatusConflict)
			return
		case ErrInvalidHabitID:
			http.Error(w, "invalid habit id", http.StatusForbidden)
			return
		default:
			http.Error(w, "internal server error", http.StatusInternalServerError)
			log.Println(err)
			return
		}
	}

	res := map[string]int{
		"current_streak": *crStreak,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Println(err)
	}
}
