package middlewares

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/natchaphonbw/usermanagement/pkg/jwt"
)

func JWTAuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")

		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return fiber.NewError(fiber.StatusUnauthorized, "Invalid or missing Authorization header")
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := jwt.VerifyAccessToken(tokenStr)
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, "Invalid token")
		}

		userID, err := uuid.Parse(claims.UserID)
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, "Invalid user ID in token")
		}

		sessionID, err := uuid.Parse(claims.SessionID)
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, "Invalid session ID in token")
		}

		c.Locals("userID", userID)
		c.Locals("sessionID", sessionID)
		c.Locals("tokenStr", tokenStr)
		return c.Next()

	}
}
