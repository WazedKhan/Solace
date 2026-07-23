package journal

import (
	"encoding/json"
	"log"
	"net/http"

	jwt_token "github.com/WazedKhan/Solace/internal/auth/token"
	"github.com/WazedKhan/Solace/internal/pagination"
	"github.com/WazedKhan/Solace/internal/storage"
	"github.com/WazedKhan/Solace/internal/utils"
)

type Handler struct {
	service *Service
	store   storage.Storage
}

func NewHandler(service *Service, store storage.Storage) *Handler {
	return &Handler{
		service: service,
		store:   store,
	}
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

	pag := pagination.ParsePagination(r)
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

	pag := pagination.ParsePagination(r)
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

func (h *Handler) GetJournalByID(w http.ResponseWriter, r *http.Request) {
	journalID := r.PathValue("id")
	if journalID == "" {
		http.Error(w, "journal id is required", http.StatusBadRequest)
		return
	}

	userID, ok := jwt_token.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}

	res, err := h.service.GetJournalByID(r.Context(), userID, journalID)
	if err != nil {
		switch err {
		case utils.ErrNotFound:
			http.Error(w, "journal not found", http.StatusNotFound)
			return

		default:
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
	}

	resp := toJournalResponse(res)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Println(err)
	}
}

func (h *Handler) UpdateJournalByID(w http.ResponseWriter, r *http.Request) {
	userID, ok := jwt_token.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}

	journalID := r.PathValue("id")
	if journalID == "" {
		http.Error(w, "journal id is required", http.StatusBadRequest)
		return
	}

	var req UpdateJournalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	res, err := h.service.UpdateJournalByID(r.Context(), userID, journalID, req)
	if err != nil {
		switch err {
		case ErrInvalidStatus:
			http.Error(w, "invalid status", http.StatusBadRequest)
			return

		case utils.ErrNotFound:
			http.Error(w, "journal not found!", http.StatusNotFound)
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
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Println(err)
	}
}

func (h *Handler) SoftDeleteJournal(w http.ResponseWriter, r *http.Request) {
	userID, ok := jwt_token.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}

	journalID := r.PathValue("id")
	if journalID == "" {
		http.Error(w, "journal id required", http.StatusBadRequest)
		return
	}

	err := h.service.SoftDeleteJournal(r.Context(), userID, journalID)
	if err != nil {
		switch err {
		case utils.ErrNotFound:
			http.Error(w, "journal is not found", http.StatusNotFound)
			return

		default:
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ConfirmUpload(w http.ResponseWriter, r *http.Request) {
	userID, ok := jwt_token.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}

	journalID := r.PathValue("id")
	if journalID == "" {
		http.Error(w, "journal id required", http.StatusBadRequest)
		return
	}

	var req ConfirmUploadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	err := h.service.ConfirmUpload(r.Context(), userID, journalID, req.Key)
	if err != nil {
		switch err {
		case ErrInvalidKey:
			http.Error(w, "image key is invalid", http.StatusBadRequest)
			return

		case ErrForbidden:
			http.Error(w, "request is not allowed", http.StatusForbidden)
			return

		case ErrImageNotFound:
			http.Error(w, "image not found with provided key", http.StatusNotFound)
			return

		default:
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}
