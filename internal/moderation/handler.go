package moderation

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

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
		r.Post("/api/v1/reports", h.CreateReport)
	})

	r.Group(func(r chi.Router) {
		r.Use(jwtMiddleware)
		r.Use(auth.RequireRole("admin", "moderator"))
		r.Get("/api/v1/admin/reports", h.ListReports)
		r.Post("/api/v1/admin/reports/{id}/resolve", h.ResolveReport)
		r.Get("/api/v1/admin/invites", h.ListInvites)
		r.Post("/api/v1/admin/invites", h.CreateInvite)
		r.Post("/api/v1/admin/invites/{id}/revoke", h.RevokeInvite)
		r.Get("/api/v1/admin/settings", h.GetSetting)
		r.Get("/api/v1/admin/instance_blocks", h.ListInstanceBlocks)
	})

	r.Group(func(r chi.Router) {
		r.Use(jwtMiddleware)
		r.Use(auth.RequireRole("admin"))
		r.Put("/api/v1/admin/settings", h.PutSetting)
		r.Post("/api/v1/admin/instance_blocks", h.BlockInstance)
		r.Delete("/api/v1/admin/instance_blocks/{id}", h.UnblockInstance)
	})
}

type CreateReportRequest struct {
	AccountID string `json:"account_id"`
	PostID    string `json:"post_id,omitempty"`
	EssayID   string `json:"essay_id,omitempty"`
	Category  string `json:"category"`
	Comment   string `json:"comment,omitempty"`
}

func (h *Handler) CreateReport(w http.ResponseWriter, r *http.Request) {
	reporterID := auth.AccountIDFromContext(r.Context())
	var req CreateReportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	if req.Category == "" {
		writeJSON(w, http.StatusUnprocessableEntity, map[string]string{"error": "category is required"})
		return
	}

	var targetAccountID *string
	if req.AccountID != "" {
		targetAccountID = &req.AccountID
	}
	var targetPostID *string
	if req.PostID != "" {
		targetPostID = &req.PostID
	}
	var targetEssayID *string
	if req.EssayID != "" {
		targetEssayID = &req.EssayID
	}

	report, err := h.svc.ReportContent(r.Context(), reporterID, targetAccountID, targetPostID, targetEssayID, req.Category, req.Comment)
	if err != nil {
		writeJSON(w, http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, report)
}

func (h *Handler) ListReports(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	if status == "" {
		status = "open"
	}
	limit := parseIntParam(r, "limit", 40)
	offset := parseIntParam(r, "offset", 0)

	reports, err := h.svc.ListReports(r.Context(), status, limit, offset)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to list reports"})
		return
	}

	writeJSON(w, http.StatusOK, reports)
}

type ResolveReportRequest struct {
	ActionTaken string `json:"action_taken"`
}

func (h *Handler) ResolveReport(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	resolverID := auth.AccountIDFromContext(r.Context())
	var req ResolveReportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	report, err := h.svc.ResolveReport(r.Context(), id, resolverID, req.ActionTaken)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "report not found"})
		return
	}

	writeJSON(w, http.StatusOK, report)
}

func (h *Handler) ListInvites(w http.ResponseWriter, r *http.Request) {
	limit := parseIntParam(r, "limit", 40)
	offset := parseIntParam(r, "offset", 0)

	invites, err := h.svc.ListInvites(r.Context(), limit, offset)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to list invites"})
		return
	}

	writeJSON(w, http.StatusOK, invites)
}

type CreateInviteRequest struct {
	MaxUses       *int `json:"max_uses,omitempty"`
	ExpiresInDays *int `json:"expires_in_days,omitempty"`
}

func (h *Handler) CreateInvite(w http.ResponseWriter, r *http.Request) {
	creatorID := auth.AccountIDFromContext(r.Context())
	var req CreateInviteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	var expiresAt *time.Time
	if req.ExpiresInDays != nil {
		t := time.Now().AddDate(0, 0, *req.ExpiresInDays)
		expiresAt = &t
	}

	invite, err := h.svc.CreateInvite(r.Context(), creatorID, req.MaxUses, expiresAt)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to create invite"})
		return
	}

	writeJSON(w, http.StatusCreated, invite)
}

func (h *Handler) RevokeInvite(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.svc.RevokeInvite(r.Context(), id); err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "invite not found"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "revoked"})
}

func (h *Handler) GetSetting(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	if key == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "key parameter is required"})
		return
	}

	var val interface{}
	if err := h.svc.GetSetting(r.Context(), key, &val); err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "setting not found"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{"key": key, "value": val})
}

type PutSettingRequest struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

func (h *Handler) PutSetting(w http.ResponseWriter, r *http.Request) {
	var req PutSettingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	if req.Key == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "key is required"})
		return
	}

	if err := h.svc.PutSetting(r.Context(), req.Key, req.Value); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "saved"})
}

func (h *Handler) ListInstanceBlocks(w http.ResponseWriter, r *http.Request) {
	limit := parseIntParam(r, "limit", 40)
	offset := parseIntParam(r, "offset", 0)

	blocks, err := h.svc.ListInstanceBlocks(r.Context(), limit, offset)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to list blocks"})
		return
	}

	writeJSON(w, http.StatusOK, blocks)
}

type BlockInstanceRequest struct {
	Domain        string `json:"domain"`
	Severity      string `json:"severity"`
	Reason        string `json:"reason,omitempty"`
	RejectMedia   bool   `json:"reject_media"`
	RejectReports bool   `json:"reject_reports"`
}

func (h *Handler) BlockInstance(w http.ResponseWriter, r *http.Request) {
	creatorID := auth.AccountIDFromContext(r.Context())
	var req BlockInstanceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	block, err := h.svc.BlockInstance(r.Context(), creatorID, req.Domain, req.Severity, req.Reason, req.RejectMedia, req.RejectReports)
	if err != nil {
		writeJSON(w, http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusCreated, block)
}

func (h *Handler) UnblockInstance(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.svc.UnblockInstance(r.Context(), id); err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "block not found"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "unblocked"})
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
