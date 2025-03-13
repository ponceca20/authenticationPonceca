package middleware

import (
	"github.com/gofiber/fiber/v2"
)

// SetupMiddleware demonstrates how to set up and use the middleware system
func SetupMiddleware(app *fiber.App) {
	// Apply central middleware configuration to all routes
	app.Use(ConfigMiddleware())

	// Example of dynamically registering a new middleware handler
	RegisterRouteMiddleware("/api/reports", createCustomMiddleware())
}

// createCustomMiddleware is an example of creating a custom middleware handler
func createCustomMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Example middleware logic
		c.Set("X-Custom-Header", "CustomValue")
		return c.Next()
	}
}
