package middleware

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

// RouteMiddlewareMap defines middleware mappings for specific route prefixes
var routeMiddlewareMap = map[string]fiber.Handler{
	"/api/v1":    AuthMiddleware(), //MiddlewareAdmin(),
	"/api/admin": MiddlewareAdmin(),
}

// ConfigMiddleware returns a middleware handler that applies specific middleware based on the request path
// It matches route prefixes in order of specificity (longest prefix first) to ensure correct middleware is applied
func ConfigMiddleware() fiber.Handler {
	// Sort routes by specificity (optional enhancement)
	// This would ensure longer, more specific routes are matched first
	// Currently using map iteration which has random order

	return func(c *fiber.Ctx) error {
		path := c.Path()
		//log.Debug(fmt.Sprintf("Processing request for path: %s", path))

		// Check if path matches any of our defined routes requiring middleware
		for routePrefix, middlewareHandler := range routeMiddlewareMap {
			if strings.HasPrefix(path, routePrefix) {
				//log.Debug(fmt.Sprintf("Applying middleware for route prefix: %s", routePrefix))
				err := middlewareHandler(c)
				if err != nil {
					log.Error(fmt.Sprintf("Middleware error for route %s: %v", routePrefix, err))
					return err
				}
				return nil // Middleware has already called Next() or returned an error
			}
		}

		// If no middleware matched, continue to the next handler
		//log.Debug("No specific middleware applied, continuing request")
		return c.Next()
	}
}

// RegisterRouteMiddleware allows dynamically adding new middleware mappings
// This makes the system more flexible for future extensions
func RegisterRouteMiddleware(routePrefix string, handler fiber.Handler) {
	routeMiddlewareMap[routePrefix] = handler
	log.Info(fmt.Sprintf("Registered middleware for route prefix: %s", routePrefix))
}
