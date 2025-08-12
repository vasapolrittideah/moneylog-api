package security

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/argon2"
)

const (
	memory      = 64 * 1024 // 64 MB
	iterations  = 3
	parallelism = 2
	saltLength  = 16
	keyLength   = 32
)

// HashPassword hashes a password using Argon2id.
func HashPassword(password string) (string, error) {
	salt := make([]byte, saltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, iterations, memory, parallelism, keyLength)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	return fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, memory, iterations, parallelism, b64Salt, b64Hash), nil
}

// VerifyPassword verifies a password against an Argon2 hash.
func VerifyPassword(password, encodedHash string) (bool, error) {
	salt, hash, err := parseHash(encodedHash)
	if err != nil {
		return false, err
	}

	otherHash := argon2.IDKey([]byte(password), salt, iterations, memory, parallelism, keyLength)
	return subtle.ConstantTimeCompare(hash, otherHash) == 1, nil
}

// parseHash parses an Argon2 hash into its salt and hash components.
func parseHash(encodedHash string) ([]byte, []byte, error) {
	var version, m, t, p int
	var saltStr, hashStr string

	_, err := fmt.Sscanf(encodedHash, "$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		&version, &m, &t, &p, &saltStr, &hashStr)
	if err != nil {
		return nil, nil, err
	}

	salt, err := base64.RawStdEncoding.DecodeString(saltStr)
	if err != nil {
		return nil, nil, err
	}

	hash, err := base64.RawStdEncoding.DecodeString(hashStr)
	if err != nil {
		return nil, nil, err
	}

	return salt, hash, nil
}
