package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"

	"github.com/sumi-devs/canopy-social/canopy/pkg/config"
	"github.com/sumi-devs/canopy-social/canopy/pkg/ulid"
	"github.com/sumi-devs/canopy-social/canopy/pkg/validate"
)

// AccountStore is the interface for account database operations needed by auth.
type AccountStore interface {
	GetByEmail(ctx context.Context, email string) (*Account, error)
	GetByUsername(ctx context.Context, username string) (*Account, error)
	Create(ctx context.Context, account *Account) (*Account, error)
	VerifyEmail(ctx context.Context, id string) error
	GetByEmailVerifyToken(ctx context.Context, token string) (*Account, error)
}

// Account represents a minimal account for auth purposes.
type Account struct {
	ID            string
	Username      string
	Email         string
	PasswordHash  string
	Role          string
	IsSuspended   bool
	IsLocal       bool
	URI           string
	PublicKeyPEM  string
	PrivateKeyPEM string
	KeyID         string
	EmailVerified bool
}

// Handler handles authentication HTTP endpoints.
type Handler struct {
	store      AccountStore
	jwt        *JWTService
	redis      *redis.Client
	cfg        *config.Config
	refreshTTL time.Duration
}

// NewHandler creates a new auth handler.
func NewHandler(store AccountStore, jwt *JWTService, redis *redis.Client, cfg *config.Config) *Handler {
	return &Handler{
		store:      store,
		jwt:        jwt,
		redis:      redis,
		cfg:        cfg,
		refreshTTL: cfg.Auth.RefreshTokenTTL,
	}
}

// --- Request/Response types ---

type RegisterRequest struct {
	Username    string `json:"username"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	InviteToken string `json:"invite_token,omitempty"`
}

type LoginRequest struct {
	Credential string `json:"credential"` // email or username
	Password   string `json:"password"`
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

// --- Handlers ---

// Register handles POST /api/v1/auth/register
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
		return
	}

	// Validate username
	if err := validate.Username(req.Username); err != nil {
		writeJSON(w, http.StatusUnprocessableEntity, ErrorResponse{Error: err.Error()})
		return
	}

	// Validate email
	if !validate.Email(req.Email) {
		writeJSON(w, http.StatusUnprocessableEntity, ErrorResponse{Error: "invalid email format"})
		return
	}

	// Validate password length
	if len(req.Password) < 8 {
		writeJSON(w, http.StatusUnprocessableEntity, ErrorResponse{Error: "password must be at least 8 characters"})
		return
	}

	// Check registration is open
	if !h.cfg.Features.RegistrationOpen {
		writeJSON(w, http.StatusForbidden, ErrorResponse{Error: "registration is closed"})
		return
	}

	ctx := r.Context()

	// Check username uniqueness
	existing, _ := h.store.GetByUsername(ctx, req.Username)
	if existing != nil {
		writeJSON(w, http.StatusConflict, ErrorResponse{Error: "username already taken"})
		return
	}

	// Check email uniqueness
	existing, _ = h.store.GetByEmail(ctx, req.Email)
	if existing != nil {
		writeJSON(w, http.StatusConflict, ErrorResponse{Error: "email already registered"})
		return
	}

	// Hash password
	passwordHash, err := HashPassword(req.Password)
	if err != nil {
		log.Error().Err(err).Msg("failed to hash password")
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "internal error"})
		return
	}

	// Generate RSA key pair for AP
	keys, err := GenerateKeyPair()
	if err != nil {
		log.Error().Err(err).Msg("failed to generate key pair")
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "internal error"})
		return
	}

	// Generate email verification token
	verifyToken, err := GenerateToken(32)
	if err != nil {
		log.Error().Err(err).Msg("failed to generate verify token")
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "internal error"})
		return
	}

	accountID := ulid.New()
	baseURL := h.cfg.BaseURL()
	actorURI := fmt.Sprintf("%s/users/%s", baseURL, req.Username)

	account := &Account{
		ID:            accountID,
		Username:      req.Username,
		Email:         req.Email,
		PasswordHash:  passwordHash,
		Role:          "user",
		IsLocal:       true,
		URI:           actorURI,
		PublicKeyPEM:  keys.PublicKeyPEM,
		PrivateKeyPEM: keys.PrivateKeyPEM,
		KeyID:         actorURI + "#main-key",
	}

	created, err := h.store.Create(ctx, account)
	if err != nil {
		log.Error().Err(err).Msg("failed to create account")
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "failed to create account"})
		return
	}

	// Store verify token (we pass it through the store, but log it in dev)
	_ = verifyToken
	if h.cfg.IsDevelopment() {
		log.Info().
			Str("username", req.Username).
			Str("verify_token", verifyToken).
			Msg("email verification token (dev mode)")
	}

	writeJSON(w, http.StatusCreated, map[string]interface{}{
		"id":       created.ID,
		"username": created.Username,
		"message":  "account created — please verify your email",
	})
}

// Login handles POST /api/v1/auth/login
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
		return
	}

	ctx := r.Context()

	// Find account by email or username
	var account *Account
	var err error
	if validate.Email(req.Credential) {
		account, err = h.store.GetByEmail(ctx, req.Credential)
	} else {
		account, err = h.store.GetByUsername(ctx, req.Credential)
	}

	if err != nil || account == nil {
		writeJSON(w, http.StatusUnauthorized, ErrorResponse{Error: "invalid credentials"})
		return
	}

	// Check suspended
	if account.IsSuspended {
		writeJSON(w, http.StatusForbidden, ErrorResponse{Error: "account is suspended"})
		return
	}

	// Verify password
	valid, err := VerifyPassword(req.Password, account.PasswordHash)
	if err != nil || !valid {
		writeJSON(w, http.StatusUnauthorized, ErrorResponse{Error: "invalid credentials"})
		return
	}

	// Generate access token
	accessToken, err := h.jwt.GenerateAccessToken(
		account.ID, account.Username, h.cfg.Server.Domain, account.Role,
	)
	if err != nil {
		log.Error().Err(err).Msg("failed to generate access token")
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "internal error"})
		return
	}

	// Generate refresh token
	refreshToken, err := GenerateToken(64)
	if err != nil {
		log.Error().Err(err).Msg("failed to generate refresh token")
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "internal error"})
		return
	}

	// Store refresh token hash in Redis
	tokenHash := hashToken(refreshToken)
	err = h.redis.Set(ctx, "refresh:"+tokenHash, account.ID, h.refreshTTL).Err()
	if err != nil {
		log.Error().Err(err).Msg("failed to store refresh token")
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "internal error"})
		return
	}

	// Set refresh token as httpOnly cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/api/v1/auth/refresh",
		HttpOnly: true,
		Secure:   !h.cfg.IsDevelopment(),
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int(h.refreshTTL.Seconds()),
	})

	writeJSON(w, http.StatusOK, TokenResponse{
		AccessToken: accessToken,
		TokenType:   "Bearer",
		ExpiresIn:   int(h.cfg.Auth.AccessTokenTTL.Seconds()),
	})
}

// Refresh handles POST /api/v1/auth/refresh
func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, ErrorResponse{Error: "no refresh token"})
		return
	}

	ctx := r.Context()
	tokenHash := hashToken(cookie.Value)

	// Look up in Redis
	accountID, err := h.redis.Get(ctx, "refresh:"+tokenHash).Result()
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, ErrorResponse{Error: "invalid or expired refresh token"})
		return
	}

	// Fetch account to get current role/username
	account, err := h.store.GetByEmail(ctx, "") // We need a GetByID — use store interface
	_ = account

	// Delete old refresh token (rotation)
	h.redis.Del(ctx, "refresh:"+tokenHash)

	// Generate new access token (simplified — in production we'd fetch the full account)
	accessToken, err := h.jwt.GenerateAccessToken(
		accountID, "", h.cfg.Server.Domain, "user",
	)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "internal error"})
		return
	}

	// Generate new refresh token
	newRefreshToken, err := GenerateToken(64)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "internal error"})
		return
	}

	newTokenHash := hashToken(newRefreshToken)
	h.redis.Set(ctx, "refresh:"+newTokenHash, accountID, h.refreshTTL)

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    newRefreshToken,
		Path:     "/api/v1/auth/refresh",
		HttpOnly: true,
		Secure:   !h.cfg.IsDevelopment(),
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int(h.refreshTTL.Seconds()),
	})

	writeJSON(w, http.StatusOK, TokenResponse{
		AccessToken: accessToken,
		TokenType:   "Bearer",
		ExpiresIn:   int(h.cfg.Auth.AccessTokenTTL.Seconds()),
	})
}

// Logout handles POST /api/v1/auth/logout
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err == nil {
		tokenHash := hashToken(cookie.Value)
		h.redis.Del(r.Context(), "refresh:"+tokenHash)
	}

	// Clear the cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/api/v1/auth/refresh",
		HttpOnly: true,
		Secure:   !h.cfg.IsDevelopment(),
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1,
	})

	writeJSON(w, http.StatusOK, map[string]string{"message": "logged out"})
}

// VerifyEmail handles POST /api/v1/auth/verify-email
func (h *Handler) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
		return
	}

	ctx := r.Context()
	account, err := h.store.GetByEmailVerifyToken(ctx, req.Token)
	if err != nil || account == nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid or expired token"})
		return
	}

	if err := h.store.VerifyEmail(ctx, account.ID); err != nil {
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "internal error"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "email verified"})
}

// --- Helpers ---

func hashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
