package themes

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

		r.Post("/api/v1/themes", h.CreatePageTheme)
		r.Get("/api/v1/themes/mine", h.GetMyTheme)
		r.Patch("/api/v1/themes/mine", h.UpdateMyTheme)
		r.Get("/api/v1/themes/mine/css", h.GetMyThemeCSS)

		r.Get("/api/v1/themes/versions", h.ListVersions)
		r.Post("/api/v1/themes/versions", h.SaveNamedVersion)
		r.Post("/api/v1/themes/versions/{id}/restore", h.RestoreVersion)
		r.Delete("/api/v1/themes/versions/{id}", h.DeleteVersion)

		r.Post("/api/v1/post_styles", h.CreatePostStyle)
		r.Get("/api/v1/post_styles", h.ListPostStyles)
		r.Get("/api/v1/post_styles/{id}", h.GetPostStyle)
		r.Delete("/api/v1/post_styles/{id}", h.DeletePostStyle)

		r.Get("/api/v1/essays/{id}/theme", h.GetEssayTheme)
		r.Patch("/api/v1/essays/{id}/theme", h.UpdateEssayTheme)
	})

	r.Get("/api/v1/themes/{accountID}", h.GetPublicTheme)
	r.Get("/api/v1/themes/{accountID}/css", h.GetPublicThemeCSS)

	r.Group(func(r chi.Router) {
		r.Use(jwtMiddleware)
		r.Use(auth.RequireRole("admin"))
		r.Get("/api/v1/admin/server_theme", h.GetServerTheme)
		r.Patch("/api/v1/admin/server_theme", h.UpdateServerTheme)
	})

	r.Get("/api/v1/fonts", h.ListAllowedFonts)
}

func (h *Handler) CreatePageTheme(w http.ResponseWriter, r *http.Request) {
	accountID := auth.AccountIDFromContext(r.Context())
	theme, err := h.svc.CreatePageTheme(r.Context(), accountID)
	if err != nil {
		writeJSON(w, http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, theme)
}

func (h *Handler) GetMyTheme(w http.ResponseWriter, r *http.Request) {
	accountID := auth.AccountIDFromContext(r.Context())
	theme, err := h.svc.GetPageTheme(r.Context(), accountID)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "no theme found"})
		return
	}
	writeJSON(w, http.StatusOK, theme)
}

func (h *Handler) UpdateMyTheme(w http.ResponseWriter, r *http.Request) {
	accountID := auth.AccountIDFromContext(r.Context())
	var req PageThemeUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	theme, err := h.svc.UpdatePageTheme(r.Context(), accountID, &req)
	if err != nil {
		writeJSON(w, http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, theme)
}

func (h *Handler) GetMyThemeCSS(w http.ResponseWriter, r *http.Request) {
	accountID := auth.AccountIDFromContext(r.Context())
	css, err := h.svc.GetThemeCSS(r.Context(), accountID)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "no theme found"})
		return
	}
	w.Header().Set("Content-Type", "text/css; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=300")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(css))
}

func (h *Handler) GetPublicTheme(w http.ResponseWriter, r *http.Request) {
	accountID := chi.URLParam(r, "accountID")
	theme, err := h.svc.GetPageTheme(r.Context(), accountID)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "no theme found"})
		return
	}
	writeJSON(w, http.StatusOK, theme)
}

func (h *Handler) GetPublicThemeCSS(w http.ResponseWriter, r *http.Request) {
	accountID := chi.URLParam(r, "accountID")
	css, err := h.svc.GetThemeCSS(r.Context(), accountID)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "no theme found"})
		return
	}
	w.Header().Set("Content-Type", "text/css; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=300")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(css))
}

func (h *Handler) ListVersions(w http.ResponseWriter, r *http.Request) {
	accountID := auth.AccountIDFromContext(r.Context())
	limit := parseIntParam(r, "limit", 40)
	offset := parseIntParam(r, "offset", 0)

	versions, err := h.svc.ListVersions(r.Context(), accountID, limit, offset)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to list versions"})
		return
	}
	writeJSON(w, http.StatusOK, versions)
}

type SaveVersionRequest struct {
	Label string `json:"label"`
}

func (h *Handler) SaveNamedVersion(w http.ResponseWriter, r *http.Request) {
	accountID := auth.AccountIDFromContext(r.Context())
	var req SaveVersionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	if req.Label == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "label is required"})
		return
	}

	version, err := h.svc.SaveNamedVersion(r.Context(), accountID, req.Label)
	if err != nil {
		writeJSON(w, http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, version)
}

func (h *Handler) RestoreVersion(w http.ResponseWriter, r *http.Request) {
	accountID := auth.AccountIDFromContext(r.Context())
	versionID := chi.URLParam(r, "id")

	theme, err := h.svc.RestoreVersion(r.Context(), accountID, versionID)
	if err != nil {
		writeJSON(w, http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, theme)
}

func (h *Handler) DeleteVersion(w http.ResponseWriter, r *http.Request) {
	accountID := auth.AccountIDFromContext(r.Context())
	versionID := chi.URLParam(r, "id")

	if err := h.svc.DeleteVersion(r.Context(), accountID, versionID); err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

func (h *Handler) CreatePostStyle(w http.ResponseWriter, r *http.Request) {
	accountID := auth.AccountIDFromContext(r.Context())
	var req CreatePostStyleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	style, err := h.svc.CreatePostStyle(r.Context(), accountID, &req)
	if err != nil {
		writeJSON(w, http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, style)
}

func (h *Handler) GetPostStyle(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	style, err := h.svc.GetPostStyle(r.Context(), id)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "post style not found"})
		return
	}
	writeJSON(w, http.StatusOK, style)
}

func (h *Handler) ListPostStyles(w http.ResponseWriter, r *http.Request) {
	accountID := auth.AccountIDFromContext(r.Context())
	limit := parseIntParam(r, "limit", 40)
	offset := parseIntParam(r, "offset", 0)

	styles, err := h.svc.ListPostStyles(r.Context(), accountID, limit, offset)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to list styles"})
		return
	}
	writeJSON(w, http.StatusOK, styles)
}

func (h *Handler) DeletePostStyle(w http.ResponseWriter, r *http.Request) {
	accountID := auth.AccountIDFromContext(r.Context())
	id := chi.URLParam(r, "id")

	if err := h.svc.DeletePostStyle(r.Context(), accountID, id); err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

func (h *Handler) GetEssayTheme(w http.ResponseWriter, r *http.Request) {
	essayID := chi.URLParam(r, "id")
	theme, err := h.svc.GetOrCreateEssayTheme(r.Context(), essayID)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "essay theme not found"})
		return
	}
	writeJSON(w, http.StatusOK, theme)
}

type UpdateEssayThemeRequest struct {
	Colors json.RawMessage `json:"colors,omitempty"`
	Fonts  json.RawMessage `json:"fonts,omitempty"`
	Layout json.RawMessage `json:"layout,omitempty"`
	BGType string          `json:"bg_type,omitempty"`
}

func (h *Handler) UpdateEssayTheme(w http.ResponseWriter, r *http.Request) {
	essayID := chi.URLParam(r, "id")
	var req UpdateEssayThemeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	theme, err := h.svc.UpdateEssayTheme(r.Context(), essayID, req.Colors, req.Fonts, req.Layout, req.BGType)
	if err != nil {
		writeJSON(w, http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, theme)
}

func (h *Handler) GetServerTheme(w http.ResponseWriter, r *http.Request) {
	theme, err := h.svc.GetServerTheme(r.Context())
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "server theme not configured"})
		return
	}
	writeJSON(w, http.StatusOK, theme)
}

type UpdateServerThemeRequest struct {
	Colors json.RawMessage `json:"colors,omitempty"`
	Fonts  json.RawMessage `json:"fonts,omitempty"`
	Layout json.RawMessage `json:"layout,omitempty"`
	BGType string          `json:"bg_type,omitempty"`
}

func (h *Handler) UpdateServerTheme(w http.ResponseWriter, r *http.Request) {
	adminID := auth.AccountIDFromContext(r.Context())
	var req UpdateServerThemeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	theme, err := h.svc.UpdateServerTheme(r.Context(), adminID, req.Colors, req.Fonts, req.Layout, req.BGType)
	if err != nil {
		writeJSON(w, http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, theme)
}

func (h *Handler) ListAllowedFonts(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"fonts": ListAllowedFonts(),
	})
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
