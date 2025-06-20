package utils

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/argon2"
)

type Argon2Config struct {
	Memory     uint32
	Time       uint32
	Threads    uint8
	KeyLength  uint32
	SaltLength uint32
}

var DefaultArgon2Config = Argon2Config{
	Memory:     64 * 1024,
	Time:       3,
	Threads:    2,
	KeyLength:  32,
	SaltLength: 16,
}

func GeneratePasswordHash(password string, config *Argon2Config) (hashBase64 string, saltBase64 string, err error) {
	salt := make([]byte, config.SaltLength)
	_, err = rand.Read(salt)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate salt: %w", err)
	}

	hash := argon2.IDKey([]byte(password), salt, config.Time, config.Memory, config.Threads, config.KeyLength)

	hashBase64 = base64.RawStdEncoding.EncodeToString(hash)
	saltBase64 = base64.RawStdEncoding.EncodeToString(salt)

	return hashBase64, saltBase64, nil
}

func VerifyPassword(password, encodedHash, encodedSalt string, config *Argon2Config) (bool, error) {
	salt, err := base64.RawStdEncoding.DecodeString(encodedSalt)
	if err != nil {
		return false, fmt.Errorf("failed to decode salt: %w", err)
	}

	expectedHash, err := base64.RawStdEncoding.DecodeString(encodedHash)
	if err != nil {
		return false, fmt.Errorf("failed to decode hash: %w", err)
	}

	newHash := argon2.IDKey([]byte(password), salt, config.Time, config.Memory, config.Threads, config.KeyLength)

	if len(newHash) != len(expectedHash) {
		return false, nil
	}
	
	match := subtle.ConstantTimeCompare(newHash, expectedHash) == 1

	return match, nil
}
