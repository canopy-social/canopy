package accounts

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

	r.Get("/api/v1/accounts/{id}", h.GetAccount)
	r.Get("/api/v1/accounts/{id}/followers", h.ListFollowers)
	r.Get("/api/v1/accounts/{id}/following", h.ListFollowing)
	r.Get("/api/v1/accounts/search", h.SearchAccounts)

	r.Group(func(r chi.Router) {
		r.Use(jwtMiddleware)
		r.Get("/api/v1/accounts/verify_credentials", h.VerifyCredentials)
		r.Patch("/api/v1/accounts/update_credentials", h.UpdateCredentials)
		r.Get("/api/v1/accounts/relationships", h.GetRelationships)
		r.Post("/api/v1/accounts/{id}/follow", h.Follow)
		r.Post("/api/v1/accounts/{id}/unfollow", h.Unfollow)
		r.Post("/api/v1/accounts/{id}/block", h.Block)
		r.Post("/api/v1/accounts/{id}/unblock", h.Unblock)
		r.Post("/api/v1/accounts/{id}/mute", h.Mute)
		r.Post("/api/v1/accounts/{id}/unmute", h.Unmute)
		r.Get("/api/v1/blocks", h.ListBlocks)
		r.Get("/api/v1/mutes", h.ListMutes)
		r.Get("/api/v1/follow_requests", h.ListFollowRequests)
		r.Post("/api/v1/follow_requests/{id}/authorize", h.AcceptFollowRequest)
		r.Post("/api/v1/follow_requests/{id}/reject", h.RejectFollowRequest)
	})

	r.Group(func(r chi.Router) {
		r.Use(jwtMiddleware)
		r.Use(auth.RequireRole("admin", "moderator"))
		r.Get("/api/v1/admin/accounts", h.AdminListAccounts)
		r.Get("/api/v1/admin/accounts/{id}", h.AdminGetAccount)
		r.Post("/api/v1/admin/accounts/{id}/action", h.AdminTakeAction)
	})
}

func (h *Handler) GetAccount(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var account *Account
	var err error
	if len(id) > 0 && id[0] == '@' {
		account, err = h.svc.GetByUsername(r.Context(), id[1:])
	} else {
		account, err = h.svc.GetByID(r.Context(), id)
	}
	if err != nil || account == nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "account not found"})
		return
	}
	writeJSON(w, http.StatusOK, account)
}

func (h *Handler) SearchAccounts(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	if q == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "query parameter 'q' is required"})
		return
	}
	limit := parseIntParam(r, "limit", 20)
	accounts, err := h.svc.SearchByUsername(r.Context(), q, limit)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "search failed"})
		return
	}
	writeJSON(w, http.StatusOK, accounts)
}

func (h *Handler) ListFollowers(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	limit := parseIntParam(r, "limit", 40)
	offset := parseIntParam(r, "offset", 0)
	accounts, err := h.svc.ListFollowers(r.Context(), id, limit, offset)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to list followers"})
		return
	}
	writeJSON(w, http.StatusOK, accounts)
}

func (h *Handler) ListFollowing(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	limit := parseIntParam(r, "limit", 40)
	offset := parseIntParam(r, "offset", 0)
	accounts, err := h.svc.ListFollowing(r.Context(), id, limit, offset)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to list following"})
		return
	}
	writeJSON(w, http.StatusOK, accounts)
}

func (h *Handler) VerifyCredentials(w http.ResponseWriter, r *http.Request) {
	accountID := auth.AccountIDFromContext(r.Context())
	account, err := h.svc.GetByID(r.Context(), accountID)
	if err != nil || account == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
		return
	}
	writeJSON(w, http.StatusOK, account)
}

func (h *Handler) UpdateCredentials(w http.ResponseWriter, r *http.Request) {
	accountID := auth.AccountIDFromContext(r.Context())
	var params UpdateProfileParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	account, err := h.svc.UpdateProfile(r.Context(), accountID, &params)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "update failed"})
		return
	}
	writeJSON(w, http.StatusOK, account)
}

func (h *Handler) Follow(w http.ResponseWriter, r *http.Request) {
	accountID := auth.AccountIDFromContext(r.Context())
	targetID := chi.URLParam(r, "id")
	rel, err := h.svc.Follow(r.Context(), accountID, targetID)
	if err != nil {
		writeJSON(w, http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, rel)
}

func (h *Handler) Unfollow(w http.ResponseWriter, r *http.Request) {
	accountID := auth.AccountIDFromContext(r.Context())
	targetID := chi.URLParam(r, "id")
	rel, err := h.svc.Unfollow(r.Context(), accountID, targetID)
	if err != nil {
		writeJSON(w, http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, rel)
}

func (h *Handler) Block(w http.ResponseWriter, r *http.Request) {
	accountID := auth.AccountIDFromContext(r.Context())
	targetID := chi.URLParam(r, "id")
	rel, err := h.svc.Block(r.Context(), accountID, targetID)
	if err != nil {
		writeJSON(w, http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, rel)
}

func (h *Handler) Unblock(w http.ResponseWriter, r *http.Request) {
	accountID := auth.AccountIDFromContext(r.Context())
	targetID := chi.URLParam(r, "id")
	rel, err := h.svc.Unblock(r.Context(), accountID, targetID)
	if err != nil {
		writeJSON(w, http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, rel)
}

func (h *Handler) Mute(w http.ResponseWriter, r *http.Request) {
	accountID := auth.AccountIDFromContext(r.Context())
	targetID := chi.URLParam(r, "id")
	rel, err := h.svc.Mute(r.Context(), accountID, targetID, true)
	if err != nil {
		writeJSON(w, http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, rel)
}

func (h *Handler) Unmute(w http.ResponseWriter, r *http.Request) {
	accountID := auth.AccountIDFromContext(r.Context())
	targetID := chi.URLParam(r, "id")
	rel, err := h.svc.Unmute(r.Context(), accountID, targetID)
	if err != nil {
		writeJSON(w, http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, rel)
}

func (h *Handler) GetRelationships(w http.ResponseWriter, r *http.Request) {
	accountID := auth.AccountIDFromContext(r.Context())
	ids := r.URL.Query()["id[]"]
	if len(ids) == 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "id[] parameter required"})
		return
	}
	var relationships []*Relationship
	for _, targetID := range ids {
		rel, err := h.svc.GetRelationship(r.Context(), accountID, targetID)
		if err != nil {
			continue
		}
		relationships = append(relationships, rel)
	}
	writeJSON(w, http.StatusOK, relationships)
}

func (h *Handler) ListBlocks(w http.ResponseWriter, r *http.Request) {
	accountID := auth.AccountIDFromContext(r.Context())
	limit := parseIntParam(r, "limit", 40)
	offset := parseIntParam(r, "offset", 0)
	accounts, err := h.svc.repo.ListBlocks(r.Context(), accountID, limit, offset)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to list blocks"})
		return
	}
	writeJSON(w, http.StatusOK, accounts)
}

func (h *Handler) ListMutes(w http.ResponseWriter, r *http.Request) {
	accountID := auth.AccountIDFromContext(r.Context())
	limit := parseIntParam(r, "limit", 40)
	offset := parseIntParam(r, "offset", 0)
	accounts, err := h.svc.repo.ListMutes(r.Context(), accountID, limit, offset)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to list mutes"})
		return
	}
	writeJSON(w, http.StatusOK, accounts)
}

func (h *Handler) ListFollowRequests(w http.ResponseWriter, r *http.Request) {
	accountID := auth.AccountIDFromContext(r.Context())
	limit := parseIntParam(r, "limit", 40)
	offset := parseIntParam(r, "offset", 0)
	accounts, err := h.svc.repo.ListPendingRequests(r.Context(), accountID, limit, offset)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to list follow requests"})
		return
	}
	writeJSON(w, http.StatusOK, accounts)
}

func (h *Handler) AcceptFollowRequest(w http.ResponseWriter, r *http.Request) {
	accountID := auth.AccountIDFromContext(r.Context())
	followerID := chi.URLParam(r, "id")
	if err := h.svc.AcceptFollowRequest(r.Context(), accountID, followerID); err != nil {
		writeJSON(w, http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
		return
	}
	rel, _ := h.svc.GetRelationship(r.Context(), accountID, followerID)
	writeJSON(w, http.StatusOK, rel)
}

func (h *Handler) RejectFollowRequest(w http.ResponseWriter, r *http.Request) {
	accountID := auth.AccountIDFromContext(r.Context())
	followerID := chi.URLParam(r, "id")
	if err := h.svc.RejectFollowRequest(r.Context(), accountID, followerID); err != nil {
		writeJSON(w, http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
		return
	}
	rel, _ := h.svc.GetRelationship(r.Context(), accountID, followerID)
	writeJSON(w, http.StatusOK, rel)
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

type AdminActionRequest struct {
	Type string `json:"type"`
	Role string `json:"role,omitempty"`
}

func (h *Handler) AdminListAccounts(w http.ResponseWriter, r *http.Request) {
	limit := parseIntParam(r, "limit", 40)
	offset := parseIntParam(r, "offset", 0)
	accs, err := h.svc.ListLocal(r.Context(), limit, offset)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to list accounts"})
		return
	}
	writeJSON(w, http.StatusOK, accs)
}

func (h *Handler) AdminGetAccount(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	account, err := h.svc.GetByID(r.Context(), id)
	if err != nil || account == nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "account not found"})
		return
	}
	writeJSON(w, http.StatusOK, account)
}

func (h *Handler) AdminTakeAction(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req AdminActionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	account, err := h.svc.GetByID(r.Context(), id)
	if err != nil || account == nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "account not found"})
		return
	}
	switch req.Type {
	case "suspend":
		err = h.svc.SetSuspended(r.Context(), id, true)
	case "unsuspend":
		err = h.svc.SetSuspended(r.Context(), id, false)
	case "silence":
		err = h.svc.SetSilenced(r.Context(), id, true)
	case "unsilence":
		err = h.svc.SetSilenced(r.Context(), id, false)
	case "change_role":
		err = h.svc.SetRole(r.Context(), id, req.Role)
	default:
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid action type"})
		return
	}
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	updated, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to fetch updated account"})
		return
	}
	writeJSON(w, http.StatusOK, updated)
}
