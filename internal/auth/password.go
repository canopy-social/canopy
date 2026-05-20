package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"

	"golang.org/x/crypto/argon2"
)

// Argon2id parameters — OWASP recommended.
const (
	argonTime    = 2
	argonMemory  = 64 * 1024 // 64 MB
	argonThreads = 4
	argonKeyLen  = 32
	saltLen      = 16
)

// HashPassword hashes a password using Argon2id.
// Returns the encoded hash string: $argon2id$v=19$m=65536,t=2,p=4$<salt>$<hash>
func HashPassword(password string) (string, error) {
	salt := make([]byte, saltLen)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("generating salt: %w", err)
	}

	hash := argon2.IDKey([]byte(password), salt, argonTime, argonMemory, argonThreads, argonKeyLen)

	encoded := fmt.Sprintf("$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s",
		argonMemory, argonTime, argonThreads,
		hex.EncodeToString(salt),
		hex.EncodeToString(hash),
	)
	return encoded, nil
}

// VerifyPassword checks a password against an Argon2id hash.
func VerifyPassword(password, encodedHash string) (bool, error) {
	var memory uint32
	var time uint32
	var threads uint8
	var saltHex, hashHex string

	_, err := fmt.Sscanf(encodedHash, "$argon2id$v=19$m=%d,t=%d,p=%d$%s",
		&memory, &time, &threads, &saltHex)
	if err != nil {
		return false, fmt.Errorf("parsing hash: %w", err)
	}

	// Split saltHex on '$' to separate salt and hash
	parts := splitOnDollar(saltHex)
	if len(parts) != 2 {
		return false, fmt.Errorf("invalid hash format")
	}
	saltHex = parts[0]
	hashHex = parts[1]

	salt, err := hex.DecodeString(saltHex)
	if err != nil {
		return false, fmt.Errorf("decoding salt: %w", err)
	}

	expectedHash, err := hex.DecodeString(hashHex)
	if err != nil {
		return false, fmt.Errorf("decoding hash: %w", err)
	}

	computed := argon2.IDKey([]byte(password), salt, time, memory, threads, uint32(len(expectedHash)))

	if len(computed) != len(expectedHash) {
		return false, nil
	}
	// Constant-time comparison
	result := byte(0)
	for i := range computed {
		result |= computed[i] ^ expectedHash[i]
	}
	return result == 0, nil
}

func splitOnDollar(s string) []string {
	var parts []string
	current := ""
	for _, c := range s {
		if c == '$' {
			parts = append(parts, current)
			current = ""
		} else {
			current += string(c)
		}
	}
	if current != "" {
		parts = append(parts, current)
	}
	return parts
}

// KeyPair holds PEM-encoded RSA keys for an actor.
type KeyPair struct {
	PrivateKeyPEM string
	PublicKeyPEM  string
}

// GenerateKeyPair generates a 2048-bit RSA key pair for an AP actor.
func GenerateKeyPair() (*KeyPair, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, fmt.Errorf("generating RSA key: %w", err)
	}

	privatePEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("marshaling public key: %w", err)
	}
	publicPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	})

	return &KeyPair{
		PrivateKeyPEM: string(privatePEM),
		PublicKeyPEM:  string(publicPEM),
	}, nil
}

// GenerateToken generates a cryptographically random hex token of the given byte length.
func GenerateToken(byteLength int) (string, error) {
	b := make([]byte, byteLength)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
