package media

import (
	"encoding/json"
	"net/http"

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
		r.Post("/api/v1/media", h.UploadMedia)
		r.Put("/api/v1/media/{id}", h.UpdateMedia)
		r.Get("/api/v1/media/{id}", h.GetMedia)
	})
}

func (h *Handler) UploadMedia(w http.ResponseWriter, r *http.Request) {
	accountID := auth.AccountIDFromContext(r.Context())
	if accountID == "" {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "failed to parse multipart form"})
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "file parameter is required"})
		return
	}
	defer file.Close()

	att, err := h.svc.ProcessUpload(r.Context(), accountID, file, header.Filename)
	if err != nil {
		writeJSON(w, http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusCreated, att)
}

func (h *Handler) UpdateMedia(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req struct {
		Description *string `json:"description"`
		AltText     *string `json:"alt_text"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	var altText *string
	if req.AltText != nil {
		altText = req.AltText
	} else if req.Description != nil {
		altText = req.Description
	}

	att, err := h.svc.UpdateAttachment(r.Context(), id, altText)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, att)
}

func (h *Handler) GetMedia(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	att, err := h.svc.GetAttachment(r.Context(), id)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "attachment not found"})
		return
	}

	writeJSON(w, http.StatusOK, att)
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
