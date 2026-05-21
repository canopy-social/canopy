package notifications

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
		r.Get("/api/v1/notifications", h.List)
		r.Post("/api/v1/notifications/{id}/read", h.MarkRead)
		r.Post("/api/v1/notifications/read", h.MarkAllRead)
		r.Delete("/api/v1/notifications/{id}", h.Dismiss)
	})
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	accountID := auth.AccountIDFromContext(r.Context())
	limit := parseIntParam(r, "limit", 20)
	if limit <= 0 || limit > 80 {
		limit = 20
	}
	maxID := r.URL.Query().Get("max_id")
	sinceID := r.URL.Query().Get("since_id")
	types := r.URL.Query()["types"]
	if len(types) == 0 {
		types = r.URL.Query()["types[]"]
	}
	excludeTypes := r.URL.Query()["exclude_types"]
	if len(excludeTypes) == 0 {
		excludeTypes = r.URL.Query()["exclude_types[]"]
	}
	list, err := h.svc.List(r.Context(), accountID, limit, maxID, sinceID, types, excludeTypes)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to list notifications"})
		return
	}
	if list == nil {
		list = make([]*Notification, 0)
	}
	writeJSON(w, http.StatusOK, list)
}

func (h *Handler) MarkRead(w http.ResponseWriter, r *http.Request) {
	accountID := auth.AccountIDFromContext(r.Context())
	id := chi.URLParam(r, "id")
	n, err := h.svc.MarkRead(r.Context(), accountID, id)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "notification not found"})
		return
	}
	writeJSON(w, http.StatusOK, n)
}

func (h *Handler) MarkAllRead(w http.ResponseWriter, r *http.Request) {
	accountID := auth.AccountIDFromContext(r.Context())
	if err := h.svc.MarkAllRead(r.Context(), accountID); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to clear notifications"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{})
}

func (h *Handler) Dismiss(w http.ResponseWriter, r *http.Request) {
	accountID := auth.AccountIDFromContext(r.Context())
	id := chi.URLParam(r, "id")
	if err := h.svc.Dismiss(r.Context(), accountID, id); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to dismiss notification"})
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
