package journal

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/WazedKhan/Solace/configs"
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

func (h *Handler) CreateJournal(w http.ResponseWriter, r *http.Request) {
	userID, ok := jwt_token.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}

	var req CreateJournalRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	res, err := h.service.CreateJournal(r.Context(), req, userID)
	if err != nil {
		switch err {
		case ErrTitleEmpty:
			http.Error(w, "title is required!", http.StatusBadRequest)
			return
		case ErrDescription:
			http.Error(w, "description is required!", http.StatusBadRequest)
			return
		case ErrInvalidStatus:
			http.Error(w, "invalid status", http.StatusBadRequest)
			return
		case utils.ErrInvalidInput:
			http.Error(w, "invalid inputs", http.StatusBadRequest)
			return
		case utils.ErrForeignKey:
			http.Error(w, "invalid mood id", http.StatusBadRequest)
			return
		default:
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
	}
	resp := JournalResponse{
		ID:          res.ID,
		Title:       res.Title,
		Description: res.Description,
		Status:      res.Status,
		ImageURL:    res.ImageURL,
		CreatedAt:   res.CreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Println(err)
	}
}

func (h *Handler) GetJournals(w http.ResponseWriter, r *http.Request) {
	userID, ok := jwt_token.UserIDFromContext(r.Context())
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
	} else if limit > configs.FetchLimit {
		limit = configs.FetchLimit
	}

	pag := pagination.QueryParams{
		Limit:  limit,
		Offset: offset,
	}
	res, err := h.service.GetJournalsByUser(r.Context(), userID, string(StatusPublished), pag)
	if err != nil {
		switch err {
		default:
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
	}

	resp := make([]JournalResponse, 0, len(res))
	for _, item := range res {
		resp = append(resp, toJournalResponse(item))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Println(err)
	}
}

func (h *Handler) GetDrafts(w http.ResponseWriter, r *http.Request) {
	userID, ok := jwt_token.UserIDFromContext(r.Context())
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
	} else if limit > configs.FetchLimit {
		limit = configs.FetchLimit
	}

	pag := pagination.QueryParams{
		Limit:  limit,
		Offset: offset,
	}
	res, err := h.service.GetJournalsByUser(r.Context(), userID, string(StatusDraft), pag)
	if err != nil {
		switch err {
		default:
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
	}

	resp := make([]JournalResponse, 0, len(res))
	for _, item := range res {
		resp = append(resp, toJournalResponse(item))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Println(err)
	}
}
