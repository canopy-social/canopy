package posts

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
		r.Use(optionalJWT)
		r.Get("/api/v1/posts/{id}", h.GetPost)
		r.Get("/api/v1/posts/{id}/context", h.GetContext)
		r.Get("/api/v1/timelines/public", h.PublicTimeline)
		r.Get("/api/v1/timelines/tag/{tag}", h.TagTimeline)
	})

	r.Group(func(r chi.Router) {
		r.Use(jwtMiddleware)
		r.Post("/api/v1/posts", h.CreatePost)
		r.Delete("/api/v1/posts/{id}", h.DeletePost)
		r.Put("/api/v1/posts/{id}", h.EditPost)
		r.Post("/api/v1/posts/{id}/like", h.LikePost)
		r.Post("/api/v1/posts/{id}/unlike", h.UnlikePost)
		r.Post("/api/v1/posts/{id}/boost", h.BoostPost)
		r.Post("/api/v1/posts/{id}/unboost", h.UnboostPost)
		r.Post("/api/v1/posts/{id}/pin", h.PinPost)
		r.Post("/api/v1/posts/{id}/unpin", h.UnpinPost)
	})
}

func (h *Handler) GetPost(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	post, err := h.svc.GetByID(r.Context(), id)
	if err != nil || post == nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "post not found"})
		return
	}
	writeJSON(w, http.StatusOK, post)
}

func (h *Handler) GetContext(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ctx, err := h.svc.GetThreadContext(r.Context(), id)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, ctx)
}

func (h *Handler) CreatePost(w http.ResponseWriter, r *http.Request) {
	accountID := auth.AccountIDFromContext(r.Context())
	var params CreatePostParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	post, err := h.svc.Create(r.Context(), accountID, &params)
	if err != nil {
		writeJSON(w, http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, post)
}

func (h *Handler) DeletePost(w http.ResponseWriter, r *http.Request) {
	accountID := auth.AccountIDFromContext(r.Context())
	postID := chi.URLParam(r, "id")
	if err := h.svc.Delete(r.Context(), postID, accountID); err != nil {
		writeJSON(w, http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "post deleted"})
}

func (h *Handler) EditPost(w http.ResponseWriter, r *http.Request) {
	accountID := auth.AccountIDFromContext(r.Context())
	postID := chi.URLParam(r, "id")
	var params UpdatePostParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	post, err := h.svc.Edit(r.Context(), postID, accountID, &params)
	if err != nil {
		writeJSON(w, http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, post)
}

func (h *Handler) LikePost(w http.ResponseWriter, r *http.Request) {
	accountID := auth.AccountIDFromContext(r.Context())
	postID := chi.URLParam(r, "id")
	if err := h.svc.Like(r.Context(), postID, accountID); err != nil {
		writeJSON(w, http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "liked"})
}

func (h *Handler) UnlikePost(w http.ResponseWriter, r *http.Request) {
	accountID := auth.AccountIDFromContext(r.Context())
	postID := chi.URLParam(r, "id")
	if err := h.svc.Unlike(r.Context(), postID, accountID); err != nil {
		writeJSON(w, http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "unliked"})
}

func (h *Handler) BoostPost(w http.ResponseWriter, r *http.Request) {
	accountID := auth.AccountIDFromContext(r.Context())
	postID := chi.URLParam(r, "id")
	if err := h.svc.Boost(r.Context(), postID, accountID); err != nil {
		writeJSON(w, http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "boosted"})
}

func (h *Handler) UnboostPost(w http.ResponseWriter, r *http.Request) {
	accountID := auth.AccountIDFromContext(r.Context())
	postID := chi.URLParam(r, "id")
	if err := h.svc.Unboost(r.Context(), postID, accountID); err != nil {
		writeJSON(w, http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "unboosted"})
}

func (h *Handler) PinPost(w http.ResponseWriter, r *http.Request) {
	accountID := auth.AccountIDFromContext(r.Context())
	postID := chi.URLParam(r, "id")
	if err := h.svc.repo.Pin(r.Context(), postID, accountID); err != nil {
		writeJSON(w, http.StatusUnprocessableEntity, map[string]string{"error": "failed to pin"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "pinned"})
}

func (h *Handler) UnpinPost(w http.ResponseWriter, r *http.Request) {
	accountID := auth.AccountIDFromContext(r.Context())
	postID := chi.URLParam(r, "id")
	if err := h.svc.repo.Unpin(r.Context(), postID, accountID); err != nil {
		writeJSON(w, http.StatusUnprocessableEntity, map[string]string{"error": "failed to unpin"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "unpinned"})
}

func (h *Handler) PublicTimeline(w http.ResponseWriter, r *http.Request) {
	limit := parseIntParam(r, "limit", 20)
	offset := parseIntParam(r, "offset", 0)
	posts, err := h.svc.ListPublicTimeline(r.Context(), limit, offset)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to load timeline"})
		return
	}
	writeJSON(w, http.StatusOK, posts)
}

func (h *Handler) TagTimeline(w http.ResponseWriter, r *http.Request) {
	tag := chi.URLParam(r, "tag")
	limit := parseIntParam(r, "limit", 20)
	offset := parseIntParam(r, "offset", 0)
	posts, err := h.svc.SearchByTag(r.Context(), tag, limit, offset)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to load tag timeline"})
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
