package middleware

import (
	"practicev2/module/authentication/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// Protected is a middleware that protects routes requiring a valid JWT.
// It extracts the token from the Authorization header, validates it,
// and stores the claims in the request context.
func Protected() fiber.Handler {
	jwtService := utils.NewJWTService()

	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return utils.SendError(c, fiber.StatusUnauthorized, "Missing or malformed JWT")
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return utils.SendError(c, fiber.StatusUnauthorized, "Missing or malformed JWT, format is 'Bearer <token>'")
		}
		tokenString := parts[1]

		claims, err := jwtService.ValidateToken(tokenString)
		if err != nil {
			return utils.SendError(c, fiber.StatusUnauthorized, "Invalid or expired JWT", err)
		}

		// Store claims in context for downstream handlers
		c.Locals("user_claims", claims)
		c.Locals("identity_id", claims.IdentityID)

		return c.Next()
	}
}
