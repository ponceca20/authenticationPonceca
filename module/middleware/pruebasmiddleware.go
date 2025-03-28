package middleware

import (
	"log"
	"strings"

	"practicev2/config"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func MiddlewareAdmin() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Next()
	}
}

// BasicTokenMiddleware realiza una validación básica del token JWT para pruebas.
func BasicTokenMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		const cookieName = "auth_token"
		// Obtener token de cookie o header
		tokenString := c.Cookies(cookieName)
		if tokenString == "" {
			authHeader := c.Get("Authorization")
			if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
				tokenString = strings.TrimPrefix(authHeader, "Bearer ")
			}
		}
		if tokenString == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token no proporcionado"})
		}

		// Obtener la clave JWT desde configuración
		jwtKey := config.GetJWTKey()

		// Validar token simple
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				log.Println("Método de firma JWT inválido:", token.Method.Alg())
				return nil, fiber.NewError(fiber.StatusUnauthorized, "Método de firma JWT inválido")
			}
			return jwtKey, nil
		})
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token inválido o expirado"})
		}

		// Validación básica exitosa, continuar con la cadena.
		return c.Next()
	}
}
