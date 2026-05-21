package essays

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/sumi-devs/canopy-social/canopy/internal/auth"
)

// Handler handles essay HTTP endpoints.
type Handler struct {
	svc *Service
}

// NewHandler creates a new essays handler.
func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// RegisterRoutes mounts all essay routes.
func (h *Handler) RegisterRoutes(r chi.Router, jwtMiddleware, optionalJWT func(http.Handler) http.Handler) {
	// Public
	r.Group(func(r chi.Router) {
		r.Use(optionalJWT)
		r.Get("/api/v1/essays/{id}", h.GetEssay)
		r.Get("/api/v1/essays/{accountID}/{slug}", h.GetEssayBySlug)
		r.Get("/api/v1/accounts/{accountID}/essays", h.ListByAccount)
	})

	// Authenticated
	r.Group(func(r chi.Router) {
		r.Use(jwtMiddleware)
		r.Post("/api/v1/essays", h.CreateEssay)
		r.Put("/api/v1/essays/{id}", h.UpdateEssay)
		r.Delete("/api/v1/essays/{id}", h.DeleteEssay)
		r.Post("/api/v1/essays/{id}/publish", h.PublishEssay)
		r.Post("/api/v1/essays/{id}/unpublish", h.UnpublishEssay)
		r.Get("/api/v1/essays/drafts", h.ListDrafts)
	})
}

func (h *Handler) GetEssay(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	essay, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "essay not found"})
		return
	}
	writeJSON(w, http.StatusOK, essay)
}

func (h *Handler) GetEssayBySlug(w http.ResponseWriter, r *http.Request) {
	accountID := chi.URLParam(r, "accountID")
	slug := chi.URLParam(r, "slug")
	essay, err := h.svc.GetBySlug(r.Context(), accountID, slug)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "essay not found"})
		return
	}
	writeJSON(w, http.StatusOK, essay)
}

func (h *Handler) ListByAccount(w http.ResponseWriter, r *http.Request) {
	accountID := chi.URLParam(r, "accountID")
	limit := parseIntParam(r, "limit", 20)
	offset := parseIntParam(r, "offset", 0)
	essays, err := h.svc.ListByAccount(r.Context(), accountID, limit, offset)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to list essays"})
		return
	}
	writeJSON(w, http.StatusOK, essays)
}

func (h *Handler) CreateEssay(w http.ResponseWriter, r *http.Request) {
	accountID := auth.AccountIDFromContext(r.Context())
	var params CreateEssayParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	essay, err := h.svc.Create(r.Context(), accountID, &params)
	if err != nil {
		writeJSON(w, http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, essay)
}

func (h *Handler) UpdateEssay(w http.ResponseWriter, r *http.Request) {
	accountID := auth.AccountIDFromContext(r.Context())
	essayID := chi.URLParam(r, "id")
	var params UpdateEssayParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	essay, err := h.svc.Update(r.Context(), essayID, accountID, &params)
	if err != nil {
		writeJSON(w, http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, essay)
}

func (h *Handler) DeleteEssay(w http.ResponseWriter, r *http.Request) {
	accountID := auth.AccountIDFromContext(r.Context())
	essayID := chi.URLParam(r, "id")
	if err := h.svc.Delete(r.Context(), essayID, accountID); err != nil {
		writeJSON(w, http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "essay deleted"})
}

func (h *Handler) PublishEssay(w http.ResponseWriter, r *http.Request) {
	accountID := auth.AccountIDFromContext(r.Context())
	essayID := chi.URLParam(r, "id")
	essay, err := h.svc.Publish(r.Context(), essayID, accountID)
	if err != nil {
		writeJSON(w, http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, essay)
}

func (h *Handler) UnpublishEssay(w http.ResponseWriter, r *http.Request) {
	accountID := auth.AccountIDFromContext(r.Context())
	essayID := chi.URLParam(r, "id")
	essay, err := h.svc.Unpublish(r.Context(), essayID, accountID)
	if err != nil {
		writeJSON(w, http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, essay)
}

func (h *Handler) ListDrafts(w http.ResponseWriter, r *http.Request) {
	accountID := auth.AccountIDFromContext(r.Context())
	limit := parseIntParam(r, "limit", 20)
	offset := parseIntParam(r, "offset", 0)
	essays, err := h.svc.ListDrafts(r.Context(), accountID, limit, offset)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to list drafts"})
		return
	}
	writeJSON(w, http.StatusOK, essays)
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
