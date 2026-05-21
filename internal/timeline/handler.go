package timeline

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

func (h *Handler) RegisterRoutes(r chi.Router, jwtMiddleware, optionalJWT func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(jwtMiddleware)
		r.Get("/api/v1/timelines/home", h.HomeTimeline)
	})

	r.Group(func(r chi.Router) {
		r.Use(optionalJWT)
		r.Get("/api/v1/timelines/public", h.PublicTimeline)
	})
}

func (h *Handler) HomeTimeline(w http.ResponseWriter, r *http.Request) {
	accountID := auth.AccountIDFromContext(r.Context())
	limit := parseIntParam(r, "limit", 20)
	maxID := r.URL.Query().Get("max_id")
	sinceID := r.URL.Query().Get("since_id")

	resp, err := h.svc.GetHomeTimeline(r.Context(), accountID, limit, maxID, sinceID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *Handler) PublicTimeline(w http.ResponseWriter, r *http.Request) {
	limit := parseIntParam(r, "limit", 20)
	offset := parseIntParam(r, "offset", 0)
	local := r.URL.Query().Get("local") == "true"

	posts, err := h.svc.GetPublicTimeline(r.Context(), local, limit, offset)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, posts)
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
