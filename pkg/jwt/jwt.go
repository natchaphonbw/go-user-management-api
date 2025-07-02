package jwt

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/matthewhartstonge/argon2"
)

var (
	accessSecretKey  = []byte(os.Getenv("JWT_SECRET"))
	refreshSecretKey = []byte(os.Getenv("JWT_REFRESH_SECRET"))
	accessTokenTTL   = 24 * time.Hour
	refreshTokenTTL  = 7 * 24 * time.Hour
)

type Claims struct {
	UserID    string `json:"user_id"`
	SessionID string `json:"session_id"`
	jwt.RegisteredClaims
}

func GenerateAccessToken(userID, sessionID string) (string, error) {
	claims := &Claims{
		UserID:    userID,
		SessionID: sessionID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// gen token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(accessSecretKey)
}

func GenerateRefreshToken(userID, sessionID string) (string, time.Time, time.Time, error) {
	issuedAt := time.Now()
	expiresAt := issuedAt.Add(refreshTokenTTL)

	claims := &Claims{
		UserID:    userID,
		SessionID: sessionID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(issuedAt),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(refreshSecretKey)
	if err != nil {
		return "", time.Time{}, time.Time{}, err
	}

	return signed, issuedAt, expiresAt, nil
}

func VerifyAccessToken(tokenStr string) (*Claims, error) {
	return verifyToken(tokenStr, accessSecretKey)
}

func VerifyRefreshToken(tokenStr string) (*Claims, error) {
	return verifyToken(tokenStr, refreshSecretKey)
}

func verifyToken(tokenStr string, secretKey []byte) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, err
}

func HashRefreshToken(token string) (string, error) {
	config := argon2.DefaultConfig()

	encoded, err := config.HashEncoded([]byte(token))
	if err != nil {
		return "", err
	}

	return string(encoded), nil
}

func VerifyRefreshTokenHash(rawToken, hashedToken string) (bool, error) {
	return argon2.VerifyEncoded([]byte(rawToken), []byte(hashedToken))
}
