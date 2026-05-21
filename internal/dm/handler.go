package dm

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/sumi-devs/canopy-social/canopy/internal/auth"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) RegisterRoutes(r chi.Router, jwtMiddleware func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(jwtMiddleware)
		r.Post("/api/v1/direct_messages", h.SendDM)
		r.Get("/api/v1/direct_messages/conversations", h.ListConversations)
		r.Get("/api/v1/direct_messages/conversations/{id}/messages", h.ListMessages)
		r.Post("/api/v1/direct_messages/conversations/{id}/read", h.MarkAsRead)
	})
}

type SendDMRequest struct {
	RecipientID string `json:"recipient_id"`
	Content     string `json:"content"`
}

func (h *Handler) SendDM(w http.ResponseWriter, r *http.Request) {
	senderID := auth.AccountIDFromContext(r.Context())
	var req SendDMRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	msg, err := h.svc.SendDM(r.Context(), senderID, req.RecipientID, req.Content)
	if err != nil {
		writeJSON(w, http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, msg)
}

func (h *Handler) ListConversations(w http.ResponseWriter, r *http.Request) {
	accountID := auth.AccountIDFromContext(r.Context())
	limit := parseIntParam(r, "limit", 20)
	offset := parseIntParam(r, "offset", 0)
	list, err := h.svc.ListConversations(r.Context(), accountID, limit, offset)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to list conversations"})
		return
	}
	if list == nil {
		list = make([]*Conversation, 0)
	}
	writeJSON(w, http.StatusOK, list)
}

func (h *Handler) ListMessages(w http.ResponseWriter, r *http.Request) {
	accountID := auth.AccountIDFromContext(r.Context())
	convID := chi.URLParam(r, "id")
	limit := parseIntParam(r, "limit", 20)
	offset := parseIntParam(r, "offset", 0)
	list, err := h.svc.ListMessages(r.Context(), accountID, convID, limit, offset)
	if err != nil {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": err.Error()})
		return
	}
	if list == nil {
		list = make([]*Message, 0)
	}
	writeJSON(w, http.StatusOK, list)
}

func (h *Handler) MarkAsRead(w http.ResponseWriter, r *http.Request) {
	accountID := auth.AccountIDFromContext(r.Context())
	convID := chi.URLParam(r, "id")
	if err := h.svc.MarkAsRead(r.Context(), accountID, convID); err != nil {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{})
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func parseIntParam(r *http.Request, key string, defaultVal int) int {
	s := r.URL.Query().Get(key)
	if s == "" {
		return defaultVal
	}
	v, err := strconv.Atoi(s)
	if err != nil || v < 0 {
		return defaultVal
	}
	return v
}
