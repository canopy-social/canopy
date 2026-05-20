package auth

import (
	"testing"
	"time"
)

func TestHashAndVerifyPassword(t *testing.T) {
	password := "mySecureP@ssw0rd"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	if hash == "" {
		t.Fatal("hash should not be empty")
	}

	// Verify correct password
	valid, err := VerifyPassword(password, hash)
	if err != nil {
		t.Fatalf("VerifyPassword failed: %v", err)
	}
	if !valid {
		t.Fatal("password should be valid")
	}

	// Verify wrong password
	valid, err = VerifyPassword("wrongpassword", hash)
	if err != nil {
		t.Fatalf("VerifyPassword failed: %v", err)
	}
	if valid {
		t.Fatal("wrong password should not be valid")
	}
}

func TestGenerateKeyPair(t *testing.T) {
	kp, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair failed: %v", err)
	}

	if kp.PrivateKeyPEM == "" {
		t.Fatal("private key should not be empty")
	}
	if kp.PublicKeyPEM == "" {
		t.Fatal("public key should not be empty")
	}

	// Ensure they start with correct PEM headers
	if kp.PrivateKeyPEM[:5] != "-----" {
		t.Fatal("private key should be PEM formatted")
	}
	if kp.PublicKeyPEM[:5] != "-----" {
		t.Fatal("public key should be PEM formatted")
	}
}

func TestJWTGenerateAndValidate(t *testing.T) {
	svc := NewJWTService("test-secret-key-that-is-long-enough", 15*time.Minute)

	token, err := svc.GenerateAccessToken("account123", "alice", "example.com", "user")
	if err != nil {
		t.Fatalf("GenerateAccessToken failed: %v", err)
	}

	if token == "" {
		t.Fatal("token should not be empty")
	}

	claims, err := svc.ValidateAccessToken(token)
	if err != nil {
		t.Fatalf("ValidateAccessToken failed: %v", err)
	}

	if claims.Subject != "account123" {
		t.Fatalf("expected subject 'account123', got '%s'", claims.Subject)
	}
	if claims.Username != "alice" {
		t.Fatalf("expected username 'alice', got '%s'", claims.Username)
	}
	if claims.Role != "user" {
		t.Fatalf("expected role 'user', got '%s'", claims.Role)
	}
}

func TestJWTExpiredToken(t *testing.T) {
	svc := NewJWTService("test-secret", -1*time.Minute) // negative TTL = already expired

	token, err := svc.GenerateAccessToken("acc1", "bob", "example.com", "user")
	if err != nil {
		t.Fatalf("GenerateAccessToken failed: %v", err)
	}

	_, err = svc.ValidateAccessToken(token)
	if err == nil {
		t.Fatal("expected error for expired token")
	}
}

func TestJWTWrongSecret(t *testing.T) {
	svc1 := NewJWTService("secret-one", 15*time.Minute)
	svc2 := NewJWTService("secret-two", 15*time.Minute)

	token, _ := svc1.GenerateAccessToken("acc1", "alice", "example.com", "user")

	_, err := svc2.ValidateAccessToken(token)
	if err == nil {
		t.Fatal("expected error for wrong secret")
	}
}

func TestGenerateToken(t *testing.T) {
	token, err := GenerateToken(32)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	if len(token) != 64 { // 32 bytes = 64 hex chars
		t.Fatalf("expected 64 hex chars, got %d", len(token))
	}

	// Two tokens should be different
	token2, _ := GenerateToken(32)
	if token == token2 {
		t.Fatal("tokens should be unique")
	}
}
